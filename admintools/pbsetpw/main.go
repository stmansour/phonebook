// apasswd = a program to set the password for a user
//
package main

import (
	"context"
	"database/sql"
	"extres"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"phonebook/db"
	"phonebook/lib"
	"time"

	_ "mysql"
)

// App is the global data structure for this app
var App struct {
	db           *sql.DB
	DBName       string
	DBUser       string
	user         string
	password     string
	fname        string
	lname        string
	usernameonly bool
}

func errcheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func getUserName() {
	n := 0
	s := fmt.Sprintf("select username from people where FirstName=\"%s\" AND LastName=\"%s\"", App.fname, App.lname)
	fmt.Println(s)
	rows, err := App.db.Query(s)
	errcheck(err)
	defer rows.Close()
	fmt.Printf("usernames for %s %s:\n", App.fname, App.lname)
	for rows.Next() {
		errcheck(rows.Scan(&App.user))
		fmt.Println(App.user)
		n++
	}
	errcheck(rows.Err())
	if n > 1 {
		fmt.Printf("There are multiple usernames for the %s %s.\n", App.fname, App.lname)
		fmt.Printf("Select appropriate username and run this program again with -u and the appropriate username\n")
		os.Exit(1)
	}
	if n == 0 {
		fmt.Printf("Database %s does not have a user named %s %s\n", App.DBName, App.fname, App.lname)
		os.Exit(1)
	}
}

func getRealName() string {
	var f, l, p string
	s := fmt.Sprintf("select FirstName,LastName,PreferredName from people where username=\"%s\"", App.user)
	errcheck(App.db.QueryRow(s).Scan(&f, &l, &p))
	if len(p) > 0 {
		f = p
	}
	return fmt.Sprintf("%s %s", f, l)
}
func readCommandLineArgs() {
	dbuPtr := flag.String("B", "ec2-user", "database user name")
	dbnmPtr := flag.String("N", "accord", "database name (accordtest, accord)")
	fnPtr := flag.String("F", "", "User first name")
	lnPtr := flag.String("L", "", "User last name")
	uPtr := flag.String("u", "username", "username")
	psPtr := flag.String("p", "", "password")
	unoPtr := flag.Bool("n", false, "if present, just dump the username, do not make password changes.")
	flag.Parse()
	App.DBName = *dbnmPtr
	App.user = *uPtr
	App.password = *psPtr
	App.DBUser = *dbuPtr
	App.fname = *fnPtr
	App.lname = *lnPtr
	App.usernameonly = *unoPtr
}

type decoy string

func main() {
	readCommandLineArgs()

	var err error
	// s := fmt.Sprintf("%s:@/%s?charset=utf8&parseTime=True", App.DBUser, App.DBName)
	// App.db, err = sql.Open("mysql", s)
	lib.ReadConfig()
	s := extres.GetSQLOpenString(App.DBName, &lib.AppConfig)
	App.db, err = sql.Open("mysql", s)
	if nil != err {
		fmt.Printf("sql.Open for database=%s, dbuser=%s: Error = %v\n", App.DBName, App.DBUser, err)
	}
	defer App.db.Close()
	err = App.db.Ping()
	if nil != err {
		fmt.Printf("App.db.Ping for database=%s, dbuser=%s: Error = %v\n", App.DBName, App.DBUser, err)
	}
	db.DB.DirDB = App.db
	db.Init()

	fmt.Printf("App.usernameonly = %t\n", App.usernameonly)

	if len(App.user) > 0 && len(App.lname) > 0 {
		getUserName()
		if App.usernameonly {
			fmt.Printf("%s\n", App.user)
			os.Exit(0)
		}
	}

	if len(App.user) > 0 && App.usernameonly {
		var a db.People
		k := decoy("key")
		ctx := context.WithValue(context.Background(), k, "Go")

		if a, err = db.GetByUsername(ctx, App.user); err != nil {
			fmt.Printf("Error getting user info: %s\n", err.Error())
			os.Exit(1)
		}
		fmt.Printf("(%d) %s %s\n", a.UID, a.FirstName, a.LastName)
		fmt.Printf("Address: %s, %s, %s, %s %s\n", a.HomeStreetAddress, a.HomeCity, a.HomeState, a.HomePostalCode, a.HomeCountry)
		fmt.Printf("Map:  %s\n", db.MapURL(a.HomeStreetAddress, a.HomeCity, a.HomeState, a.HomePostalCode, a.HomeCountry))
		os.Exit(0)
	}

	if len(App.password) == 0 {
		App.password = lib.RandPasswordStringRunes(8)
	}

	err = lib.UpdateUserPassword(App.user, App.password, App.db)
	if nil != err {
		fmt.Printf("error = %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s\nusername: %s\npassword: %s\nOK\n", getRealName(), App.user, App.password)
}
