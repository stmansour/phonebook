package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func saveAdminEditClassHandler(w http.ResponseWriter, r *http.Request) {
	var sess *session
	var ui uiSupport
	sess = nil
	if 0 < initHandlerSession(sess, &ui, w, r) {
		return
	}
	sess = ui.X

	// SECURITY
	if !sess.elemPermsAny(ELEMCLASS, PERMMOD) {
		ulog("Permissions refuse saveAdminEditClass page on userid=%d (%s), role=%s\n", sess.UID, sess.Firstname, sess.Urole.Name)
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}

	var c class
	path := "/saveAdminEditClass/"
	ClassCodestr := r.RequestURI[len(path):]
	if len(ClassCodestr) == 0 {
		fmt.Fprintf(w, "The RequestURI needs the person's Company Code. It was not found on the URI:  %s\n", r.RequestURI)
		return
	}
	ClassCode, err := strconv.Atoi(ClassCodestr)
	if err != nil {
		fmt.Fprintf(w, "Error converting Company Code to a number: %v. URI: %s\n", err, r.RequestURI)
		return
	}

	c.ClassCode = ClassCode

	c.Name = r.FormValue("Name")
	c.Designation = r.FormValue("Designation")
	c.Description = r.FormValue("Description")

	if 0 == ClassCode {
		insert, err := Phonebook.db.Prepare("INSERT INTO classes (Name,Designation,Description) VALUES(?,?,?)")
		errcheck(err)
		_, err = insert.Exec(c.Name, c.Designation, c.Description)
		errcheck(err)

		// read this record back to get the ClassCode...
		rows, err := Phonebook.db.Query("select ClassCode from classes where Name=? and Designation=?", c.Name, c.Designation)
		errcheck(err)
		defer rows.Close()
		nClassCode := 0 // quick way to handle multiple matches... in this case, largest ClassCode wins, it hast to be the latest class added
		for rows.Next() {
			errcheck(rows.Scan(&ClassCode))
			if ClassCode > nClassCode {
				nClassCode = ClassCode
			}
		}
		errcheck(rows.Err())
		ClassCode = nClassCode
		c.ClassCode = ClassCode
	} else {
		update, err := Phonebook.db.Prepare("update classes set Name=?,Designation=?,Description=? where ClassCode=?")
		errcheck(err)
		_, err = update.Exec(c.Name, c.Designation, c.Description, ClassCode)
		if nil != err {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	http.Redirect(w, r, fmt.Sprintf("/class/%d", ClassCode), http.StatusFound)
}
