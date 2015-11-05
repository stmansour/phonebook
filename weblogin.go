package main

import (
	"crypto/sha512"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
)

func webloginHandler(w http.ResponseWriter, r *http.Request) {
	n := 0 //error number associated with this login attempt
	loggedIn := false
	fmt.Printf("webloginHandler\n")
	myusername := strings.ToLower(r.FormValue("username"))
	password := []byte(r.FormValue("password"))
	sha := sha512.Sum512(password)
	mypasshash := fmt.Sprintf("%x", sha)

	var passhash string
	var uid int
	err := Phonebook.db.QueryRow("select uid,passhash from people where username=?", myusername).Scan(&uid, &passhash)
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
