package main

import (
	"fmt"
	"net/http"
	"phonebook/authz"
	"phonebook/sess"
	"strconv"
	"text/template"
)

func adminEditHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var ssn *sess.Session
	var ui uiSupport
	ssn = nil
	if 0 < initHandlerSession(ssn, &ui, w, r) {
		return
	}
	ssn = ui.X

	var d personDetail
	d.Reports = make([]person, 0)
	path := "/adminEdit/"
	uidstr := r.RequestURI[len(path):]
	if len(uidstr) == 0 {
		fmt.Fprintf(w, "the RequestURI needs to know the person's uid. It was not found on the URI:  %s\n", r.RequestURI)
		return
	}
	uid, err := strconv.Atoi(uidstr)
	if err != nil {
		fmt.Fprintf(w, "Error converting uid to a number: %v. URI: %s\n", err, r.RequestURI)
		return
	}
	d.UID = uid

	adminReadDetails(&d)
	breadcrumbAdd(ssn, "AdminEdit Person", fmt.Sprintf("/adminEdit/%d", uid))

	//---------------------------------------------------------------------
	// SECURITY
	//		Access to the screen requires authz.PERMMOD permission.  The data
	//		in personDetail includes those fields with VIEW and MOD perms
	//---------------------------------------------------------------------
	if !ssn.ElemPermsAny(authz.ELEMPERSON, authz.PERMMOD) {
		fmt.Printf("adminEditHandler:  ssn.ElemPermsAny(authz.ELEMPERSON, authz.PERMVIEW|authz.PERMMOD) returned 0\n")
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}
	d.filterSecurityRead(ssn, authz.PERMVIEW|authz.PERMMOD)
	ui.D = &d
	t, _ := template.New("adminEdit.html").Funcs(funcMap).ParseFiles("adminEdit.html")
	initUIData(&ui)
	// fmt.Printf("AdminEditHandler: d = %#v\n", d)
	err = t.Execute(w, &ui)
	if nil != err {
		errmsg := fmt.Sprintf("adminEditHandler: err = %v\n", err)
		ulog(errmsg)
		fmt.Println(errmsg)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
