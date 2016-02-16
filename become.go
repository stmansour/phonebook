package main

import (
	"net/http"
	"strconv"
)

func adminBecomeHandler(w http.ResponseWriter, r *http.Request) {
	var sess *session
	var ui uiSupport
	sess = nil
	if 0 < initHandlerSession(sess, &ui, w, r) {
		return
	}
	sess = ui.X
	Phonebook.ReqCountersMem <- 1    // ask to access the shared mem, blocks until granted
	<-Phonebook.ReqCountersMemAck    // make sure we got it
	Counters.ViewPerson++            // initialize our data
	Phonebook.ReqCountersMemAck <- 1 // tell Dispatcher we're done with the data

	path := "/become/"
	uidstr := r.RequestURI[len(path):]
	uid := 0
	if len(uidstr) > 0 {
		uid, _ = strconv.Atoi(uidstr)
	}

	// fmt.Printf("Current session = %s\n", sess.ToString())
	var tmp session
	var d personDetail
	d.Reports = make([]person, 0)
	d.UID = sess.UIDorig
	adminReadDetails(&d)
	getRoleInfo(d.RID, &tmp)
	// fmt.Printf("UIDorig = %d,  role name = %s, RID = %d\n", sess.UIDorig, tmp.Urole.Name, tmp.Urole.RID)

	//============================================================
	// SECURITY
	//============================================================
	if tmp.Urole.Name != "Administrator" {
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}

	sess.sessionBecome(uid)
	searchHandler(w, r)
}
