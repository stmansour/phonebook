package main

import (
	"fmt"
	"net/http"
	"phonebook/db"
	"phonebook/sess"
	"text/template"
)

func searchClassHandler(w http.ResponseWriter, r *http.Request) {
	var s string
	w.Header().Set("Content-Type", "text/html")
	var ssn *sess.Session
	var ui uiSupport
	ssn = nil
	if 0 < initHandlerSession(ssn, &ui, w, r) {
		return
	}
	ssn = ui.X
	breadcrumbReset(ssn, "Search Business Units", "/searchcl/")
	Phonebook.ReqCountersMem <- 1    // ask to access the shared mem, blocks until granted
	<-Phonebook.ReqCountersMemAck    // make sure we got it
	Counters.SearchClasses++         // initialize our data
	Phonebook.ReqCountersMemAck <- 1 // tell Dispatcher we're done with the data

	var d searchClassResults

	d.Query = r.FormValue("searchstring")
	if len(d.Query) > 0 {
		s = "select ClassCode,Name,Designation,Description from classes where "
		s += fmt.Sprintf("Name like \"%%%s%%\" or Designation like \"%%%s%%\" or Description like \"%%%s%%\"",
			d.Query, d.Query, d.Query)
		s += fmt.Sprintf("order by Designation")
		// fmt.Printf("query = %s\n", s)
	} else {
		d.Query = "  "
		s = "select ClassCode,Name,Designation,Description from classes order by Designation"
	}
	rows, err := Phonebook.db.Query(s)
	errcheck(err)
	defer rows.Close()

	for rows.Next() {
		var c db.Class
		errcheck(rows.Scan(&c.ClassCode, &c.Name, &c.Designation, &c.Description))
		d.Matches = append(d.Matches, c)
	}
	errcheck(rows.Err())

	t, _ := template.New("searchClass.html").Funcs(funcMap).ParseFiles("searchClass.html")

	ui.L = &d
	err = t.Execute(w, &ui)

	if nil != err {
		errmsg := fmt.Sprintf("searchClassHandler: err = %v\n", err)
		ulog(errmsg)
		fmt.Println(errmsg)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
