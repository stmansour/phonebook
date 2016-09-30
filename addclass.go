package main

import (
	"fmt"
	"net/http"
	"text/template"
)

func adminAddClassHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var sess *session
	var ui uiSupport
	sess = nil

	loadCompanies()
	if 0 < initHandlerSession(sess, &ui, w, r) {
		return
	}
	sess = ui.X

	// SECURITY
	if !sess.elemPermsAny(ELEMPERSON, PERMCREATE) {
		ulog("Permissions refuse AddClass page on userid=%d (%s), role=%s\n", sess.UID, sess.Firstname, sess.Urole.Name)
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}
	breadcrumbAdd(sess, "Add Business Unit", "/adminAddClass/")

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
