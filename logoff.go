package main

import (
	"net/http"
	"phonebook/sess"
	"time"
)

func logoffHandler(w http.ResponseWriter, r *http.Request) {
	Phonebook.ReqCountersMem <- 1    // ask to access the shared mem, blocks until granted
	<-Phonebook.ReqCountersMemAck    // make sure we got it
	Counters.Logoff++                // initialize our data
	Phonebook.ReqCountersMemAck <- 1 // tell Dispatcher we're done with the data

	var ok bool
	w.Header().Set("Content-Type", "text/html")
	var ssn *sess.Session
	var ui uiSupport
	ssn = nil
	if 0 < initHandlerSession(ssn, &ui, w, r) {
		return
	}
	ssn = ui.X

	cookie, err := r.Cookie(sess.SessionCookieName)
	if nil != cookie && err == nil {
		ssn, ok = sess.SessionGet(cookie.Value)
		if ok {
			sess.SessionDelete(ssn)
		}
		// force the cookie to expire
		cookie.Expires = time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
		cookie.Path = "/"
		http.SetCookie(w, cookie)
		r.AddCookie(cookie) // need this so that the redirect to search finds the cookie
	}
	http.Redirect(w, r, "/signin/", http.StatusFound)
}
