package main

import (
	"net/http"
	"text/template"
)

func adminHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var sess *session
	var ui uiSupport
	sess = nil
	if 0 < initHandlerSession(sess, &ui, w, r) {
		return
	}
	sess = ui.X

	// SECURITY
	if !(sess.elemPermsAny(ELEMPERSON, PERMCREATE) ||
		sess.elemPermsAny(ELEMCOMPANY, PERMCREATE) ||
		sess.elemPermsAny(ELEMCLASS, PERMCREATE)) {
		ulog("Permissions refuse admin page on userid=%d (%s), role=%s\n", sess.UID, sess.Firstname, sess.Urole.Name)
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}

	t, _ := template.New("admin.html").Funcs(funcMap).ParseFiles("admin.html")

	err := t.Execute(w, &ui)
	if nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
