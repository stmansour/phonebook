package main

import (
	"fmt"
	"net/http"
	"phonebook/authz"
	"phonebook/sess"
	"strconv"
	"text/template"
)

func adminEditClassHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var ssn *sess.Session
	var ui uiSupport
	ssn = nil
	if 0 < initHandlerSession(ssn, &ui, w, r) {
		return
	}
	ssn = ui.X

	// SECURITY
	if !ssn.ElemPermsAny(authz.ELEMCLASS, authz.PERMMOD) {
		ulog("Permissions refuse adminEditClass page on userid=%d (%s), role=%s\n", ssn.UID, ssn.Firstname, ssn.Urole.Name)
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}

	var d class
	path := "/adminEditClass/"
	uidstr := r.RequestURI[len(path):]
	if len(uidstr) == 0 {
		fmt.Fprintf(w, "the RequestURI needs to know the classcode. It was not found on the URI:  %s\n", r.RequestURI)
		return
	}
	ClassCode, err := strconv.Atoi(uidstr)
	if err != nil {
		fmt.Fprintf(w, "Error converting classcode to a number: %v. URI: %s\n", err, r.RequestURI)
		return
	}
	breadcrumbAdd(ssn, "AdminEdit Class", fmt.Sprintf("/adminEditClass/%d", ClassCode))
	d.ClassCode = ClassCode
	getClassInfo(ClassCode, &d)
	ui.A = &d
	ui.A.filterSecurityRead(ssn, authz.PERMVIEW|authz.PERMMOD)

	t, _ := template.New("adminEditClass.html").Funcs(funcMap).ParseFiles("adminEditClass.html")
	initUIData(&ui)

	// this interface needs the complete list of companies
	for i := 0; i < len(PhonebookUI.CompanyList); i++ {
		ui.CompanyList = append(ui.CompanyList, PhonebookUI.CompanyList[i])
	}

	//fmt.Printf("ui.A = %#v\n", ui.A)
	err = t.Execute(w, &ui)
	if nil != err {
		errmsg := fmt.Sprintf("adminEditClassHandler: err = %v\n", err)
		ulog(errmsg)
		fmt.Println(errmsg)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
