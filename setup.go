package main

import (
	"fmt"
	"net/http"
	"text/template"
)

func setupHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("entered setup handler\n")
	w.Header().Set("Content-Type", "text/html")
	var sess *session
	var ui uiSupport
	sess = nil
	if 0 < initHandlerSession(sess, &ui, w, r) {
		return
	}
	sess = ui.X
	breadcrumbReset(sess, "Setup", "/setup/")

	// SECURITY
	if !(sess.elemPermsAny(ELEMPBSVC, PERMEXEC)) {
		ulog("Permissions refuse setup page on userid=%d (%s), role=%s\n", sess.UID, sess.Firstname, sess.Urole.Name)
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}

	t, _ := template.New("setup.html").Funcs(funcMap).ParseFiles("setup.html")

	err := t.Execute(w, &ui)
	if nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
