package ws

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"phonebook/db"
	"phonebook/lib"
	"strings"
	"time"
)

// AuthenticateData is the struct with the username and password
// used for authentication.  This is the data that the user sends
// to the phonebook server.
type AuthenticateData struct {
	User       string `json:"user"`
	Pass       string `json:"pass"`
	FLAGS      uint64 `json:"flags"`
	UserAgent  string `json:"useragent"`
	RemoteAddr string `json:"remoteaddr"`
}

// AuthSuccessResponse will be the response structure used when
// authentication is successful.
type AuthSuccessResponse struct {
	Status   string `json:"status"`
	UID      int64  `json:"uid"`
	Name     string `json:"Name"`
	ImageURL string `json:"ImageURL"`
	Token    string `json:"Token"`
	Expire   string `json:"Expire"` // DATETIMEFMT in this format "2006-01-02T15:04 "
}

// ValidateCookie describes the data sent by an AIR app to check
// whether or not a db.Cookie value is valid.
type ValidateCookie struct {
	CookieVal string `json:"cookieval"`
	IP        string `json:"ip"`
	UserAgent string `json:"useragent"`
	FLAGS     uint64 `json:"flags"`
}

const (
	// JSONDATETIME is format roller and others use for datetime over json
	JSONDATETIME = "2006-01-02T15:04:00Z"
)

func getImageURL(UID int64) (string, error) {
	fname := ""
	if UID > 0 {
		m, err := filepath.Glob(fmt.Sprintf("./pictures/%d.*", UID))
		if nil != err {
			return fname, err
		}
		if len(m) > 0 {
			fname = m[0]
		}
	}
	return fname, nil
}

// NewSessionFromCookie is a wrapper around the sess version of this so that
// no db access is needed
func NewSessionFromCookie(c *db.SessionCookie) *db.Session {
	s := db.NewSessionFromCookie(c)
	//s.ImageURL = db.GetImageLocation(uid)
	return s
}

// SvcAuthenticate generates a password hash from the supplied POST info and
//     along with the user name compares it to what is in the database. If
//     there is a match, then the response is {status: success}.  If it fails
//     then the response is {status: failed}.  No indication will be given
//     indicating whether the username is not recognized or the password for
//     the supplied username is not correct.
//
// INPUTS:
//     w = file descriptor to write result
//     r = http requrest
//     d = pointer to data parsed by service dispatcher
//
// RETURNS:
//     nothing at this time
//-----------------------------------------------------------------------------
func SvcAuthenticate(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	var funcname = "SvcAuthenticate"
	var err error
	var foo AuthenticateData
	var UID int64
	var Name string

	// lib.Console("Entered %s\n", funcname)

	data := []byte(d.data)
	if err = json.Unmarshal(data, &foo); err != nil {
		e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
		SvcErrorReturn(w, e, funcname)
		goto exit
	}

	lib.Console("User = %s, Pass = %s\n", foo.User, foo.Pass)
	lib.Console("IP = %s, UserAgent = %s\n", foo.RemoteAddr, foo.UserAgent)

	//------------------------------------------------------------------------
	// We can detect the forwarded-for value, but it is not used. In this
	// case, another server has sent the login on behalf of a client further
	// back. So, here, we want to use exactly the value that has been sent
	// by the requester.
	//------------------------------------------------------------------------
	// fwdaddr := r.Header.Get("X-Forwarded-For")
	// lib.Console("X-Forwarded-For value = %q\n", fwdaddr)

	lib.Console("svcAuth A\n")
	UID, Name, err = DoAuthentication(foo.User, foo.Pass)
	if err != nil {
		lib.Console("svcAuth B\n")
		SvcErrorReturn(w, err, funcname)
		goto exit
	}
	lib.Console("svcAuth C\n")

	//----------------------------------------------------------------------------------
	// If UID > 0 then the username and password match.  So, we get the user a session.
	// Use an existing setup if the user already has a session from the endpoint that
	// submitted this request. Otherwise, create a new one.
	//----------------------------------------------------------------------------------
	if UID > 0 {
		lib.Console("svcAuth D\n")
		imageProfilePath := db.GetImageLocation(UID) // we need this in multiple cases

		//------------------------------------------------------------------------
		// Before generating a new db.Cookie, see if this user / useragent / ip
		// combination already has a valid db.Cookie.
		//------------------------------------------------------------------------
		c, err := db.FindMatchingSessionCookie(foo.User, foo.RemoteAddr, foo.UserAgent)
		if err != nil {
			lib.Console("svcAuth E\n")
			err := fmt.Errorf("error finding db.Cookie: %s", err.Error())
			SvcErrorReturn(w, err, funcname)
			goto exit
		}
		lib.Console("svcAuth F....   len(c.Cookie) = %d, c.UserName = %s\n", len(c.Cookie), c.UserName)
		if len(c.Cookie) > 0 && foo.User == c.UserName {
			lib.Console("svcAuth G\n")
			//-----------------------------------------------------------------------
			// This user already has a db.Cookie in the same useragent. Just update
			// the existing info and return it...
			//-----------------------------------------------------------------------
			g := AuthSuccessResponse{
				Status:   "success",
				UID:      UID,
				Name:     Name,
				ImageURL: imageProfilePath,
				Token:    c.Cookie,
				Expire:   c.Expire.In(db.SessionManager.ZoneUTC).Format(JSONDATETIME),
			}
			lib.Console("svcAuth H\n")
			//--------------------------------
			// get the associated session...
			//--------------------------------
			s, ok := db.Sessions[c.Cookie]
			if !ok { // this could possibly happen if the timeing is *just* right, but we need to create it
				lib.Console("svcAuth I\n")
				s = db.NewSessionFromCookie(&c)
			}
			//----------------------------------------------------
			// update its timeout now that it has been used...
			//----------------------------------------------------
			lib.Console("svcAuth J\n")
			db.ReUpCookieTime(s)
			db.UpdateSessionCookieDB(s)
			g.Expire = s.Expire.In(db.SessionManager.ZoneUTC).Format(JSONDATETIME)

			//----------------------------------------------------
			// And now we're done... return the response
			//----------------------------------------------------
			lib.Console("g = %#v\n", g)
			SvcWriteResponse(&g, w)
			lib.Ulog("user %s successfully piggybacked on existing session\n", foo.User)
			goto exit
		}
		lib.Console("svcAuth K\n")

		//---------------------------------------------------------------------------
		// If we hit this point, it means that there currently is no entry in the
		// session table for the this user. Create one...
		//---------------------------------------------------------------------------
		c = db.GenerateSessionCookie(UID, foo.User, foo.UserAgent, foo.RemoteAddr)
		db.NewSessionFromCookie(&c) // we don't need the return value, we just need the session to be put into memory

		g := AuthSuccessResponse{
			Status:   "success",
			UID:      c.UID,
			Name:     Name,
			ImageURL: imageProfilePath,
			Token:    c.Cookie,
			Expire:   c.Expire.In(db.SessionManager.ZoneUTC).Format(JSONDATETIME),
		}

		lib.Console("svcAuth L  (username: %s, user agent: %s)\n", c.UserName, c.UserAgent)
		lib.Console("g = %#v\n", g)
		SvcWriteResponse(&g, w)
		lib.Ulog("user %s successfully logged in\n", foo.User)

		err = db.InsertSessionCookie(c.UID, c.UserName, c.Cookie, &c.Expire, c.UserAgent, c.IP)
		if err != nil {
			lib.Console("svcAuth M\n")
			err = fmt.Errorf("error inserting db.Cookie into sessiondb: %s", err.Error())
			SvcErrorReturn(w, err, funcname)
			goto exit
		}
		lib.Console("svcAuth N\n")
		goto exit
	}
	lib.Console("svcAuth O\n")
	err = fmt.Errorf("login failed")
	SvcErrorReturn(w, err, funcname)

exit:
	lib.Console("svcAuth P\n")
	db.DumpSessions()
	db.DumpSessionCookies()
}

// DoAuthentication builds a password hash out of the supplied user and
// password information. It then looks up the user in the database. If the
// password hashes match, then the login is successful
//
// INPUTS:
//  User = username
//  Pass = user's password
//
// RETURNS:
//  int64 =  UID if the login was successful, 0 otherwise
//  name  = user's first name (or preferred name if it exists)
//  error = any error encountered
//-----------------------------------------------------------------------------
func DoAuthentication(User, Pass string) (int64, string, error) {
	myusername := strings.ToLower(User)
	password := []byte(Pass)
	sha := sha512.Sum512(password)
	mypasshash := fmt.Sprintf("%x", sha)
	var err error

	// lookup the user
	q := fmt.Sprintf("SELECT UID,FirstName,PreferredName,Status,passhash FROM people WHERE UserName=%q", myusername)
	var passhash string
	var UID int64
	var Status int
	var first, preferred string
	err = SvcCtx.db.QueryRow(q).Scan(&UID, &first, &preferred, &Status, &passhash)
	if err != nil {
		return int64(0), first, err
	}
	if len(preferred) > 0 {
		first = preferred
	}
	switch Status {
	case 0:
		err = fmt.Errorf("account is inactive")
		return UID, first, err
	case 1:
		if passhash != mypasshash {
			err = fmt.Errorf("login failed")
			return int64(0), first, err
		}
		return UID, first, nil // login is successful
	default:
		err = fmt.Errorf("unrecognized user status: %d", Status)
		return UID, first, err
	}
}

// SvcValidateCookie is called by an AIR app when it finds the air db.Cookie
// but has no associated session.  If the db.Cookie is in our sessions table
// then we send back a success response with the same info we include in
// a successful login. Otherwise, we send an appropriate error response
//
// FLAGS of the data refine the operation:
//
//     1<<0  -  if this bit is set it means just send back success
//              or failure, do not send back all other information
//              associated with the session containing the db.Cookie.
//              The response will come back with Status: "success"
//              if the db.Cookie was found, "failure" if the db.Cookie
//              was not found, or "error" if an error was encountered.
//
//     1<<1  =  Update the timeout time for the db.Cookie and session -
//              Increments by phonebook session manager timeout time.
//
// INPUTS:
//     w = file descriptor to write result
//     r = http requrest
//     d = pointer to data parsed by service dispatcher
//
// RETURNS:
//     nothing at this time
//-----------------------------------------------------------------------------
func SvcValidateCookie(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	var funcname = "SvcValidateCookie"
	var err error
	var foo ValidateCookie
	var c db.SessionCookie
	var imageProfilePath string
	var resp string
	var g AuthSuccessResponse

	lib.Console("Entered %s\n", funcname)

	data := []byte(d.data)
	if err = json.Unmarshal(data, &foo); err != nil {
		e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
		SvcErrorReturn(w, e, funcname)
		goto exit1
	}

	lib.Console("request for session db.Cookie:  %s, IP = %s, UserAgent = %s\n", foo.CookieVal, foo.IP, foo.UserAgent)
	c, err = db.GetSessionCookie(foo.CookieVal)
	if err != nil {
		lib.Ulog("signinHandler: error getting session db.Cookie: %s\n", err.Error())
	}
	lib.Console("Found session db.Cookie: UID=%d, UserName=%s, Expire=%s\n", c.UID, c.UserName, c.Expire.Format("JSONDATETIME"))
	lib.Console("                      IP = %s,  UserAgent = %s\n", c.IP, c.UserAgent)

	resp = "failure"
	if c.UID > 0 {
		lib.Console("%s: db.Cookie found:  c.UID = %d\n", funcname, c.UID)
		resp = "success"
		//------------------------------------------------------------------
		// if the request calls for the timestamp to be updated, do so now
		// that we know it exists.
		//------------------------------------------------------------------
		if foo.FLAGS&2 > 0 {
			s, ok := db.Sessions[c.Cookie]
			if !ok {
				//----------------------------------------------------------------
				// This means that the db.Cookie was found in the database but not
				// in memory. The most likely reason for this is that phonebook
				// was restarted.  In any case, we need to add this db.Cookie to the
				// in memory sessions...
				//----------------------------------------------------------------
				s = db.NewSessionFromCookie(&c) // we don't need the return value, we just need the session to be put into memory
				lib.Console("Session was not found in memory.  Adding it to memory.  db.Cookie = %s\n", c.Cookie)

				// err = fmt.Errorf("*** UNEXPECTED STATE: session with db.Cookie %s was not found in db.Sessions", c.Cookie)
				// lib.Console("%s\n", err.Error())
				// lib.Ulog("%s\n", err.Error())
				// SvcErrorReturn(w, err, funcname)
			}
			s.Expire = s.Expire.Add(db.SessionManager.SessionTimeout * time.Minute)
			if err = db.UpdateSessionCookieDB(s); err != nil {
				lib.Console("Error updating session db.Cookie = %s\n", err.Error())
				lib.Ulog("%s: could not update session db.Cookie: %s\n", funcname, err.Error())
			}
			lib.Console("UPDATED SESSION db.Cookie TIMEOUT TIME\n")
		}
	}

	//------------------------------------------------------------------
	// if the request was to ONLY verify the existence of the db.Cookie...
	//------------------------------------------------------------------
	lib.Console("D\n")
	if foo.FLAGS&1 > 0 {
		lib.Console("D1 - verify existence of db.Cookie only\n")
		g = AuthSuccessResponse{Status: resp}
		SvcWriteResponse(&g, w)
		goto exit1
	}

	lib.Console("E - full reply: Status = %s, UID = %d, username = %s\n", resp, c.UID, c.UserName)

	//------------------------------------------------------------------
	// add the known information to the response
	//------------------------------------------------------------------
	imageProfilePath = db.GetImageLocation(c.UID)
	g = AuthSuccessResponse{
		Status:   resp,
		UID:      c.UID,
		Name:     c.UserName,
		ImageURL: imageProfilePath,
		Token:    c.Cookie,
		Expire:   c.Expire.In(db.SessionManager.ZoneUTC).Format(JSONDATETIME),
	}
	SvcWriteResponse(&g, w)

exit1:
	db.DumpSessions()
	db.DumpSessionCookies()
}

// SvcLogoff removes a session from the
func SvcLogoff(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	var funcname = "SvcLogoff"
	var err error
	var foo ValidateCookie

	lib.Console("Entered %s\n", funcname)
	// lib.Console("svcLogoff: A\n")
	data := []byte(d.data)
	if err = json.Unmarshal(data, &foo); err != nil {
		// lib.Console("svcLogoff: B\n")
		e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
		SvcErrorReturn(w, e, funcname)
		return
	}
	// lib.Console("svcLogoff: C\n")
	lib.Console("unmarshaled request.  db.Cookie value = %s\n", foo.CookieVal)

	ssn, ok := db.SessionGet(foo.CookieVal)
	if ok {
		// lib.Console("svcLogoff: D\n")
		lib.Console("found session with that db.Cookie. Deleting.\n")
		db.SessionDelete(ssn)
	} else {
		// lib.Console("svcLogoff: E\n")
		lib.Console("No session with that db.Cookie found in memory.\n")
	}
	// lib.Console("svcLogoff: F\n")
	if err := db.DeleteSessionCookie(foo.CookieVal); err != nil {
		// lib.Console("svcLogoff: G\n")
		lib.Ulog("Error deleteing session db.Cookie: %s\n", err.Error())
	}
	// lib.Console("svcLogoff: H\n")
	SvcWriteSuccessResponse(w)

	db.DumpSessions()
	db.DumpSessionCookies()

}
