package main

import (
	"fmt"
	"net/http"
	"text/template"
)

func searchCompaniesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var sess *session
	var ui uiSupport
	sess = nil
	if 0 < initHandlerSession(sess, &ui, w, r) {
		return
	}
	sess = ui.X

	var s string
	var d searchCoResults

	d.Query = r.FormValue("searchstring")
	if len(d.Query) > 0 {
		s = "select CoCode,LegalName,CommonName,Phone,Fax,Email,Designation from companies where "
		s += fmt.Sprintf("LegalName like \"%%%s%%\" or CommonName like \"%%%s%%\" or Phone like \"%%%s%%\" or Fax like \"%%%s%%\" or email like \"%%%s%%\" or designation like \"%%%s%%\" ",
			d.Query, d.Query, d.Query, d.Query, d.Query, d.Query)
		s += fmt.Sprintf("order by CommonName,LegalName")
		// fmt.Printf("query = %s\n", s)
	} else {
		s = "select CoCode,LegalName,CommonName,Phone,Fax,Email,Designation from companies order by CommonName,LegalName"
		d.Query = " "
	}
	rows, err := Phonebook.db.Query(s)
	errcheck(err)
	defer rows.Close()

	for rows.Next() {
		var c company
		errcheck(rows.Scan(&c.CoCode, &c.LegalName, &c.CommonName, &c.Phone, &c.Fax, &c.Email, &c.Designation))
		d.Matches = append(d.Matches, c)
	}
	errcheck(rows.Err())
	t, _ := template.New("searchco.html").ParseFiles("searchco.html")
	ui.T = &d
	err = t.Execute(w, &ui)

	if nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
