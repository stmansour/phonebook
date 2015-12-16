package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func companyInit(c *company) {
	c.LegalName = ""
	c.CommonName = ""
	c.Address = ""
	c.Address2 = ""
	c.City = ""
	c.State = ""
	c.PostalCode = ""
	c.Country = ""
	c.Phone = ""
	c.Fax = ""
	c.Email = ""
	c.Designation = ""
	c.Active = 0
}

func (c *company) filterSecurityRead(sess *session, permRequired int) {
	filterSecurityRead(c, ELEMCOMPANY, sess, permRequired, 0)
}

func getCompanyInfo(cocode int, c *company) {
	s := fmt.Sprintf("select cocode,LegalName,CommonName,Address,Address2,City,State,PostalCode,Country,Phone,Fax,Email,Designation,Active,EmploysPersonnel from companies where cocode=%d", cocode)
	rows, err := Phonebook.db.Query(s)
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		errcheck(rows.Scan(&c.CoCode, &c.LegalName, &c.CommonName, &c.Address, &c.Address2, &c.City, &c.State, &c.PostalCode, &c.Country, &c.Phone, &c.Fax, &c.Email, &c.Designation, &c.Active, &c.EmploysPersonnel))
	}
	errcheck(rows.Err())
}

func companyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var sess *session
	var ui uiSupport
	sess = nil
	if 0 < initHandlerSession(sess, &ui, w, r) {
		return
	}
	sess = ui.X
	Phonebook.ReqCountersMem <- 1    // ask to access the shared mem, blocks until granted
	<-Phonebook.ReqCountersMemAck    // make sure we got it
	Counters.ViewCompany++           // initialize our data
	Phonebook.ReqCountersMemAck <- 1 // tell Dispatcher we're done with the data

	// SECURITY
	if !sess.elemPermsAny(ELEMCOMPANY, PERMVIEW) {
		ulog("Permissions refuse company view page on userid=%d (%s), role=%s\n", sess.UID, sess.Firstname, sess.Urole.Name)
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}

	var c company
	path := "/company/"
	costr := r.RequestURI[len(path):]
	if len(costr) > 0 {
		cocode, _ := strconv.Atoi(costr)
		breadcrumbAdd(sess, "Company", fmt.Sprintf("/company/%d", cocode))
		getCompanyInfo(cocode, &c)
		t, _ := template.New("company.html").Funcs(funcMap).ParseFiles("company.html")
		ui.C = &c
		ui.C.filterSecurityRead(sess, PERMVIEW)
		err := t.Execute(w, &ui)
		if nil != err {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		fmt.Fprintf(w, "cocode = %s\nCould not convert to number\n", costr)
	}
}
