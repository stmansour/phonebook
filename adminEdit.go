package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

func adminEditHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var sess *session
	var ui uiSupport
	sess = nil
	if 0 < initHandlerSession(sess, &ui, w, r) {
		return
	}
	sess = ui.X

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
	breadcrumbAdd(sess, "AdminEdit Person", fmt.Sprintf("/adminEdit/%d", uid))

	//---------------------------------------------------------------------
	// SECURITY
	//		Access to the screen requires PERMMOD permission.  The data
	//		in personDetail includes those fields with VIEW and MOD perms
	//---------------------------------------------------------------------
	if !sess.elemPermsAny(ELEMPERSON, PERMMOD) {
		fmt.Printf("adminEditHandler:  sess.elemPermsAny(ELEMPERSON, PERMVIEW|PERMMOD) returned 0\n")
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}
	d.filterSecurityRead(sess, PERMVIEW|PERMMOD)
	ui.D = &d
	t, _ := template.New("adminEdit.html").Funcs(funcMap).ParseFiles("adminEdit.html")
	initUIData(&ui)
	// fmt.Printf("AdminEditHandler: d = %#v\n", d)
	err = t.Execute(w, &ui)
	if nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
