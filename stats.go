package main

import (
	"fmt"
	"net/http"
	"phonebook/sess"
	"text/template"
)

func statsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var mysession *sess.Session
	var ui uiSupport
	mysession = nil
	if 0 < initHandlerSession(mysession, &ui, w, r) {
		return
	}
	mysession = ui.X
	breadcrumbAdd(mysession, "Stats", "/stats/")

	var MyCounters UsageCounters
	var MyiCounters UsageCounters
	Phonebook.ReqCountersMem <- 1 // ask to access the shared mem, blocks until granted
	<-Phonebook.ReqCountersMemAck // make sure we got it
	MyCounters = TotCounters
	MyiCounters = Counters
	Phonebook.ReqCountersMemAck <- 1 // tell Dispatcher we're done with the data

	ui.K = &MyCounters
	ui.Ki = &MyiCounters

	sess.SessionManager.ReqSessionMem <- 1 // ask to access the shared mem, blocks until granted
	<-sess.SessionManager.ReqSessionMemAck // make sure we got it
	var p []sess.Session
	for _, v := range sess.Sessions {
		s := sess.Session{}
		s.Token = v.Token
		s.Firstname = v.Firstname
		s.ImageURL = v.ImageURL
		s.UID = v.UID
		s.Username = v.Username
		s.Expire = v.Expire
		p = append(p, s)
	}
	sess.SessionManager.ReqSessionMemAck <- 1 // tell SessionDispatcher we're done with the data

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
