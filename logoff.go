package main

import (
	"net/http"
	"phonebook/db"
	"time"
)

func logoffHandler(w http.ResponseWriter, r *http.Request) {
	Phonebook.ReqCountersMem <- 1    // ask to access the shared mem, blocks until granted
	<-Phonebook.ReqCountersMemAck    // make sure we got it
	Counters.Logoff++                // initialize our data
	Phonebook.ReqCountersMemAck <- 1 // tell Dispatcher we're done with the data

	var ok bool
	// var ui uiSupport
	var ssn *db.Session

	// ssn = nil
	// if 0 < initHandlerSession(ssn, &ui, w, r) {
	// 	return
	// }
	// ssn = ui.X
	cookie, err := r.Cookie(db.SessionCookieName)

	w.Header().Set("Content-Type", "text/html")
	if nil != cookie && err == nil {
		ssn, ok = db.SessionGet(cookie.Value)
		if ok {
			db.SessionDelete(ssn)
		}
		// force the cookie to expire
		cookie.Expires = time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
		cookie.Path = "/"
		http.SetCookie(w, cookie)
		r.AddCookie(cookie) // need this so that the redirect to search finds the cookie
	}
	http.Redirect(w, r, "/signin/", http.StatusFound)
}
