// add a user
//   needs firstname, lastname, username, passwork, role

package main

import (
	"crypto/sha512"
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

import _ "github.com/go-sql-driver/mysql"

// Role defines a collection of FieldPerms that can be assigned to a person
type Role struct {
	RID  int    // assigned by DB
	Name string // role name
}

// KeyVal is a struct def for generic key/value string pairs
type KeyVal struct {
	key   string
	value string
}

//--------------------------------------------------------------------
//  FINANCE
//--------------------------------------------------------------------
type company struct {
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
}

type myComp struct {
	CompCode int    // code for this comp type
	Name     string // name for this code
	HaveIt   int    // 0 = does not have it, 1 = has it
}

type aDeduction struct {
	DCode  int    // code for this deduction
	Name   string // name for this deduction
	HaveIt int    // 0 = does not have it, 1 = has it
}

//--------------------------------------------------------------------
//  ROLE SECURITY
//--------------------------------------------------------------------
const (
	PERMNONE       = 0      // no permissions to see, view, modify, delete, print, or anything to this field
	PERMVIEW       = 1 << 0 // OK to view   this field for any element (Person, Company, Class)
	PERMCREATE     = 1 << 1 // OK to create   "
	PERMMOD        = 1 << 2 // OK to modify   "
	PERMDEL        = 1 << 3 // OK to delete   "
	PERMPRINT      = 1 << 4 // OK to print    "
	PERMOWNERVIEW  = 1 << 5 // OK for the owner to view this field  (applies to Person elements)
	PERMOWNERMOD   = 1 << 6 // OK for the owner to modify this field
	PERMOWNERPRINT = 1 << 7 // OK for the owner to modify this field
	PERMEXEC       = 1 << 8 // OK to execute

	ELEMPERSON  = 1 // people
	ELEMCOMPANY = 2 // companies
	ELEMCLASS   = 3 // classes
	ELEMPBSVC   = 4 // the executable service
)

// FieldPerm defines how a specific element field can be accessed
type FieldPerm struct {
	Elem  int    // Element: Person, Company, or Class
	Field string // field within the Element
	Perm  int    // 'logical or' of all permissions on this field
	Descr string // description of the field
}

type class struct {
	ClassCode   int
	Name        string
	Designation string
	Description string
}

//--------------------------------------------------------------------
//  PEOPLE-RELATED STRUCTURES
//  used in Phonebook searches and as direct reports
//--------------------------------------------------------------------
type person struct {
	UID           int
	LastName      string
	FirstName     string
	PreferredName string
	PrimaryEmail  string
	JobCode       int
	OfficePhone   string
	CellPhone     string
	DeptCode      int
	DeptName      string
	Employer      string
}

// personDetail is a structure with the basic information
// describing a virtual user. It is the same personDetail
// structure as in the Phonebook program, but updated with
// testing specific info... like Profile
type personDetail struct {
	Pro           *Profile
	SessionCookie *http.Cookie
	// below this point, it is exactly the struct defined in Phonebook
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
	Company                 company
	CoCode                  int
	MgrUID                  int
	JobTitle                string
	Class                   string
	ClassCode               int
	MgrName                 string
	Image                   string // ptr to image -- URI
	Reports                 []person
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
	MyComps                 []myComp
	MyDeductions            []aDeduction
}

// App is the global data structure for this app
var App struct {
	Seed             int64
	DBName           string
	DBUser           string
	Host             string
	Port             int
	TestIterations   int             // number of iterations (mutually exclusive with TestDuration)
	TestUsers        int             // number of users to test with
	TestDuration     int             // time in minutes
	Debug            bool            // show debug information
	Peeps            []*personDetail // the people to use for test users
	FirstNames       []string        // array of first names
	LastNames        []string        // array of last names
	Streets          []string        // array of street names
	Cities           []string        // array of cities
	States           []string        // array of states
	CoCodeToName     map[int]string  // map from company code to company name
	NameToCoCode     map[string]int  // map from company name to company code
	NameToJobCode    map[string]int  // jobtitle to jobcode
	AcceptCodeToName map[int]string  // Acceptance to jobcode
	NameToDeptCode   map[string]int  // department name to dept code
	NameToClassCode  map[string]int  // class designation to classcode
	ClassCodeToName  map[int]string  // index by classcode to get the name
	Months           []string        // a map for month number to month name
	Roles            []Role          // the roles saved in the database
	JCLo, JCHi       int             // lo and high indeces for jobcode
	DeptLo, DeptHi   int             // lo and high indeces for department
	db               *sql.DB
}

func createUser(v *personDetail) {
	v.RID = 1
	Nlast := len(App.LastNames)
	Nfirst := len(App.FirstNames)
	v.FirstName = strings.ToLower(App.FirstNames[rand.Intn(Nfirst)])
	v.LastName = strings.ToLower(App.LastNames[rand.Intn(Nlast)])
	v.UserName = getUsername(v.FirstName, v.LastName)
	v.Status = 1
	v.OfficePhone = randomPhoneNumber()
	v.CellPhone = randomPhoneNumber()
	v.OfficeFax = randomPhoneNumber()
	v.HomeStreetAddress = randomAddress()
	v.HomeCity = App.Cities[rand.Intn(len(App.Cities))]
	v.HomeState = App.States[rand.Intn(len(App.States))]
	v.HomePostalCode = fmt.Sprintf("%05d", rand.Intn(99999))
	v.HomeCountry = "USA"
	v.DeptCode = rand.Intn(1 + rand.Intn(App.DeptHi-App.DeptLo))
	v.JobCode = rand.Intn(App.JCLo + rand.Intn(App.JCHi-App.JCLo))

	sha := sha512.Sum512([]byte("accord"))
	passhash := fmt.Sprintf("%x", sha)

	stmt, err := App.db.Prepare("INSERT INTO people (UserName,passhash,FirstName,LastName,RID,Status," + //6
		"OfficePhone,CellPhone,OfficeFax," + //9
		"HomeStreetAddress,HomeCity,HomeState,HomePostalCode,HomeCountry," +
		"DeptCode,JobCode) " + //14
		//           1                 10
		" VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if nil != err {
		fmt.Printf("error = %v\n", err)
		os.Exit(1)
	}
	_, err = stmt.Exec(v.UserName, passhash, v.FirstName, v.LastName, v.RID, v.Status,
		v.OfficePhone, v.CellPhone, v.OfficeFax,
		v.HomeStreetAddress, v.HomeCity, v.HomeState, v.HomePostalCode, v.HomeCountry,
		v.DeptCode, v.JobCode)
	if nil != err {
		fmt.Printf("error = %v\n", err)
	}
	v.Pro = &Tester
	fmt.Printf("Added user to database %s:  username: %s, access role: %d\n", App.DBName, v.UserName, v.RID)
}

func loadUsers() {
	rows, err := App.db.Query("select uid from people")
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		var m personDetail
		var uid int
		errcheck(rows.Scan(&uid))
		m.UID = uid
		adminReadDetails(&m)
		App.Peeps = append(App.Peeps, &m)
	}
	errcheck(rows.Err())
}

func main() {
	readCommandLineArgs()

	if App.TestUsers > 100 {
		fmt.Printf("Maximum users per simulation is 100.  You specified %d. Please reduce user count.\n", App.TestUsers)
		os.Exit(1)
	}

	var err error
	s := fmt.Sprintf("%s:@/%s?charset=utf8&parseTime=True", App.DBUser, App.DBName)
	App.db, err = sql.Open("mysql", s)
	if nil != err {
		fmt.Printf("sql.Open: Error = %v\n", err)
	}
	defer App.db.Close()
	err = App.db.Ping()
	if nil != err {
		fmt.Printf("App.db.Ping: Error = %v\n", err)
	}
	readAccessRoles()
	loadNames()
	loadMaps()
	App.Peeps = make([]*personDetail, 0)
	initProfiles()

	loadUsers()
	if App.TestUsers > len(App.Peeps) {
		for i := 0; i < App.TestUsers-len(App.Peeps); i++ {
			var v personDetail
			createUser(&v)
		}
		App.Peeps = make([]*personDetail, 0)
		loadUsers()
	}

	executeSimulation()
}
