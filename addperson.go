package main

import (
	"fmt"
	"net/http"
	"phonebook/db"
)

func adminAddPersonHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var ssn *db.Session
	var ui uiSupport
	ssn = nil
	if 0 < initHandlerSession(ssn, &ui, w, r) {
		return
	}
	ssn = ui.X

	// SECURITY
	if !ssn.ElemPermsAny(db.ELEMPERSON, db.PERMCREATE) {
		ulog("Permissions refuse AddPerson page on userid=%d (%s), role=%s\n", ssn.UID, ssn.Firstname, ssn.PMap.Urole.Name)
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}
	breadcrumbAdd(ssn, "Add Person", "/adminAddPerson/")

	var d db.PersonDetail
	d.Reports = make([]db.Person, 0)
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
	d.CoCode = ssn.CoCode // if this is done by an HR, the person should default to the same company
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
	d.RID = 4 // Viewer
	ui.D = &d

	err := renderTemplate(w, ui, "adminEdit.html")
	if nil != err {
		errmsg := fmt.Sprintf("adminAddPersonHandler:  err = %v\n", err)
		ulog(errmsg)
		fmt.Println(errmsg)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
