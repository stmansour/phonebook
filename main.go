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
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

import _ "github.com/go-sql-driver/mysql"

type company struct {
	CoCode      int
	Name        string
	Designation string
	Address     string
	Address2    string
	City        string
	State       string
	PostalCode  string
	Country     string
	Phone       string
	Fax         string
	Email       string
}

type person struct {
	UID          int
	LastName     string
	FirstName    string
	PrimaryEmail string
	JobCode      int
	OfficePhone  string
	CellPhone    string
	DeptCode     int
	DeptName     string
	Employer     string
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
	LastReview              string
	NextReview              string
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
	MgrName                 string
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

type searchResults struct {
	Query   string
	Matches []person
}

// uiSupport is an umbrella structure in which we can pass many useful
// data objects to the UI
type uiSupport struct {
	CoCodeToName     map[int]string // map from company code to company name
	NameToCoCode     map[string]int // map from company name to company code
	NameToJobCode    map[string]int // jobtitle to jobcode
	AcceptCodeToName map[int]string // Acceptance to jobcode
	NameToDeptCode   map[string]int // department name to dept code
	Months           []string       // a map for month number to month name
	D                *personDetail
}

// PhonebookUI is the instance of uiSupport used by this app
var PhonebookUI uiSupport

// Phonebook is the global application structure providing
// information that any function might need.
var Phonebook struct {
	db *sql.DB
}

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
		http.Redirect(w, r, "/search/", http.StatusFound)
	}
}

func loadMaps() {
	var code int
	var name string

	PhonebookUI.CoCodeToName = make(map[int]string)
	PhonebookUI.NameToCoCode = make(map[string]int)
	PhonebookUI.AcceptCodeToName = make(map[int]string)

	rows, err := Phonebook.db.Query("select cocode,name from companies")
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
	for i := ACPTUNKNOWN; i <= ACPTLAST; i++ {
		PhonebookUI.AcceptCodeToName[i] = acceptIntToString(i)
	}

	PhonebookUI.Months = []string{
		"January", "February", "March", "April",
		"May", "June", "July", "August",
		"September", "October", "November", "December",
	}
}

func shutdownHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<html><body><h1>shutting down in 5 seconds!</h1></body></html>")
	go func() {
		time.Sleep(time.Duration(5 * time.Second))
		os.Exit(0)
	}()
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

var chttp = http.NewServeMux()

func main() {
	db, err := sql.Open("mysql", "ec2-user:@/accord?charset=utf8&parseTime=True")
	if nil != err {
		fmt.Printf("sql.Open: Error = %v\n", err)
	}
	defer db.Close()

	err = db.Ping()
	if nil != err {
		fmt.Printf("db.Ping: Error = %v\n", err)
	}
	Phonebook.db = db
	loadMaps()

	chttp.Handle("/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/search/", searchHandler)
	http.HandleFunc("/detail/", detailHandler)
	http.HandleFunc("/editDetail/", editDetailHandler)
	http.HandleFunc("/savePersonDetails/", savePersonDetailsHandler)
	http.HandleFunc("/adminEdit/", adminEditHandler)
	http.HandleFunc("/adminView/", adminViewHandler)
	http.HandleFunc("/saveAdminEdit/", saveAdminEditHandler)
	http.HandleFunc("/company/", companyHandler)
	http.HandleFunc("/status/", statusHandler)
	http.HandleFunc("/shutdown/", shutdownHandler)

	err = http.ListenAndServe(":8250", nil)
	if nil != err {
		fmt.Printf("*** Error on http.ListenAndServe: %v\n", err)
		os.Exit(1)
	}
}
