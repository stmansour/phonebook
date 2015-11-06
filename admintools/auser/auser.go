// auser  a program to set the username for a person in the accord database
//        based on their UID
package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
)

import _ "github.com/go-sql-driver/mysql"

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
	if len(os.Args) < 3 {
		fmt.Println("usage: auser <uid> <username>\nExample:  auser 211 sman\n")
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

	uid := strToInt(os.Args[1])
	if uid < 1 {
		fmt.Printf("invalid uid: %s\n", os.Args[1])
	}
	username := os.Args[2]
	update, err := db.Prepare("update people set username=? where uid=?")
	if nil != err {
		fmt.Printf("error = %v\n", err)
		os.Exit(1)
	}
	_, err = update.Exec(username, uid)
	if nil != err {
		fmt.Printf("error = %v\n", err)
	} else {
		fmt.Printf("password for uid %d has been updated\n", uid)
	}

}
