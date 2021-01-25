package db

import (
	"database/sql"
	"extres"
	"math/rand"
	"phonebook/lib"
	"time"
)

// DB is the context structure of data for the DB infrastructure
var DB struct {
	DirDB    *sql.DB
	noAuth   bool // is authrization needed to access the db?
	Config   extres.ExternalResources
	DBFields map[string]string // map of db table fields DBFields[tablename] = field list
	Zone     *time.Location    // what timezone should the server use?
	Key      []byte            // crypto key
	Rand     *rand.Rand        // for generating Reference Numbers or other UniqueIDs
}

// MyComp describes the MyComp struct
type MyComp struct {
	CompCode int64  // code for this comp type
	Name     string // name for this code
	HaveIt   int64  // 0 = does not have it, 1 = has it
}

// ADeduction describes the ADeduction struct
type ADeduction struct {
	DCode  int64  // code for this deduction
	Name   string // name for this deduction
	HaveIt int64  // 0 = does not have it, 1 = has it
}

// Class defines a business unit within a company
//--------------------------------------------------------------------
type Class struct {
	ClassCode   int64  // uid of Business Unit
	CoCode      int64  // uid of parent company
	Name        string // name of Business Unit
	Designation string // business unit designation
	Description string
	LastModTime time.Time
	LastModBy   int64
	C           Company // parent company (just a holder for convenience)
}

// Company defines the structure of data for a company
//--------------------------------------------------------------------
type Company struct {
	CoCode           int64
	LegalName        string
	CommonName       string
	Address          string
	Address2         string
	City             string
	State            string
	PostalCode       string
	Country          string
	Phone            string
	Fax              string
	Email            string
	Designation      string
	Active           int64
	EmploysPersonnel int64
	C                []Class // an array of classes for the business units of this
}

// Person defines a low-details version of the Person table
//--------------------------------------------------------------------
type Person struct {
	UID              int64
	LastName         string
	FirstName        string
	PreferredName    string
	PrimaryEmail     string
	JobCode          int64
	OfficePhone      string
	CellPhone        string
	OfficeFax        string
	DeptCode         int64
	DeptName         string
	Employer         string
	ProfileImageURL  string
	ProfileImagePath string
}

// PeopleTypeDown is the struct needed to match names in typedown controls
//--------------------------------------------------------------------
type PeopleTypeDown struct {
	Recid int64 `json:"recid"` // this will hold the UID
	UID   int64
	Name  string
}

// PersonDetail defines all details version of the Person table
//--------------------------------------------------------------------
type PersonDetail struct {
	UID                     int64
	UserName                string
	LastName                string
	FirstName               string
	PrimaryEmail            string
	JobCode                 int64
	OfficePhone             string
	CellPhone               string
	DeptName                string
	MiddleName              string
	Salutation              string
	Status                  int64
	PositionControlNumber   string
	OfficeFax               string
	SecondaryEmail          string
	EligibleForRehire       int64
	LastReview              time.Time
	NextReview              time.Time
	Birthdate               string
	BirthMonth              int64
	BirthDOM                int64
	HomeStreetAddress       string
	HomeStreetAddress2      string
	HomeCity                string
	HomeState               string
	HomePostalCode          string
	HomeCountry             string
	StateOfEmployment       string
	CountryOfEmployment     string
	PreferredName           string
	Comps                   []int64 // an array of CompensationType values (ints)
	RID                     int64   // security role assigned to this person
	CompensationStr         string  //used in the admin edit interface
	DeptCode                int64
	Company                 Company
	CoCode                  int64
	MgrUID                  int64
	JobTitle                string
	Class                   string
	ClassCode               int64
	MgrName                 string
	Image                   string // ptr to image -- URI
	Reports                 []Person
	Deductions              []int64
	DeductionsStr           string
	EmergencyContactName    string
	EmergencyContactPhone   string
	AcceptedHealthInsurance int64
	AcceptedDentalInsurance int64
	Accepted401K            int64
	Hire                    time.Time
	Termination             time.Time
	NameToCoCode            map[string]int64
	NameToJobCode           map[string]int64
	AcceptCodeToName        map[int64]string
	NameToDeptCode          map[string]int64 // department name to dept code
	MyComps                 []MyComp
	MyDeductions            []ADeduction
	ProfileImageURL         string
	ProfileImagePath        string
}

// Init initializes the database infrastructure
//
// INPUTS
//  name  - name of the db to load.  This name overrides the one in the config file
//          if its length is > 0
// RETURNS
//  error - any error encountered
//-----------------------------------------------------------------------------
func Init() error {
	CreatePreparedStmts()
	return nil
}

// BUInfo contains the information about a business unit needed for typedown
//--------------------------------------------------------------------
type BUInfo struct {
	Recid       int64  `json:"recid"`
	ClassCode   int64  // business unit uid
	CoCode      int64  // parent company uid
	Name        string // BU name
	Designation string // bu designation
}

// GetBUByBUD searches for the Business Unit by the supplied BUD
//
// INPUTS
//   s = the BUD to search for
//
// RETURNS
//   the Business Unit record
//   err Any errors encountered
//-----------------------------------------------------------------------------
func GetBUByBUD(s string) (Class, error) {
	//funcname := "GetBUByBUD"
	var p Class
	row := PrepStmts.GetBUByBUD.QueryRow(s)
	err := row.Scan(&p.ClassCode, &p.CoCode, &p.Name, &p.Designation, &p.Description)
	SkipSQLNoRowsError(&err)
	return p, err
}

// GetBUTypeDown returns a slice of session cookies
//
// INPUTS
//
// RETURNS
//  array of BUInfo structs matching the string s1
//
//  err      Any errors encountered
//-----------------------------------------------------------------------------
func GetBUTypeDown(s1 string, limit int64) ([]BUInfo, error) {
	funcname := "GetBUTypeDown"
	var m []BUInfo
	s := "%" + s1 + "%"
	// lib.Console("s = %q\n", s)
	rows, err := PrepStmts.GetBUTypeDown.Query(s, limit)
	if err != nil {
		lib.Ulog("%s: error getting rows: %s\n", funcname, err.Error())
		return m, err
	}
	defer rows.Close()

	for rows.Next() {
		var p BUInfo
		if err := rows.Scan(&p.ClassCode, &p.CoCode, &p.Name, &p.Designation); err != nil {
			lib.Ulog("%s: error getting row:  %v\n", funcname, err)
			return m, err
		}
		p.Recid = p.ClassCode // a unique identifier for this class
		m = append(m, p)
	}

	return m, nil
}

// GetSessionCookieDB searches the session table for the speified cookie.
//
// INPUTS
//  cookie - the Web Cookie value string. If err == nil then it is filled
//           with all the info associated with the session table record.
//           If it is not found, then len(c.Cookie) == 0
//
// RETURNS
//  SessionCookie - if err == nil then a SessionCookie filled out with the
//           information in the session table record. If err != nil, then
//           the SessionCookie value will have len() == 0
//
//  err      Any errors encountered
//-----------------------------------------------------------------------------
func GetSessionCookieDB(cookie string) (SessionCookie, error) {
	var c SessionCookie
	err := PrepStmts.GetSessionCookie.QueryRow(cookie).Scan(&c.UID, &c.UserName, &c.Cookie, &c.Expire, &c.UserAgent, &c.IP)
	if nil != err {
		if !lib.IsSQLNoResultsError(err) {
			lib.Ulog("GetSessionCookie: error getting cookie:  %v\n", err)
			lib.Ulog("cookie = %s\n", cookie)
			return c, err
		}
	}
	return c, nil
}

// GetAllSessionCookies returns a slice of session cookies
//
// INPUTS
//
// RETURNS
//  []SessionCookie - a slice with all the rows in the sessions table.
//
//  err      Any errors encountered
//-----------------------------------------------------------------------------
func GetAllSessionCookies() ([]SessionCookie, error) {
	funcname := "GetAllSessionCookies"
	var m []SessionCookie
	rows, err := PrepStmts.GetAllSessionCookies.Query()
	if err != nil {
		lib.Ulog("%s: error getting rows: %s\n", funcname, err.Error())
		return m, err
	}
	defer rows.Close()

	for rows.Next() {
		var c SessionCookie
		err := rows.Scan(&c.UID, &c.UserName, &c.Cookie, &c.Expire, &c.UserAgent, &c.IP)
		if err != nil {
			lib.Ulog("%s: error getting row:  %v\n", funcname, err)
			return m, err
		}
		m = append(m, c)
	}

	return m, nil
}

// FindMatchingSessionCookie searches the session table for the speified cookie.
//
// INPUTS
//  user   - this is the user making the request. It is assumed this user
//           has already been authenticated.
//  ip     - the user's ip address
//  ua     - the user agent making the request
//
// RETURNS
//  SessionCookie - the session cookie found matching the params supplied.
//           If no match was made the session cookie will have a zero length
//			 value for the .Cookie field
//  err      Any errors encountered
//-----------------------------------------------------------------------------
func FindMatchingSessionCookie(user, ip, ua string) (SessionCookie, error) {
	var c SessionCookie
	err := PrepStmts.FindMatchingSessionCookie.QueryRow(user, ip, ua).Scan(&c.UID, &c.UserName, &c.Cookie, &c.Expire, &c.UserAgent, &c.IP)
	if nil != err {
		if !lib.IsSQLNoResultsError(err) {
			lib.Ulog("FindMatchingSessionCookie: error finding match:  %v\n", err)
			lib.Ulog("user, ip, ua = %s, %s, %s\n", user, ip, ua)
			return c, err
		}
	}
	return c, nil
}

// DeleteSessionCookie updates the specified cookie with the new expire time
//-----------------------------------------------------------------------------
func DeleteSessionCookie(cookie string) error {
	_, err := PrepStmts.DeleteSessionCookie.Exec(cookie)
	if nil != err {
		lib.Ulog("DeleteSessionCookie: error deleting cookie:  %v\n", err)
		lib.Ulog("cookie = %s\n", cookie)
	}
	return err
}

// InsertSessionCookie inserts a new session cookie into the sessions table
//-----------------------------------------------------------------------------
func InsertSessionCookie(UID int64, user string, cookie string, dt *time.Time, ua, ip string) error {
	lib.Console("InsertSessionCookie: %d, %s, ua = %s, ip = %s\n", UID, user, ua, ip)
	_, err := PrepStmts.InsertSessionCookie.Exec(UID, user, cookie, *dt, ua, ip)
	if nil != err {
		lib.Ulog("InsertSessionCookie: error inserting Cookie:  %v\n", err)
		lib.Ulog("UID = %d, user = %s, ip = %s cookie = %s, ua = %s\n", UID, user, ip, cookie, ua)
	}
	return err
}

// UpdateSessionCookie inserts a new session cookie into the sessions table
//-----------------------------------------------------------------------------
func UpdateSessionCookie(cookie string, dt *time.Time) error {
	_, err := PrepStmts.UpdateSessionCookie.Exec(*dt, cookie)
	if nil != err {
		lib.Ulog("UpdateSessionCookie: error updating Cookie:  %v\n", err)
	}
	return err
}
