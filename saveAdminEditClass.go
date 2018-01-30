package main

import (
	"fmt"
	"net/http"
	"phonebook/authz"
	"phonebook/sess"
	"strconv"
	"strings"
)

func saveAdminEditClassHandler(w http.ResponseWriter, r *http.Request) {
	var ssn *sess.Session
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
	if !ssn.ElemPermsAny(authz.ELEMCLASS, authz.PERMMOD) {
		ulog("Permissions refuse saveAdminEditClass page on userid=%d (%s), role=%s\n", ssn.UID, ssn.Firstname, ssn.Urole.Name)
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
		c.CoCode, err = strconv.Atoi(r.FormValue("CoCode"))
		if err != nil {
			c.CoCode = 0
		}

		//-------------------------------
		// SECURITY
		//-------------------------------
		var co class                                   // container for current information
		co.ClassCode = ClassCode                       // initialize
		getClassInfo(ClassCode, &co)                   // get the rest of the info
		co.filterSecurityMerge(ssn, authz.PERMMOD, &c) // merge new info based on permissions

		if 0 == ClassCode {
			_, err = Phonebook.prepstmt.insertClass.Exec(co.CoCode, co.Name, co.Designation, co.Description, ssn.UID)
			errcheck(err)

			// read this record back to get the ClassCode...
			rows, err := Phonebook.prepstmt.classReadBack.Query(co.Name, co.Designation)
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
			loadClasses() // This is a new class, we've saved it, now we need to reload our company list...
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
