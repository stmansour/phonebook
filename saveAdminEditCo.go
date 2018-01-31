package main

import (
	"fmt"
	"net/http"
	"phonebook/authz"
	"phonebook/db"
	"phonebook/sess"
	"strconv"
	"strings"
)

// LEGALNAMESIZE is the size of the varchar for the company legal name
var LEGALNAMESIZE = 50

// COMMONNAMESIZE is the sql size
var COMMONNAMESIZE = 50

func saveAdminEditCoHandler(w http.ResponseWriter, r *http.Request) {
	var ssn *sess.Session
	var ui uiSupport
	ssn = nil
	if 0 < initHandlerSession(ssn, &ui, w, r) {
		return
	}
	ssn = ui.X
	Phonebook.ReqCountersMem <- 1    // ask to access the shared mem, blocks until granted
	<-Phonebook.ReqCountersMemAck    // make sure we got it
	Counters.AdminEditCompany++      // initialize our data
	Phonebook.ReqCountersMemAck <- 1 // tell Dispatcher we're done with the data

	// SECURITY
	if !ssn.ElemPermsAny(authz.ELEMCOMPANY, authz.PERMMOD) {
		ulog("Permissions refuse saveAdminEditCo page on userid=%d (%s), role=%s\n", ssn.UID, ssn.Firstname, ssn.PMap.Urole.Name)
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}

	var c db.Company
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

		if len(c.LegalName) > LEGALNAMESIZE {
			c.LegalName = c.LegalName[0:LEGALNAMESIZE]
		}
		if len(c.CommonName) > COMMONNAMESIZE {
			c.CommonName = c.CommonName[0:COMMONNAMESIZE]
		}

		//-------------------------------
		// SECURITY
		//-------------------------------
		var co db.Company           // container for current information
		co.CoCode = CoCode          // initialize
		getCompanyInfo(CoCode, &co) // fetch all its data

		// func (c *company) filterSecurityMerge(ssn *sess.Session, permRequired int, cNew *company) {
		// 	filterSecurityMerge(c, ssn, authz.ELEMCOMPANY, permRequired, cNew, 0)
		// }
		filterSecurityMerge(&co, ssn, authz.ELEMCOMPANY, authz.PERMMOD, &c, 0) // merge

		if 0 == CoCode {
			_, err = Phonebook.prepstmt.insertCompany.Exec(c.LegalName, c.CommonName, c.Designation,
				c.Email, c.Phone, c.Fax, c.Active, c.EmploysPersonnel,
				c.Address, c.Address2, c.City, c.State, c.PostalCode, c.Country, ssn.UID)
			errcheck(err)

			// read this record back to get the CoCode...
			rows, err := Phonebook.prepstmt.companyReadback.Query(c.CommonName, c.LegalName)
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
		} else {
			_, err = Phonebook.prepstmt.updateCompany.Exec(c.LegalName, c.CommonName, c.Designation, c.Email, c.Phone,
				c.Fax, c.EmploysPersonnel, c.Active, c.Address, c.Address2, c.City, c.State,
				c.PostalCode, c.Country, ssn.UID, CoCode)
			if nil != err {
				errmsg := fmt.Sprintf("saveAdminEditCoHandler: Phonebook.prepstmt.updateCompany.Exec: err = %v\n", err)
				ulog(errmsg)
				fmt.Println(errmsg)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

	}
	loadCompanies() // It may be a new company, or its active/inactive status may have changed.
	http.Redirect(w, r, breadcrumbBack(ssn, 2), http.StatusFound)
}
