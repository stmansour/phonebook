package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (c *class) filterSecurityRead(sess *session, permRequired int) {
	filterSecurityRead(c, ELEMCLASS, sess, permRequired, 0)
}

func getClassInfo(classcode int, c *class) {
	s := fmt.Sprintf("select classcode,Name,Designation,Description from classes where classcode=%d", classcode)
	rows, err := Phonebook.db.Query(s)
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		errcheck(rows.Scan(&c.ClassCode, &c.Name, &c.Designation, &c.Description))
	}
	errcheck(rows.Err())
}

func classHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var sess *session
	var ui uiSupport
	sess = nil
	if 0 < initHandlerSession(sess, &ui, w, r) {
		return
	}
	sess = ui.X

	// SECURITY
	if !sess.elemPermsAny(ELEMCLASS, PERMVIEW) {
		ulog("Permissions refuse class page on userid=%d (%s), role=%s\n", sess.UID, sess.Firstname, sess.Urole.Name)
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}

	var c class
	path := "/class/"
	costr := r.RequestURI[len(path):]
	if len(costr) > 0 {
		classcode, _ := strconv.Atoi(costr)
		getClassInfo(classcode, &c)
		ui.A = &c
		ui.A.filterSecurityRead(sess, PERMVIEW)
		t, _ := template.New("class.html").Funcs(funcMap).ParseFiles("class.html")
		err := t.Execute(w, &ui)
		if nil != err {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		fmt.Fprintf(w, "classcode = %s\nCould not convert to number\n", costr)
	}
}
