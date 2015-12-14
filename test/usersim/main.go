// add a user
//   needs firstname, lastname, username, passwork, role

package main

import (
	"bytes"
	"compress/gzip"
	"crypto/sha512"
	"database/sql"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
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

// VUser is a structure with the basic information
// describing a virtual user
type VUser struct {
	UID                     int
	UserName                string
	FirstName               string
	LastName                string
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
	CoCode                  int
	MgrUID                  int
	JobTitle                string
	Class                   string
	ClassCode               int
	MgrName                 string
	Image                   string // ptr to image -- URI
	Deductions              []int
	DeductionsStr           string
	EmergencyContactName    string
	EmergencyContactPhone   string
	AcceptedHealthInsurance int
	AcceptedDentalInsurance int
	Accepted401K            int
	Hire                    time.Time
	Termination             time.Time
	//Company                 company
	//Reports                 []person
}

// App is the global data structure for this app
var App struct {
	seed             int64
	db               *sql.DB
	DBName           string
	DBUser           string
	FirstNames       []string
	LastNames        []string
	Streets          []string
	Cities           []string
	States           []string
	CoCodeToName     map[int]string // map from company code to company name
	NameToCoCode     map[string]int // map from company name to company code
	NameToJobCode    map[string]int // jobtitle to jobcode
	AcceptCodeToName map[int]string // Acceptance to jobcode
	NameToDeptCode   map[string]int // department name to dept code
	NameToClassCode  map[string]int // class designation to classcode
	ClassCodeToName  map[int]string // index by classcode to get the name
	Months           []string       // a map for month number to month name
	Roles            []Role         // the roles saved in the database
	JCLo, JCHi       int            // lo and high indeces for jobcode
	DeptLo, DeptHi   int            // lo and high indeces for department
	host             string         // phonebook server host
	port             int            // phonebook server port
}

func createUser(v *VUser) {
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
	v.DeptCode = rand.Intn(App.DeptLo + rand.Intn(App.DeptHi-App.DeptLo))
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
	fmt.Printf("Added user to database %s:  username: %s, access role: %d\n", App.DBName, v.UserName, v.RID)
}

func login(v *VUser) {
	URL := fmt.Sprint("http://%s:%d/signin", App.host, App.port)
	hc := http.Client{}

	form := url.Values{}
	form.Add("username", v.UserName)
	form.Add("password", "accord")
	// req.PostForm = form
	// req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	req, err := http.NewRequest("POST", URL, bytes.NewBufferString(form.Encode()))
	errcheck(err)

	// if App.debug {
	// 	dump, err := httputil.DumpRequest(req, false)
	// 	errcheck(err)
	// 	fmt.Printf("\n\ndumpRequest = %s\n", string(dump))
	// }

	resp, err := hc.Do(req)
	errcheck(err)
	defer resp.Body.Close()

	// if App.debug {
	// 	dump, err := httputil.DumpResponse(resp, true)
	// 	errcheck(err)
	// 	fmt.Printf("\n\ndumpResponse = %s\n", string(dump))
	// }

	// Check that the server actually sent compressed data
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		defer reader.Close()
	default:
		reader = resp.Body
	}

	fmt.Printf("response: %s\n", resp.Body)
}

func main() {
	readCommandLineArgs()

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
	for i := 0; i < 5; i++ {
		var v VUser
		createUser(&v)
		login(&v)
	}
}
