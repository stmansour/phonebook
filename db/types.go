package db

import (
	"database/sql"
	"phonebook/lib"
	"time"
)

// DB is the context structure of data for the DB infrastructure
var DB struct {
	DirDB *sql.DB
}

// MyComp describes the MyComp struct
type MyComp struct {
	CompCode int    // code for this comp type
	Name     string // name for this code
	HaveIt   int    // 0 = does not have it, 1 = has it
}

// ADeduction describes the ADeduction struct
type ADeduction struct {
	DCode  int    // code for this deduction
	Name   string // name for this deduction
	HaveIt int    // 0 = does not have it, 1 = has it
}

// Class defines a business unit within a company
//--------------------------------------------------------------------
type Class struct {
	ClassCode   int
	CoCode      int
	Name        string
	Designation string
	Description string
	LastModTime time.Time
	LastModBy   int
	C           Company // parent company
}

// Company defines the structure of data for a company
//--------------------------------------------------------------------
type Company struct {
	CoCode           int
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
	Active           int
	EmploysPersonnel int
	C                []Class // an array of classes for the business units of this
}

// Person defines a low-details version of the Person table
//--------------------------------------------------------------------
type Person struct {
	UID              int
	LastName         string
	FirstName        string
	PreferredName    string
	PrimaryEmail     string
	JobCode          int
	OfficePhone      string
	CellPhone        string
	OfficeFax        string
	DeptCode         int
	DeptName         string
	Employer         string
	ProfileImageURL  string
	ProfileImagePath string
}

// PersonDetail defines all details version of the Person table
//--------------------------------------------------------------------
type PersonDetail struct {
	UID                     int
	UserName                string
	LastName                string
	FirstName               string
	PrimaryEmail            string
	JobCode                 int
	OfficePhone             string
	CellPhone               string
	DeptName                string
	MiddleName              string
	Salutation              string
	Status                  int
	PositionControlNumber   string
	OfficeFax               string
	SecondaryEmail          string
	EligibleForRehire       int
	LastReview              time.Time
	NextReview              time.Time
	Birthdate               string
	BirthMonth              int
	BirthDOM                int
	HomeStreetAddress       string
	HomeStreetAddress2      string
	HomeCity                string
	HomeState               string
	HomePostalCode          string
	HomeCountry             string
	StateOfEmployment       string
	CountryOfEmployment     string
	PreferredName           string
	Comps                   []int  // an array of CompensationType values (ints)
	RID                     int    // security role assigned to this person
	CompensationStr         string //used in the admin edit interface
	DeptCode                int
	Company                 Company
	CoCode                  int
	MgrUID                  int
	JobTitle                string
	Class                   string
	ClassCode               int
	MgrName                 string
	Image                   string // ptr to image -- URI
	Reports                 []Person
	Deductions              []int
	DeductionsStr           string
	EmergencyContactName    string
	EmergencyContactPhone   string
	AcceptedHealthInsurance int
	AcceptedDentalInsurance int
	Accepted401K            int
	Hire                    time.Time
	Termination             time.Time
	NameToCoCode            map[string]int
	NameToJobCode           map[string]int
	AcceptCodeToName        map[int]string
	NameToDeptCode          map[string]int // department name to dept code
	MyComps                 []MyComp
	MyDeductions            []ADeduction
	ProfileImageURL         string
	ProfileImagePath        string
}

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

// PrepStmts are the sql prepared statements
var PrepStmts struct {
	DeleteSessionCookie  *sql.Stmt
	DeleteExpiredCookies *sql.Stmt
	GetSessionCookie     *sql.Stmt
	InsertSessionCookie  *sql.Stmt
	UpdateSessionCookie  *sql.Stmt
	LoginInfo            *sql.Stmt
	GetImagePath         *sql.Stmt
}

// CreatePreparedStmts creates prepared sql statements
func CreatePreparedStmts() {
	var err error
	var flds string
	flds = "UID,UserName,Cookie,DtExpire"
	PrepStmts.InsertSessionCookie, err = DB.DirDB.Prepare("INSERT INTO sessions (" + flds + ") VALUES(?,?,?,?)")
	lib.Errcheck(err)
	PrepStmts.GetSessionCookie, err = DB.DirDB.Prepare("SELECT " + flds + " FROM sessions WHERE Cookie=?")
	lib.Errcheck(err)
	PrepStmts.UpdateSessionCookie, err = DB.DirDB.Prepare("UPDATE sessions SET DtExpire=? WHERE Cookie=?")
	lib.Errcheck(err)
	PrepStmts.DeleteSessionCookie, err = DB.DirDB.Prepare("DELETE FROM sessions WHERE Cookie=?")
	lib.Errcheck(err)
	PrepStmts.DeleteExpiredCookies, err = DB.DirDB.Prepare("DELETE FROM sessions WHERE DtExpire <= ?")
	lib.Errcheck(err)

	PrepStmts.LoginInfo, err = DB.DirDB.Prepare("SELECT uid,firstname,preferredname,PrimaryEmail,passhash,rid FROM people WHERE UserName=?")
	lib.Errcheck(err)

	// get image path from the people table
	PrepStmts.GetImagePath, err = DB.DirDB.Prepare("SELECT ImagePath from people WHERE UID=?")
	lib.Errcheck(err)
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

// GetSessionCookie searches the session table for the speified cookie.
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
func GetSessionCookie(cookie string) (SessionCookie, error) {
	var c SessionCookie
	err := PrepStmts.GetSessionCookie.QueryRow(cookie).Scan(&c.UID, &c.UserName, &c.Cookie, &c.Expire, &c.UserAgent, &c.IP)
	if nil != err {
		if !lib.IsSQLNoResultsError(err) {
			lib.Ulog("UpdateSessionCookie: error updating expire time:  %v\n", err)
			lib.Ulog("cookie = %s\n", cookie)
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
		lib.Ulog("UpdateSessionCookie: error updating expire time:  %v\n", err)
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
