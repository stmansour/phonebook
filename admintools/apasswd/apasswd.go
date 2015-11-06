// apasswd = a program to set the password for a user
//
package main

import (
	"crypto/sha512"
	"database/sql"
	"fmt"
	"os"
)

import _ "github.com/go-sql-driver/mysql"

func main() {
	if len(os.Args) < 3 {
		fmt.Println("usage: apasswd username newpassword")
		os.Exit(2)
	}

	dbopenparms := "ec2-user:@/accord?charset=utf8&parseTime=True"
	db, err := sql.Open("mysql", dbopenparms)
	if nil != err {
		fmt.Printf("sql.Open: Error = %v\n", err)
	}
	defer db.Close()

	err = db.Ping()
	if nil != err {
		fmt.Printf("db.Ping: Error = %v\n", err)
	}
	fmt.Printf("MySQL database opened with \"%s\"\n", dbopenparms)

	username := os.Args[1]
	passwd := os.Args[2]
	sha := sha512.Sum512([]byte(passwd))
	passhash := fmt.Sprintf("%x", sha)
	update, err := db.Prepare("update people set passhash=? where username=?")
	if nil != err {
		fmt.Printf("error = %v\n", err)
		os.Exit(1)
	}
	_, err = update.Exec(passhash, username)
	if nil != err {
		fmt.Printf("error = %v\n", err)
	} else {
		fmt.Printf("password for user %s has been updated\n", username)
	}

}
