// auser  a program to set the username for a person in the accord database
//        based on their UID
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"phonebook/lib"
	"strconv"
	"strings"
)

import _ "github.com/go-sql-driver/mysql"

// App is the global data structure for this app
var App struct {
	db       *sql.DB
	DBName   string
	DBUser   string
	UID      int
	username string
}

func readCommandLineArgs() {
	dbuPtr := flag.String("B", "ec2-user", "database user name")
	dbnmPtr := flag.String("N", "accord", "database name (accordtest, accord)")
	nPtr := flag.String("n", "", "new user name")
	uPtr := flag.Int("u", 0, "user's UID")
	flag.Parse()
	App.DBName = *dbnmPtr
	App.DBUser = *dbuPtr
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

	if "" == App.username {
		fmt.Printf("You must supply the username using the -n option\n")
		os.Exit(2)
	}
	if 0 == App.UID {
		fmt.Printf("You must supply the uid of the user for whom you wish to change the username using the -u option\n")
		os.Exit(2)
	}
	var err error
	// s := fmt.Sprintf("ec2-user:@/%s?charset=utf8&parseTime=True", App.DBName)
	lib.ReadConfig()
	s := lib.GetSQLOpenString(App.DBUser, App.DBName)
	App.db, err = sql.Open("mysql", s)
	if nil != err {
		fmt.Printf("sql.Open for database=%s, dbuser=%s: Error = %v\n", App.DBName, App.DBUser, err)
	}
	defer App.db.Close()
	err = App.db.Ping()
	if nil != err {
		fmt.Printf("App.db.Ping for database=%s, dbuser=%s: Error = %v\n", App.DBName, App.DBUser, err)
	}

	var uid int
	err = App.db.QueryRow("select uid from people where username=?", App.username).Scan(&uid)
	switch {
	case err == sql.ErrNoRows:
		fmt.Printf("username = %s is available for use in database %s\n", App.username, App.DBName)
	case err != nil:
		fmt.Printf("error with QueryRow selecting username: %s,  error = %v\n", App.username, err)
		os.Exit(1)
	default:
		fmt.Printf("username %s is already being used in database %s. UID = %d\n", App.username, App.DBName, uid)
		os.Exit(2)
	}

	update, err := App.db.Prepare("update people set username=? where uid=?")
	if nil != err {
		fmt.Printf("error = %v\n", err)
		os.Exit(1)
	}
	c, err := update.Exec(App.username, App.UID)
	if nil != err {
		switch err {
		case sql.ErrNoRows:
			fmt.Printf("Database %s does not contain a user with uid=%d\n", App.DBName, App.UID)
			os.Exit(1)
		default:
			fmt.Printf("error = %v\n", err)
			os.Exit(1)
		}
	} else {
		n, _ := c.RowsAffected()
		if n == 0 {
			fmt.Printf("Database %s does not contain a user with uid=%d\n", App.DBName, App.UID)
			os.Exit(1)
		}
		fmt.Printf("uid %d now has username %s\n", App.UID, App.username)
	}
}
