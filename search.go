package main

import (
	"fmt"
	"net/http"
	"phonebook/authz"
	"phonebook/db"
	"phonebook/sess"
	"strings"
)

func searchHandler(w http.ResponseWriter, r *http.Request) {
	var ssn *sess.Session
	var ui uiSupport
	ssn = nil

	if 0 < initHandlerSession(ssn, &ui, w, r) {
		return
	}

	Phonebook.ReqCountersMem <- 1    // ask to access the shared mem, blocks until granted
	<-Phonebook.ReqCountersMemAck    // make sure we got it
	Counters.SearchPeople++          // initialize our data
	Phonebook.ReqCountersMemAck <- 1 // tell Dispatcher we're done with the data

	ssn = ui.X
	breadcrumbReset(ssn, "Search People", "/search/")
	w.Header().Set("Content-Type", "text/html")

	var d searchResults
	var s string
	d.Query = strings.TrimSpace(r.FormValue("searchstring"))
	inclterms := "" != r.FormValue("inclterms")

	searchTerms := strings.Split(d.Query, " ")

	//===========================================================
	//  First, determine the deptcodes that match this query...
	//===========================================================
	var dca []int
	l := len(d.Query)
	for deptname, deptcode := range ui.NameToDeptCode {
		if l > 0 {
			if strings.Contains(strings.ToLower(deptname), strings.ToLower(d.Query)) {
				dca = append(dca, deptcode)
			}
		} else {
			dca = append(dca, deptcode)
		}
	}

	// Here are the major search fields
	s = "select uid,lastname,firstname,preferredname,jobcode,primaryemail,officephone,officefax,cellphone,deptcode from people where "

	// if the user has access and wants to include terminated employees...
	if !inclterms {
		s += "status>0 and "
	}

	// here are the general conditions
	s += "("
	if l > 0 {
		switch len(searchTerms) {
		case 2:
			s += fmt.Sprintf(`((firstname like "%%%s%%" and lastname like "%%%s%%") or (PreferredName like "%%%s%%" and lastname like "%%%s%%") `,
				searchTerms[0], searchTerms[1], searchTerms[0], searchTerms[1])
		case 3:
			s += fmt.Sprintf(`((firstname like "%%%s%%" and middlename like "%%%s%%" and lastname like "%%%s%%") or (PreferredName like "%%%s%%" and middlename like "%%%s%%" and lastname like "%%%s%%") `,
				searchTerms[0], searchTerms[1], searchTerms[2], searchTerms[0], searchTerms[1], searchTerms[2])
		default:
			s += fmt.Sprintf("(lastname like \"%%%s%%\" or firstname like \"%%%s%%\" or PreferredName like \"%%%s%%\" ",
				d.Query, d.Query, d.Query)

		}
		s += fmt.Sprintf("or primaryemail like \"%%%s%%\" or cellphone like \"%%%s%%\" or OfficePhone like \"%%%s%%\" or OfficeFax like \"%%%s%%\") ",
			d.Query, d.Query, d.Query, d.Query)
	}

	// include departments...
	if len(dca) > 0 {
		if l > 0 {
			s += "or "
		}
		s += "("
		for i := 0; i < len(dca); i++ {
			if i > 0 {
				s += " or "
			}
			s += fmt.Sprintf("deptcode=%d", dca[i])
		}
		s += ") "
	}
	s += fmt.Sprintf(") order by lastname,firstname LIMIT 75")

	// fmt.Printf("query = %s\n", s)
	rows, err := Phonebook.db.Query(s)
	errcheck(err)
	defer rows.Close()

	for rows.Next() {
		var m db.Person
		errcheck(rows.Scan(&m.UID, &m.LastName, &m.FirstName, &m.PreferredName, &m.JobCode, &m.PrimaryEmail, &m.OfficePhone, &m.OfficeFax, &m.CellPhone, &m.DeptCode))
		m.DeptName = getDepartmentFromDeptCode(m.DeptCode)
		pm := &m
		// func (d *person) filterSecurityRead(sess *session, permRequired int) {
		// 	filterSecurityRead(d, ELEMPERSON, sess, permRequired, d.UID)
		// }
		filterSecurityRead(pm, authz.ELEMPERSON, ssn, authz.PERMVIEW|authz.PERMMOD, m.UID)
		d.Matches = append(d.Matches, m)
	}
	errcheck(rows.Err())
	if l == 0 {
		d.Query = " "
	}
	ui.R = &d

	err = renderTemplate(w, ui, "search.html")

	if nil != err {
		errmsg := fmt.Sprintf("searchHandler: err = %v\n", err)
		ulog(errmsg)
		fmt.Println(errmsg)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
