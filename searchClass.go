package main

import (
	"fmt"
	"net/http"
	"text/template"
)

func searchClassHandler(w http.ResponseWriter, r *http.Request) {
	var s string
	w.Header().Set("Content-Type", "text/html")
	var sess *session
	var ui uiSupport
	sess = nil
	if 0 < initHandlerSession(sess, &ui, w, r) {
		return
	}
	sess = ui.X
	breadcrumbReset(sess, "Search Classes", "/searchcl/")

	var d searchClassResults

	d.Query = r.FormValue("searchstring")
	if len(d.Query) > 0 {
		s = "select ClassCode,Name,Designation,Description from classes where "
		s += fmt.Sprintf("Name like \"%%%s%%\" or Designation like \"%%%s%%\" or Description like \"%%%s%%\"",
			d.Query, d.Query, d.Query)
		s += fmt.Sprintf("order by Name,Designation")
		// fmt.Printf("query = %s\n", s)
	} else {
		d.Query = "  "
		s = "select ClassCode,Name,Designation,Description from classes order by Name,Designation"
	}
	rows, err := Phonebook.db.Query(s)
	errcheck(err)
	defer rows.Close()

	for rows.Next() {
		var c class
		errcheck(rows.Scan(&c.ClassCode, &c.Name, &c.Designation, &c.Description))
		d.Matches = append(d.Matches, c)
	}
	errcheck(rows.Err())

	t, _ := template.New("searchClass.html").Funcs(funcMap).ParseFiles("searchClass.html")

	ui.L = &d
	err = t.Execute(w, &ui)

	if nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
