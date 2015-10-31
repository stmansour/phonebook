package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func getCompanyInfo(cocode int, c *company) {
	s := fmt.Sprintf("select cocode,LegalName,CommonName,Address,Address2,City,State,PostalCode,Country,Phone,Fax,Email,Designation,Active,EmploysPersonnel from companies where cocode=%d", cocode)
	rows, err := Phonebook.db.Query(s)
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		errcheck(rows.Scan(&c.CoCode, &c.LegalName, &c.CommonName, &c.Address, &c.Address2, &c.City, &c.State, &c.PostalCode, &c.Country, &c.Phone, &c.Fax, &c.Email, &c.Designation, &c.Active, &c.EmploysPersonnel))
	}
	errcheck(rows.Err())
}

func companyHandler(w http.ResponseWriter, r *http.Request) {
	var c company
	costr := r.RequestURI[9:]
	if len(costr) > 0 {
		cocode, _ := strconv.Atoi(costr)
		getCompanyInfo(cocode, &c)
		t, _ := template.New("company.html").ParseFiles("company.html")
		err := t.Execute(w, &c)
		if nil != err {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		fmt.Fprintf(w, "cocode = %s\nCould not convert to number\n", costr)
	}
}
