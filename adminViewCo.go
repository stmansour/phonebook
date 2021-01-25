package main

import (
	"fmt"
	"net/http"
	"phonebook/db"
	"strconv"
)

func adminViewCompanyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("entered adminViewCompanyHandler\n")
	w.Header().Set("Content-Type", "text/html")
	var ssn *db.Session
	var ui uiSupport
	ssn = nil
	if 0 < initHandlerSession(ssn, &ui, w, r) {
		return
	}
	ssn = ui.X
	Phonebook.ReqCountersMem <- 1    // ask to access the shared mem, blocks until granted
	<-Phonebook.ReqCountersMemAck    // make sure we got it
	Counters.ViewCompany++           // initialize our data
	Phonebook.ReqCountersMemAck <- 1 // tell Dispatcher we're done with the data

	// SECURITY
	if !ssn.ElemPermsAny(db.ELEMCOMPANY, db.PERMVIEW|db.PERMMOD) {
		ulog("Permissions refuse adminViewCo page on userid=%d (%s), role=%s\n", ssn.UID, ssn.Firstname, ssn.PMap.Urole.Name)
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}

	var c db.Company
	path := "/adminViewCo/"
	CoCodeStr := r.RequestURI[len(path):]
	if len(CoCodeStr) == 0 {
		fmt.Fprintf(w, "the RequestURI needs to know the Company Code. It was not found on the URI:  %s\n", r.RequestURI)
		return
	}
	CoCode, err := strconv.ParseInt(CoCodeStr, 10, 64)
	if err != nil {
		fmt.Fprintf(w, "Error converting Company Code to a number: %v. URI: %s\n", err, r.RequestURI)
		return
	}
	breadcrumbAdd(ssn, "AdminView Company", fmt.Sprintf("/adminViewCo/%d", CoCode))
	getCompanyInfo(CoCode, &c)
	ui.C = &c
	filterSecurityRead(ui.C, db.ELEMCOMPANY, ssn, db.PERMVIEW, 0)

	err = renderTemplate(w, ui, "adminViewCo.html")

	if nil != err {
		errmsg := fmt.Sprintf("adminViewCompanyHandler: err = %v\n", err)
		ulog(errmsg)
		fmt.Println(errmsg)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
