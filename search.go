package main

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"
)

func searchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var sess *session
	var ui uiSupport
	sess = nil
	if 0 < initHandlerSession(sess, &ui, w, r) {
		return
	}
	sess = ui.X
	breadcrumbReset(sess, "Search People", "/search/")

	var d searchResults
	var s string
	d.Query = strings.TrimSpace(r.FormValue("searchstring"))
	inclterms := "" != r.FormValue("inclterms")

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
	s = "select uid,lastname,firstname,preferredname,jobcode,primaryemail,officephone,cellphone,deptcode from people where "

	// if the user has access and wants to include terminated employees...
	if !inclterms {
		s += "status>0 and "
	}

	// here are the general conditions
	s += "("
	if l > 0 {
		s += fmt.Sprintf("(lastname like \"%%%s%%\" or firstname like \"%%%s%%\" or PreferredName like \"%%%s%%\" or primaryemail like \"%%%s%%\" or cellphone like \"%%%s%%\" or OfficePhone like \"%%%s%%\" or OfficeFax like \"%%%s%%\") ",
			d.Query, d.Query, d.Query, d.Query, d.Query, d.Query, d.Query)
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
	s += fmt.Sprintf(") order by lastname,firstname")

	// fmt.Printf("query = %s\n", s)
	rows, err := Phonebook.db.Query(s)
	errcheck(err)
	defer rows.Close()

	for rows.Next() {
		var m person
		errcheck(rows.Scan(&m.UID, &m.LastName, &m.FirstName, &m.PreferredName, &m.JobCode, &m.PrimaryEmail, &m.OfficePhone, &m.CellPhone, &m.DeptCode))
		m.DeptName = getDepartmentFromDeptCode(m.DeptCode)
		pm := &m
		pm.filterSecurityRead(sess, PERMVIEW|PERMMOD)
		d.Matches = append(d.Matches, m)
	}
	errcheck(rows.Err())
	if l == 0 {
		d.Query = " "
	}
	ui.R = &d
	t, _ := template.New("search.html").Funcs(funcMap).ParseFiles("search.html")
	err = t.Execute(w, &ui)

	if nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
