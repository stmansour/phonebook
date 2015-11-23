package main

import "net/http"

func logoffHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var sess *session
	var ui uiSupport
	sess = nil
	if 0 < initHandlerSession(sess, &ui, w, r) {
		return
	}
	sess = ui.X
	cookie, err := r.Cookie("accord")
	if nil != cookie && err == nil {
		sess = sessionGet(cookie.Value)
		sessionDelete(sess)
	}
	http.Redirect(w, r, "/signin/", http.StatusFound)
}
