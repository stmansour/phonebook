package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

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

	var c class
	path := "/class/"
	costr := r.RequestURI[len(path):]
	if len(costr) > 0 {
		classcode, _ := strconv.Atoi(costr)
		getClassInfo(classcode, &c)
		t, _ := template.New("class.html").ParseFiles("class.html")
		ui.A = &c
		err := t.Execute(w, &ui)
		if nil != err {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		fmt.Fprintf(w, "classcode = %s\nCould not convert to number\n", costr)
	}
}
