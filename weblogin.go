package main

import (
	"crypto/md5"
	"crypto/sha512"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// initHandlerSession validates the session cookie and redirects if necessary.
// it also initializes the uiSession variable
// RETURNS:  0 = no problems
//           1 = redirected
func initHandlerSession(sess *session, ui *uiSupport, w http.ResponseWriter, r *http.Request) int {
	var ok bool
	cookie, err := r.Cookie("accord")
	if nil != cookie && err == nil {
		sess, ok = sessionGet(cookie.Value)
		if !ok || sess == nil {
			http.Redirect(w, r, "/signin/", http.StatusFound)
			return 1
		}
		sess.refresh(w, r)
	} else {
		//fmt.Printf("REDIRECT to signin\n")
		http.Redirect(w, r, "/signin/", http.StatusFound)
		return 1
	}
	ui.X = sess
	Phonebook.ReqMem <- 1    // ask to access the shared mem, blocks until granted
	<-Phonebook.ReqMemAck    // make sure we got it
	initUIData(ui)           // initialize our data
	Phonebook.ReqMemAck <- 1 // tell Dispatcher we're done with the data
	return 0
}

func webloginHandler(w http.ResponseWriter, r *http.Request) {
	n := 0 //error number associated with this login attempt
	loggedIn := false
	myusername := strings.ToLower(r.FormValue("username"))
	password := []byte(r.FormValue("password"))
	sha := sha512.Sum512(password)
	mypasshash := fmt.Sprintf("%x", sha)

	var passhash, firstname, preferredname string
	var uid, RID int
	err := Phonebook.db.QueryRow("select uid,firstname,preferredname,passhash,rid from people where username=?", myusername).Scan(&uid, &firstname, &preferredname, &passhash, &RID)
	switch {
	case err == sql.ErrNoRows:
		ulog("No user with username = %s\n", myusername)
		n = 1
	case err != nil:
		ulog("login username: %s,  error = %v\n", myusername, err)
		n = 2
	default:
		// ulog("found username %s in database. UID = %d\n", myusername, uid)
	}

	if passhash == mypasshash {
		loggedIn = true
		ulog("user %s logged in\n", myusername)
		expiration := time.Now().Add(10 * time.Minute)
		//=================================================================================
		// There could be multiple sessions from the same user on different browsers.
		// These could be on the same or separate machines. We need the IP and the browser
		// to guarantee uniqueness...
		//=================================================================================
		key := myusername + r.Header.Get("User-Agent") + r.RemoteAddr
		cval := fmt.Sprintf("%x", md5.Sum([]byte(key)))
		name := firstname
		if len(preferredname) > 0 {
			name = preferredname
		}

		s := sessionNew(cval, myusername, name, uid, RID, "/images/anon.png")
		cookie := http.Cookie{Name: "accord", Value: s.Token, Expires: expiration}
		cookie.Path = "/"
		http.SetCookie(w, &cookie)
	} else {
		ulog("user name or password did not match for: %s\n", myusername)
		n = 1
	}

	if !loggedIn {
		http.Redirect(w, r, fmt.Sprintf("/signin/%d", n), http.StatusFound)
	} else {
		http.Redirect(w, r, "/search/", http.StatusFound)
	}
}
