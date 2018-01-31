package main

import (
	"net/http"
	"phonebook/authz"
	"phonebook/db"
	"phonebook/sess"
	"strconv"
)

func adminBecomeHandler(w http.ResponseWriter, r *http.Request) {
	var ssn *sess.Session
	var ui uiSupport
	ssn = nil
	if 0 < initHandlerSession(ssn, &ui, w, r) {
		return
	}
	ssn = ui.X
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

	// fmt.Printf("Current session = %s\n", ssn.ToString())
	var tmp sess.Session
	var d db.PersonDetail
	d.Reports = make([]db.Person, 0)
	d.UID = ssn.UIDorig
	adminReadDetails(&d)
	authz.GetRoleInfo(d.RID, &tmp.PMap)
	// fmt.Printf("UIDorig = %d,  role name = %s, RID = %d\n", ssn.UIDorig, tmp.Urole.Name, tmp.Urole.RID)

	//============================================================
	// SECURITY
	//============================================================
	if tmp.PMap.Urole.Name != "Administrator" {
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}

	sessionBecome(ssn, uid)
	searchHandler(w, r)
}
