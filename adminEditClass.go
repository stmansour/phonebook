package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

func adminEditClassHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var sess *session
	var ui uiSupport
	sess = nil
	if 0 < initHandlerSession(sess, &ui, w, r) {
		return
	}
	sess = ui.X

	// SECURITY
	if !sess.elemPermsAny(ELEMCLASS, PERMMOD) {
		ulog("Permissions refuse adminEditClass page on userid=%d (%s), role=%s\n", sess.UID, sess.Firstname, sess.Urole.Name)
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}

	var d class
	path := "/adminEditClass/"
	uidstr := r.RequestURI[len(path):]
	if len(uidstr) == 0 {
		fmt.Fprintf(w, "the RequestURI needs to know the classcode. It was not found on the URI:  %s\n", r.RequestURI)
		return
	}
	ClassCode, err := strconv.Atoi(uidstr)
	if err != nil {
		fmt.Fprintf(w, "Error converting classcode to a number: %v. URI: %s\n", err, r.RequestURI)
		return
	}
	d.ClassCode = ClassCode
	getClassInfo(ClassCode, &d)
	ui.A = &d
	ui.A.filterSecurityRead(sess, PERMVIEW|PERMMOD)

	t, _ := template.New("adminEditClass.html").Funcs(funcMap).ParseFiles("adminEditClass.html")
	initUIData(&ui)
	//fmt.Printf("ui.A = %#v\n", ui.A)
	err = t.Execute(w, &ui)
	if nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
