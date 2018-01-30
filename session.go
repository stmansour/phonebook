package main

import (
	"fmt"
	"phonebook/authz"
	"phonebook/sess"
	"phonebook/ui"
)

func sessionInit() {
	sess.Sessions = make(map[string]*sess.Session)
}

func sessionGet(token string) (*sess.Session, bool) {
	s, ok := sess.Sessions[token]
	return s, ok
}

func hasAccess(s *sess.Session, el int, fieldName string, access int) bool {
	var perm int
	var ok bool

	sess.SessionManager.ReqSessionMem <- 1 // ask to access the shared mem, blocks until granted
	<-sess.SessionManager.ReqSessionMemAck // make sure we got it
	switch el {
	case authz.ELEMPERSON:
		perm, ok = s.Pp[fieldName] // here's the permission we have
	case authz.ELEMCOMPANY:
		perm, ok = s.Pco[fieldName] // here's the permission we have
	case authz.ELEMCLASS:
		perm, ok = s.Pcl[fieldName] // here's the permission we have
	case authz.ELEMPBSVC:
		perm, ok = s.Ppr[fieldName] // here's the permission we have
	}
	sess.SessionManager.ReqSessionMemAck <- 1 // tell SessionDispatcher we're done with the data
	ok = (0 != perm&access)
	// fmt.Printf("hasFieldAccess: access to el: %d, field %s, access 0x%02x: %v\n", el, fieldName, access, ok)
	return ok // could be true or false

}

func hasFieldAccess(token string, el int, fieldName string, access int) bool {
	s, ok := sess.Sessions[token]
	if !ok {
		fmt.Printf("hasFieldAccess:  Could not find sess.Session for %s\n", token)
		return false
	}
	return hasAccess(s, el, fieldName, access)
}

func hasPERMMODaccess(token string, el int, fieldName string) bool {
	return hasFieldAccess(token, el, fieldName, authz.PERMMOD)
}

//=====================================================================================
// SYNOPSIS:
// 		hasAdminScreenAccess scans the permissions of the supplied element's fields
// 		in the sess.Session associated with the logged in user. If the at least one of the
// 		fields has the requested permission, this function returns true. Otherwise it
// 		returns false
// PARAMS:
//		token - sess.Session token
//		el - check data fields for this element type. One of authz.ELEMPERSON, authz.ELEMCOMPANY,
//			 authz.ELEMCLASS
//      perm - logical OR of the required permissions
// RETURNS:
//		true  - if the user with this sess.Session has the required permissions to see the
//			    admin screen
//      false - if the user does not have the required permissions
//=====================================================================================
func pvtHasAdminScreenAccess(s *sess.Session, el int, perm int) bool {
	var p int
	var ok bool
	for i := 0; i < len(adminScreenFields); i++ {
		if adminScreenFields[i].Elem == el {
			if (el == authz.ELEMPERSON && adminScreenFields[i].AdminScreen) || (el != authz.ELEMPERSON) {
				switch el {
				case authz.ELEMPERSON:
					p, ok = s.Pp[adminScreenFields[i].FieldName] // here's the permission we have
				case authz.ELEMCOMPANY:
					p, ok = s.Pco[adminScreenFields[i].FieldName] // here's the permission we have
				case authz.ELEMCLASS:
					p, ok = s.Pcl[adminScreenFields[i].FieldName] // here's the permission we have
				case authz.ELEMPBSVC:
					p, ok = s.Ppr[adminScreenFields[i].FieldName] // here's the permission we have
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
	s, ok := sess.Sessions[token]
	if !ok {
		fmt.Printf("hasAdminScreenAccess:  Could not find sess.Session for %s\n", token)
		return false
	}
	sess.SessionManager.ReqSessionMem <- 1    // ask to access the shared mem, blocks until granted
	<-sess.SessionManager.ReqSessionMemAck    // make sure we got it
	ok = pvtHasAdminScreenAccess(s, el, perm) //we have the memory, do the work
	sess.SessionManager.ReqSessionMemAck <- 1 // tell SessionDispatcher we're done with the data
	return ok
}

//=====================================================================================
// SYNOPSIS:
// 		showAdminButton determines whether or not the Admin button needs to appear
// 		on the menu.
// PARAMS:
//		token - sess.Session token
// RETURNS:
//		true  - if the admin button should be shown
//      false - if it should not
//=====================================================================================
func pvtShowAdminButton(s *sess.Session) bool {
	for i := 0; i < len(s.Urole.Perms); i++ {
		if s.Urole.Perms[i].Perm&authz.PERMCREATE != 0 {
			return true
		}
	}
	return false
}
func showAdminButton(token string) bool {
	s, ok := sess.Sessions[token]
	if !ok {
		fmt.Printf("showAdminButton:  Could not find sess.Session for %s\n", token)
		return false
	}
	sess.SessionManager.ReqSessionMem <- 1    // ask to access the shared mem, blocks until granted
	<-sess.SessionManager.ReqSessionMemAck    // make sure we got it
	ok = pvtShowAdminButton(s)                //we have the memory, do the work
	sess.SessionManager.ReqSessionMemAck <- 1 // tell SessionDispatcher we're done with the data
	return ok
}

func getRoleInfo(rid int, s *sess.Session) {
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
	s.Ppr = make(map[string]int)

	for i := 0; i < len(r.Perms); i++ {
		var f authz.FieldPerm
		f.Elem = r.Perms[i].Elem
		f.Field = r.Perms[i].Field
		if found < 0 {
			f.Perm = authz.PERMVIEW
		} else {
			f.Perm = r.Perms[i].Perm
		}
		s.Urole.Perms = append(s.Urole.Perms, f)

		// fast access maps:
		switch f.Elem {
		case authz.ELEMPERSON:
			s.Pp[f.Field] = f.Perm
		case authz.ELEMCOMPANY:
			s.Pco[f.Field] = f.Perm
		case authz.ELEMCLASS:
			s.Pcl[f.Field] = f.Perm
		case authz.ELEMPBSVC:
			s.Ppr[f.Field] = f.Perm
		}
	}
}

func sessionNew(token, username, firstname string, uid int, rid int) *sess.Session {
	s := new(sess.Session)
	s.Token = token
	s.Username = username
	s.Firstname = firstname
	s.UID = uid
	s.UIDorig = uid
	s.ImageURL = getImageFilename(uid)
	s.Breadcrumbs = make([]ui.Crumb, 0)
	getRoleInfo(rid, s)

	if Phonebook.SecurityDebug {
		for i := 0; i < len(s.Urole.Perms); i++ {
			ulog("f: %s,  perm: %02x\n", s.Urole.Perms[i].Field, s.Urole.Perms[i].Perm)
		}
	}

	var d personDetail
	d.UID = uid
	//getSecurityList(&d)

	// err := Phonebook.prepstmt.getUserCoCode.QueryRow("select cocode from people where uid=?", uid).Scan(&s.CoCode)
	err := Phonebook.db.QueryRow(fmt.Sprintf("select cocode from people where uid=%d", uid)).Scan(&s.CoCode)
	if nil != err {
		ulog("Unable to read CoCode for userid=%d,  err = %v\n", uid, err)
	}

	sess.SessionManager.ReqSessionMem <- 1 // ask to access the shared mem, blocks until granted
	<-sess.SessionManager.ReqSessionMemAck // make sure we got it
	sess.Sessions[token] = s
	sess.SessionManager.ReqSessionMemAck <- 1 // tell SessionDispatcher we're done with the data
	sulog("New Session: %s\n", s.ToString())
	sulog("sess.Session.Urole.perms = %+v\n", s.Urole.Perms)
	return s
}

// Privileged function allowing one user to become another user. This is meant
// to be used by Administrators or User Support personnel.
func sessionBecome(s *sess.Session, uid int) {
	var d personDetail
	d.Reports = make([]person, 0)
	d.UID = uid
	adminReadDetails(&d)

	s.Firstname = d.FirstName
	if 0 < len(d.PreferredName) {
		s.Firstname = d.PreferredName
	}
	s.UID = uid
	s.Username = d.UserName
	s.ImageURL = getImageFilename(uid)
	getRoleInfo(d.RID, s)

	if Phonebook.SecurityDebug {
		for i := 0; i < len(s.Urole.Perms); i++ {
			ulog("f: %s,  perm: %02x\n", s.Urole.Perms[i].Field, s.Urole.Perms[i].Perm)
		}
	}

	ulog("user %d to BECOME user %d", s.UIDorig, s.UID)
}

// remove the supplied sess.Session.
// if there is a better idiomatic way to do this, please let me know.
func sessionDelete(s *sess.Session) {
	// fmt.Printf("Session being deleted: %s\n", s.ToString())
	// fmt.Printf("sess.Sessions before delete:\n")
	// dumpSessions()

	ss := make(map[string]*sess.Session, 0)

	sess.SessionManager.ReqSessionMem <- 1 // ask to access the shared mem, blocks until granted
	<-sess.SessionManager.ReqSessionMemAck // make sure we got it
	for k, v := range sess.Sessions {
		if s.Token != k {
			ss[k] = v
		}
	}
	sess.Sessions = ss
	sess.SessionManager.ReqSessionMemAck <- 1 // tell SessionDispatcher we're done with the data
	// fmt.Printf("sess.Sessions after delete:\n")
	// dumpSessions()
}
