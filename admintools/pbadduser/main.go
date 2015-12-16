// add a user
//   needs firstname, lastname, username, passwork, role

package main

import (
	"crypto/sha512"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"
)

import _ "github.com/go-sql-driver/mysql"

// Role defines a collection of FieldPerms that can be assigned to a person
type Role struct {
	RID  int    // assigned by DB
	Name string // role name
}

// App is the global data structure for this app
var App struct {
	db        *sql.DB
	DBName    string
	DBUser    string
	username  string
	firstname string
	lastname  string
	passwd    string
	rname     string
	RID       int
	DumpRoles bool
	Roles     []Role // the roles saved in the database
}

func errcheck(err error) {
	if nil != err {
		fmt.Printf("Error = %v\n", err)
		os.Exit(1)
	}
}

func readAccessRoles() {
	rows, err := App.db.Query("select RID,Name from Roles")
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		var r Role
		errcheck(rows.Scan(&r.RID, &r.Name))
		App.Roles = append(App.Roles, r)
	}
	errcheck(rows.Err())
}

func readCommandLineArgs() {
	dbuPtr := flag.String("B", "ec2-user", "database user name")
	dbnmPtr := flag.String("N", "accordtest", "database name (accordtest, accord)")
	uPtr := flag.String("u", "", "username")
	pPtr := flag.String("p", "accord", "password")
	fPtr := flag.String("f", "", "first or given name")
	lPtr := flag.String("l", "", "last or surname name")
	rPtr := flag.String("r", "Viewer", "sets the user's role")
	RPtr := flag.Bool("R", false, "dump roles to stdout")
	flag.Parse()
	App.DBName = *dbnmPtr
	App.DBUser = *dbuPtr
	App.username = *uPtr
	App.firstname = *fPtr
	App.lastname = *lPtr
	App.passwd = *pPtr
	App.rname = *rPtr
	App.DumpRoles = *RPtr
}

func getUsername() {
	//============================================
	// generate a unique username...
	//============================================
	App.username = strings.ToLower(App.firstname[0:1] + App.lastname)
	if len(App.username) > 17 {
		App.username = App.username[0:17]
	}
	UserName := App.username
	var xx int
	nUID := 0
	for {
		found := false
		rows, err := App.db.Query("select uid from people where UserName=?", UserName)
		errcheck(err)
		defer rows.Close()
		for rows.Next() {
			errcheck(rows.Scan(&xx))
			nUID++
			found = true
			UserName = fmt.Sprintf("%s%d", App.username, nUID)
		}
		if !found {
			break
		}
	}
	App.username = UserName
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
	if App.DumpRoles {
		for i := 0; i < len(App.Roles); i++ {
			fmt.Printf("%d - %s\n", App.Roles[i].RID, App.Roles[i].Name)
		}
		os.Exit(0)
	}

	App.RID = 0

	for i := 0; i < len(App.Roles); i++ {
		if App.Roles[i].Name == App.rname {
			App.RID = App.Roles[i].RID
		}
	}

	if 0 == App.RID {
		fmt.Printf("Could not find role named: %s\n", App.rname)
		os.Exit(0)
	}

	getUsername()
	sha := sha512.Sum512([]byte(App.passwd))
	passhash := fmt.Sprintf("%x", sha)

	stmt, err := App.db.Prepare("INSERT INTO people (UserName,passhash,FirstName,LastName,RID) VALUES(?,?,?,?,?)")
	if nil != err {
		fmt.Printf("error = %v\n", err)
		os.Exit(1)
	}
	_, err = stmt.Exec(App.username, passhash, App.firstname, App.lastname, App.RID)
	if nil != err {
		fmt.Printf("error = %v\n", err)
	} else {
		fmt.Printf("Added user to database %s:  username: %s, access role: %s\n", App.DBName, App.username, App.rname)
	}
}
