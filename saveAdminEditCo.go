package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func saveAdminEditCoHandler(w http.ResponseWriter, r *http.Request) {
	var sess *session
	var ui uiSupport
	sess = nil
	if 0 < initHandlerSession(sess, &ui, w, r) {
		return
	}
	sess = ui.X
	Phonebook.ReqCountersMem <- 1    // ask to access the shared mem, blocks until granted
	<-Phonebook.ReqCountersMemAck    // make sure we got it
	Counters.AdminEditCompany++      // initialize our data
	Phonebook.ReqCountersMemAck <- 1 // tell Dispatcher we're done with the data

	// SECURITY
	if !sess.elemPermsAny(ELEMCOMPANY, PERMMOD) {
		ulog("Permissions refuse saveAdminEditCo page on userid=%d (%s), role=%s\n", sess.UID, sess.Firstname, sess.Urole.Name)
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}

	var c company
	path := "/saveAdminEditCo/"
	CoCodestr := r.RequestURI[len(path):]
	if len(CoCodestr) == 0 {
		fmt.Fprintf(w, "The RequestURI needs the person's Company Code. It was not found on the URI:  %s\n", r.RequestURI)
		return
	}
	CoCode, err := strconv.Atoi(CoCodestr)
	if err != nil {
		fmt.Fprintf(w, "Error converting Company Code to a number: %v. URI: %s\n", err, r.RequestURI)
		return
	}

	c.CoCode = CoCode
	action := strings.ToLower(r.FormValue("action"))
	if "delete" == action {
		url := fmt.Sprintf("/delCompany/%d", CoCode)
		http.Redirect(w, r, url, http.StatusFound)
		return
	}

	if "save" == action {
		c.LegalName = r.FormValue("LegalName")
		c.CommonName = r.FormValue("CommonName")
		c.Designation = r.FormValue("Designation")
		c.Email = r.FormValue("Email")
		c.Phone = r.FormValue("Phone")
		c.Fax = r.FormValue("Fax")
		c.Active = activeToInt(r.FormValue("Active")) // 5
		c.EmploysPersonnel = yesnoToInt(r.FormValue("EmploysPersonnel"))
		c.Address = r.FormValue("Address")
		c.Address2 = r.FormValue("Address2") //10
		c.City = r.FormValue("City")
		c.State = r.FormValue("State")
		c.PostalCode = r.FormValue("PostalCode")
		c.Country = r.FormValue("Country")

		//-------------------------------
		// SECURITY
		//-------------------------------
		var co company                            // container for current information
		co.CoCode = CoCode                        // initialize
		getCompanyInfo(CoCode, &co)               // fetch all its data
		co.filterSecurityMerge(sess, PERMMOD, &c) // merge

		if 0 == CoCode {
			insert, err := Phonebook.db.Prepare("INSERT INTO companies (LegalName,CommonName,Designation," +
				"Email,Phone,Fax,Active,EmploysPersonnel,Address,Address2,City,State,PostalCode,Country,lastmodby) " +
				//      1                 10                  20                  30
				"VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
			errcheck(err)
			_, err = insert.Exec(c.LegalName, c.CommonName, c.Designation,
				c.Email, c.Phone, c.Fax, c.Active, c.EmploysPersonnel,
				c.Address, c.Address2, c.City, c.State, c.PostalCode, c.Country, sess.UID)
			errcheck(err)

			// read this record back to get the CoCode...
			rows, err := Phonebook.db.Query("select CoCode from companies where CommonName=? and LegalName=?", c.CommonName, c.LegalName)
			errcheck(err)
			defer rows.Close()
			nCoCode := 0 // quick way to handle multiple matches... in this case, largest CoCode wins, it hast to be the latest person added
			for rows.Next() {
				errcheck(rows.Scan(&CoCode))
				if CoCode > nCoCode {
					nCoCode = CoCode
				}
			}
			errcheck(rows.Err())
			CoCode = nCoCode
			c.CoCode = CoCode
			loadCompanies() // This is a new company, we've saved it, now we need to reload our company list...
		} else {
			update, err := Phonebook.db.Prepare("update companies set LegalName=?,CommonName=?,Designation=?,Email=?,Phone=?,Fax=?,EmploysPersonnel=?,Active=?,Address=?,Address2=?,City=?,State=?,PostalCode=?,Country=?,lastmodby=? where CoCode=?")
			errcheck(err)
			_, err = update.Exec(c.LegalName, c.CommonName, c.Designation, c.Email, c.Phone,
				c.Fax, c.EmploysPersonnel, c.Active, c.Address, c.Address2, c.City, c.State,
				c.PostalCode, c.Country, sess.UID, CoCode)
			if nil != err {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

	}
	http.Redirect(w, r, breadcrumbBack(sess, 2), http.StatusFound)
}
