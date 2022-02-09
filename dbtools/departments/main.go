package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"phonebook/lib"
	"strings"

	_ "mysql"
)

// App is the global data structure for this app
var App struct {
	db       *sql.DB
	DBName   string
	DBUser   string
	FileName string
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
	InsertJT, err := db.Prepare("INSERT INTO departments (name) VALUES(?)")
	errcheck(err)
	f, err := os.Open(App.FileName)
	if err != nil {
		fmt.Printf("Error opening file: %s :: %s\n", App.FileName, err.Error())
		os.Exit(1)
	}
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
	jobtPtr := flag.String("d", "depts.csv", "The file containing the job titles to load")
	flag.Parse()
	App.DBName = *dbnmPtr
	App.FileName = *jobtPtr
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

	fmt.Printf("loadDepartments\n")

	loadDepartments(App.db)
	fmt.Printf("Completed!\n")
}
