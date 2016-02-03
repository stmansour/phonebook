package main

import (
	"fmt"
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

	Phonebook.ReqSessionMem <- 1 // ask to access the shared mem, blocks until granted
	<-Phonebook.ReqSessionMemAck // make sure we got it
	var p []session
	for _, v := range sessions {
		s := session{}
		s.Token = v.Token
		s.Firstname = v.Firstname
		s.ImageURL = v.ImageURL
		s.UID = v.UID
		s.Username = v.Username
		s.Expire = v.Expire
		p = append(p, s)
	}
	Phonebook.ReqSessionMemAck <- 1 // tell SessionDispatcher we're done with the data

	ui.N = p

	t, _ := template.New("stats.html").Funcs(funcMap).ParseFiles("stats.html")
	err := t.Execute(w, &ui)

	if nil != err {
		errmsg := fmt.Sprintf("statsHandler: err = %v\n", err)
		ulog(errmsg)
		fmt.Println(errmsg)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
