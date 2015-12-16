package main

import "net/http"

func logoffHandler(w http.ResponseWriter, r *http.Request) {
	var ok bool
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
	Counters.Logoff++                // initialize our data
	Phonebook.ReqCountersMemAck <- 1 // tell Dispatcher we're done with the data

	cookie, err := r.Cookie("accord")
	if nil != cookie && err == nil {
		sess, ok = sessionGet(cookie.Value)
		if ok {
			sessionDelete(sess)
		}
	}
	http.Redirect(w, r, "/signin/", http.StatusFound)
}
