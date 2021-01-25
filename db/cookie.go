package db

import (
	"crypto/md5"
	"fmt"
	"phonebook/lib"
	"time"
)

// SessionCookieName is a string holding the cookie name for browser cookies
// throughout the air suite.
//---------------------------------------------------------------------------
var SessionCookieName = string("air")

// Cookie management for web client and web services

// GenerateSessionCookie - create a new cookie
//
// INPUTS
//  username   - user's login name
//  useragent  - the user's client. If this is being called by another server
//               via a web service, the server should pass in the user's agent
//               used to make the request. This helps differentiate between the
//               same user using multiple clients (perhaps different browsers
//               at the same time).
//  remoteaddr - the user's IP address in string form
//
// RETURNS
//  string     - a unique key identifier for the user
//-----------------------------------------------------------------------------
func GenerateSessionCookie(UID int64, username, useragent, remoteaddr string) SessionCookie {
	lib.Console("Entered GenerateSessionCookie:  ua = %s, ip = %s\n", useragent, remoteaddr)
	var c SessionCookie
	key := username + useragent + remoteaddr
	c.Cookie = fmt.Sprintf("%x", md5.Sum([]byte(key)))
	c.UID = UID
	c.UserName = username
	c.Expire = time.Now().Add(SessionManager.SessionTimeout * time.Minute)
	c.UserAgent = useragent
	c.IP = remoteaddr
	lib.Console("GenerateSessionCookie    %s : %s : %s  --> %s\n", username, useragent, remoteaddr, c.Cookie)
	lib.Console("   PAddr : User Agent    %s : %s\n", c.IP, c.UserAgent)
	return c
}

// GetSessionCookie - try to find the supplied cookie
//
// INPUTS
//  s           - the session cookie
//
// RETURNS
//  error       - any errors encountered, or nil if no errors
//-----------------------------------------------------------------------------
func GetSessionCookie(s string) (SessionCookie, error) {
	return GetSessionCookieDB(s)
}

// InsertSessionCookieDB - add a new cookie to the session db table
//
// INPUTS
//  s           - the session containing the cookie
//
// RETURNS
//  error       - any errors encountered, or nil if no errors
//-----------------------------------------------------------------------------
func InsertSessionCookieDB(s *Session) error {
	return InsertSessionCookie(s.UID, s.Username, s.Token, &s.Expire, s.UserAgent, s.IP)
}

// UpdateSessionCookieDB - update the expire time of an existing cookie. It
// is assumed that the expire time in the session is correct, that is the
// value that will be written to the database.
//
// INPUTS
//  s           - the session containing the cookie
//
// RETURNS
//  error       - any errors encountered, or nil if no errors
//-----------------------------------------------------------------------------
func UpdateSessionCookieDB(s *Session) error {
	lib.Console("Entered UpdateSessionCookie: token = %s, Expire = %v\n", s.Token, s.Expire)
	err := UpdateSessionCookie(s.Token, &s.Expire)
	return err
}

// DeleteSessionCookieDB - remove the cookie from the session list. This is called
//                when the user explicitly logs out rather
//
// INPUTS
//  s           - the session containing the cookie
//
// RETURNS
//  error       - any errors encountered, or nil if no errors
//-----------------------------------------------------------------------------
func DeleteSessionCookieDB(s *Session) error {
	return nil
}

// ExpiredCookieCleaner removes sessions that have timed out
//-----------------------------------------------------------------------------
func ExpiredCookieCleaner() {
	for {
		select {
		case <-time.After(1 * time.Minute):
			now := time.Now()
			_, err := PrepStmts.DeleteExpiredCookies.Exec(now)
			if err != nil {
				lib.Ulog("Error removing expired coockies = %s\n", err.Error())
			}
		}
	}
}

// DumpSessionCookies is a convenience function to dump the current
//     session table to the console.
//
// INPUTS
//
// RETURNS
//  err      Any errors encountered
//-----------------------------------------------------------------------------
func DumpSessionCookies() error {
	fmt.Printf("DB SESSIONS COOKIE TABLE\n")
	m, err := GetAllSessionCookies()
	if err != nil {
		return err
	}
	for i := 0; i < len(m); i++ {
		fmt.Printf("%3d. %s  %s (%d) ip(%s), %s, ua = %s\n", i, m[i].Expire.Format(lib.JSONDATETIME), m[i].UserName, m[i].UID, m[i].IP, m[i].Cookie, m[i].UserAgent)
	}
	return nil
}
