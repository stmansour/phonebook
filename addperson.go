package main

import (
	"net/http"
	"text/template"
)

func adminAddPersonHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var sess *session
	var ui uiSupport
	sess = nil
	if 0 < initHandlerSession(sess, &ui, w, r) {
		return
	}

	var d personDetail
	d.Reports = make([]person, 0)
	initMyDeductions(&d)
	initMyComps(&d)

	d.Status = ACTIVE
	d.EligibleForRehire = YES
	d.UID = 0
	d.Salutation = ""
	d.FirstName = ""
	d.MiddleName = ""
	d.LastName = ""
	d.PrimaryEmail = ""
	d.OfficePhone = ""
	d.CellPhone = ""
	d.PreferredName = ""
	d.EmergencyContactName = ""
	d.EmergencyContactPhone = ""
	d.CoCode = 0
	d.JobCode = 0
	d.DeptCode = 0
	d.PositionControlNumber = ""
	d.HomeStreetAddress = ""
	d.HomeStreetAddress2 = ""
	d.HomeCity = ""
	d.HomeState = ""
	d.HomePostalCode = ""
	d.HomeCountry = "USA"
	d.PrimaryEmail = ""
	d.SecondaryEmail = ""
	d.OfficePhone = ""
	d.OfficeFax = ""
	d.CellPhone = ""
	d.DeptName = ""
	d.LastReview = stringToDate("")
	d.NextReview = stringToDate("")
	d.BirthDOM = 0
	d.BirthMonth = 0
	d.MgrUID = 0
	d.Accepted401K = ACPTUNKNOWN
	d.AcceptedDentalInsurance = ACPTUNKNOWN
	d.AcceptedHealthInsurance = ACPTUNKNOWN
	d.Hire = stringToDate("")
	d.Termination = stringToDate("")
	d.StateOfEmployment = ""
	d.CountryOfEmployment = "USA"

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

	t, _ := template.New("adminEdit.html").Funcs(funcMap).ParseFiles("adminEdit.html")
	ui.D = &d
	err := t.Execute(w, &ui)
	if nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
