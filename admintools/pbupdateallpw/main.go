// pbupdateallpw = a program to set the password for all users
//
package main

import (
	"crypto/sha512"
	"database/sql"
	"flag"
	"fmt"
	"os"
)

import _ "github.com/go-sql-driver/mysql"

// App is the global data structure for this app
var App struct {
	db       *sql.DB
	DBName   string
	DBUser   string
	password string
}

func readCommandLineArgs() {
	dbuPtr := flag.String("B", "ec2-user", "database user name")
	dbnmPtr := flag.String("N", "accordtest", "database name (accordtest, accord)")
	psPtr := flag.String("p", "accord", "password")
	flag.Parse()
	App.DBName = *dbnmPtr
	App.password = *psPtr
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

	sha := sha512.Sum512([]byte(App.password))
	passhash := fmt.Sprintf("%x", sha)
	update, err := App.db.Prepare("update people set passhash=?")
	if nil != err {
		fmt.Printf("error = %v\n", err)
		os.Exit(1)
	}
	_, err = update.Exec(passhash)
	if nil != err {
		fmt.Printf("error = %v\n", err)
	} else {
		fmt.Printf("password for all users has been set to \"%s\"\n", App.password)
	}

}
