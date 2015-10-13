package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func strToInt(s string) int {
	if len(s) == 0 {
		return 0
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		fmt.Printf("Error converting %s to a number: %v\n", s, err)
		return 0
	}
	return n
}

func saveAdminEditHandler(w http.ResponseWriter, r *http.Request) {
	var d personDetail
	path := "/saveAdminEdit/"
	uidstr := r.RequestURI[len(path):]
	if len(uidstr) == 0 {
		fmt.Fprintf(w, "The RequestURI needs the person's uid. It was not found on the URI:  %s\n", r.RequestURI)
		return
	}
	uid, err := strconv.Atoi(uidstr)
	if err != nil {
		fmt.Fprintf(w, "Error converting uid to a number: %v. URI: %s\n", err, r.RequestURI)
		return
	}
	d.UID = uid
	d.Salutation = r.FormValue("Salutation")
	d.FirstName = r.FormValue("FirstName")
	d.MiddleName = r.FormValue("MiddleName")
	d.LastName = r.FormValue("LastName")
	d.PrimaryEmail = r.FormValue("PrimaryEmail")
	d.OfficePhone = r.FormValue("OfficePhone")
	d.CellPhone = r.FormValue("CellPhone")
	d.PreferredName = r.FormValue("PreferredName")
	d.EmergencyContactName = r.FormValue("EmergencyContactName")
	d.EmergencyContactPhone = r.FormValue("EmergencyContactPhone")
	d.CoCode = strToInt(r.FormValue("CoCode"))
	d.JobCode = strToInt(r.FormValue("JobCode"))
	d.DeptCode = strToInt(r.FormValue("DeptCode"))
	d.PositionControlNumber = r.FormValue("PositionControlNumber")
	d.HomeStreetAddress = r.FormValue("HomeStreetAddress")
	d.HomeStreetAddress2 = r.FormValue("HomeStreetAddress2")
	d.HomeCity = r.FormValue("HomeCity")
	d.HomeState = r.FormValue("HomeState")
	d.HomePostalCode = r.FormValue("HomePostalCode")
	d.HomeCountry = r.FormValue("HomeCountry")
	d.PrimaryEmail = r.FormValue("PrimaryEmail")
	d.SecondaryEmail = r.FormValue("SecondaryEmail")
	d.OfficePhone = r.FormValue("OfficePhone")
	d.OfficeFax = r.FormValue("OfficeFax")
	d.CellPhone = r.FormValue("CellPhone")
	d.Department = r.FormValue("Department")
	d.Status = strToInt(r.FormValue("Status"))
	d.EligibleForRehire = r.FormValue("EligibleForRehire")
	d.LastReview = r.FormValue("LastReview")
	d.NextReview = r.FormValue("NextReview")
	d.Birthdate = r.FormValue("Birthdate")
	d.CompensationType = r.FormValue("CompensationType")
	d.MgrUID = strToInt(r.FormValue("MgrUID"))

	fmt.Printf("d = %+v\n", d)
	update, err := Phonebook.db.Prepare("update people set Salutation=?,FirstName=?,MiddleName=?,LastName=?,PreferredName=?,EmergencyContactName=?,EmergencyContactPhone=?,PrimaryEmail=?,SecondaryEmail=?,OfficePhone=?,OfficeFax=?,CellPhone=?,CoCode=?,JobCode=?,PositionControlNumber=?,DeptCode=?,HomeStreetAddress=?,HomeStreetAddress2=?,HomeCity=?,HomeState=?,HomePostalCode=?,HomeCountry=?,status=? where people.uid=?")
	errcheck(err)
	_, err = update.Exec(
		d.Salutation, d.FirstName, d.MiddleName, d.LastName, d.PreferredName,
		d.EmergencyContactName, d.EmergencyContactPhone,
		d.PrimaryEmail, d.SecondaryEmail, d.OfficePhone, d.OfficeFax, d.CellPhone, d.CoCode, d.JobCode,
		d.PositionControlNumber, d.DeptCode,
		d.HomeStreetAddress, d.HomeStreetAddress2, d.HomeCity, d.HomeState, d.HomePostalCode, d.HomeCountry,
		d.Status,
		uid)

	if nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, fmt.Sprintf("/adminView/%d", uid), http.StatusFound)
}
