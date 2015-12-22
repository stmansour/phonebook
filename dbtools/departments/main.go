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
	DBUser     string
	JTFileName string
}

func errcheck(err error) {
	if nil != err {
		fmt.Printf("Error = %v\n", err)
		os.Exit(1)
	}
}

func loadDepartments(db *sql.DB) {
	//--------------------------------------------------------------------------
	// Populate the JobTitles table
	//--------------------------------------------------------------------------
	InsertJT, err := db.Prepare("INSERT INTO Departments (name) VALUES(?)")
	errcheck(err)
	filename := "depts.csv"
	f, err := os.Open(filename)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		da := strings.Split(line, ",")
		_, err := InsertJT.Exec(strings.Trim(da[0], " \n\r"))
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

	loadDepartments(App.db)
}
