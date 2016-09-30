package main

import (
	"net/http"
	"strconv"
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
		action == "add business unit" || action == "add company" || action == "stats" || action == "setup" {
		url := r.FormValue("url")
		// fmt.Printf("action = %s,  url = %s\n", action, url)
		http.Redirect(w, r, url, http.StatusFound)
	} else if action == "shutdown" {
		http.Redirect(w, r, r.FormValue("url"), http.StatusFound)
	} else if action == "restart" {
		http.Redirect(w, r, r.FormValue("url"), http.StatusFound)
	} else {
		ulog("adminViewBtnHandler: unrecognized action: %s\n", action)
		http.Redirect(w, r, "/search/", http.StatusFound)
	}
}

func popHandler(w http.ResponseWriter, r *http.Request) {
	var sess *session
	var ui uiSupport
	sess = nil
	if 0 < initHandlerSession(sess, &ui, w, r) {
		return
	}
	sess = ui.X
	path := "/pop/"
	uidstr := r.RequestURI[len(path):]
	if len(uidstr) > 0 {
		n, _ := strconv.Atoi(uidstr)
		http.Redirect(w, r, breadcrumbBack(sess, n), http.StatusFound)
	} else {
		http.Redirect(w, r, "/search/", http.StatusFound)
	}
}
