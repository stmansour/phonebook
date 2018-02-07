package ws

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"phonebook/db"
	"phonebook/lib"
	"phonebook/sess"
	"strings"
)

// AuthenticateData is the struct with the username and password
// used for authentication
type AuthenticateData struct {
	User       string `json:"user"`
	Pass       string `json:"pass"`
	FLAGS      uint64 `json:"flags"`
	UserAgent  string `json:"useragent"`
	RemoteAddr string `json:"remoteaddr"`
}

// AuthSuccessResponse will be the response structure used when
// authentication is successful
type AuthSuccessResponse struct {
	Status   string `json:"status"`
	UID      int64  `json:"uid"`
	Name     string `json:"Name"`
	ImageURL string `json:"ImageURL"`
	Token    string `json:"Token"`
	Expire   string `json:"Expire"` // DATETIMEFMT in this format "2006-01-02T15:04 "
}

// ValidateCookie describes the data sent by an AIR app to check
// whether or not a cookie value is valid.
type ValidateCookie struct {
	CookieVal string `json:"cookieval"`
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

	lib.Console("Entered %s\n", funcname)

	data := []byte(d.data)
	if err = json.Unmarshal(data, &foo); err != nil {
		e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
		SvcErrorReturn(w, e, funcname)
		return
	}

	lib.Console("User = %s\n", foo.User)
	lib.Console("Pass = %s\n", foo.Pass)

	UID, Name, err := DoAuthentication(foo.User, foo.Pass)
	if err != nil {
		SvcErrorReturn(w, err, funcname)
		return
	}
	if UID > 0 {
		fname, err := getImageURL(UID)
		if err != nil {
			err := fmt.Errorf("login failed")
			SvcErrorReturn(w, err, funcname)
		}
		c := sess.GenerateSessionCookie(UID, foo.User, foo.UserAgent, foo.RemoteAddr)

		g := AuthSuccessResponse{
			Status:   "success",
			UID:      UID,
			Name:     Name,
			ImageURL: fname,
			Token:    c.Cookie,
			Expire:   c.Expire.In(sess.SessionManager.ZoneUTC).Format(JSONDATETIME),
		}
		lib.Console("g = %#v\n", g)
		SvcWriteResponse(&g, w)
		lib.Ulog("user %s successfully logged in\n", foo.User)
		err = db.InsertSessionCookie(c.UID, c.UserName, c.Cookie, &c.Expire)
		if err == nil {
			return
		}
	} else {
		err := fmt.Errorf("login failed")
		SvcErrorReturn(w, err, funcname)
	}
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

	// lookup the user
	q := fmt.Sprintf("SELECT UID,FirstName,PreferredName,passhash FROM people WHERE UserName=%q", myusername)
	var passhash string
	var UID int64
	var first, preferred string
	err := SvcCtx.db.QueryRow(q).Scan(&UID, &first, &preferred, &passhash)
	if err != nil {
		return int64(0), first, err
	}
	if passhash != mypasshash {
		err := fmt.Errorf("login failed")
		return int64(0), first, err
	}
	if len(preferred) > 0 {
		first = preferred
	}
	return UID, first, nil // login is successful
}

// SvcValidateCookie is called by an AIR app when it finds the air cookie
// but has no associated session.  If the cookie is in our sessions table
// then we send back a success response with the same info we include in
// a successful login. Otherwise, we send an appropriate error response
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

	lib.Console("Entered %s\n", funcname)

	data := []byte(d.data)
	if err = json.Unmarshal(data, &foo); err != nil {
		e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
		SvcErrorReturn(w, e, funcname)
		return
	}

	lib.Console("request for session cookie:  %s\n", foo.CookieVal)
	c, err := sess.GetSessionCookie(foo.CookieVal)
	if err != nil {
		lib.Ulog("signinHandler: error getting session cookie: %s\n", err.Error())
	}
	fname, err := getImageURL(c.UID)
	if err != nil {
		err := fmt.Errorf("login failed")
		SvcErrorReturn(w, err, funcname)
	}
	g := AuthSuccessResponse{
		Status:   "success",
		UID:      c.UID,
		Name:     c.UserName,
		ImageURL: fname,
		Token:    c.Cookie,
		Expire:   c.Expire.In(sess.SessionManager.ZoneUTC).Format(JSONDATETIME),
	}
	SvcWriteResponse(&g, w)
}

// SvcLogoff removes a session from the
func SvcLogoff(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	var funcname = "SvcLogoff"
	var err error
	var foo ValidateCookie

	lib.Console("Entered %s\n", funcname)

	data := []byte(d.data)
	if err = json.Unmarshal(data, &foo); err != nil {
		e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
		SvcErrorReturn(w, e, funcname)
		return
	}
	lib.Console("unmarshaled request.  cookie value = %s\n", foo.CookieVal)

	ssn, ok := sess.SessionGet(foo.CookieVal)
	if ok {
		lib.Console("found session with that cookie. Deleting.\n")
		sess.SessionDelete(ssn)
	} else {
		lib.Console("No session with that cookie. Deleting the cooki.\n")
		if err := db.DeleteSessionCookie(foo.CookieVal); err != nil {
			lib.Ulog("Error deleteing session cookie: %s\n", err.Error())
		}
	}
	lib.Console("Done!\n")
	SvcWriteSuccessResponse(w)
}
