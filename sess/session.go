package sess

import (
	"database/sql"
	"fmt"
	"net/http"
	"phonebook/authz"
	"phonebook/db"
	"phonebook/lib"
	"phonebook/ui"
	"time"
)

// SessionManager is the struct containing key values for the Session
// management infrastructure
var SessionManager struct {
	ReqSessionMem      chan int // request to access Session data memory
	ReqSessionMemAck   chan int // done with Session datamemory
	SessionCleanupTime time.Duration
	SecurityDebug      bool
	SessionTimeout     time.Duration
	db                 *sql.DB        // the database connection
	ZoneUTC            *time.Location // what timezone should the server use?
}

// Session is the generic Session
type Session struct {
	Token        string         // this is the md5 hash, unique id
	Username     string         // associated username
	Firstname    string         // user's first name
	UID          int64          // user's db uid
	UIDorig      int64          // original uid (for use with method sessionBecome())
	UsernameOrig string         // original username
	CoCode       int            // logged in user's company
	ImageURL     string         // user's picture
	Expire       time.Time      // when does the cookie expire
	Breadcrumbs  []ui.Crumb     // where is the user in the screen hierarchy
	PMap         authz.PermMaps // user's role and associated maps
}

// Sessions is the map of Session structs indexed by the SessionKey (the browser cookie value)
var Sessions map[string]*Session

// SessionGet returns the in memory session with the supplied token
func SessionGet(token string) (*Session, bool) {
	s, ok := Sessions[token]
	return s, ok
}

// InitSessionManager initializes the Session infrastructure
//
// INPUTS
//  none
//
// RETURNS
//  nothing
//-----------------------------------------------------------------------------
func InitSessionManager(clean, timeout time.Duration, db *sql.DB, debug bool) {
	var err error
	SessionManager.ReqSessionMem = make(chan int)
	SessionManager.ReqSessionMemAck = make(chan int)
	SessionManager.SessionCleanupTime = clean
	SessionManager.SessionTimeout = timeout
	Sessions = make(map[string]*Session)
	SessionManager.SecurityDebug = debug
	SessionManager.db = db
	SessionManager.ZoneUTC, err = time.LoadLocation("UTC")
	if err != nil {
		lib.Ulog("InitSessionManager: error reading timezone: ", err.Error())
	}
	go SessionDispatcher()
	go SessionCleanup()
	go ExpiredCookieCleaner()
}

// SessionDispatcher controls access to shared memory.
//-----------------------------------------------------------------------------
func SessionDispatcher() {
	for {
		select {
		case <-SessionManager.ReqSessionMem:
			SessionManager.ReqSessionMemAck <- 1 // tell caller go ahead
			<-SessionManager.ReqSessionMemAck    // block until caller is done with mem
		}
	}
}

// SessionCleanup periodically spins through the list of Sessions
// and removes any which have timed out.
//-----------------------------------------------------------------------------
func SessionCleanup() {
	for {
		select {
		case <-time.After(SessionManager.SessionCleanupTime * time.Minute):
			SessionManager.ReqSessionMem <- 1  // ask to access the shared mem, blocks until granted
			<-SessionManager.ReqSessionMemAck  // make sure we got it
			ss := make(map[string]*Session, 0) // here's the new Session list
			n := 0                             // total number removed
			for k, v := range Sessions {       // look at every Session
				if time.Now().After(v.Expire) { // if it's still active...
					n++ // removed another
				} else {
					ss[k] = v // ...copy it to the new list
				}
			}
			Sessions = ss                        // set the new list
			SessionManager.ReqSessionMemAck <- 1 // tell SessionDispatcher we're done with the data
			//fmt.Printf("SessionCleanup completed. %d removed. Current Session list size = %d\n", n, len(Sessions))
		}
	}
}

// ToString is the stringger for Session variables
//-----------------------------------------------------------------------------
func (s *Session) ToString() string {
	if nil == s {
		return "nil"
	}
	return fmt.Sprintf("User(%s) Name(%s) UID(%d) Token(%s)  Role(%s)",
		s.Username, s.Firstname, s.UID, s.Token, s.PMap.Urole.Name)
}

// DumpSessions prints out the session map for debugging
//-----------------------------------------------------------------------------
func DumpSessions() {
	i := 0
	for _, v := range Sessions {
		fmt.Printf("%2d. %s\n", i, v.ToString())
		i++
	}
}

// Refresh updates the cookie and Session with a new expire time.
//-----------------------------------------------------------------------------
func (s *Session) Refresh(w http.ResponseWriter, r *http.Request) int {
	lib.Console("Entered Session.Refresh\n")
	cookie, err := r.Cookie(SessionCookieName)
	if nil != cookie && err == nil {
		lib.Console("Cookie found: %s\n", cookie.Value)
		cookie.Expires = time.Now().Add(SessionManager.SessionTimeout * time.Minute)
		lib.Console("Setting expire time to: %v\n", cookie.Expires)
		SessionManager.ReqSessionMem <- 1    // ask to access the shared mem, blocks until granted
		<-SessionManager.ReqSessionMemAck    // make sure we got it
		s.Expire = cookie.Expires            // update the Session information
		SessionManager.ReqSessionMemAck <- 1 // tell SessionDispatcher we're done with the data
		cookie.Path = "/"
		http.SetCookie(w, cookie)
		lib.Console("Session.Expire = %v\n", s.Expire)
		UpdateSessionCookie(s)
		return 0
	}
	return 1
}

// NewSessionFromCookie - Creates a new in-memory session based on a cookie
// that exists in the session table. There are several circumstances which
// cause us to get here:
//		a) the login may have come from a separate running instance of
//		   this server
//		b) this server may have been restarted
//		c) the user may have logged into another AIR app in the suite
//
// RETURNS
//  *Session - it will be empty if there was any problem. Otherwise it will
//		have all required session information
//-----------------------------------------------------------------------------
func NewSessionFromCookie(c *db.SessionCookie) *Session {
	var email, passhash, firstname, preferredname string
	var uid, RID int

	err := db.PrepStmts.LoginInfo.QueryRow(c.UserName).Scan(&uid, &firstname, &preferredname, &email, &passhash, &RID)
	if err != nil {
		s := new(Session)
		lib.Ulog("Error reading person with username %s: %s", c.UserName, err.Error())
		return s // it's empty because of the error
	}
	if len(preferredname) > 0 {
		firstname = preferredname
	}
	return pvtNewSession(c, firstname, RID, false)
}

// NewSession returns a new session.  This entry point requires an update
// to the session table.
//-----------------------------------------------------------------------------
func NewSession(c *db.SessionCookie, firstname string, rid int) *Session {
	return pvtNewSession(c, firstname, rid, true)
}

// pvtNewSession creates a new session, updates the session table if necessary,
// adds the new session to the in-memory session table, and returns the session
//-----------------------------------------------------------------------------
func pvtNewSession(c *db.SessionCookie, firstname string, rid int, updateSessionTable bool) *Session {
	// lib.Ulog("Entering NewSession: %s (%d)\n", username, uid)
	uid := int(c.UID)
	s := new(Session)
	s.Token = c.Cookie
	s.Username = c.UserName
	s.Firstname = firstname
	s.UID = c.UID
	s.UIDorig = c.UID
	s.ImageURL = ui.GetImageFilename(uid)
	s.Breadcrumbs = make([]ui.Crumb, 0)
	s.Expire = c.Expire
	authz.GetRoleInfo(rid, &s.PMap)

	if authz.Authz.SecurityDebug {
		for i := 0; i < len(s.PMap.Urole.Perms); i++ {
			lib.Ulog("f: %s,  perm: %02x\n", s.PMap.Urole.Perms[i].Field, s.PMap.Urole.Perms[i].Perm)
		}
	}

	var d db.PersonDetail
	d.UID = uid

	err := SessionManager.db.QueryRow(fmt.Sprintf("SELECT CoCode FROM people WHERE UID=%d", uid)).Scan(&s.CoCode)
	if nil != err {
		lib.Ulog("Unable to read CoCode for userid=%d,  err = %v\n", uid, err)
	}

	if updateSessionTable {
		err = InsertSessionCookie(s)
		if err != nil {
			lib.Ulog("Unable to save session for UID = %d to database,  err = %s\n", uid, err.Error())
		}
	}

	SessionManager.ReqSessionMem <- 1 // ask to access the shared mem, blocks until granted
	<-SessionManager.ReqSessionMemAck // make sure we got it
	Sessions[c.Cookie] = s
	SessionManager.ReqSessionMemAck <- 1 // tell SessionDispatcher we're done with the data

	return s
}

// SessionDelete removes the supplied sess.Session.
// If there is a better idiomatic way to do this, please let me know.
// It also removes the session from the db sessions table.
//-----------------------------------------------------------------------------
func SessionDelete(s *Session) {
	// fmt.Printf("Session being deleted: %s\n", s.ToString())
	// fmt.Printf("sess.Sessions before delete:\n")
	// dumpSessions()

	if err := db.DeleteSessionCookie(s.Token); err != nil {
		lib.Ulog("Error deleteing session cookie: %s\n", err.Error())
	}

	ss := make(map[string]*Session, 0)

	SessionManager.ReqSessionMem <- 1 // ask to access the shared mem, blocks until granted
	<-SessionManager.ReqSessionMemAck // make sure we got it
	for k, v := range Sessions {
		if s.Token != k {
			ss[k] = v
		}
	}
	Sessions = ss
	SessionManager.ReqSessionMemAck <- 1 // tell SessionDispatcher we're done with the data
	// fmt.Printf("sess.Sessions after delete:\n")
	// dumpSessions()
}

//=====================================================================================
// pvtElemPermsAny determines whether or not the Session has permissions to perform the
// requested operations.  NOTE:  This interface does check the UID to fully cover
// permissions authz.PERMOWNERVIEW or authz.PERMOWNERMOD. This must be done at a higher level.
//
// ARGS:
//   ent  = which element? authz.ELEMPERSON, authz.ELEMCOMPANY, authz.ELEMCLASS, ...
//   perm = logical or of the desired permissions.  Example authz.PERMCREATE | authz.PERMMOD
//
// RETURNS:
//   true if there are ANY fields for the specified element for
//   with the requested permission.
//=====================================================================================
func pvtElemPermsAny(s *Session, elem int, perm int) bool {
	// lib.Ulog("elemPermsAny:  elem=%d, perm = 0x%02x\n", elem, perm)
	for i := 0; i < len(s.PMap.Urole.Perms); i++ {
		// lib.Ulog("s.PMap.Urole.Perms[%d].Elem = %d\n", i, s.PMap.Urole.Perms[i].Elem)
		if s.PMap.Urole.Perms[i].Elem == elem {
			res := s.PMap.Urole.Perms[i].Perm & perm
			// lib.Ulog("fieldname: %s  s.PMap.Urole.Perms[%d].Perm = 0x%02x, s.PMap.Urole.Perms[%d].Perm & perm = 0x%02x\n", s.PMap.Urole.Perms[i].Field, i, s.PMap.Urole.Perms[i].Elem, i, res)
			if res != 0 { // if any of the permissions exist
				// lib.Ulog("return true") // we're good to go for this check
				return true
			}
		}
	}
	// lib.Ulog("return false")
	return false
}

// ElemPermsAny returns true if the session as permissions to perform any of the
// requested actions. Otherwise it return s false
//-----------------------------------------------------------------------------
func (s *Session) ElemPermsAny(elem int, perm int) bool {
	SessionManager.ReqSessionMem <- 1    // ask to access the shared mem, blocks until granted
	<-SessionManager.ReqSessionMemAck    // make sure we got it
	ok := pvtElemPermsAny(s, elem, perm) // look for perms
	SessionManager.ReqSessionMemAck <- 1 // tell SessionDispatcher we're done with the data
	return ok
}

//=====================================================================================
// pvtElemPermsAll determines whether or not the Session has permissions to perform all
// requested operations.  NOTE: This interface does check the UID to fully cover
// permissions authz.PERMOWNERVIEW or authz.PERMOWNERMOD. This must be done at a higher level.
//
// ARGS:
//   ent  = which element? authz.ELEMPERSON, authz.ELEMCOMPANY, authz.ELEMCLASS, ...
//   perm = logical or of the desired permissions.  Example authz.PERMCREATE | authz.PERMMOD
//
// RETURNS:
//   true if ALL permission fields for the specified element are present
//=====================================================================================
func pvtElemPermsAll(s *Session, elem int, perm int) bool {
	// lib.Ulog("elemPermsAll:  elem=%d, perm = 0x%02x\n", elem, perm)
	for i := 0; i < len(s.PMap.Urole.Perms); i++ {
		// lib.Ulog("s.PMap.Urole.Perms[%d].Elem = %d\n", i, s.PMap.Urole.Perms[i].Elem)
		if s.PMap.Urole.Perms[i].Elem == elem {
			res := s.PMap.Urole.Perms[i].Perm & perm
			// lib.Ulog("fieldname: %s  s.PMap.Urole.Perms[%d].Perm = 0x%02x, s.PMap.Urole.Perms[%d].Perm & perm = 0x%02x\n", s.PMap.Urole.Perms[i].Field, i, s.PMap.Urole.Perms[i].Elem, i, res)
			if res == perm { // if all bits are present, res will match perm
				// lib.Ulog("return true") // we're good to go for this check
				return true
			}
		}
	}
	// lib.Ulog("return false")
	return false
}

// ElemPermsAll returns true if the supplied session has permissions to perform all requested
// operations. Otherwise it returns false.
//----------------------------------------------------------------------------------------------
func (s *Session) ElemPermsAll(elem int, perm int) bool {
	SessionManager.ReqSessionMem <- 1    // ask to access the shared mem, blocks until granted
	<-SessionManager.ReqSessionMemAck    // make sure we got it
	ok := pvtElemPermsAll(s, elem, perm) // look for perms
	SessionManager.ReqSessionMemAck <- 1 // tell SessionDispatcher we're done with the data
	return ok
}
