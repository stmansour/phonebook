package main

import (
	"fmt"
	"net/http"
	"text/template"
)

func helpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var sess *session
	var ui uiSupport
	sess = nil
	if 0 < initHandlerSession(sess, &ui, w, r) {
		return
	}
	sess = ui.X
	Phonebook.ReqCountersMem <- 1    // ask to access the shared mem, blocks until granted
	<-Phonebook.ReqCountersMemAck    // make sure we got it
	Counters.ViewClass++             // initialize our data
	Phonebook.ReqCountersMemAck <- 1 // tell Dispatcher we're done with the data
	breadcrumbAdd(sess, "Help", "/help/")
	t, _ := template.New("help.html").Funcs(funcMap).ParseFiles("help.html")
	err := t.Execute(w, &ui)
	if nil != err {
		errmsg := fmt.Sprintf("helpHandler: err = %v\n", err)
		ulog(errmsg)
		fmt.Println(errmsg)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
