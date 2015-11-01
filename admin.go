package main

import (
	"net/http"
	"text/template"
)

func adminHandler(w http.ResponseWriter, r *http.Request) {

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

	t, _ := template.New("admin.html").Funcs(funcMap).ParseFiles("admin.html")
	err := t.Execute(w, &PhonebookUI)
	if nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
