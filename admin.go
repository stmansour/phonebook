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
	var ui uiSupport
	Phonebook.ReqMem <- 1    // ask to access the shared mem, blocks until granted
	<-Phonebook.ReqMemAck    // make sure we got it
	initUIData(&ui)          // initialize our data
	Phonebook.ReqMemAck <- 1 // tell Dispatcher we're done with the data
	err := t.Execute(w, &ui)
	if nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
