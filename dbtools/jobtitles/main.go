package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"
)

import _ "github.com/go-sql-driver/mysql"

// App is the global data structure for this app
var App struct {
	db         *sql.DB
	DBName     string
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
	InsertJT, err := db.Prepare("INSERT INTO JobTitles (title) VALUES(?)")
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
	dbnmPtr := flag.String("N", "accord", "database name (accordtest, accord)")
	jobtPtr := flag.String("j", "jobtitles.csv", "The file containing the job titles to load")
	flag.Parse()
	App.DBName = *dbnmPtr
	App.JTFileName = *jobtPtr
}

func main() {
	readCommandLineArgs()

	var err error
	s := fmt.Sprintf("sman:@/%s?charset=utf8&parseTime=True", App.DBName)
	App.db, err = sql.Open("mysql", s)
	if nil != err {
		fmt.Printf("sql.Open: Error = %v\n", err)
	}
	defer App.db.Close()
	err = App.db.Ping()
	if nil != err {
		fmt.Printf("App.db.Ping: Error = %v\n", err)
	}

	CreateJobTitlesTable(App.db)
}
