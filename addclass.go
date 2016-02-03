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
	breadcrumbAdd(sess, "Add Class", "/adminAddClass/")

	var c class
	c.ClassCode = 0
	c.Designation = ""
	c.Description = ""

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

	t, _ := template.New("adminEditClass.html").Funcs(funcMap).ParseFiles("adminEditClass.html")

	ui.A = &c
	err := t.Execute(w, &ui)
	if nil != err {
		errmsg := fmt.Sprintf("adminAddClassHandler: err = %v\n", err)
		ulog(errmsg)
		fmt.Println(errmsg)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
