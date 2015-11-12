package main

import (
	"fmt"
	"net/http"
	"time"
)

type session struct {
	Token     string    // this is the md5 hash, unique id
	Username  string    // associated username
	Firstname string    // user's first name
	UID       int       // user's db uid
	ImageURL  string    // user's picture
	Expire    time.Time // when does the cookie expire
}

var sessions map[string]*session

func sessionInit() {
	sessions = make(map[string]*session)
}
func sessionGet(token string) *session {
	return sessions[token]
}

func (s *session) ToString() string {
	if nil == s {
		return "nil"
	}
	return fmt.Sprintf("User(%s) Name(%s) UID(%d) Token(%s)", s.Username, s.Firstname, s.UID, s.Token)
}

func dumpSessions() {
	i := 0
	for _, v := range sessions {
		fmt.Printf("%2d. %s\n", i, v.ToString())
		i++
	}
}

func sessionNew(token, username, firstname string, uid int, image string) *session {
	s := new(session)
	s.Token = token
	s.Username = username
	s.Firstname = firstname
	s.UID = uid
	s.ImageURL = image
	sessions[token] = s
	return s
}

func (s *session) refresh(w http.ResponseWriter, r *http.Request) int {
	cookie, err := r.Cookie("accord")
	if nil != cookie && err == nil {
		cookie.Expires = time.Now().Add(10 * time.Minute)
		s.Expire = cookie.Expires
		cookie.Path = "/"
		http.SetCookie(w, cookie)
		return 0
	}
	return 1
}

// remove the supplied session.
// if there is a better idiomatic way to do this, please let me know.
func sessionDelete(s *session) {
	// fmt.Printf("Session being deleted: %s\n", s.ToString())
	// fmt.Printf("sessions before delete:\n")
	// dumpSessions()

	ss := make(map[string]*session, 0)
	for k, v := range sessions {
		if s.Token != k {
			ss[k] = v
		}
	}
	sessions = ss
	// fmt.Printf("sessions after delete:\n")
	// dumpSessions()
}
