package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func delCheckError(caller string, sess *session, err error, s string, w http.ResponseWriter, r *http.Request) bool {
	if nil != err {
		ulog("%s: \"%s\"  err = %v\n", caller, s, err)
		fmt.Printf("%s: \"%s\"  err = %v\n", caller, s, err)
		http.Redirect(w, r, breadcrumbBack(sess, 2), http.StatusFound)
		return true
	}
	return false
}

func delPersonHandler(w http.ResponseWriter, r *http.Request) {
	var sess *session
	var ui uiSupport
	sess = nil
	if 0 < initHandlerSession(sess, &ui, w, r) {
		return
	}
	sess = ui.X

	// SECURITY
	if !hasAccess(sess, ELEMPERSON, "ElemEntity", PERMDEL) {
		ulog("Permissions refuse delPersonHandler page on userid=%d (%s), role=%s\n", sess.UID, sess.Firstname, sess.Urole.Name)
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}

	// var d personDetail
	c := "delPersonHandler"
	m := strings.Split(r.RequestURI, "/")
	uidstr := m[len(m)-1]
	if len(uidstr) == 0 {
		fmt.Fprintf(w, "The RequestURI needs the person's uid. It was not found on the URI:  %s\n", r.RequestURI)
		return
	}
	uid, err := strconv.Atoi(uidstr)
	if err != nil {
		fmt.Fprintf(w, "Error converting uid to a number: %v. URI: %s\n", err, r.RequestURI)
		return
	}

	//===============================================================
	// in order to delete a person, we must delete all references
	// to the person in the following database tables:
	//		deductions
	//		compensation
	//===============================================================
	s := fmt.Sprintf("DELETE FROM people WHERE UID=%d", uid)
	stmt, err := Phonebook.db.Prepare(s)
	if delCheckError(c, sess, err, s, w, r) {
		return
	}
	_, err = stmt.Exec()
	if delCheckError(c, sess, err, s, w, r) {
		return
	}

	s = fmt.Sprintf("DELETE FROM deductions WHERE UID=%d", uid)
	stmt, err = Phonebook.db.Prepare(s)
	if delCheckError(c, sess, err, s, w, r) {
		return
	}
	_, err = stmt.Exec()
	if delCheckError(c, sess, err, s, w, r) {
		return
	}

	s = fmt.Sprintf("DELETE FROM compensation WHERE UID=%d", uid)
	stmt, err = Phonebook.db.Prepare(s)
	if delCheckError(c, sess, err, s, w, r) {
		return
	}
	_, err = stmt.Exec()
	if delCheckError(c, sess, err, s, w, r) {
		return
	}
	http.Redirect(w, r, "/search/", http.StatusFound)
}

func delClassHandler(w http.ResponseWriter, r *http.Request) {
	var sess *session
	var ui uiSupport
	sess = nil
	if 0 < initHandlerSession(sess, &ui, w, r) {
		return
	}
	sess = ui.X

	// SECURITY
	if !hasAccess(sess, ELEMCLASS, "ElemEntity", PERMDEL) {
		ulog("Permissions refuse delCoHandler page on userid=%d (%s), role=%s\n", sess.UID, sess.Firstname, sess.Urole.Name)
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}

	c := "delClassHandler"
	m := strings.Split(r.RequestURI, "/")
	uidstr := m[len(m)-1]
	if len(uidstr) == 0 {
		fmt.Fprintf(w, "The RequestURI needs the class code. It was not found on the URI:  %s\n", r.RequestURI)
		return
	}
	ClassCode, err := strconv.Atoi(uidstr)
	if err != nil {
		fmt.Fprintf(w, "Error converting class code to a number: %v. URI: %s\n", err, r.RequestURI)
		return
	}

	//===============================================================
	// in order to delete a person, we must delete all references
	// to the person in the following database tables:
	//		deductions
	//		compensationc,
	//===============================================================
	s := fmt.Sprintf("DELETE FROM classes WHERE ClassCode=%d", ClassCode)
	stmt, err := Phonebook.db.Prepare(s)
	if delCheckError(c, sess, err, s, w, r) {
		return
	}
	_, err = stmt.Exec()
	if delCheckError(c, sess, err, s, w, r) {
		return
	}
	// we've deleted it, now we need to reload our class list...
	loadClasses()
	http.Redirect(w, r, "/searchcl/", http.StatusFound)
}

func delCoHandler(w http.ResponseWriter, r *http.Request) {
	var sess *session
	var ui uiSupport
	sess = nil
	if 0 < initHandlerSession(sess, &ui, w, r) {
		return
	}
	sess = ui.X

	// SECURITY
	if !hasAccess(sess, ELEMCOMPANY, "ElemEntity", PERMDEL) {
		ulog("Permissions refuse delCoHandler page on userid=%d (%s), role=%s\n", sess.UID, sess.Firstname, sess.Urole.Name)
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}

	c := "delCoHandler"
	m := strings.Split(r.RequestURI, "/")
	uidstr := m[len(m)-1]
	if len(uidstr) == 0 {
		fmt.Fprintf(w, "The RequestURI needs the company code. It was not found on the URI:  %s\n", r.RequestURI)
		return
	}
	CoCode, err := strconv.Atoi(uidstr)
	if err != nil {
		fmt.Fprintf(w, "Error converting company code to a number: %v. URI: %s\n", err, r.RequestURI)
		return
	}

	//===============================================================
	// in order to delete a person, we must delete all references
	// to the person in the following database tables:
	//		deductions
	//		compensation
	//===============================================================
	s := fmt.Sprintf("DELETE FROM companies WHERE CoCode=%d", CoCode)
	stmt, err := Phonebook.db.Prepare(s)
	if delCheckError(c, sess, err, s, w, r) {
		return
	}
	_, err = stmt.Exec()
	if delCheckError(c, sess, err, s, w, r) {
		return
	}
	// we've deleted it, now we need to reload our company list...
	loadCompanies()
	http.Redirect(w, r, "/searchco/", http.StatusFound)
}
