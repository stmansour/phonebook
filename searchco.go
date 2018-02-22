package main

import (
	"fmt"
	"net/http"
	"phonebook/authz"
	"phonebook/db"
	"phonebook/sess"
)

func searchCompaniesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var ssn *sess.Session
	var ui uiSupport
	ssn = nil
	if 0 < initHandlerSession(ssn, &ui, w, r) {
		return
	}
	ssn = ui.X
	breadcrumbReset(ssn, "Search Companies", "/searchco/")
	Phonebook.ReqCountersMem <- 1    // ask to access the shared mem, blocks until granted
	<-Phonebook.ReqCountersMemAck    // make sure we got it
	Counters.SearchCompanies++       // initialize our data
	Phonebook.ReqCountersMemAck <- 1 // tell Dispatcher we're done with the data

	var s string
	var d searchCoResults

	d.Query = r.FormValue("searchstring")
	if len(d.Query) > 0 {
		s = "select CoCode,LegalName,CommonName,Phone,Fax,Email,Designation from companies where "
		s += fmt.Sprintf("LegalName like \"%%%s%%\" or CommonName like \"%%%s%%\" or Phone like \"%%%s%%\" or Fax like \"%%%s%%\" or email like \"%%%s%%\" or designation like \"%%%s%%\" ",
			d.Query, d.Query, d.Query, d.Query, d.Query, d.Query)
		s += fmt.Sprintf("order by Designation")
		// fmt.Printf("query = %s\n", s)
	} else {
		s = "select CoCode,LegalName,CommonName,Phone,Fax,Email,Designation from companies order by Designation"
		d.Query = " "
	}
	rows, err := Phonebook.db.Query(s)
	errcheck(err)
	defer rows.Close()

	for rows.Next() {
		var c db.Company
		errcheck(rows.Scan(&c.CoCode, &c.LegalName, &c.CommonName, &c.Phone, &c.Fax, &c.Email, &c.Designation))
		pc := &c
		// func (c *db.Company) filterSecurityRead(ssn *sess.Session, permRequired int) {
		// 	filterSecurityRead(c, authz.ELEMCOMPANY, ssn, permRequired, 0)
		// }
		filterSecurityRead(pc, authz.ELEMCOMPANY, ssn, authz.PERMVIEW, 0)
		d.Matches = append(d.Matches, c)
	}
	errcheck(rows.Err())

	ui.T = &d
	err = renderTemplate(w, ui, "searchco.html")

	if nil != err {
		errmsg := fmt.Sprintf("searchCompaniesHandler: err = %v\n", err)
		ulog(errmsg)
		fmt.Println(errmsg)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
