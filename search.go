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

	var d searchResults
	d.Query = r.FormValue("searchstring")
	if len(d.Query) > 0 {
		//===========================================================
		//  First, determine the deptcodes that match this query...
		//===========================================================
		var dca []int
		for deptname, deptcode := range ui.NameToDeptCode {
			if strings.Contains(strings.ToLower(deptname), strings.ToLower(d.Query)) {
				dca = append(dca, deptcode)
			}
		}

		s := "select uid,lastname,firstname,preferredname,jobcode,primaryemail,officephone,cellphone,deptcode from people "
		s += fmt.Sprintf("where status>0 and ((lastname like \"%%%s%%\" or firstname like \"%%%s%%\" or PreferredName like \"%%%s%%\" or primaryemail like \"%%%s%%\" or cellphone like \"%%%s%%\" or OfficePhone like \"%%%s%%\" or OfficeFax like \"%%%s%%\") ",
			d.Query, d.Query, d.Query, d.Query, d.Query, d.Query, d.Query)
		if len(dca) > 0 {
			s += "or ("
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
			d.Matches = append(d.Matches, m)
		}
		errcheck(rows.Err())
	}

	ui.R = &d
	t, _ := template.New("search.html").ParseFiles("search.html")
	err := t.Execute(w, &ui)

	if nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
