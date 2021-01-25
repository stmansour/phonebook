package main

import (
	"fmt"
	"net/http"
	"phonebook/db"
	"phonebook/lib"
	"phonebook/ws"
	"strconv"
	"text/template"
)

// ErrMsgs are a list of string we can convey to the UI to indicate when an error occurs
var ErrMsgs = []string{
	"",                               // 0
	"Username or password not found", // 1
	"System error",                   // 2
}

// normal call:  http://host:8250/search/
// login failed: http://host:8250/search/1     (the error number is in the filepath)
func signinHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	// lib.Console("Entered signinHandler.  Looking for cookie named: %s\n", db.SessionCookieName)
	cookie, _ := r.Cookie(db.SessionCookieName)
	if nil != cookie {
		// lib.Console("Cookie named %s found.  value = %s\n", db.SessionCookieName, cookie.Value)
		s, ok := db.SessionGet(cookie.Value)
		if ok {
			if s.Token == cookie.Value {
				// fmt.Printf("FOUND session, redirecting\n")
				http.Redirect(w, r, "/search/", http.StatusFound)
				return
			}
		}
		//----------------------------------------------------------------
		// The login may have come from another AIR application. Create a
		// session
		//----------------------------------------------------------------
		// lib.Console("Cookie named %s was not found in the session table\n", db.SessionCookieName)
		// lib.Console("Creating new session...\n")
		c, err := db.GetSessionCookie(cookie.Value)
		// lib.Console("session cookie:  %#v\n", c)
		if err != nil {
			lib.Ulog("signinHandler: error getting session cookie: %s\n", err.Error())
		} else if len(c.Cookie) > 0 {
			s := ws.NewSessionFromCookie(&c)
			// lib.Console("Creating new session from cookie. s.Username = %s\n", s.Username)
			if len(s.Username) > 0 {
				http.Redirect(w, r, "/search/", http.StatusFound)
				return
			}
		}
	}

	var err error
	n := 0
	path := "/signin/"
	nstr := r.RequestURI[len(path):]
	if len(nstr) > 0 {
		if '0' <= nstr[0] && nstr[0] <= '9' {
			// path == errno
			n, err = strconv.Atoi(nstr)
			if err != nil {
				ulog("signinHandler: Error converting uid to a number: %v. URI: %s\n", err, r.RequestURI)
			}
		} else {
			// path == return to this path after login
			// TODO: work out how to make it happen
		}
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
		errmsg := fmt.Sprintf("signinHandler: err = %v\n", err)
		ulog(errmsg)
		fmt.Println(errmsg)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
