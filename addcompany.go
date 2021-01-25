package main

import (
	"fmt"
	"net/http"
	"phonebook/db"
)

func adminAddCompanyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var ssn *db.Session
	var ui uiSupport
	ssn = nil
	if 0 < initHandlerSession(ssn, &ui, w, r) {
		return
	}
	ssn = ui.X

	// SECURITY
	if !ssn.ElemPermsAny(db.ELEMCOMPANY, db.PERMCREATE) {
		ulog("Permissions refuse AddCompany page on userid=%d (%s), role=%s\n", ssn.UID, ssn.Firstname, ssn.PMap.Urole.Name)
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}
	breadcrumbAdd(ssn, "Add Company", "/adminAddCompany/")

	var c db.Company
	c.CoCode = 0
	c.Active = YES
	c.EmploysPersonnel = NO
	c.LegalName = ""
	c.CommonName = ""
	c.Address = ""
	c.Address2 = ""
	c.City = ""
	c.State = ""
	c.PostalCode = ""
	c.Country = "USA"
	c.Phone = ""
	c.Fax = ""
	c.Email = ""
	c.Designation = ""

	// funcMap := template.FuncMap{
	// 	"compToString":      compensationTypeToString,
	// 	"acceptIntToString": acceptIntToString,
	// 	"dateToString":      dateToString,
	// 	"dateYear":          dateYear,
	// 	"monthStringToInt":  monthStringToInt,
	// 	"add":               add,
	// 	"sub":               sub,
	// 	"rmd":               rmd,
	// 	"mul":               mul,
	// 	"div":               div,
	// }

	ui.C = &c
	err := renderTemplate(w, ui, "adminEditCo.html")
	if nil != err {
		errmsg := fmt.Sprintf("adminAddCompanyHandler: err = %v\n", err)
		ulog(errmsg)
		fmt.Println(errmsg)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
