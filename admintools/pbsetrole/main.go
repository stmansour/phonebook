// auser  a program to set the role for user in the accord database
//        based on their UID
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
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
	rname     string
	DumpRoles bool
	RID       int
	Roles     []Role // the roles saved in the database
}

func errcheck(err error) {
	if nil != err {
		fmt.Printf("Error = %v\n", err)
		os.Exit(1)
	}
}

func readAccessRoles() {
	rows, err := App.db.Query("select RID,Name from roles")
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
	dbnmPtr := flag.String("N", "accord", "database name (accordtest, accord)")
	nPtr := flag.String("u", "", "username")
	RPtr := flag.Bool("R", false, "dump roles to stdout")
	rPtr := flag.String("r", "Viewer", "sets the user's role")
	flag.Parse()
	App.DBName = *dbnmPtr
	App.username = *nPtr
	App.rname = *rPtr
	App.DumpRoles = *RPtr
	App.DBUser = *dbuPtr
}

func main() {
	readCommandLineArgs()

	var err error
	s := fmt.Sprintf("%s:@/%s?charset=utf8&parseTime=True", App.DBUser, App.DBName)
	App.db, err = sql.Open("mysql", s)
	if nil != err {
		fmt.Printf("sql.Open for database=%s, dbuser=%s: Error = %v\n", App.DBName, App.DBUser, err)
	}
	defer App.db.Close()
	err = App.db.Ping()
	if nil != err {
		fmt.Printf("App.db.Ping for database=%s, dbuser=%s: Error = %v\n", App.DBName, App.DBUser, err)
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

	uid := 0
	err = App.db.QueryRow("select uid from people where username=?", App.username).Scan(&uid)
	if 0 == uid {
		fmt.Printf("Database %s does not have a user with username = %s\n", App.DBName, App.username)
		os.Exit(1)
	}
	update, err := App.db.Prepare("update people set RID=? where uid=?")
	if nil != err {
		fmt.Printf("error = %v\n", err)
		os.Exit(1)
	}
	_, err = update.Exec(App.RID, uid)
	if nil != err {
		fmt.Printf("error = %v\n", err)
	} else {
		fmt.Printf("password for uid %d has been updated to %s\n", uid, App.rname)
	}

}
