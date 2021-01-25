package main

import (
	"fmt"
	"phonebook/db"
)

func hasAccess(s *db.Session, el int, fieldName string, access int64) bool {
	var perm int64
	var ok bool

	db.SessionManager.ReqSessionMem <- 1 // ask to access the shared mem, blocks until granted
	<-db.SessionManager.ReqSessionMemAck // make sure we got it
	switch el {
	case db.ELEMPERSON:
		perm, ok = s.PMap.Pp[fieldName] // here's the permission we have
	case db.ELEMCOMPANY:
		perm, ok = s.PMap.Pco[fieldName] // here's the permission we have
	case db.ELEMCLASS:
		perm, ok = s.PMap.Pcl[fieldName] // here's the permission we have
	case db.ELEMPBSVC:
		perm, ok = s.PMap.Ppr[fieldName] // here's the permission we have
	}
	db.SessionManager.ReqSessionMemAck <- 1 // tell SessionDispatcher we're done with the data
	ok = (0 != perm&access)
	// fmt.Printf("hasFieldAccess: access to el: %d, field %s, access 0x%02x: %v\n", el, fieldName, access, ok)
	return ok // could be true or false

}

func hasFieldAccess(token string, el int, fieldName string, access int64) bool {
	s, ok := db.Sessions[token]
	if !ok {
		fmt.Printf("hasFieldAccess:  Could not find db.Session for %s\n", token)
		return false
	}
	return hasAccess(s, el, fieldName, access)
}

func hasPERMMODaccess(token string, el int, fieldName string) bool {
	return hasFieldAccess(token, el, fieldName, db.PERMMOD)
}

//=====================================================================================
// SYNOPSIS:
// 		hasAdminScreenAccess scans the permissions of the supplied element's fields
// 		in the db.Session associated with the logged in user. If the at least one of the
// 		fields has the requested permission, this function returns true. Otherwise it
// 		returns false
// PARAMS:
//		token - db.Session token
//		el - check data fields for this element type. One of db.ELEMPERSON, db.ELEMCOMPANY,
//			 db.ELEMCLASS
//      perm - logical OR of the required permissions
// RETURNS:
//		true  - if the user with this db.Session has the required permissions to see the
//			    admin screen
//      false - if the user does not have the required permissions
//=====================================================================================
func pvtHasAdminScreenAccess(s *db.Session, el int, perm int64) bool {
	var p int64
	var ok bool
	for i := 0; i < len(adminScreenFields); i++ {
		if adminScreenFields[i].Elem == el {
			if (el == db.ELEMPERSON && adminScreenFields[i].AdminScreen) || (el != db.ELEMPERSON) {
				switch el {
				case db.ELEMPERSON:
					p, ok = s.PMap.Pp[adminScreenFields[i].FieldName] // here's the permission we have
				case db.ELEMCOMPANY:
					p, ok = s.PMap.Pco[adminScreenFields[i].FieldName] // here's the permission we have
				case db.ELEMCLASS:
					p, ok = s.PMap.Pcl[adminScreenFields[i].FieldName] // here's the permission we have
				case db.ELEMPBSVC:
					p, ok = s.PMap.Ppr[adminScreenFields[i].FieldName] // here's the permission we have
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

func hasAdminScreenAccess(token string, el int, perm int64) bool {
	s, ok := db.Sessions[token]
	if !ok {
		fmt.Printf("hasAdminScreenAccess:  Could not find db.Session for %s\n", token)
		return false
	}
	db.SessionManager.ReqSessionMem <- 1      // ask to access the shared mem, blocks until granted
	<-db.SessionManager.ReqSessionMemAck      // make sure we got it
	ok = pvtHasAdminScreenAccess(s, el, perm) //we have the memory, do the work
	db.SessionManager.ReqSessionMemAck <- 1   // tell SessionDispatcher we're done with the data
	return ok
}

//=====================================================================================
// SYNOPSIS:
// 		showAdminButton determines whether or not the Admin button needs to appear
// 		on the menu.
// PARAMS:
//		token - db.Session token
// RETURNS:
//		true  - if the admin button should be shown
//      false - if it should not
//=====================================================================================
func pvtShowAdminButton(s *db.Session) bool {
	for i := 0; i < len(s.PMap.Urole.Perms); i++ {
		if s.PMap.Urole.Perms[i].Perm&db.PERMCREATE != 0 {
			return true
		}
	}
	return false
}

func showAdminButton(token string) bool {
	s, ok := db.Sessions[token]
	if !ok {
		fmt.Printf("showAdminButton:  Could not find db.Session for %s\n", token)
		return false
	}
	db.SessionManager.ReqSessionMem <- 1    // ask to access the shared mem, blocks until granted
	<-db.SessionManager.ReqSessionMemAck    // make sure we got it
	ok = pvtShowAdminButton(s)              //we have the memory, do the work
	db.SessionManager.ReqSessionMemAck <- 1 // tell SessionDispatcher we're done with the data
	return ok
}

// Privileged function allowing one user to become another user. This is meant
// to be used by Administrators or User Support personnel.
func sessionBecome(s *db.Session, uid int64) {
	var d db.PersonDetail
	d.Reports = make([]db.Person, 0)
	d.UID = uid
	adminReadDetails(&d)

	s.Firstname = d.FirstName
	if 0 < len(d.PreferredName) {
		s.Firstname = d.PreferredName
	}
	s.UID = int64(uid)
	s.Username = d.UserName
	s.ImageURL = db.GetImageLocation(uid)
	db.GetRoleInfo(d.RID, &s.PMap)

	if db.Authz.SecurityDebug {
		for i := 0; i < len(s.PMap.Urole.Perms); i++ {
			ulog("f: %s,  perm: %02x\n", s.PMap.Urole.Perms[i].Field, s.PMap.Urole.Perms[i].Perm)
		}
	}

	ulog("user %d to BECOME user %d", s.UIDorig, s.UID)
}
