package main

import (
	"fmt"
	"net/http"
	"text/template"
)

func adminAddCompanyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var sess *session
	var ui uiSupport
	sess = nil
	if 0 < initHandlerSession(sess, &ui, w, r) {
		return
	}
	sess = ui.X

	// SECURITY
	if !sess.elemPermsAny(ELEMCOMPANY, PERMCREATE) {
		ulog("Permissions refuse AddCompany page on userid=%d (%s), role=%s\n", sess.UID, sess.Firstname, sess.Urole.Name)
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}
	breadcrumbAdd(sess, "Add Company", "/adminAddCompany/")

	var c company
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

	t, _ := template.New("adminEditCo.html").Funcs(funcMap).ParseFiles("adminEditCo.html")

	ui.C = &c
	err := t.Execute(w, &ui)
	if nil != err {
		errmsg := fmt.Sprintf("adminAddCompanyHandler: err = %v\n", err)
		ulog(errmsg)
		fmt.Println(errmsg)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
