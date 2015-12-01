package main

import (
	"net/http"
	"strings"
)

func adminViewBtnHandler(w http.ResponseWriter, r *http.Request) {
	var sess *session
	var ui uiSupport
	sess = nil
	if 0 < initHandlerSession(sess, &ui, w, r) {
		return
	}
	sess = ui.X

	action := strings.ToLower(r.FormValue("action"))
	// fmt.Printf("action = %s\n", action)

	if action == "done" {
		s := breadcrumbBack(sess, 2)
		// fmt.Printf("breadcrumbBack redirects to: %s\n", s)
		http.Redirect(w, r, s, http.StatusFound)
	} else if action == "adminedit" || action == "adminview" || action == "add person" ||
		action == "add class" || action == "add company" {
		url := r.FormValue("url")
		// fmt.Printf("action = %s,  url = %s\n", action, url)
		http.Redirect(w, r, url, http.StatusFound)
	} else {
		ulog("adminViewBtnHandler: unrecognized action: %s\n", action)
		http.Redirect(w, r, "/search/", http.StatusFound)
	}
}
