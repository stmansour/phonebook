package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

func adminEditCompanyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var sess *session
	var ui uiSupport
	sess = nil
	if 0 < initHandlerSession(sess, &ui, w, r) {
		return
	}
	sess = ui.X

	// SECURITY
	if !sess.elemPermsAny(ELEMCOMPANY, PERMMOD) {
		ulog("Permissions refuse adminEditCo page on userid=%d (%s), role=%s\n", sess.UID, sess.Firstname, sess.Urole.Name)
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}

	var c company
	path := "/adminEditCo/"
	CoCodestr := r.RequestURI[len(path):]
	if len(CoCodestr) == 0 {
		fmt.Fprintf(w, "the RequestURI needs the Company Code. It was not found on the URI:  %s\n", r.RequestURI)
		return
	}
	CoCode, err := strconv.Atoi(CoCodestr)
	if err != nil {
		fmt.Fprintf(w, "Error converting Company Code to a number: %v. URI: %s\n", err, r.RequestURI)
		return
	}
	breadcrumbAdd(sess, "AdminEdit Company", fmt.Sprintf("/adminEditCo/%d", CoCode))
	getCompanyInfo(CoCode, &c)
	ui.C = &c
	ui.C.filterSecurityRead(sess, PERMMOD)

	t, _ := template.New("adminEditCo.html").Funcs(funcMap).ParseFiles("adminEditCo.html")
	err = t.Execute(w, &ui)
	if nil != err {
		errmsg := fmt.Sprintf("adminEditCompanyHandler: err = %v\n", err)
		ulog(errmsg)
		fmt.Println(errmsg)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
