package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
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

func restartHandler(w http.ResponseWriter, r *http.Request) {
	var sess *session
	var ui uiSupport
	sess = nil
	if 0 < initHandlerSession(sess, &ui, w, r) {
		return
	}
	sess = ui.X

	perm, ok := sess.Ppr["Restart"]
	// fmt.Printf("restartHandler: perm=0x%02x\n", perm)
	if ok {
		if perm&PERMEXEC != 0 {
			ulog("restart invoked by UID %d, %s\n", sess.UID, sess.Username)
			cmd := "restart"
			out, err := exec.Command("./activate.sh", cmd).Output()
			if err != nil {
				ulog("Error executing 'activate.sh restart' = %v\n", err)
				http.Redirect(w, r, "/search/", http.StatusFound)
				return
			}
			ulog("Output from './activate.sh restart':\n%s\n", string(out))
			return
		}
	}
	http.Redirect(w, r, "/search/", http.StatusFound)
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
			ulog("shutdownHandler successfully invoked\n")
			extAdminShutdown(w, r)
			return
		}
	}
	http.Redirect(w, r, "/search/", http.StatusFound)
}

func extAdminShutdown(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<html><body><h1>shutting down in 5 seconds!</h1></body></html>")
	ulog("Shutdown initiated\n")
	go func() {
		time.Sleep(time.Duration(5 * time.Second))
		os.Exit(0)
	}()
	return
}
