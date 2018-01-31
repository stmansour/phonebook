package db

import (
	"database/sql"
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
	UID           int
	LastName      string
	FirstName     string
	PreferredName string
	PrimaryEmail  string
	JobCode       int
	OfficePhone   string
	CellPhone     string
	OfficeFax     string
	DeptCode      int
	DeptName      string
	Employer      string
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
}

// Init initializes the database infrastructure
//
// INPUTS
//  name  - name of the db to load.  This name overrides the one in the config file
//          if its length is > 0
// RETURNS
//  error - any error encountered
// //-----------------------------------------------------------------------------
// func Init(overridename string) error {
// 	name := lib.AppConfig.Dbname
// 	if len(overridename) > 0 {
// 		name = overridename
// 	}
// 	dbopenparms := extres.GetSQLOpenString(name, &lib.AppConfig)
// 	db, err := sql.Open("mysql", dbopenparms)
// 	lib.Errcheck(err)
// 	defer db.Close()
// 	err = db.Ping()
// 	if nil != err {
// 		lib.Ulog("db.Ping: Error = %v\n", err)
// 		s := fmt.Sprintf("Could not establish database connection to db: %s, dbuser: %s\n", name, lib.AppConfig.Dbuser)
// 		lib.Ulog(s)
// 		fmt.Println(s)
// 		os.Exit(2)
// 	}
// 	lib.Ulog("MySQL database opened with \"%s\"\n", dbopenparms)
// 	DB.DirDB = db
// 	return nil
// }
