// auser  a program to set the username for a person in the accord database
//        based on their UID
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

import _ "github.com/go-sql-driver/mysql"

// App is the global data structure for this app
var App struct {
	db       *sql.DB
	DBName   string
	UID      int
	username string
}

func readCommandLineArgs() {
	dbnmPtr := flag.String("N", "accordtest", "database name (accordtest, accord)")
	uPtr := flag.Int("u", 0, "user's UID")
	nPtr := flag.String("n", "", "new user name")
	flag.Parse()
	App.DBName = *dbnmPtr
	App.UID = *uPtr
	App.username = *nPtr
}

func strToInt(s string) int {
	if len(s) == 0 {
		return 0
	}
	s = strings.Trim(s, " \n\r")
	n, err := strconv.Atoi(s)
	if err != nil {
		fmt.Printf("Error converting %s to a number: %v\n", s, err)
		return 0
	}
	return n
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

	update, err := App.db.Prepare("update people set username=? where uid=?")
	if nil != err {
		fmt.Printf("error = %v\n", err)
		os.Exit(1)
	}
	_, err = update.Exec(App.username, App.UID)
	if nil != err {
		fmt.Printf("error = %v\n", err)
	} else {
		fmt.Printf("password for uid %d has been updated\n", App.UID)
	}

}
