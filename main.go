// Phonebook - a temporary directory interface

package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"
)

import _ "github.com/go-sql-driver/mysql"

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

// Role defines a collection of FieldPerms that can be assigned to a person
type Role struct {
	RID   int         // assigned by DB
	Name  string      // role name
	Descr string      // role description
	Perms []FieldPerm // permissions for all fields, all entities
}

//--------------------------------------------------------------------
//  PEOPLE-RELATED STRUCTURES
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

type personDetail struct {
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

// dataFields lists all the field names and other information
// about the field:
// 		- its description
//		- whether the field is only accessible via an Administration screen
type dataFields struct {
	Elem        int
	FieldName   string
	AdminScreen bool
	Description string
}

var adminScreenFields = []dataFields{
	{ELEMPERSON, "Status", false, "Indicates whether the person is an active employee."},
	{ELEMPERSON, "EligibleForRehire", false, "Indicates whether a past employee can be rehired."},
	{ELEMPERSON, "UID", false, "A unique identifier associated with the employee. Once created, it never changes."},
	{ELEMPERSON, "Salutation", false, "'Mr.', 'Mrs.', 'Ms.', etc."},
	{ELEMPERSON, "FirstName", false, "The person's common name."},
	{ELEMPERSON, "MiddleName", false, "The person's middle name."},
	{ELEMPERSON, "LastName", false, "The person's surname or last name."},
	{ELEMPERSON, "PreferredName", false, "Less formal name but more commonly used, for example 'Mike' rather than 'Michael'."},
	{ELEMPERSON, "PrimaryEmail", false, "The primary email address to use for this person."},
	{ELEMPERSON, "OfficePhone", false, "This person's office telephone number."},
	{ELEMPERSON, "CellPhone", false, "This person's cellphone number."},
	{ELEMPERSON, "EmergencyContactName", true, "Name of someone to contact in the event of an emergency."},
	{ELEMPERSON, "EmergencyContactPhone", true, "Phone number for the emergency contact."},
	{ELEMPERSON, "HomeStreetAddress", true, "def"},
	{ELEMPERSON, "HomeStreetAddress2", true, "def"},
	{ELEMPERSON, "HomeCity", true, "def"},
	{ELEMPERSON, "HomeState", true, "def"},
	{ELEMPERSON, "HomePostalCode", true, "def"},
	{ELEMPERSON, "HomeCountry", true, "def"},
	{ELEMPERSON, "PrimaryEmail", true, "def"},
	{ELEMPERSON, "SecondaryEmail", true, "def"},
	{ELEMPERSON, "OfficePhone", true, "def"},
	{ELEMPERSON, "OfficeFax", true, "def"},
	{ELEMPERSON, "CellPhone", true, "def"},
	{ELEMPERSON, "BirthDOM", true, "def"},
	{ELEMPERSON, "BirthMonth", true, "def"},
	{ELEMPERSON, "CoCode", true, "The company code associated with this user."},
	{ELEMPERSON, "JobCode", true, "def"},
	{ELEMPERSON, "ClassCode", true, "def"},
	{ELEMPERSON, "DeptCode", true, "def"},
	{ELEMPERSON, "PositionControlNumber", true, "def"},
	{ELEMPERSON, "MgrUID", true, "def"},
	{ELEMPERSON, "Accepted401K", true, "def"},
	{ELEMPERSON, "AcceptedDentalInsurance", true, "def"},
	{ELEMPERSON, "AcceptedHealthInsurance", true, "def"},
	{ELEMPERSON, "Hire", true, "def"},
	{ELEMPERSON, "Termination", true, "def"},
	{ELEMPERSON, "LastReview", true, "def"},
	{ELEMPERSON, "NextReview", true, "def"},
	{ELEMPERSON, "StateOfEmployment", false, "def"},
	{ELEMPERSON, "CountryOfEmployment", false, "def"},
	{ELEMPERSON, "Comps", true, "def"},
	{ELEMPERSON, "Deductions", true, "def"},
	{ELEMPERSON, "MyDeductions", true, "def"},
	{ELEMPERSON, "RID", true, "role identifier"},
	{ELEMPERSON, "ElemEntity", true, "The entire entity"},
	{ELEMCOMPANY, "CoCode", false, "def"},
	{ELEMCOMPANY, "LegalName", false, "def"},
	{ELEMCOMPANY, "CommonName", false, "def"},
	{ELEMCOMPANY, "Address", false, "def"},
	{ELEMCOMPANY, "Address2", false, "def"},
	{ELEMCOMPANY, "City", false, "def"},
	{ELEMCOMPANY, "State", false, "def"},
	{ELEMCOMPANY, "PostalCode", false, "def"},
	{ELEMCOMPANY, "Country", false, "def"},
	{ELEMCOMPANY, "Phone", false, "def"},
	{ELEMCOMPANY, "Fax", false, "def"},
	{ELEMCOMPANY, "Email", false, "def"},
	{ELEMCOMPANY, "Designation", false, "def"},
	{ELEMCOMPANY, "Active", false, "def"},
	{ELEMCOMPANY, "EmploysPersonnel", false, "def"},
	{ELEMCOMPANY, "ElemEntity", true, "The entire entity"},
	{ELEMCLASS, "ClassCode", false, "def"},
	{ELEMCLASS, "Name", false, "def"},
	{ELEMCLASS, "Designation", false, "def"},
	{ELEMCLASS, "Description", false, "def"},
	{ELEMCLASS, "ElemEntity", true, "The entire entity"},
	{ELEMPBSVC, "Shutdown", true, "Shut down the running Phonebook service"},
	{ELEMPBSVC, "Restart", true, "Restart the running Phonebook service"},
}

type class struct {
	ClassCode   int
	Name        string
	Designation string
	Description string
}

type searchResults struct {
	Query   string
	Matches []person
}

type searchCoResults struct {
	Query   string
	Matches []company
}

type searchClassResults struct {
	Query   string
	Matches []class
}

type signin struct {
	ErrNo  int    // 0 = no error, otherwise signin error
	ErrMsg string // err message string for user
}

//--------------------------------------------------------------------
// uiSupport is an umbrella structure in which we can pass many useful
// data objects to the UI
//--------------------------------------------------------------------
type uiSupport struct {
	CoCodeToName     map[int]string // map from company code to company name
	NameToCoCode     map[string]int // map from company name to company code
	NameToJobCode    map[string]int // jobtitle to jobcode
	AcceptCodeToName map[int]string // Acceptance to jobcode
	NameToDeptCode   map[string]int // department name to dept code
	NameToClassCode  map[string]int // class designation to classcode
	ClassCodeToName  map[int]string // index by classcode to get the name
	Months           []string       // a map for month number to month name
	Roles            []Role         // list of roles -- fields are not initialized
	C                *company
	A                *class
	D                *personDetail
	R                *searchResults
	S                *signin
	T                *searchCoResults
	L                *searchClassResults
	X                *session
	K                *UsageCounters
	N                []session
}

// PhonebookUI is the instance of uiSupport used by this app
var PhonebookUI uiSupport

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

// Phonebook is the global application structure providing
// information that any function might need.
var Phonebook struct {
	Port               int           // port on which we listen
	db                 *sql.DB       // the database connection
	prepstmt           PrepSQL       // struct of prepared sql statements
	DBName             string        // name of database to use
	DBUser             string        // user phonebook should use for accessing db
	LogFile            *os.File      // where to log messages
	Roles              []Role        // the roles saved in the database
	ReqMem             chan int      // request to access UI data memory
	ReqMemAck          chan int      // done with memory
	ReqSessionMem      chan int      // request to access Session data memory
	ReqSessionMemAck   chan int      // done with Session datamemory
	ReqCountersMem     chan int      // request to access counters
	ReqCountersMemAck  chan int      // done with counters mem
	DebugToScreen      bool          // show logged messages to screen
	Debug              bool          // push debug log messages to the logfile
	SecurityDebug      bool          // push security debug messages to the logfile
	SessionTimeout     time.Duration // timeout in minutes
	SessionCleanupTime time.Duration // time in minutes
	CountersUpdateTime int           // time in minutes
}

// UsageCounters defines the type of stats phonebook stores
type UsageCounters struct {
	SearchPeople     int64
	SearchClasses    int64
	SearchCompanies  int64
	EditPerson       int64
	ViewPerson       int64
	ViewClass        int64
	ViewCompany      int64
	AdminEditPerson  int64
	AdminEditClass   int64
	AdminEditCompany int64
	DeletePerson     int64
	DeleteClass      int64
	DeleteCompany    int64
	SignIn           int64
	Logoff           int64
}

// Counters stores the number of times each function was executed
// The numbers are cumulative over server restarts and maintained
// in the database.
var Counters UsageCounters

var funcMap map[string]interface{}

var chttp = http.NewServeMux()

func errcheck(err error) {
	if err != nil {
		fmt.Printf("error = %v\n", err)
		log.Fatal(err)
	}
}

// HomeHandler serves static http content such as the .css files
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, ".") {
		chttp.ServeHTTP(w, r)
	} else {
		http.Redirect(w, r, "/signin/", http.StatusFound)
	}
}

func initUIData(u *uiSupport) {
	u.CoCodeToName = make(map[int]string, len(PhonebookUI.CoCodeToName))
	u.NameToCoCode = make(map[string]int, len(PhonebookUI.NameToCoCode))
	for k, v := range PhonebookUI.CoCodeToName {
		u.CoCodeToName[k] = v
		u.NameToCoCode[v] = k
	}
	u.AcceptCodeToName = make(map[int]string, len(PhonebookUI.AcceptCodeToName))
	for k, v := range PhonebookUI.AcceptCodeToName {
		u.AcceptCodeToName[k] = v
	}
	u.NameToDeptCode = make(map[string]int, len(PhonebookUI.NameToDeptCode))
	for k, v := range PhonebookUI.NameToDeptCode {
		u.NameToDeptCode[k] = v
	}
	u.NameToJobCode = make(map[string]int, len(PhonebookUI.NameToJobCode))
	for k, v := range PhonebookUI.NameToJobCode {
		u.NameToJobCode[k] = v
	}
	u.NameToClassCode = make(map[string]int, len(PhonebookUI.NameToClassCode))
	u.ClassCodeToName = make(map[int]string, len(PhonebookUI.NameToClassCode))
	for k, v := range PhonebookUI.NameToClassCode {
		u.NameToClassCode[k] = v
		u.ClassCodeToName[v] = k
	}
	u.Months = make([]string, len(PhonebookUI.Months))
	for i := 0; i < len(PhonebookUI.Months); i++ {
		u.Months[i] = PhonebookUI.Months[i]
	}
	u.Roles = make([]Role, len(Phonebook.Roles))
	for i := 0; i < len(Phonebook.Roles); i++ {
		u.Roles[i] = Role{}
		u.Roles[i].Name = Phonebook.Roles[i].Name
		u.Roles[i].RID = Phonebook.Roles[i].RID
	}
}

// Dispatcher controls access to shared resources.
func Dispatcher() {
	for {
		select {
		case <-Phonebook.ReqMem:
			Phonebook.ReqMemAck <- 1 // tell caller go ahead
			<-Phonebook.ReqMemAck    // block until caller is done with mem
		}
	}
}

func loadCompanies() {
	var code int
	var name string
	PhonebookUI.CoCodeToName = make(map[int]string)
	PhonebookUI.NameToCoCode = make(map[string]int)

	rows, err := Phonebook.db.Query("select cocode,CommonName from companies")
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		errcheck(rows.Scan(&code, &name))
		PhonebookUI.CoCodeToName[code] = name
		PhonebookUI.NameToCoCode[name] = code
	}
	errcheck(rows.Err())
}

func loadClasses() {
	var code int
	var name string

	PhonebookUI.NameToClassCode = make(map[string]int)
	PhonebookUI.ClassCodeToName = make(map[int]string)
	rows, err := Phonebook.db.Query("select classcode,designation from classes")
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		errcheck(rows.Scan(&code, &name))
		PhonebookUI.NameToClassCode[name] = code
		PhonebookUI.ClassCodeToName[code] = name
	}
	// for k, v := range Phonebook.NameToClassCode {
	// 	fmt.Printf("%s %d\n", k, v)
	// }
	errcheck(rows.Err())
}

func loadMaps() {
	var code int
	var name string

	funcMap = template.FuncMap{
		"compToString":         compensationTypeToString,
		"acceptIntToString":    acceptIntToString,
		"dateToString":         dateToString,
		"dateYear":             dateYear,
		"monthStringToInt":     monthStringToInt,
		"add":                  add,
		"sub":                  sub,
		"rmd":                  rmd,
		"mul":                  mul,
		"div":                  div,
		"hasFieldAccess":       hasFieldAccess,
		"hasPERMMODaccess":     hasPERMMODaccess,
		"hasAdminScreenAccess": hasAdminScreenAccess,
		"showAdminButton":      showAdminButton,
		"getBreadcrumb":        getBreadcrumb,
		"getHTMLBreadcrumb":    getHTMLBreadcrumb,
		"datetimeToString":     datetimeToString,
	}
	loadCompanies()
	loadClasses()

	PhonebookUI.NameToJobCode = make(map[string]int)
	rows, err := Phonebook.db.Query("select jobcode,title from jobtitles")
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		errcheck(rows.Scan(&code, &name))
		PhonebookUI.NameToJobCode[name] = code
	}
	errcheck(rows.Err())

	PhonebookUI.NameToDeptCode = make(map[string]int)
	rows, err = Phonebook.db.Query("select deptcode,name from departments order by name")
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		errcheck(rows.Scan(&code, &name))
		PhonebookUI.NameToDeptCode[name] = code
	}
	errcheck(rows.Err())

	PhonebookUI.AcceptCodeToName = make(map[int]string)
	for i := ACPTUNKNOWN; i <= ACPTLAST; i++ {
		PhonebookUI.AcceptCodeToName[i] = acceptIntToString(i)
	}

	PhonebookUI.Months = make([]string, len(fmtMonths))
	for i := 0; i < len(fmtMonths); i++ {
		PhonebookUI.Months[i] = fmtMonths[i]
	}

	errcheck(Phonebook.db.QueryRow("select SearchPeople,SearchClasses,SearchCompanies,"+
		"EditPerson,ViewPerson,ViewClass,ViewCompany,"+
		"AdminEditPerson,AdminEditClass,AdminEditCompany,"+
		"DeletePerson,DeleteClass,DeleteCompany,SignIn,Logoff from counters").Scan(
		&Counters.SearchPeople, &Counters.SearchClasses, &Counters.SearchCompanies,
		&Counters.EditPerson, &Counters.ViewPerson, &Counters.ViewClass, &Counters.ViewCompany,
		&Counters.AdminEditPerson, &Counters.AdminEditClass, &Counters.AdminEditCompany,
		&Counters.DeletePerson, &Counters.DeleteClass, &Counters.DeleteCompany,
		&Counters.SignIn, &Counters.Logoff))
}

func initHTTP() {
	chttp.Handle("/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/admin/", adminHandler)
	http.HandleFunc("/adminAddClass/", adminAddClassHandler)
	http.HandleFunc("/adminAddCompany/", adminAddCompanyHandler)
	http.HandleFunc("/adminAddPerson/", adminAddPersonHandler)
	http.HandleFunc("/adminEdit/", adminEditHandler)
	http.HandleFunc("/adminEditClass/", adminEditClassHandler)
	http.HandleFunc("/adminEditCo/", adminEditCompanyHandler)
	http.HandleFunc("/adminView/", adminViewHandler)
	http.HandleFunc("/adminViewBtn/", adminViewBtnHandler)
	http.HandleFunc("/class/", classHandler)
	http.HandleFunc("/company/", companyHandler)
	http.HandleFunc("/delClass/", delClassHandler)
	http.HandleFunc("/delClassRefErr/", delClassRefErr)
	http.HandleFunc("/delCompany/", delCoHandler)
	http.HandleFunc("/delCoRefErr/", delCoRefErr)
	http.HandleFunc("/delPerson/", delPersonHandler)
	http.HandleFunc("/delPersonRefErr/", delPersonRefErrHandler)
	http.HandleFunc("/detail/", detailHandler)
	http.HandleFunc("/detailpop/", detailpopHandler)
	http.HandleFunc("/editDetail/", editDetailHandler)
	http.HandleFunc("/extAdminShutdown/", extAdminShutdown)
	http.HandleFunc("/help/", helpHandler)
	http.HandleFunc("/inactivatePerson/", inactivatePersonHandler)
	http.HandleFunc("/logoff/", logoffHandler)
	http.HandleFunc("/pop/", popHandler)
	http.HandleFunc("/restart/", restartHandler)
	http.HandleFunc("/saveAdminEdit/", saveAdminEditHandler)
	http.HandleFunc("/saveAdminEditClass/", saveAdminEditClassHandler)
	http.HandleFunc("/saveAdminEditCo/", saveAdminEditCoHandler)
	http.HandleFunc("/savePersonDetails/", savePersonDetailsHandler)
	http.HandleFunc("/search/", searchHandler)
	http.HandleFunc("/searchcl/", searchClassHandler)
	http.HandleFunc("/searchco/", searchCompaniesHandler)
	http.HandleFunc("/setup/", setupHandler)
	http.HandleFunc("/shutdown/", shutdownHandler)
	http.HandleFunc("/signin/", signinHandler)
	http.HandleFunc("/status/", statusHandler)
	http.HandleFunc("/stats/", statsHandler)
	http.HandleFunc("/weblogin/", webloginHandler)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func readCommandLineArgs() {
	dbusPtr := flag.String("B", "ec2-user", "database username")
	cntrPtr := flag.Int("c", 5, "counter update period in minutes")
	dbugPtr := flag.Bool("d", false, "debug mode - includes debug info in logfile")
	dtscPtr := flag.Bool("D", false, "LogToScreen mode - prints log messages to stdout")
	dbnmPtr := flag.String("N", "accord", "database name")
	portPtr := flag.Int("p", 8250, "port on which Phonebook listens")
	sbugPtr := flag.Bool("s", false, "security debug mode - includes security debugging info in logfile")

	flag.Parse()

	Phonebook.Port = *portPtr
	Phonebook.Debug = *dbugPtr
	Phonebook.SecurityDebug = *sbugPtr
	Phonebook.DebugToScreen = *dtscPtr
	Phonebook.DBName = *dbnmPtr
	Phonebook.DBUser = *dbusPtr
	Phonebook.CountersUpdateTime = *cntrPtr
}

func main() {
	//=============================
	//  Hardcoded defaults...
	//=============================
	Phonebook.ReqMem = make(chan int)
	Phonebook.ReqMemAck = make(chan int)
	Phonebook.ReqSessionMem = make(chan int)
	Phonebook.ReqSessionMemAck = make(chan int)
	Phonebook.ReqCountersMem = make(chan int)
	Phonebook.ReqCountersMemAck = make(chan int)
	Phonebook.Roles = make([]Role, 0)
	Phonebook.SessionTimeout = 120    // minutes
	Phonebook.SessionCleanupTime = 10 // minutes

	//==============================================
	// There may be some command line overrides...
	//==============================================
	readCommandLineArgs()

	//==============================================
	// Now open the logfile and the database...
	//==============================================
	var err error
	Phonebook.LogFile, err = os.OpenFile("Phonebook.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer Phonebook.LogFile.Close()
	log.SetOutput(Phonebook.LogFile)
	ulog("*** Accord PHONEBOOK ***\n")

	//==============================================
	// And the database...
	//==============================================
	dbopenparms := fmt.Sprintf("%s:@/%s?charset=utf8&parseTime=True", Phonebook.DBUser, Phonebook.DBName)
	db, err := sql.Open("mysql", dbopenparms)
	if nil != err {
		ulog("sql.Open: Error = %v\n", err)
	}
	defer db.Close()
	err = db.Ping()
	if nil != err {
		ulog("db.Ping: Error = %v\n", err)
		os.Exit(2)
	}
	ulog("MySQL database opened with \"%s\"\n", dbopenparms)
	Phonebook.db = db
	buildPreparedStatements()

	//==============================================
	// Load some of the database info...
	//==============================================
	loadMaps()
	readAccessRoles()
	if Phonebook.Debug {
		dumpAccessRoles()
	}

	//==============================================
	// On with the show...
	//==============================================
	go Dispatcher()
	go CounterDispatcher()
	go UpdateCounters()
	go SessionDispatcher()
	go SessionCleanup()

	initHTTP()
	sessionInit()

	ulog("Phonebook initiating HTTP service on port %d\n", Phonebook.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", Phonebook.Port), nil)
	if nil != err {
		fmt.Printf("*** Error on http.ListenAndServe: %v\n", err)
		ulog("*** Error on http.ListenAndServe: %v\n", err)
		os.Exit(1)
	}
}
