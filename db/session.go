package db

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"phonebook/lib"
	"strings"
	"time"
)

// SessionCookie defines the struct for the database table where session
// cookies are managed.
type SessionCookie struct {
	UID       int64     // uid of the user
	UserName  string    // username for the user
	Cookie    string    // the cookie value
	Expire    time.Time // that timestamp when it expires
	UserAgent string    // client identifier
	IP        string    // end user's IP address
}

// SessionManager is the struct containing key values for the Session
// management infrastructure
var SessionManager struct {
	ReqSessionMem      chan int64 // request to access Session data memory
	ReqSessionMemAck   chan int64 // done with Session datamemory
	SessionCleanupTime time.Duration
	SecurityDebug      bool
	SessionTimeout     time.Duration
	db                 *sql.DB        // the database connection
	ZoneUTC            *time.Location // what timezone should the server use?
}

// Session is the generic Session
type Session struct {
	Token       string      // this is the md5 hash, unique id, the cookie value
	Username    string      // associated username
	Firstname   string      // user's first name
	UID         int64       // user's db uid
	UIDorig     int64       // original uid (for use with method sessionBecome())
	CoCode      int64       // logged in user's company
	ImageURL    string      // user's picture
	Expire      time.Time   // when does the cookie expire
	Breadcrumbs []lib.Crumb // where is the user in the screen hierarchy
	PMap        PermMaps    // user's role and associated maps
	IP          string      // user's IP address
	UserAgent   string      // the user's client
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
	SessionManager.ReqSessionMem = make(chan int64)
	SessionManager.ReqSessionMemAck = make(chan int64)
	SessionManager.SessionCleanupTime = clean
	SessionManager.SessionTimeout = timeout
	Sessions = make(map[string]*Session)
	SessionManager.SecurityDebug = debug
	SessionManager.db = db
	SessionManager.ZoneUTC, err = time.LoadLocation("UTC")
	if err != nil {
		lib.Ulog("InitSessionManager: error reading timezone: %s", err.Error())
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
	return fmt.Sprintf("User(%s) Name(%s) UID(%d) IP(%s) Token(%s)  Role(%s)",
		s.Username, s.Firstname, s.UID, s.IP, s.Token, s.PMap.Urole.Name)
}

// DumpSessions prints out the session map for debugging
//-----------------------------------------------------------------------------
func DumpSessions() {
	fmt.Printf("\nDIRECTORY INTERNAL SESSION TABLE\n")
	i := 0
	for _, v := range Sessions {
		fmt.Printf("%2d. %s\n", i, v.ToString())
		i++
	}
	fmt.Printf("END\n")

}

// ReUpCookieTime updates the timeout time associated with a session cookie.
//-----------------------------------------------------------------------------
func ReUpCookieTime(s *Session) {
	t := time.Now().Add(SessionManager.SessionTimeout * time.Minute)
	SessionManager.ReqSessionMem <- 1    // ask to access the shared mem, blocks until granted
	<-SessionManager.ReqSessionMemAck    // make sure we got it
	s.Expire = t                         // update the Session information
	SessionManager.ReqSessionMemAck <- 1 // tell SessionDispatcher we're done with the data
}

// Refresh updates the cookie and Session with a new expire time.
//-----------------------------------------------------------------------------
func (s *Session) Refresh(w http.ResponseWriter, r *http.Request) int64 {
	// lib.Console("Entered Session.Refresh\n")
	cookie, err := r.Cookie(SessionCookieName)
	if nil != cookie && err == nil {
		// lib.Console("Cookie found: %s\n", cookie.Value)
		ReUpCookieTime(s)
		cookie.Expires = s.Expire
		// lib.Console("Setting expire time to: %v\n", cookie.Expires)
		cookie.Path = "/"
		http.SetCookie(w, cookie)
		// lib.Console("Session.Expire = %v\n", s.Expire)
		UpdateSessionCookieDB(s)
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
func NewSessionFromCookie(c *SessionCookie) *Session {
	var email, passhash, firstname, preferredname string
	var uid, RID int64

	err := PrepStmts.LoginInfo.QueryRow(c.UserName).Scan(&uid, &firstname, &preferredname, &email, &passhash, &RID)
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
func NewSession(c *SessionCookie, firstname string, rid int64) *Session {
	return pvtNewSession(c, firstname, rid, true)
}

// pvtNewSession creates a new session, updates the session table if necessary,
// adds the new session to the in-memory session table, and returns the session
//-----------------------------------------------------------------------------
func pvtNewSession(c *SessionCookie, firstname string, rid int64, updateSessionTable bool) *Session {
	// lib.Ulog("Entering NewSession: %s (%d)\n", username, uid)
	uid := int64(c.UID)
	s := new(Session)
	s.Token = c.Cookie
	s.Username = c.UserName
	s.Firstname = firstname
	s.UID = c.UID
	s.UIDorig = c.UID
	s.ImageURL = GetImageLocation(uid)
	s.Breadcrumbs = make([]lib.Crumb, 0)
	s.Expire = c.Expire
	s.IP = c.IP
	s.UserAgent = c.UserAgent
	GetRoleInfo(rid, &s.PMap)

	if Authz.SecurityDebug {
		for i := 0; i < len(s.PMap.Urole.Perms); i++ {
			lib.Ulog("f: %s,  perm: %02x\n", s.PMap.Urole.Perms[i].Field, s.PMap.Urole.Perms[i].Perm)
		}
	}

	var d PersonDetail
	d.UID = uid

	err := SessionManager.db.QueryRow(fmt.Sprintf("SELECT CoCode FROM people WHERE UID=%d", uid)).Scan(&s.CoCode)
	if nil != err {
		lib.Ulog("Unable to read CoCode for userid=%d,  err = %v\n", uid, err)
	}

	if updateSessionTable {
		lib.Console("JUST BEFORE InsertSessionCookie: s.IP = %s, s.UserAgent = %s\n", s.IP, s.UserAgent)
		err = InsertSessionCookieDB(s)
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
	fmt.Printf("Session being deleted: %s\n", s.ToString())
	// fmt.Printf("sess.Sessions before delete:\n")
	// DumpSessions()

	if err := DeleteSessionCookie(s.Token); err != nil {
		lib.Ulog("Error deleting session cookie: %s\n", err.Error())
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
	// DumpSessions()
}

//=====================================================================================
// pvtElemPermsAny determines whether or not the Session has permissions to perform the
// requested operations.  NOTE:  This interface does check the UID to fully cover
// permissions db.PERMOWNERVIEW or db.PERMOWNERMOD. This must be done at a higher level.
//
// ARGS:
//   ent  = which element? db.ELEMPERSON, db.ELEMCOMPANY, db.ELEMCLASS, ...
//   perm = logical or of the desired permissions.  Example db.PERMCREATE | db.PERMMOD
//
// RETURNS:
//   true if there are ANY fields for the specified element for
//   with the requested permission.
//=====================================================================================
func pvtElemPermsAny(s *Session, elem int64, perm int64) bool {
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
func (s *Session) ElemPermsAny(elem int64, perm int64) bool {
	SessionManager.ReqSessionMem <- 1    // ask to access the shared mem, blocks until granted
	<-SessionManager.ReqSessionMemAck    // make sure we got it
	ok := pvtElemPermsAny(s, elem, perm) // look for perms
	SessionManager.ReqSessionMemAck <- 1 // tell SessionDispatcher we're done with the data
	return ok
}

//=====================================================================================
// pvtElemPermsAll determines whether or not the Session has permissions to perform all
// requested operations.  NOTE: This interface does check the UID to fully cover
// permissions db.PERMOWNERVIEW or db.PERMOWNERMOD. This must be done at a higher level.
//
// ARGS:
//   ent  = which element? db.ELEMPERSON, db.ELEMCOMPANY, db.ELEMCLASS, ...
//   perm = logical or of the desired permissions.  Example db.PERMCREATE | db.PERMMOD
//
// RETURNS:
//   true if ALL permission fields for the specified element are present
//=====================================================================================
func pvtElemPermsAll(s *Session, elem int64, perm int64) bool {
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
func (s *Session) ElemPermsAll(elem int64, perm int64) bool {
	SessionManager.ReqSessionMem <- 1    // ask to access the shared mem, blocks until granted
	<-SessionManager.ReqSessionMemAck    // make sure we got it
	ok := pvtElemPermsAll(s, elem, perm) // look for perms
	SessionManager.ReqSessionMemAck <- 1 // tell SessionDispatcher we're done with the data
	return ok
}

// ValidateSessionCookie verifies the existence of a cookie and will update its
// timeout value if update=true
//
// INPUTS
//  cval   = value of cookie string for "air"
//  update = if true, increase the timeout value by the SessionTimeout amount
//
// RETURNS
//  any errors encountered
//----------------------------------------------------------------------------------------------
func ValidateSessionCookie(r *http.Request, update bool) (bool, error) {
	// funcname := "ValidateCookie"
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		return false, err
	}
	cval := cookie.Value
	// lib.Console("validate Cookie:  %s\n", cval)

	c, err := GetSessionCookie(cval)
	if err != nil {
		return false, err
	}

	// lib.Console("Found session Cookie: UID=%d, UserName=%s, Expire=%v\n", c.UID, c.UserName, c.Expire)

	if c.UID < 1 {
		return false, nil
	}

	// lib.Console("%s: Cookie found:  c.UID = %d\n", funcname, c.UID)
	//------------------------------------------------------------------
	// if the request calls for the timestamp to be updated, do so now
	// that we know it exists.
	//------------------------------------------------------------------
	if update {
		s, ok := Sessions[c.Cookie]
		if !ok {
			//----------------------------------------------------------------
			// This means that the Cookie was found in the database but not
			// in memory. The most likely reason for this is that phonebook
			// was restarted.  In any case, we need to add this Cookie to the
			// in memory sessions...
			//----------------------------------------------------------------
			s = NewSessionFromCookie(&c) // we don't need the return value, we just need the session to be put into memory
			// lib.Console("Session was not found in memory.  Adding it to memory.  Cookie = %s\n", c.Cookie)
		}
		s.Expire = s.Expire.Add(SessionManager.SessionTimeout * time.Minute)
		if err = UpdateSessionCookieDB(s); err != nil {
			return true, err
		}
		// lib.Console("UPDATED SESSION Cookie TIMEOUT TIME\n")
	}
	return true, nil

}

// GetSession returns the session based on the cookie in the supplied
// HTTP connection. If the "air" cookie is valid, it will either find the
// existing session or create a new session.
//
// INPUT
//  ctx database context
//  w - http writer to client
//  r - the request where we look for the cookie
//  c - pointer to the validate cookie response.  If UID > 0 it means that
//      the cookie has already been validated and that the other fields are
//      valid -- we don't need to make another call to the directory server.
//
// RETURNS
//  session - pointer to the new session
//  error   - any error encountered
//-----------------------------------------------------------------------------
func GetSession(ctx context.Context, w http.ResponseWriter, r *http.Request) (*Session, error) {
	// funcname := "GetSession"
	// var b AIRAuthenticateResponse
	var ok bool
	var sess *Session

	// util.Console("GetSession 1\n")
	// util.Console("\nSession Table:\n")
	// DumpSessions()
	// util.Console("\n")
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		// util.Console("GetSession 2\n")
		if strings.Contains(err.Error(), "cookie not present") {
			// util.Console("GetSession 3\n")
			return nil, nil
		}
		// util.Console("GetSession 4\n")
		return nil, err
	}
	// util.Console("GetSession 5\n")
	sess, ok = Sessions[cookie.Value]
	// if !ok || sess == nil {
	// 	var b ValidateCookieResponse
	// 	util.Console("GetSession 6\n")
	//
	// 	b, err = ValidateSessionCookie(cookie.Value, true)
	// 	if err != nil {
	// 		util.Console("GetSession 7\n")
	// 		return sess, err
	// 	}
	// 	util.Console("ValidateSessionCookie returned b = %#v\n", b)
	// 	util.Console("Directory Service Expire time = %s\n", time.Time(b.Expire).Format(util.RRDATETIMEINPFMT))
	// 	sess, err = CreateSession(ctx, &b)
	// 	util.Console("GetSession 8\n")
	// 	if err != nil {
	// 		util.Console("GetSession 9\n")
	// 		return nil, err
	// 	}
	// 	util.Console("*** NEW SESSION CREATED ***\n")
	// }
	if ok && sess != nil {
		sess.Refresh(w, r)
	}
	return sess, nil
}
