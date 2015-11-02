package main

import (
	"net/http"
	"text/template"
)

func adminAddCompanyHandler(w http.ResponseWriter, r *http.Request) {
	var c company
	c.CoCode = 0
	c.Active = YES
	c.EmploysPersonnel = NO
	c.LegalName = ""
	c.CommonName = ""
	c.Address = ""
	c.Address2 = ""
	c.City = ""
	c.State = ""
	c.PostalCode = ""
	c.Country = "USA"
	c.Phone = ""
	c.Fax = ""
	c.Email = ""
	c.Designation = ""

	funcMap := template.FuncMap{
		"compToString":      compensationTypeToString,
		"acceptIntToString": acceptIntToString,
		"dateToString":      dateToString,
		"dateYear":          dateYear,
		"monthStringToInt":  monthStringToInt,
		"add":               add,
		"sub":               sub,
		"rmd":               rmd,
		"mul":               mul,
		"div":               div,
	}

	t, _ := template.New("adminEditCo.html").Funcs(funcMap).ParseFiles("adminEditCo.html")
	PhonebookUI.C = &c
	err := t.Execute(w, &PhonebookUI)
	if nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
