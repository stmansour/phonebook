package main

import (
	"fmt"
	"net/http"
	"time"
)

type session struct {
	Token       string         // this is the md5 hash, unique id
	Username    string         // associated username
	Firstname   string         // user's first name
	UID         int            // user's db uid
	Urole       Role           // user's role for permissions
	CoCode      int            // logged in user's company
	ImageURL    string         // user's picture
	Expire      time.Time      // when does the cookie expire
	Pp          map[string]int // quick way to reference person permissions based on field name
	Pco         map[string]int // quick way to reference company permissions based on field name
	Pcl         map[string]int // quick way to reference class permissions based on field name
	Breadcrumbs []Crumb        // where is the user in the screen hierarchy
}

var sessions map[string]*session

// SessionDispatcher controls access to shared memory.
func SessionDispatcher() {
	for {
		select {
		case <-Phonebook.ReqSessionMem:
			Phonebook.ReqSessionMemAck <- 1 // tell caller go ahead
			<-Phonebook.ReqSessionMemAck    // block until caller is done with mem
		}
	}
}

// SessionCleanup periodically spins through the list of sessions
// and removes any which have timed out.
func SessionCleanup() {
	for {
		select {
		case <-time.After(Phonebook.SessionCleanupTime * time.Minute):
			Phonebook.ReqSessionMem <- 1       // ask to access the shared mem, blocks until granted
			<-Phonebook.ReqSessionMemAck       // make sure we got it
			ss := make(map[string]*session, 0) // here's the new session list
			n := 0                             // total number removed
			for k, v := range sessions {       // look at every session
				if time.Now().After(v.Expire) { // if it's still active...
					n++ // removed another
					fmt.Printf("v.Expire = %s, t.Now() = %s\n",
						v.Expire.Format(time.RFC3339), time.Now().Format(time.RFC3339))
				} else {
					ss[k] = v // ...copy it to the new list
				}
			}
			sessions = ss                   // set the new list
			Phonebook.ReqSessionMemAck <- 1 // tell SessionDispatcher we're done with the data
			fmt.Printf("SessionCleanup completed. %d removed. Current session list size = %d\n", n, len(sessions))
		}
	}
}

func sessionInit() {
	sessions = make(map[string]*session)
}
func sessionGet(token string) (*session, bool) {
	s, ok := sessions[token]
	return s, ok
}

func (s *session) ToString() string {
	if nil == s {
		return "nil"
	}
	return fmt.Sprintf("User(%s) Name(%s) UID(%d) Token(%s)  Role(%s)",
		s.Username, s.Firstname, s.UID, s.Token, s.Urole.Name)
}

func dumpSessions() {
	i := 0
	for _, v := range sessions {
		fmt.Printf("%2d. %s\n", i, v.ToString())
		i++
	}
}

func hasPERMMODaccess(token string, el int, fieldName string) bool {
	var perm int
	s, ok := sessions[token]
	if !ok {
		fmt.Printf("hasPERMMODaccess:  Could not find session for %s\n", token)
		return false
	}

	Phonebook.ReqSessionMem <- 1 // ask to access the shared mem, blocks until granted
	<-Phonebook.ReqSessionMemAck // make sure we got it
	switch el {
	case ELEMPERSON:
		perm, ok = s.Pp[fieldName] // here's the permission we have
	case ELEMCOMPANY:
		perm, ok = s.Pco[fieldName] // here's the permission we have
	case ELEMCLASS:
		perm, ok = s.Pcl[fieldName] // here's the permission we have
	}
	Phonebook.ReqSessionMemAck <- 1 // tell SessionDispatcher we're done with the data
	ok = (0 != perm&PERMMOD)
	dulog("%v\n", ok)
	return ok // could be true or false
}

//=====================================================================================
// SYNOPSIS:
// 		hasAdminScreenAccess scans the permissions of the supplied element's fields
// 		in the session associated with the logged in user. If the at least one of the
// 		fields has the requested permission, this function returns true. Otherwise it
// 		returns false
// PARAMS:
//		token - session token
//		el - check data fields for this element type. One of ELEMPERSON, ELEMCOMPANY,
//			 ELEMCLASS
//      perm - logical OR of the required permissions
// RETURNS:
//		true  - if the user with this session has the required permissions to see the
//			    admin screen
//      false - if the user does not have the required permissions
//=====================================================================================
func pvtHasAdminScreenAccess(s *session, el int, perm int) bool {
	var p int
	var ok bool
	for i := 0; i < len(adminScreenFields); i++ {
		if adminScreenFields[i].Elem == el {
			if (el == ELEMPERSON && adminScreenFields[i].AdminScreen) || (el != ELEMPERSON) {
				switch el {
				case ELEMPERSON:
					p, ok = s.Pp[adminScreenFields[i].FieldName] // here's the permission we have
				case ELEMCOMPANY:
					p, ok = s.Pco[adminScreenFields[i].FieldName] // here's the permission we have
				case ELEMCLASS:
					p, ok = s.Pcl[adminScreenFields[i].FieldName] // here's the permission we have
				}
				if ok { // if we have a permission for the field name
					// fmt.Printf("p = 0x%02x\n", p)
					pcheck := p & perm // AND it with the required permission
					if 0 != pcheck {   // if the result is non-zero...
						// fmt.Printf("granted\n")
						return true // ... we have the permission to view the screen
					}
				}
			}
		}
	}
	// fmt.Printf("not granted\n")
	return false
}

func hasAdminScreenAccess(token string, el int, perm int) bool {
	s, ok := sessions[token]
	if !ok {
		fmt.Printf("hasAdminScreenAccess:  Could not find session for %s\n", token)
		return false
	}
	Phonebook.ReqSessionMem <- 1              // ask to access the shared mem, blocks until granted
	<-Phonebook.ReqSessionMemAck              // make sure we got it
	ok = pvtHasAdminScreenAccess(s, el, perm) //we have the memory, do the work
	Phonebook.ReqSessionMemAck <- 1           // tell SessionDispatcher we're done with the data
	return ok
}

//=====================================================================================
// SYNOPSIS:
// 		showAdminButton determines whether or not the Admin button needs to appear
// 		on the menu.
// PARAMS:
//		token - session token
// RETURNS:
//		true  - if the admin button should be shown
//      false - if it should not
//=====================================================================================
func pvtShowAdminButton(s *session) bool {
	for i := 0; i < len(s.Urole.Perms); i++ {
		if s.Urole.Perms[i].Perm&PERMCREATE != 0 {
			return true
		}
	}
	return false
}
func showAdminButton(token string) bool {
	s, ok := sessions[token]
	if !ok {
		fmt.Printf("showAdminButton:  Could not find session for %s\n", token)
		return false
	}
	Phonebook.ReqSessionMem <- 1    // ask to access the shared mem, blocks until granted
	<-Phonebook.ReqSessionMemAck    // make sure we got it
	ok = pvtShowAdminButton(s)      //we have the memory, do the work
	Phonebook.ReqSessionMemAck <- 1 // tell SessionDispatcher we're done with the data
	return ok
}

func getRoleInfo(rid int, s *session) {
	found := -1
	idx := -1

	// try to find the requested index
	sulog("len(Phonebook.Roles)=%d\n", len(Phonebook.Roles))
	sulog("getRoleInfo - looking for rid=%d\n", rid)
	for i := 0; i < len(Phonebook.Roles); i++ {
		sulog("Phonebook.Roles[%d] = %+v\n", i, Phonebook.Roles[i])
		if rid == Phonebook.Roles[i].RID {
			found = i
			idx = i
			s.Urole.Name = Phonebook.Roles[i].Name
			s.Urole.RID = rid
			break
		}
	}

	if found < 0 {
		idx = 0
		ulog("Did not find rid == %d, all permissions set to read-only\n", rid)
	}

	r := Phonebook.Roles[idx]
	s.Pp = make(map[string]int)
	s.Pco = make(map[string]int)
	s.Pcl = make(map[string]int)

	for i := 0; i < len(r.Perms); i++ {
		var f FieldPerm
		f.Elem = r.Perms[i].Elem
		f.Field = r.Perms[i].Field
		if found < 0 {
			f.Perm = PERMVIEW
		} else {
			f.Perm = r.Perms[i].Perm
		}
		s.Urole.Perms = append(s.Urole.Perms, f)

		// fast access maps:
		switch f.Elem {
		case ELEMPERSON:
			s.Pp[f.Field] = f.Perm
		case ELEMCOMPANY:
			s.Pco[f.Field] = f.Perm
		case ELEMCLASS:
			s.Pcl[f.Field] = f.Perm
		}
	}

}

func sessionNew(token, username, firstname string, uid int, rid int, image string) *session {
	s := new(session)
	s.Token = token
	s.Username = username
	s.Firstname = firstname
	s.UID = uid
	s.ImageURL = image
	s.Breadcrumbs = make([]Crumb, 0)
	getRoleInfo(rid, s)

	if Phonebook.SecurityDebug {
		for i := 0; i < len(s.Urole.Perms); i++ {
			ulog("f: %s,  perm: %02x\n", s.Urole.Perms[i].Field, s.Urole.Perms[i].Perm)
		}
	}

	var d personDetail
	d.UID = uid
	//getSecurityList(&d)
	err := Phonebook.db.QueryRow("select cocode from people where uid=?", uid).Scan(&s.CoCode)
	if nil != err {
		ulog("Unable to read CoCode for userid=%d,  err = %v\n", uid, err)
	}

	Phonebook.ReqSessionMem <- 1 // ask to access the shared mem, blocks until granted
	<-Phonebook.ReqSessionMemAck // make sure we got it
	sessions[token] = s
	Phonebook.ReqSessionMemAck <- 1 // tell SessionDispatcher we're done with the data
	sulog("New Session: %s\n", s.ToString())
	sulog("session.Urole.perms = %+v\n", s.Urole.Perms)
	return s
}

//=====================================================================================
// refresh updates the cookie and session with a new expire time.
//=====================================================================================
func (s *session) refresh(w http.ResponseWriter, r *http.Request) int {
	cookie, err := r.Cookie("accord")
	if nil != cookie && err == nil {
		cookie.Expires = time.Now().Add(Phonebook.SessionTimeout * time.Minute)
		Phonebook.ReqSessionMem <- 1    // ask to access the shared mem, blocks until granted
		<-Phonebook.ReqSessionMemAck    // make sure we got it
		s.Expire = cookie.Expires       // update the session information
		Phonebook.ReqSessionMemAck <- 1 // tell SessionDispatcher we're done with the data
		cookie.Path = "/"
		http.SetCookie(w, cookie)
		return 0
	}
	return 1
}

//=====================================================================================
// elemPermsAny determines whether or not the session has permissions to perform the
// requested operations.  NOTE:  This interface does check the UID to fully cover
// permissions PERMOWNERVIEW or PERMOWNERMOD. This must be done at a higher level.
//
// ARGS:
//   ent  = which element? ELEMPERSON, ELEMCOMPANY, ELEMCLASS, ...
//   perm = logical or of the desired permissions.  Example PERMCREATE | PERMMOD
//
// RETURNS:
//   true if there are ANY fields for the specified element for
//   with the requested permission.
//=====================================================================================
func pvtElemPermsAny(s *session, elem int, perm int) bool {
	sulog("elemPermsAny:  elem=%d, perm = 0x%02x\n", elem, perm)
	for i := 0; i < len(s.Urole.Perms); i++ {
		sulog("s.Urole.Perms[%d].Elem = %d\n", i, s.Urole.Perms[i].Elem)
		if s.Urole.Perms[i].Elem == elem {
			res := s.Urole.Perms[i].Perm & perm
			sulog("fieldname: %s  s.Urole.Perms[%d].Perm = 0x%02x, s.Urole.Perms[%d].Perm & perm = 0x%02x\n",
				s.Urole.Perms[i].Field, i, s.Urole.Perms[i].Elem, i, res)
			if res != 0 { // if any of the permissions exist
				sulog("return true") // we're good to go for this check
				return true
			}
		}
	}
	sulog("return false")
	return false
}

func (s *session) elemPermsAny(elem int, perm int) bool {
	Phonebook.ReqSessionMem <- 1         // ask to access the shared mem, blocks until granted
	<-Phonebook.ReqSessionMemAck         // make sure we got it
	ok := pvtElemPermsAny(s, elem, perm) // look for perms
	Phonebook.ReqSessionMemAck <- 1      // tell SessionDispatcher we're done with the data
	return ok
}

//=====================================================================================
// elemPermsAll determines whether or not the session has permissions to perform all
// requested operations.  NOTE: This interface does check the UID to fully cover
// permissions PERMOWNERVIEW or PERMOWNERMOD. This must be done at a higher level.
//
// ARGS:
//   ent  = which element? ELEMPERSON, ELEMCOMPANY, ELEMCLASS, ...
//   perm = logical or of the desired permissions.  Example PERMCREATE | PERMMOD
//
// RETURNS:
//   true if ALL permission fields for the specified element are present
//=====================================================================================
func pvtElemPermsAll(s *session, elem int, perm int) bool {
	sulog("elemPermsAll:  elem=%d, perm = 0x%02x\n", elem, perm)
	for i := 0; i < len(s.Urole.Perms); i++ {
		sulog("s.Urole.Perms[%d].Elem = %d\n", i, s.Urole.Perms[i].Elem)
		if s.Urole.Perms[i].Elem == elem {
			res := s.Urole.Perms[i].Perm & perm
			sulog("fieldname: %s  s.Urole.Perms[%d].Perm = 0x%02x, s.Urole.Perms[%d].Perm & perm = 0x%02x\n",
				s.Urole.Perms[i].Field, i, s.Urole.Perms[i].Elem, i, res)
			if res == perm { // if all bits are present, res will match perm
				sulog("return true") // we're good to go for this check
				return true
			}
		}
	}
	sulog("return false")
	return false
}
func (s *session) elemPermsAll(elem int, perm int) bool {
	Phonebook.ReqSessionMem <- 1         // ask to access the shared mem, blocks until granted
	<-Phonebook.ReqSessionMemAck         // make sure we got it
	ok := pvtElemPermsAll(s, elem, perm) // look for perms
	Phonebook.ReqSessionMemAck <- 1      // tell SessionDispatcher we're done with the data
	return ok
}

// remove the supplied session.
// if there is a better idiomatic way to do this, please let me know.
func sessionDelete(s *session) {
	// fmt.Printf("Session being deleted: %s\n", s.ToString())
	// fmt.Printf("sessions before delete:\n")
	// dumpSessions()

	ss := make(map[string]*session, 0)

	Phonebook.ReqSessionMem <- 1 // ask to access the shared mem, blocks until granted
	<-Phonebook.ReqSessionMemAck // make sure we got it
	for k, v := range sessions {
		if s.Token != k {
			ss[k] = v
		}
	}
	sessions = ss
	Phonebook.ReqSessionMemAck <- 1 // tell SessionDispatcher we're done with the data
	// fmt.Printf("sessions after delete:\n")
	// dumpSessions()
}
