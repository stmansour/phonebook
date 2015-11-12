// Phonebook - a temporary directory interface
//  TODO:
//    change Status to a boolean
//    add employers addresses - add google maps insert
//    no one uses middle name?
//    Salutation ... needed, but no data
//    Preferred name  (ex: Joe rather than Joseph, Steve rather than Steven)
//    Email - delete?  just use PrimaryEmail and secondaryEmail
//    LastReviewDate, NextReviewDate ??  maybe just multivalued table of review dates?
//    Birthday -- maybe
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

import _ "github.com/go-sql-driver/mysql"

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

type personDetail struct {
	UID                     int
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

// uiSupport is an umbrella structure in which we can pass many useful
// data objects to the UI
type uiSupport struct {
	CoCodeToName     map[int]string // map from company code to company name
	NameToCoCode     map[string]int // map from company name to company code
	NameToJobCode    map[string]int // jobtitle to jobcode
	AcceptCodeToName map[int]string // Acceptance to jobcode
	NameToDeptCode   map[string]int // department name to dept code
	NameToClassCode  map[string]int // class designation to classcode
	ClassCodeToName  map[int]string // index by classcode to get the name
	Months           []string       // a map for month number to month name
	C                *company
	A                *class
	D                *personDetail
	R                *searchResults
	S                *signin
	T                *searchCoResults
	L                *searchClassResults
	X                *session
}

// PhonebookUI is the instance of uiSupport used by this app
var PhonebookUI uiSupport

// Phonebook is the global application structure providing
// information that any function might need.
var Phonebook struct {
	Port          int // port on which we listen
	db            *sql.DB
	LogFile       *os.File
	ReqMem        chan int // request to access UI data memory
	ReqMemAck     chan int // done with memory
	DebugToScreen bool
	Debug         bool
}

var chttp = http.NewServeMux()

func errcheck(err error) {
	if err != nil {
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

func loadMaps() {
	var code int
	var name string

	PhonebookUI.CoCodeToName = make(map[int]string)
	PhonebookUI.NameToCoCode = make(map[string]int)
	PhonebookUI.AcceptCodeToName = make(map[int]string)

	rows, err := Phonebook.db.Query("select cocode,CommonName from companies")
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		errcheck(rows.Scan(&code, &name))
		PhonebookUI.CoCodeToName[code] = name
		PhonebookUI.NameToCoCode[name] = code
	}
	errcheck(rows.Err())

	PhonebookUI.NameToJobCode = make(map[string]int)
	rows, err = Phonebook.db.Query("select jobcode,title from jobtitles")
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		errcheck(rows.Scan(&code, &name))
		PhonebookUI.NameToJobCode[name] = code
	}
	errcheck(rows.Err())

	PhonebookUI.NameToDeptCode = make(map[string]int)
	rows, err = Phonebook.db.Query("select deptcode,name from departments")
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		errcheck(rows.Scan(&code, &name))
		PhonebookUI.NameToDeptCode[name] = code
	}
	errcheck(rows.Err())

	PhonebookUI.NameToClassCode = make(map[string]int)
	PhonebookUI.ClassCodeToName = make(map[int]string)
	rows, err = Phonebook.db.Query("select classcode,designation from classes")
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		errcheck(rows.Scan(&code, &name))
		PhonebookUI.NameToClassCode[name] = code
		PhonebookUI.ClassCodeToName[code] = name
	}
	// for k, v := range PhonebookUI.NameToClassCode {
	// 	fmt.Printf("%s %d\n", k, v)
	// }
	errcheck(rows.Err())

	for i := ACPTUNKNOWN; i <= ACPTLAST; i++ {
		PhonebookUI.AcceptCodeToName[i] = acceptIntToString(i)
	}

	PhonebookUI.Months = make([]string, len(fmtMonths))
	for i := 0; i < len(fmtMonths); i++ {
		PhonebookUI.Months[i] = fmtMonths[i]

	}
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
	http.HandleFunc("/class/", classHandler)
	http.HandleFunc("/company/", companyHandler)
	http.HandleFunc("/detail/", detailHandler)
	http.HandleFunc("/editDetail/", editDetailHandler)
	http.HandleFunc("/logoff/", logoffHandler)
	http.HandleFunc("/saveAdminEdit/", saveAdminEditHandler)
	http.HandleFunc("/saveAdminEditClass/", saveAdminEditClassHandler)
	http.HandleFunc("/saveAdminEditCo/", saveAdminEditCoHandler)
	http.HandleFunc("/savePersonDetails/", savePersonDetailsHandler)
	http.HandleFunc("/search/", searchHandler)
	http.HandleFunc("/searchcl/", searchClassHandler)
	http.HandleFunc("/searchco/", searchCompaniesHandler)
	http.HandleFunc("/shutdown/", shutdownHandler)
	http.HandleFunc("/signin/", signinHandler)
	http.HandleFunc("/status/", statusHandler)
	http.HandleFunc("/weblogin/", webloginHandler)
}

func shutdownHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<html><body><h1>shutting down in 5 seconds!</h1></body></html>")
	ulog("Shutdown initiated from web service\n")
	go func() {
		time.Sleep(time.Duration(5 * time.Second))
		os.Exit(0)
	}()
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func readCommandLineArgs() {
	portPtr := flag.Int("p", 8250, "port on which Phonebook listens")
	dbugPtr := flag.Bool("d", false, "debug mode - includes debug info in logfile")
	dtscPtr := flag.Bool("D", false, "LogToScreen mode - prints log messages to stdout")
	flag.Parse()

	Phonebook.Port = *portPtr
	Phonebook.Debug = *dbugPtr
	Phonebook.DebugToScreen = *dtscPtr
}

func main() {
	var err error
	Phonebook.LogFile, err = os.OpenFile("Phonebook.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer Phonebook.LogFile.Close()
	log.SetOutput(Phonebook.LogFile)
	ulog("*** Accord PHONEBOOK ***\n")

	dbopenparms := "ec2-user:@/accord?charset=utf8&parseTime=True"
	db, err := sql.Open("mysql", dbopenparms)
	if nil != err {
		ulog("sql.Open: Error = %v\n", err)
	}
	defer db.Close()

	err = db.Ping()
	if nil != err {
		ulog("db.Ping: Error = %v\n", err)
	}
	ulog("MySQL database opened with \"%s\"\n", dbopenparms)

	Phonebook.db = db
	Phonebook.ReqMem = make(chan int)
	Phonebook.ReqMemAck = make(chan int)
	loadMaps()

	readCommandLineArgs()

	go Dispatcher()

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
