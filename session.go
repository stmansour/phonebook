package main

import (
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
	sessions = make(map[string]*session, 0)
}
func sessionGet(token string) *session {
	return sessions[token]
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
		cookie.Path = "/"
		http.SetCookie(w, cookie)
		return 0
	}
	return 1
}
