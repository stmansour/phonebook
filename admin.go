package main

import (
	"fmt"
	"net/http"
	"os"
	"text/template"
	"time"
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
	breadcrumbReset(sess, "Admin", "/admin/")

	// SECURITY
	if !(sess.elemPermsAny(ELEMPERSON, PERMCREATE) ||
		sess.elemPermsAny(ELEMCOMPANY, PERMCREATE) ||
		sess.elemPermsAny(ELEMCLASS, PERMCREATE) ||
		sess.elemPermsAny(ELEMPBSVC, PERMEXEC)) {
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

func shutdownHandler(w http.ResponseWriter, r *http.Request) {
	var sess *session
	var ui uiSupport
	sess = nil
	if 0 < initHandlerSession(sess, &ui, w, r) {
		return
	}
	sess = ui.X

	perm, ok := sess.Ppr["Shutdown"]
	if ok {
		if perm&PERMEXEC != 0 {
			fmt.Fprintf(w, "<html><body><h1>shutting down in 5 seconds!</h1></body></html>")
			ulog("Shutdown initiated from web service\n")
			go func() {
				time.Sleep(time.Duration(5 * time.Second))
				os.Exit(0)
			}()
			return
		}
	}
	http.Redirect(w, r, "/search/", http.StatusFound)
}
