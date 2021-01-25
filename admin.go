package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"phonebook/db"
	"time"
)

func adminHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var ssn *db.Session
	var ui uiSupport
	ssn = nil
	if 0 < initHandlerSession(ssn, &ui, w, r) {
		return
	}
	ssn = ui.X
	breadcrumbReset(ssn, "Admin", "/admin/")

	// SECURITY
	if !(ssn.ElemPermsAny(db.ELEMPERSON, db.PERMCREATE) ||
		ssn.ElemPermsAny(db.ELEMCOMPANY, db.PERMCREATE) ||
		ssn.ElemPermsAny(db.ELEMCLASS, db.PERMCREATE) ||
		ssn.ElemPermsAny(db.ELEMPBSVC, db.PERMEXEC)) {
		ulog("Permissions refuse admin page on userid=%d (%s), role=%s\n", ssn.UID, ssn.Firstname, ssn.PMap.Urole.Name)
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}

	err := renderTemplate(w, ui, "admin.html")

	if nil != err {
		errmsg := fmt.Sprintf("adminHandler: err = %v\n", err)
		ulog(errmsg)
		fmt.Println(errmsg)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func restartHandler(w http.ResponseWriter, r *http.Request) {
	var ssn *db.Session
	var ui uiSupport
	ssn = nil
	if 0 < initHandlerSession(ssn, &ui, w, r) {
		return
	}
	ssn = ui.X

	perm, ok := ssn.PMap.Ppr["Restart"]
	// fmt.Printf("restartHandler: perm=0x%02x\n", perm)
	if ok {
		if perm&db.PERMEXEC != 0 {
			ulog("restart invoked by UID %d, %s\n", ssn.UID, ssn.Username)
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
	var ssn *db.Session
	var ui uiSupport
	ssn = nil
	if 0 < initHandlerSession(ssn, &ui, w, r) {
		return
	}
	ssn = ui.X

	perm, ok := ssn.PMap.Ppr["Shutdown"]
	if ok {
		if perm&db.PERMEXEC != 0 {
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
