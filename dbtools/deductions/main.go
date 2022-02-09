package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
)

import _ "mysql"

// App is the global data structure for this app
var App struct {
	db         *sql.DB
	DBName     string
	DBUser     string
	JTFileName string
}

func errcheck(err error) {
	if nil != err {
		fmt.Printf("Error = %v\n", err)
		os.Exit(1)
	}
}

// DDUNKNOWN through DDTAXES are constants to represent
// the enumerations for Deductions
const (
	DDUNKNOWN      = iota // an unknown deduction
	DD401K                // 401K deduction
	DD401KLOAN            // 401K loan deduction
	DDCHILDSUPPORT        // Child Support deduction
	DDDENTAL              // dental coverage deduction
	DDFSA                 // FSA
	DDGARN                // garnished wages
	DDGROUPLIFE           // group life insurance
	DDHOUSING             // housing deduction
	DDMEDICAL             // medical insurance deducrtion
	DDMISCDED             // misc
	DDTAXES               // taxes
	DDEND                 // end of deduction list
)

// func deductionStringToInt(s string) int {
// 	var i int
// 	s = strings.ToUpper(s)
// 	s = strings.Replace(s, " ", "", -1)
// 	switch {
// 	case s == "401K":
// 		i = DD401K
// 	case s == "401KLOAN":
// 		i = DD401KLOAN
// 	case s == "CHILDSUPPORT":
// 		i = DDCHILDSUPPORT
// 	case s == "DENTAL":
// 		i = DDDENTAL
// 	case s == "FSA":
// 		i = DDFSA
// 	case s == "GARN":
// 		i = DDGARN
// 	case s == "GROUPLIFE":
// 		i = DDGROUPLIFE
// 	case s == "HOUSING":
// 		i = DDHOUSING
// 	case s == "MEDICAL":
// 		i = DDMEDICAL
// 	case s == "MISCDED":
// 		i = DDMISCDED
// 	case s == "TAXES":
// 		i = DDTAXES
// 	default:
// 		fmt.Printf("Unknown compensation type: %s\n", s)
// 		i = DDUNKNOWN
// 	}
// 	return i
// }

func deductionToString(i int) string {
	var s string
	switch {
	case i == DDUNKNOWN:
		s = "Unknown"
	case i == DD401K:
		s = "401K"
	case i == DD401KLOAN:
		s = "401K Loan"
	case i == DDCHILDSUPPORT:
		s = "Child Support"
	case i == DDFSA:
		s = "FSA"
	case i == DDGARN:
		s = "GARN"
	case i == DDGROUPLIFE:
		s = "Group Life"
	case i == DDHOUSING:
		s = "Housing"
	case i == DDDENTAL:
		s = "Dental"
	case i == DDMEDICAL:
		s = "Medical"
	case i == DDMISCDED:
		s = "Miscded"
	case i == DDTAXES:
		s = "Taxes"
	default:
		s = "UKNOWN COMPENSATION TYPE"
	}
	return s
}

func createDeductionsList(db *sql.DB) {
	Insrt, err := db.Prepare("INSERT INTO deductionlist (dcode,name) VALUES(?,?)")
	errcheck(err)

	for i := 0; i < DDEND; i++ {
		_, err := Insrt.Exec(i, deductionToString(i))
		errcheck(err)
	}
}

func readCommandLineArgs() {
	dbuPtr := flag.String("B", "ec2-user", "database user name")
	dbnmPtr := flag.String("N", "accordtest", "database name (accordtest, accord)")
	jobtPtr := flag.String("j", "jobtitles.csv", "The file containing the job titles to load")
	flag.Parse()
	App.DBName = *dbnmPtr
	App.JTFileName = *jobtPtr
	App.DBUser = *dbuPtr
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

	createDeductionsList(App.db)
}
