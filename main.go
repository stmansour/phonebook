// Phonebook - a temporary directory interface

package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"phonebook/authz"
	"phonebook/db"
	"phonebook/lib"
	"phonebook/sess"
	"phonebook/ws"
	"runtime/debug"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

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
	{authz.ELEMPERSON, "Status", false, "Indicates whether the person is an active employee."},
	{authz.ELEMPERSON, "EligibleForRehire", false, "Indicates whether a past employee can be rehired."},
	{authz.ELEMPERSON, "UID", false, "A unique identifier associated with the employee. Once created, it never changes."},
	{authz.ELEMPERSON, "Salutation", false, "'Mr.', 'Mrs.', 'Ms.', etc."},
	{authz.ELEMPERSON, "FirstName", false, "The person's common name."},
	{authz.ELEMPERSON, "MiddleName", false, "The person's middle name."},
	{authz.ELEMPERSON, "LastName", false, "The person's surname or last name."},
	{authz.ELEMPERSON, "PreferredName", false, "Less formal name but more commonly used, for example 'Mike' rather than 'Michael'."},
	{authz.ELEMPERSON, "PrimaryEmail", false, "The primary email address to use for this person."},
	{authz.ELEMPERSON, "OfficePhone", false, "This person's office telephone number."},
	{authz.ELEMPERSON, "CellPhone", false, "This person's cellphone number."},
	{authz.ELEMPERSON, "EmergencyContactName", true, "Name of someone to contact in the event of an emergency."},
	{authz.ELEMPERSON, "EmergencyContactPhone", true, "Phone number for the emergency contact."},
	{authz.ELEMPERSON, "HomeStreetAddress", true, "def"},
	{authz.ELEMPERSON, "HomeStreetAddress2", true, "def"},
	{authz.ELEMPERSON, "HomeCity", true, "def"},
	{authz.ELEMPERSON, "HomeState", true, "def"},
	{authz.ELEMPERSON, "HomePostalCode", true, "def"},
	{authz.ELEMPERSON, "HomeCountry", true, "def"},
	{authz.ELEMPERSON, "PrimaryEmail", true, "def"},
	{authz.ELEMPERSON, "SecondaryEmail", true, "def"},
	{authz.ELEMPERSON, "OfficePhone", true, "def"},
	{authz.ELEMPERSON, "OfficeFax", true, "def"},
	{authz.ELEMPERSON, "CellPhone", true, "def"},
	{authz.ELEMPERSON, "BirthDOM", true, "def"},
	{authz.ELEMPERSON, "BirthMonth", true, "def"},
	{authz.ELEMPERSON, "CoCode", true, "The company code associated with this user."},
	{authz.ELEMPERSON, "JobCode", true, "def"},
	{authz.ELEMPERSON, "ClassCode", true, "def"},
	{authz.ELEMPERSON, "DeptCode", true, "def"},
	{authz.ELEMPERSON, "PositionControlNumber", true, "def"},
	{authz.ELEMPERSON, "MgrUID", true, "def"},
	{authz.ELEMPERSON, "Accepted401K", true, "def"},
	{authz.ELEMPERSON, "AcceptedDentalInsurance", true, "def"},
	{authz.ELEMPERSON, "AcceptedHealthInsurance", true, "def"},
	{authz.ELEMPERSON, "Hire", true, "def"},
	{authz.ELEMPERSON, "Termination", true, "def"},
	{authz.ELEMPERSON, "LastReview", true, "def"},
	{authz.ELEMPERSON, "NextReview", true, "def"},
	{authz.ELEMPERSON, "StateOfEmployment", false, "def"},
	{authz.ELEMPERSON, "CountryOfEmployment", false, "def"},
	{authz.ELEMPERSON, "Comps", true, "def"},
	{authz.ELEMPERSON, "Deductions", true, "def"},
	{authz.ELEMPERSON, "MyDeductions", true, "def"},
	{authz.ELEMPERSON, "RID", true, "role identifier"},
	{authz.ELEMPERSON, "ElemEntity", true, "The entire entity"},
	{authz.ELEMCOMPANY, "CoCode", false, "def"},
	{authz.ELEMCOMPANY, "LegalName", false, "def"},
	{authz.ELEMCOMPANY, "CommonName", false, "def"},
	{authz.ELEMCOMPANY, "Address", false, "def"},
	{authz.ELEMCOMPANY, "Address2", false, "def"},
	{authz.ELEMCOMPANY, "City", false, "def"},
	{authz.ELEMCOMPANY, "State", false, "def"},
	{authz.ELEMCOMPANY, "PostalCode", false, "def"},
	{authz.ELEMCOMPANY, "Country", false, "def"},
	{authz.ELEMCOMPANY, "Phone", false, "def"},
	{authz.ELEMCOMPANY, "Fax", false, "def"},
	{authz.ELEMCOMPANY, "Email", false, "def"},
	{authz.ELEMCOMPANY, "Designation", false, "def"},
	{authz.ELEMCOMPANY, "Active", false, "def"},
	{authz.ELEMCOMPANY, "EmploysPersonnel", false, "def"},
	{authz.ELEMCOMPANY, "ElemEntity", true, "The entire entity"},
	{authz.ELEMCLASS, "ClassCode", false, "def"},
	{authz.ELEMCLASS, "Name", false, "def"},
	{authz.ELEMCLASS, "Designation", false, "def"},
	{authz.ELEMCLASS, "Description", false, "def"},
	{authz.ELEMCLASS, "ElemEntity", true, "The entire entity"},
	{authz.ELEMPBSVC, "Shutdown", true, "Shut down the running Phonebook service"},
	{authz.ELEMPBSVC, "Restart", true, "Restart the running Phonebook service"},
}

type searchResults struct {
	Query   string
	Matches []db.Person
}

type searchCoResults struct {
	Query   string
	Matches []db.Company
}

type searchClassResults struct {
	Query   string
	Matches []db.Class
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
	CoCodeToName     map[int]string    // map from company code to company name
	NameToCoCode     map[string]int    // map from company name to company code
	NameToJobCode    map[string]int    // jobtitle to jobcode
	AcceptCodeToName map[int]string    // Acceptance to jobcode
	NameToDeptCode   map[string]int    // department name to dept code
	NameToClassCode  map[string]int    // db.Class designation to classcode
	ClassCodeToName  map[int]string    // index by classcode to get the name
	Months           []string          // a map for month number to month name
	Roles            []authz.Role      // list of roles -- fields are not initialized
	Images           map[string]string // interface images
	CompanyList      []db.Company      // list of all company structs
	C                *db.Company
	A                *db.Class
	D                *db.PersonDetail
	R                *searchResults
	S                *signin
	T                *searchCoResults
	L                *searchClassResults
	X                *sess.Session
	K                *UsageCounters
	Ki               *UsageCounters
	N                []sess.Session
	ErrMsg           template.HTML // if the caller wants to convey an error message
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
	classInfo          *sql.Stmt // get db.Class attributes
	companyInfo        *sql.Stmt // company attributes
	countersUpdate     *sql.Stmt // feature usage counters update
	delClass           *sql.Stmt // deletes a db.Class
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
	insertClass        *sql.Stmt // adding a new db.Class
	classReadBack      *sql.Stmt // read back newly written db.Class
	updateClass        *sql.Stmt // update a db.Class
	insertCompany      *sql.Stmt // insert a new company
	companyReadback    *sql.Stmt // read back newly written company
	updateCompany      *sql.Stmt // update a company
	updateMyDetails    *sql.Stmt // person updating their own details
	updatePasswd       *sql.Stmt // person updating their passwd
	readFieldPerms     *sql.Stmt // read field permissions
	accessRoles        *sql.Stmt // read access roles
	getUserCoCode      *sql.Stmt // read the cocode for a person
	loginInfo          *sql.Stmt // read info for login
	CompanyClasses     *sql.Stmt // read a list of classes that belong to a company
	GetAllCompanies    *sql.Stmt // query to select all companies
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
	ReqMem             chan int      // request to access UI data memory
	ReqMemAck          chan int      // done with memory
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
// in the database. Counters is the - incremental update -- the number
// of operations done since the last save.
var Counters UsageCounters

// TotCounters is the total number of the counters across all servers
var TotCounters UsageCounters

var funcMap map[string]interface{}

var chttp = http.NewServeMux()

func errcheck(err error) {
	if err != nil {
		fmt.Printf("error = %v\n", err)
		debug.PrintStack()
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

// find the first filename that matches the base filename in /images/
func findMatchingFilename(base string, backup string) string {

	m, err := filepath.Glob(fmt.Sprintf("./images/%s.*", base))
	if nil != err {
		ulog("filepath.Glob returned error: %v\n", err)
	}
	if len(m) == 0 {
		// copy the default file into /images/
		fmt.Printf("COPY DEFAULT FILE TO IMAGES\n")
		return backup
	}
	return m[0]
}

// load the branding images...
var uiDflt = []string{
	"logo",
	"search",
	"searchco",
	"searchcl",
	"detail",
	"company",
	"class",
	"signin",
	"signinlogo",
	"adminView",
	"adminEdit",
	"adminEditClass",
	"adminEditCo",
	"stats",
	"setup",
	"delPersonRefErr",
	"delCoRefErr",
	"delClassRefErr",
}

func initUI() {
	PhonebookUI.Images = make(map[string]string)
	for i := 0; i < len(uiDflt); i++ {
		PhonebookUI.Images[uiDflt[i]] = findMatchingFilename(uiDflt[i], uiDflt[i]+".png")
	}

	// for k, v := range PhonebookUI.Images {
	// 	fmt.Printf("%s -> %s\n", k, v)
	// }
}

// func bugCheck(u *uiSupport) {
// 	var m []int
// 	for k := range PhonebookUI.CoCodeToName {
// 		m = append(m, k)
// 	}
// 	sort.Ints(m)

// 	k := m[0]
// 	for i := 1; i < len(m); i++ {
// 		if k+1 != m[i] {
// 			fmt.Printf("k=%d, i=%d, m[i]=%d\n", k, i, m[i])
// 			fmt.Printf("m = %#v\n", m)
// 			os.Exit(1)
// 		}
// 		k = m[i]
// 	}
// }

func initUIData(u *uiSupport) {
	u.Images = make(map[string]string, len(PhonebookUI.Images))
	for k, v := range PhonebookUI.Images {
		u.Images[k] = v
	}
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
	u.Roles = make([]authz.Role, len(authz.Authz.Roles))
	for i := 0; i < len(authz.Authz.Roles); i++ {
		u.Roles[i] = authz.Role{}
		u.Roles[i].Name = authz.Authz.Roles[i].Name
		u.Roles[i].RID = authz.Authz.Roles[i].RID
	}
	u.ErrMsg = ""
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
	PhonebookUI.CoCodeToName = make(map[int]string)
	PhonebookUI.NameToCoCode = make(map[string]int)
	PhonebookUI.CompanyList = PhonebookUI.CompanyList[:0] // empty it first

	rows, err := Phonebook.prepstmt.GetAllCompanies.Query()
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		var c db.Company
		errcheck(rows.Scan(&c.CoCode, &c.LegalName, &c.CommonName, &c.Address, &c.Address2, &c.City, &c.State, &c.PostalCode, &c.Country, &c.Phone, &c.Fax, &c.Email, &c.Designation, &c.Active, &c.EmploysPersonnel))
		PhonebookUI.CompanyList = append(PhonebookUI.CompanyList, c)
		if c.EmploysPersonnel != 0 {
			PhonebookUI.CoCodeToName[c.CoCode] = c.LegalName
			PhonebookUI.NameToCoCode[c.LegalName] = c.CoCode
		}

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

func getVer() string {
	return lib.GetVersionNo()
}
func getBTime() string {
	return lib.GetBuildTime()
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
		"smrand":               smrand,
		"hasFieldAccess":       hasFieldAccess,
		"hasPERMMODaccess":     hasPERMMODaccess,
		"hasAdminScreenAccess": hasAdminScreenAccess,
		"showAdminButton":      showAdminButton,
		"getBreadcrumb":        getBreadcrumb,
		"getHTMLBreadcrumb":    getHTMLBreadcrumb,
		"datetimeToString":     datetimeToString,
		"phoneURL":             phoneURL,
		"mapURL":               mapURL,
		"GetVersionNo":         getVer,
		"GetBuildTime":         getBTime,
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

	ReadTotalCounters()

	Counters.SearchPeople = 0
	Counters.SearchClasses = 0
	Counters.SearchCompanies = 0
	Counters.EditPerson = 0
	Counters.ViewPerson = 0
	Counters.ViewClass = 0
	Counters.ViewCompany = 0
	Counters.AdminEditPerson = 0
	Counters.AdminEditClass = 0
	Counters.AdminEditCompany = 0
	Counters.DeletePerson = 0
	Counters.DeleteClass = 0
	Counters.DeleteCompany = 0
	Counters.SignIn = 0
	Counters.Logoff = 0

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
	http.HandleFunc("/become/", adminBecomeHandler)
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
	http.HandleFunc("/resetpw/", resetpwHandler)
	http.HandleFunc("/restart/", restartHandler)
	http.HandleFunc("/saveAdminEdit/", saveAdminEditHandler)
	http.HandleFunc("/saveAdminEditClass/", saveAdminEditClassHandler)
	http.HandleFunc("/saveAdminEditCo/", saveAdminEditCoHandler)
	http.HandleFunc("/savePersonDetails/", savePersonDetailsHandler)
	http.HandleFunc("/saveSetup/", saveSetupHandler)
	http.HandleFunc("/search/", searchHandler)
	http.HandleFunc("/searchcl/", searchClassHandler)
	http.HandleFunc("/searchco/", searchCompaniesHandler)
	http.HandleFunc("/setup/", setupHandler)
	http.HandleFunc("/shutdown/", shutdownHandler)
	http.HandleFunc("/signin/", signinHandler)
	http.HandleFunc("/status/", statusHandler)
	http.HandleFunc("/stats/", statsHandler)
	http.HandleFunc("/weblogin/", webloginHandler)
	http.HandleFunc("/v1/", ws.V1ServiceHandler)

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
	vPtr := flag.Bool("v", false, "version request - dumps version to stdout")

	flag.Parse()

	if *vPtr {
		fmt.Printf("Version: %s\n", lib.GetVersionNo())
		os.Exit(0)
	}

	Phonebook.Port = *portPtr
	Phonebook.Debug = *dbugPtr
	Phonebook.SecurityDebug = *sbugPtr
	Phonebook.DebugToScreen = *dtscPtr
	Phonebook.DBName = *dbnmPtr
	Phonebook.DBUser = *dbusPtr
	Phonebook.CountersUpdateTime = *cntrPtr
}

func main() {
	rand.Seed(time.Now().UnixNano()) // We need this for generating random passwords, probably other random things too.

	//=============================
	//  Hardcoded defaults...
	//=============================
	Phonebook.ReqMem = make(chan int)
	Phonebook.ReqMemAck = make(chan int)
	Phonebook.ReqCountersMem = make(chan int)
	Phonebook.ReqCountersMemAck = make(chan int)
	Phonebook.SessionTimeout = 15     // minutes
	Phonebook.SessionCleanupTime = 10 // minutes
	authz.Init(Phonebook.SecurityDebug)

	//==============================================
	// There may be some command line overrides...
	//==============================================
	readCommandLineArgs()

	//==============================================
	// Now open the logfile and the database...
	//==============================================
	var err error
	Phonebook.LogFile, err = os.OpenFile("Phonebook.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	lib.Errcheck(err)
	defer Phonebook.LogFile.Close()
	log.SetOutput(Phonebook.LogFile)
	ulog("*** Accord PHONEBOOK ***\n")

	// lib.ReadConfig()
	// db.Init(Phonebook.DBName)
	// Phonebook.db = db.DB.DirDB
	// buildPreparedStatements()
	lib.ReadConfig()
	dbopenparms := lib.GetSQLOpenString(Phonebook.DBUser, Phonebook.DBName)
	pbdb, err := sql.Open("mysql", dbopenparms)
	lib.Errcheck(err)
	defer pbdb.Close()
	err = pbdb.Ping()
	if nil != err {
		ulog("pbdb.Ping: Error = %v\n", err)
		s := fmt.Sprintf("Could not establish database connection to pbdb: %s, dbuser: %s\n", Phonebook.DBName, Phonebook.DBUser)
		ulog(s)
		fmt.Println(s)
		os.Exit(2)
	}
	ulog("MySQL database opened with \"%s\"\n", dbopenparms)
	Phonebook.db = pbdb
	db.DB.DirDB = pbdb
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
	initUI()
	sess.InitSessionManager(Phonebook.SessionCleanupTime, Phonebook.SessionTimeout, pbdb, Phonebook.SecurityDebug)
	go Dispatcher()
	go CounterDispatcher()
	go UpdateCounters()

	initHTTP()
	ws.InitServices(Phonebook.db)

	ulog("Phonebook initiating HTTP service on port %d\n", Phonebook.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", Phonebook.Port), nil)
	if nil != err {
		fmt.Printf("*** Error on http.ListenAndServe: %v\n", err)
		ulog("*** Error on http.ListenAndServe: %v\n", err)
		os.Exit(1)
	}
}
