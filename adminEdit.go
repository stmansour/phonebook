package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

func adminEditHandler(w http.ResponseWriter, r *http.Request) {
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

	adminReadDetails(uid, &d)

	t, _ := template.New("adminEdit.html").ParseFiles("adminEdit.html")
	err = t.Execute(w, &d)
	if nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
