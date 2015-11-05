package main

import (
	"net/http"
	"strconv"
	"text/template"
)

// ErrMsgs are a list of string we can convey to the UI to indicate when an error occurs
var ErrMsgs = []string{
	"", // 0
	"Username or password not found", // 1
	"System error",                   // 2
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	n := 0
	path := "/signin/"
	nstr := r.RequestURI[len(path):]
	if len(nstr) > 0 {
		n, err = strconv.Atoi(nstr)
		if err != nil {
			ulog("signinHandler: Error converting uid to a number: %v. URI: %s\n", err, r.RequestURI)
		}
	}

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

	t, _ := template.New("signin.html").Funcs(funcMap).ParseFiles("signin.html")
	var ui uiSupport
	Phonebook.ReqMem <- 1    // ask to access the shared mem, blocks until granted
	<-Phonebook.ReqMemAck    // make sure we got it
	initUIData(&ui)          // initialize our data
	Phonebook.ReqMemAck <- 1 // tell Dispatcher we're done with the data

	var S signin
	S.ErrNo = n
	S.ErrMsg = ErrMsgs[n]

	ui.S = &S

	err = t.Execute(w, &ui)
	if nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
