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
	"strings"
)

import _ "github.com/go-sql-driver/mysql"

// Phonebook is the global application structure providing
// information that any function might need.
var Phonebook struct {
	db *sql.DB
}

type company struct {
	CoCode     int
	Company    string
	Address    string
	Address2   string
	City       string
	State      string
	PostalCode string
	Country    string
	Phone      string
	Fax        string
	Email      string
}

type person struct {
	UID          int
	LastName     string
	FirstName    string
	PrimaryEmail string
	JobCode      int
	OfficePhone  string
	CellPhone    string
	Department   string
	Employer     string
}

type personDetail struct {
	UID                   int
	LastName              string
	FirstName             string
	PrimaryEmail          string
	JobCode               int
	OfficePhone           string
	CellPhone             string
	Department            string
	MiddleName            string
	Salutation            string
	Status                string
	PositionControlNumber string
	OfficeFax             string
	SecondaryEmail        string
	EligibleForRehire     string
	LastReview            string
	NextReview            string
	Birthdate             string
	HomeStreetAddress     string
	HomeStreetAddress2    string
	HomeCity              string
	HomeState             string
	HomePostalCode        string
	HomeCountry           string
	StateOfEmployment     string
	CountryOfEmployment   string
	PreferredName         string
	CompensationType      string
	DeptCode              int
	Company               company
	CoCode                int
	MgrUID                int
	JobTitle              string
	CostCenter            string
	MgrName               string
	Reports               []person
	EmergencyContactName  string
	EmergencyContactPhone string
	// HealthInsuranceAccepted string
	// DentalInsuranceAccepted string
	// Accepted401K string
}

type searchResults struct {
	Query   string
	Matches []person
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

var chttp = http.NewServeMux()

func main() {
	db, err := sql.Open("mysql", "sman:@/smtest")
	if nil != err {
		fmt.Printf("sql.Open: Error = %v\n", err)
	}
	defer db.Close()

	err = db.Ping()
	if nil != err {
		fmt.Printf("db.Ping: Error = %v\n", err)
	}
	Phonebook.db = db

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

	http.ListenAndServe(":8250", nil)
}
