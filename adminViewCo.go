package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

func adminViewCompanyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("entered adminViewCompanyHandler\n")
	w.Header().Set("Content-Type", "text/html")
	var sess *session
	var ui uiSupport
	sess = nil
	if 0 < initHandlerSession(sess, &ui, w, r) {
		return
	}
	sess = ui.X

	// SECURITY
	if !sess.elemPermsAny(ELEMCOMPANY, PERMVIEW|PERMMOD) {
		ulog("Permissions refuse adminViewCo page on userid=%d (%s), role=%s\n", sess.UID, sess.Firstname, sess.Urole.Name)
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}

	var c company
	path := "/adminViewCo/"
	CoCodeStr := r.RequestURI[len(path):]
	if len(CoCodeStr) == 0 {
		fmt.Fprintf(w, "the RequestURI needs to know the Company Code. It was not found on the URI:  %s\n", r.RequestURI)
		return
	}
	CoCode, err := strconv.Atoi(CoCodeStr)
	if err != nil {
		fmt.Fprintf(w, "Error converting Company Code to a number: %v. URI: %s\n", err, r.RequestURI)
		return
	}
	getCompanyInfo(CoCode, &c)

	fmt.Printf("adminViewCompanyHandler:  A\n")
	funcMap := template.FuncMap{
		"compToString":      compensationTypeToString,
		"acceptIntToString": acceptIntToString,
		"dateToString":      dateToString,
		"dateYear":          dateYear,
		"monthStringToInt":  monthStringToInt,
		"add":               add,
		"sub":               sub,
		"rmd":               rmd,
		"mul":               mul,
		"div":               div,
	}

	t, _ := template.New("adminViewCo.html").Funcs(funcMap).ParseFiles("adminViewCo.html")
	ui.C = &c
	fmt.Printf("adminViewCompanyHandler:  Z\n")
	err = t.Execute(w, &ui)
	if nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
