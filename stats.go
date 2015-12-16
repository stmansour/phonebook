package main

import (
	"net/http"
	"text/template"
)

func statsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var sess *session
	var ui uiSupport
	sess = nil
	if 0 < initHandlerSession(sess, &ui, w, r) {
		return
	}
	sess = ui.X
	breadcrumbAdd(sess, "Stats", "/stats/")

	var MyCounters UsageCounters
	Phonebook.ReqCountersMem <- 1 // ask to access the shared mem, blocks until granted
	<-Phonebook.ReqCountersMemAck // make sure we got it
	MyCounters = Counters
	Phonebook.ReqCountersMemAck <- 1 // tell Dispatcher we're done with the data

	ui.K = &MyCounters

	t, _ := template.New("stats.html").Funcs(funcMap).ParseFiles("stats.html")
	err := t.Execute(w, &ui)

	if nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
