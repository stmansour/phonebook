package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func saveAdminEditCoHandler(w http.ResponseWriter, r *http.Request) {
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

	if 0 == CoCode {
		insert, err := Phonebook.db.Prepare("INSERT INTO companies (LegalName,CommonName,Designation," +
			"Email,Phone,Fax,Active,EmploysPersonnel,Address,Address2,City,State,PostalCode,Country) " +
			//      1                 10                  20                  30
			"VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
		errcheck(err)
		_, err = insert.Exec(c.LegalName, c.CommonName, c.Designation,
			c.Email, c.Phone, c.Fax, c.Active, c.EmploysPersonnel,
			c.Address, c.Address2, c.City, c.State, c.PostalCode, c.Country)
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
	} else {
		update, err := Phonebook.db.Prepare("update companies set LegalName=?,CommonName=?,Designation=?,Email=?,Phone=?,Fax=?,EmploysPersonnel=?,Active=?,Address=?,Address2=?,City=?,State=?,PostalCode=?,Country=? where CoCode=?")
		errcheck(err)
		_, err = update.Exec(c.LegalName, c.CommonName, c.Designation, c.Email, c.Phone,
			c.Fax, c.EmploysPersonnel, c.Active, c.Address, c.Address2, c.City, c.State,
			c.PostalCode, c.Country, CoCode)
		if nil != err {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	http.Redirect(w, r, fmt.Sprintf("/company/%d", CoCode), http.StatusFound)
}
