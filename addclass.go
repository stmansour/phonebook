package main

import (
	"fmt"
	"net/http"
	"phonebook/authz"
	"phonebook/sess"
	"text/template"
)

func adminAddClassHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var ssn *sess.Session
	var ui uiSupport
	ssn = nil

	loadCompanies()
	if 0 < initHandlerSession(ssn, &ui, w, r) {
		return
	}
	ssn = ui.X

	// SECURITY
	if !ssn.ElemPermsAny(authz.ELEMPERSON, authz.PERMCREATE) {
		ulog("Permissions refuse AddClass page on userid=%d (%s), role=%s\n", ssn.UID, ssn.Firstname, ssn.Urole.Name)
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}
	breadcrumbAdd(ssn, "Add Business Unit", "/adminAddClass/")

	var c class
	c.ClassCode = 0
	c.Designation = ""
	c.Description = ""

	t, _ := template.New("adminEditClass.html").Funcs(funcMap).ParseFiles("adminEditClass.html")

	ui.A = &c
	ui.CompanyList = PhonebookUI.CompanyList
	err := t.Execute(w, &ui)
	if nil != err {
		errmsg := fmt.Sprintf("adminAddClassHandler: err = %v\n", err)
		ulog(errmsg)
		fmt.Println(errmsg)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
