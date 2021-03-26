package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"phonebook/lib"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

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

// CreateJobTitlesTable not only creates the JobTitles table, it makes a pass through
// the people table, replaces the  Title field with the appropriate deptcode field.
func CreateJobTitlesTable(db *sql.DB) {
	//--------------------------------------------------------------------------
	// Populate the JobTitles table
	//--------------------------------------------------------------------------
	InsertJT, err := db.Prepare("INSERT INTO jobtitles (title) VALUES(?)")
	errcheck(err)
	jobtitles := "jobtitles.csv"
	f, err := os.Open(jobtitles)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		s := strings.Trim(line, " \n\r")
		_, err := InsertJT.Exec(s)

		errcheck(err)
	}
}

func readCommandLineArgs() {
	dbuPtr := flag.String("B", "ec2-user", "database user name")
	dbnmPtr := flag.String("N", "accord", "database name (accordtest, accord)")
	jobtPtr := flag.String("j", "jobtitles.csv", "The file containing the job titles to load")
	flag.Parse()
	App.DBName = *dbnmPtr
	App.JTFileName = *jobtPtr
	App.DBUser = *dbuPtr
}

func main() {
	var err error
	fmt.Printf("reading command line args\n")
	readCommandLineArgs()
	fmt.Printf("ReadConfig\n")
	lib.ReadConfig()
	fmt.Printf("GetSQLOpenString(%q,%q)\n", App.DBUser, App.DBName)
	dbopenparms := lib.GetSQLOpenString(App.DBUser, App.DBName)
	fmt.Printf("sql.Open(), dbopenparms = %s\n", dbopenparms)
	App.db, err = sql.Open("mysql", dbopenparms)
	lib.Errcheck(err)
	defer App.db.Close()
	fmt.Printf("Ping\n")
	err = App.db.Ping()
	if nil != err {
		fmt.Printf("App.db.Ping: Error = %v\n", err)
		s := fmt.Sprintf("Could not establish database connection to pbdb: %s, dbuser: %s\n", App.DBName, App.DBUser)
		// ulog(s)
		fmt.Println(s)
		os.Exit(2)
	}

	// var err error
	// s := fmt.Sprintf("%s:@/%s?charset=utf8&parseTime=True", App.DBUser, App.DBName)
	// App.db, err = sql.Open("mysql", s)
	// if nil != err {
	// 	fmt.Printf("sql.Open: Error = %v\n", err)
	// }
	// defer App.db.Close()
	// err = App.db.Ping()
	// if nil != err {
	// 	fmt.Printf("App.db.Ping: Error = %v\n", err)
	// }
	fmt.Printf("CreateJobTitlesTable\n")

	CreateJobTitlesTable(App.db)
}
