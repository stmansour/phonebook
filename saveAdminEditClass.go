package main

import (
	"fmt"
	"net/http"
	"phonebook/db"
	"strconv"
	"strings"
)

func saveAdminEditClassHandler(w http.ResponseWriter, r *http.Request) {
	var ssn *db.Session
	var ui uiSupport
	ssn = nil
	if 0 < initHandlerSession(ssn, &ui, w, r) {
		return
	}
	ssn = ui.X
	Phonebook.ReqCountersMem <- 1    // ask to access the shared mem, blocks until granted
	<-Phonebook.ReqCountersMemAck    // make sure we got it
	Counters.AdminEditClass++        // initialize our data
	Phonebook.ReqCountersMemAck <- 1 // tell Dispatcher we're done with the data

	// SECURITY
	if !ssn.ElemPermsAny(db.ELEMCLASS, db.PERMMOD) {
		ulog("Permissions refuse saveAdminEditClass page on userid=%d (%s), role=%s\n", ssn.UID, ssn.Firstname, ssn.PMap.Urole.Name)
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}

	var c db.Class
	path := "/saveAdminEditClass/"
	ClassCodestr := r.RequestURI[len(path):]
	if len(ClassCodestr) == 0 {
		fmt.Fprintf(w, "The RequestURI needs the person's Company Code. It was not found on the URI:  %s\n", r.RequestURI)
		return
	}
	ClassCode, err := strconv.ParseInt(ClassCodestr, 10, 64)
	if err != nil {
		fmt.Fprintf(w, "Error converting Company Code to a number: %v. URI: %s\n", err, r.RequestURI)
		return
	}

	c.ClassCode = ClassCode
	action := strings.ToLower(r.FormValue("action"))
	if "delete" == action {
		url := fmt.Sprintf("/delClass/%d", ClassCode)
		http.Redirect(w, r, url, http.StatusFound)
		return
	}
	if "save" == action {
		c.Name = r.FormValue("Name")
		c.Designation = r.FormValue("Designation")
		c.Description = r.FormValue("Description")
		if len(c.Designation) > 3 {
			c.Designation = c.Designation[0:3]
		}
		c.CoCode, err = strconv.ParseInt(r.FormValue("CoCode"), 10, 64)
		if err != nil {
			c.CoCode = 0
		}

		//-------------------------------
		// SECURITY
		//-------------------------------
		var co db.Class              // container for current information
		co.ClassCode = ClassCode     // initialize
		getClassInfo(ClassCode, &co) // get the rest of the info

		// func (c *db.Class) filterSecurityMerge(ssn *sess.Session, permRequired int, cNew *db.Class) {
		// 	filterSecurityMerge(c, ssn, db.ELEMCLASS, permRequired, cNew, 0)
		// }
		filterSecurityMerge(&co, ssn, db.ELEMCLASS, db.PERMMOD, &c, 0) // merge new info based on permissions

		if 0 == ClassCode {
			_, err = Phonebook.prepstmt.insertClass.Exec(co.CoCode, co.Name, co.Designation, co.Description, ssn.UID)
			errcheck(err)

			// read this record back to get the ClassCode...
			rows, err := Phonebook.prepstmt.classReadBack.Query(co.Name, co.Designation)
			errcheck(err)
			defer rows.Close()
			nClassCode := int64(0) // quick way to handle multiple matches... in this case, largest ClassCode wins, it hast to be the latest db.Class added
			for rows.Next() {
				errcheck(rows.Scan(&ClassCode))
				if ClassCode > nClassCode {
					nClassCode = ClassCode
				}
			}
			errcheck(rows.Err())
			ClassCode = nClassCode
			c.ClassCode = ClassCode
			loadClasses() // This is a new db.Class, we've saved it, now we need to reload our company list...
		} else {
			_, err = Phonebook.prepstmt.updateClass.Exec(co.CoCode, co.Name, co.Designation, co.Description, ssn.UID, ClassCode)
			if nil != err {
				errmsg := fmt.Sprintf("saveAdminEditClassHandler: Phonebook.prepstmt.adminUpdatePerson.Exec: err = %v\n", err)
				ulog(errmsg)
				fmt.Println(errmsg)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
	http.Redirect(w, r, breadcrumbBack(ssn, 2), http.StatusFound)
}
