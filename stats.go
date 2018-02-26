package main

import (
	"fmt"
	"net/http"
	"phonebook/sess"
)

func statsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var mysession *sess.Session
	var uis uiSupport
	mysession = nil
	if 0 < initHandlerSession(mysession, &uis, w, r) {
		return
	}
	mysession = uis.X
	breadcrumbAdd(mysession, "Stats", "/stats/")

	var MyCounters UsageCounters
	var MyiCounters UsageCounters
	Phonebook.ReqCountersMem <- 1 // ask to access the shared mem, blocks until granted
	<-Phonebook.ReqCountersMemAck // make sure we got it
	MyCounters = TotCounters
	MyiCounters = Counters
	Phonebook.ReqCountersMemAck <- 1 // tell Dispatcher we're done with the data

	uis.K = &MyCounters
	uis.Ki = &MyiCounters

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

	uis.N = p

	err := renderTemplate(w, uis, "stats.html")

	if nil != err {
		errmsg := fmt.Sprintf("statsHandler: err = %v\n", err)
		ulog(errmsg)
		fmt.Println(errmsg)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
