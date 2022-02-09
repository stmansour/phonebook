// add a user
//   needs firstname, lastname, username, passwork, role

package main

import (
	"crypto/sha512"
	"database/sql"
	"extres"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"phonebook/lib"
	"strings"
	"time"

	_ "mysql"
)

// Role defines a collection of FieldPerms that can be assigned to a person
type Role struct {
	RID  int64  // assigned by DB
	Name string // role name
}

// KeyVal is a struct def for generic key/value string pairs
type KeyVal struct {
	key   string
	value string
}

type testContext struct {
	d        *personDetail
	co       *company
	cl       *class
	testtype int64
}

// ProductName is the name of the product that appears in the header of all html.
var ProductName = string("AIR Directory")

//--------------------------------------------------------------------
//  FINANCE
//--------------------------------------------------------------------
type company struct {
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
}

type myComp struct {
	CompCode int64  // code for this comp type
	Name     string // name for this code
	HaveIt   int64  // 0 = does not have it, 1 = has it
}

type aDeduction struct {
	DCode  int64  // code for this deduction
	Name   string // name for this deduction
	HaveIt int64  // 0 = does not have it, 1 = has it
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
	Elem  int64  // Element: Person, Company, or Class
	Field string // field within the Element
	Perm  int64  // 'logical or' of all permissions on this field
	Descr string // description of the field
}

type class struct {
	ClassCode   int64
	CoCode      int64
	Name        string
	Designation string
	Description string
}

//--------------------------------------------------------------------
//  PEOPLE-RELATED STRUCTURES
//  used in Phonebook searches and as direct reports
//--------------------------------------------------------------------
type person struct {
	UID           int64
	LastName      string
	FirstName     string
	PreferredName string
	PrimaryEmail  string
	JobCode       int64
	OfficePhone   string
	CellPhone     string
	DeptCode      int64
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
	Company                 company
	CoCode                  int64
	MgrUID                  int64
	JobTitle                string
	Class                   string
	ClassCode               int64
	MgrName                 string
	Image                   string // ptr to image -- URI
	Reports                 []person
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
	MyComps                 []myComp
	MyDeductions            []aDeduction
}

// PrepSQL is the data type holding all prepared statements
// for use within phonebook
type PrepSQL struct {
	deductList         *sql.Stmt // deduction list names and id vals
	getComps           *sql.Stmt // compensations associated with a user
	myDeductions       *sql.Stmt // deductions for a specific user
	adminPersonDetails *sql.Stmt // for AdminView and AdminEdit
	classInfo          *sql.Stmt // get class attributes
	companyInfo        *sql.Stmt // company attributes
	countersUpdate     *sql.Stmt // feature usage counters update
	delClass           *sql.Stmt // deletes a class
	delCompany         *sql.Stmt // deletes a company
	delPerson          *sql.Stmt // deletes a person
	delPersonComp      *sql.Stmt // part of delperson
	delPersonDeduct    *sql.Stmt // part of delperson
	getJobTitle        *sql.Stmt // title associated with a job code
	nameFromUID        *sql.Stmt // name lookup
	deptName           *sql.Stmt // name from DeptCode
	directReports      *sql.Stmt // folks who report to an individual
	personDetail       *sql.Stmt // get a bunch of user attributes
	adminInsertPerson  *sql.Stmt // insert a new person
	adminReadBack      *sql.Stmt // read back newly inserted person
	adminUpdatePerson  *sql.Stmt // admin update person
	insertComp         *sql.Stmt // part of admin update person
	insertDeduct       *sql.Stmt // part of admin update person
	insertClass        *sql.Stmt // adding a new class
	classReadBack      *sql.Stmt // read back newly written class
	updateClass        *sql.Stmt // update a class
	insertCompany      *sql.Stmt // insert a new company
	companyReadback    *sql.Stmt // read back newly written company
	updateCompany      *sql.Stmt // update a company
	updateMyDetails    *sql.Stmt // person updating their own details
	updatePasswd       *sql.Stmt // person updating their passwd
	readFieldPerms     *sql.Stmt // read field permissions
	accessRoles        *sql.Stmt // read access roles
	getUserCoCode      *sql.Stmt // read the cocode for a person
	loginInfo          *sql.Stmt // read info for login
}

// App is the global data structure for this app
var App struct {
	Seed             int64
	DBName           string
	DBUser           string
	Host             string
	Port             int64
	FirstUserIndex   int64            // index of first user to test
	db               *sql.DB          // database connection
	prepstmt         PrepSQL          // struct of prepared sql statements
	TestIterations   int64            // number of iterations (mutually exclusive with TestDuration)
	TestUsers        int64            // number of users to test with
	TestDurationMins int64            // duration in minutes
	TestDurationHrs  int64            // duration in hours
	TestDuration     time.Duration    // duration of tests
	Debug            bool             // show debug information
	ShowTestMatching bool             // to debug when matches fail
	UpdateDBOnly     bool             // just update the db as needed, don't run the simulation
	TotalCompanies   int64            // number of companies to create
	TotalClasses     int64            // number of classes to create
	Peeps            []*personDetail  // the people to use for test users
	FirstNames       []string         // array of first names
	LastNames        []string         // array of last names
	Streets          []string         // array of street names
	Cities           []string         // array of cities
	States           []string         // array of states
	Companies        []string         // array of random company names
	RandClasses      []string         // array of random class names
	CoCodeToName     map[int64]string // map from company code to company name
	NameToCoCode     map[string]int64 // map from company name to company code
	NameToJobCode    map[string]int64 // jobtitle to jobcode
	AcceptCodeToName map[int64]string // Acceptance to jobcode
	NameToDeptCode   map[string]int64 // department name to dept code
	NameToClassCode  map[string]int64 // class designation to classcode
	ClassCodeToName  map[int64]string // index by classcode to get the name
	Months           []string         // a map for month number to month name
	Roles            []Role           // the roles saved in the database
	JCLo, JCHi       int64            // lo and high indeces for jobcode
	DeptLo, DeptHi   int64            // lo and high indeces for department
	CompanyList      []company        // array of company structs for all defined companies
	LogFile          *os.File         // where to log messages
}

func fillUserFields(v *personDetail) {
	Nlast := len(App.LastNames)
	Nfirst := len(App.FirstNames)
	v.FirstName = strings.ToLower(App.FirstNames[rand.Intn(Nfirst)])
	v.LastName = strings.ToLower(App.LastNames[rand.Intn(Nlast)])
	v.UserName = getUsername(v.FirstName, v.LastName)
	v.PreferredName = strings.ToLower(App.FirstNames[rand.Intn(Nfirst)])
	v.Status = 1
	v.OfficePhone = randomPhoneNumber()
	v.CellPhone = randomPhoneNumber()
	v.OfficeFax = randomPhoneNumber()
	v.HomeStreetAddress = randomAddress()
	v.HomeCity = App.Cities[rand.Intn(len(App.Cities))]
	v.HomeState = App.States[rand.Intn(len(App.States))]
	v.HomePostalCode = fmt.Sprintf("%05d", rand.Intn(99999))
	v.HomeCountry = "USA"

	var w = int(App.DeptHi - App.DeptLo - 1)
	var x = int(1 + rand.Intn(w))
	v.DeptCode = int64(1 + rand.Intn(x))

	v.JobCode = 1 + rand.Int63n(1+rand.Int63n(App.JCHi-App.JCLo-1))
	v.EmergencyContactName = strings.ToLower(App.FirstNames[rand.Intn(Nfirst)]) + strings.ToLower(App.LastNames[rand.Intn(Nlast)])
	v.EmergencyContactPhone = randomPhoneNumber()
	v.PrimaryEmail = randomEmail(v.LastName, v.FirstName)
	v.SecondaryEmail = randomEmail(v.LastName, v.FirstName)

	v.ClassCode = int64(1 + rand.Intn(len(App.NameToClassCode)))

	// select a random company, but be sure it employs people.
	for {
		i := rand.Intn(len(App.CompanyList))
		if App.CompanyList[i].EmploysPersonnel > 0 {
			v.CoCode = App.CompanyList[i].CoCode
			break
		}
	}
}

func createUser(v *personDetail) {
	v.RID = 1
	fillUserFields(v)

	sha := sha512.Sum512([]byte("accord"))
	passhash := fmt.Sprintf("%x", sha)

	stmt, err := App.db.Prepare("INSERT INTO people (UserName,passhash,FirstName,LastName,RID,Status," + //6
		"OfficePhone,CellPhone,OfficeFax," + //9
		"HomeStreetAddress,HomeCity,HomeState,HomePostalCode,HomeCountry," + // 14
		"DeptCode,JobCode,PreferredName,EmergencyContactName,EmergencyContactPhone," + // 19
		"PrimaryEmail,SecondaryEmail,ClassCode,CoCode) " + // 23
		//       1                 10                  20
		" VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if nil != err {
		fmt.Printf("error = %v\n", err)
		os.Exit(1)
	}
	_, err = stmt.Exec(v.UserName, passhash, v.FirstName, v.LastName, v.RID, v.Status,
		v.OfficePhone, v.CellPhone, v.OfficeFax,
		v.HomeStreetAddress, v.HomeCity, v.HomeState, v.HomePostalCode, v.HomeCountry,
		v.DeptCode, v.JobCode, v.PreferredName, v.EmergencyContactName, v.EmergencyContactPhone,
		v.PrimaryEmail, v.SecondaryEmail, v.ClassCode, v.CoCode)
	if nil != err {
		fmt.Printf("error = %v\n", err)
	}
	v.Pro = &Tester
	//fmt.Printf("Added user to database %s:  username: %s, access role: %d\n", App.DBName, v.UserName, v.RID)

}

func loadUsers() {
	rows, err := App.db.Query("select uid from people")
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		var m personDetail
		var uid int64
		errcheck(rows.Scan(&uid))
		m.UID = uid
		adminReadDetails(&m)
		App.Peeps = append(App.Peeps, &m)
	}
	errcheck(rows.Err())
}

func loadCompanyList() {
	rows, err := App.db.Query("SELECT CoCode,LegalName,CommonName,Address,Address2,City,State,PostalCode,Country,Phone,Fax,Email,Designation,Active,EmploysPersonnel FROM companies")
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		var m company
		errcheck(rows.Scan(&m.CoCode, &m.LegalName, &m.CommonName, &m.Address, &m.Address2, &m.City, &m.State, &m.PostalCode, &m.Country, &m.Phone, &m.Fax, &m.Email, &m.Designation, &m.Active, &m.EmploysPersonnel))
		App.CompanyList = append(App.CompanyList, m)
	}
	errcheck(rows.Err())
}

func main() {
	readCommandLineArgs()

	//==============================================
	// Now open the logfile and the database...
	//==============================================
	var err error
	App.LogFile, err = os.OpenFile("usersim.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	lib.Errcheck(err)
	defer App.LogFile.Close()
	log.SetOutput(App.LogFile)
	lib.Ulog("*** usersim ***\n")

	if App.TestUsers > 100 && !App.UpdateDBOnly {
		fmt.Printf("Maximum users per simulation is 100.  You specified %d. Please reduce user count.\n", App.TestUsers)
		os.Exit(1)
	}

	// s := fmt.Sprintf("%s:@/%s?charset=utf8&parseTime=True", App.DBUser, App.DBName)
	lib.ReadConfig()
	s := extres.GetSQLOpenString(App.DBName, &lib.AppConfig)
	App.db, err = sql.Open("mysql", s)
	if nil != err {
		fmt.Printf("sql.Open: Error = %v\n", err)
	}
	defer App.db.Close()
	err = App.db.Ping()
	if nil != err {
		fmt.Printf("App.db.Ping: Error = %v\n", err)
	}

	buildPreparedStatements()
	readAccessRoles()
	loadNames()
	loadMaps()

	App.Peeps = make([]*personDetail, 0)
	initProfiles()

	loadUsers()
	loadCompanyList()
	if App.TestUsers > int64(len(App.Peeps)) {
		for i := int64(0); i < App.TestUsers-int64(len(App.Peeps)); i++ {
			var v personDetail
			createUser(&v)
		}
		App.Peeps = make([]*personDetail, 0)
		loadUsers()
	}

	if !App.UpdateDBOnly {
		lib.Ulog("Database Name: %s\n", App.DBName)
		lib.Ulog("Database User: %s\n", App.DBUser)
		lib.Ulog("User simulation count: %d\n", App.TestUsers)
		lib.Ulog("Initiate simulation\n")
		executeSimulation()
		lib.Ulog("Simulation completed\n")
	}
}
