package main

import (
	"fmt"
	"net/http"
	"phonebook/authz"
	"phonebook/db"
	"phonebook/sess"
	"strconv"
)

// func (c *db.Class) filterSecurityRead(ssn *sess.Session, permRequired int) {
// 	filterSecurityRead(c, authz.ELEMCLASS, ssn, permRequired, 0)
// }

func getClassInfo(classcode int, c *db.Class) {
	Phonebook.ReqCountersMem <- 1    // ask to access the shared mem, blocks until granted
	<-Phonebook.ReqCountersMemAck    // make sure we got it
	Counters.ViewClass++             // initialize our data
	Phonebook.ReqCountersMemAck <- 1 // tell Dispatcher we're done with the data
	// s := fmt.Sprintf("select classcode,Name,Designation,Description from classes where classcode=%d", classcode)
	rows, err := Phonebook.prepstmt.classInfo.Query(classcode)
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		errcheck(rows.Scan(&c.ClassCode, &c.CoCode, &c.Name, &c.Designation, &c.Description))
	}
	errcheck(rows.Err())

	if c.CoCode > 0 {
		errcheck(Phonebook.prepstmt.companyInfo.QueryRow(c.CoCode).Scan(&c.C.CoCode, &c.C.LegalName, &c.C.CommonName, &c.C.Address, &c.C.Address2, &c.C.City, &c.C.State, &c.C.PostalCode, &c.C.Country, &c.C.Phone, &c.C.Fax, &c.C.Email, &c.C.Designation, &c.C.Active, &c.C.EmploysPersonnel))
	}
}

func classHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var ssn *sess.Session
	var ui uiSupport
	ssn = nil
	if 0 < initHandlerSession(ssn, &ui, w, r) {
		return
	}
	ssn = ui.X

	// SECURITY
	if !ssn.ElemPermsAny(authz.ELEMCLASS, authz.PERMVIEW) {
		ulog("Permissions refuse class page on userid=%d (%s), role=%s\n", ssn.UID, ssn.Firstname, ssn.PMap.Urole.Name)
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}

	var c db.Class
	path := "/class/"
	costr := r.RequestURI[len(path):]
	if len(costr) > 0 {
		classcode, _ := strconv.Atoi(costr)
		breadcrumbAdd(ssn, "Class", fmt.Sprintf("/class/%d", classcode))
		getClassInfo(classcode, &c)
		ui.A = &c
		filterSecurityRead(ui.A, authz.ELEMCLASS, ssn, authz.PERMVIEW, 0)

		err := renderTemplate(w, ui, "class.html")

		if nil != err {
			errmsg := fmt.Sprintf("classHandler: err = %v\n", err)
			ulog(errmsg)
			fmt.Println(errmsg)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		fmt.Fprintf(w, "classcode = %s\nCould not convert to number\n", costr)
	}
}
