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
	Phonebook.ReqCountersMem <- 1    // ask to access the shared mem, blocks until granted
	<-Phonebook.ReqCountersMemAck    // make sure we got it
	Counters.ViewClass++             // initialize our data
	Phonebook.ReqCountersMemAck <- 1 // tell Dispatcher we're done with the data
	// s := fmt.Sprintf("select classcode,Name,Designation,Description from classes where classcode=%d", classcode)
	rows, err := Phonebook.prepstmt.classInfo.Query(classcode)
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		errcheck(rows.Scan(&c.ClassCode, &c.CoCode, &c.Name, &c.Designation, &c.Description))
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
		breadcrumbAdd(sess, "Class", fmt.Sprintf("/class/%d", classcode))
		getClassInfo(classcode, &c)
		ui.A = &c
		ui.A.filterSecurityRead(sess, PERMVIEW)
		t, _ := template.New("class.html").Funcs(funcMap).ParseFiles("class.html")
		err := t.Execute(w, &ui)
		if nil != err {
			errmsg := fmt.Sprintf("classHandler: err = %v\n", err)
			ulog(errmsg)
			fmt.Println(errmsg)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		fmt.Fprintf(w, "classcode = %s\nCould not convert to number\n", costr)
	}
}
