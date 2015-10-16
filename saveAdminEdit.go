package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func strToInt(s string) int {
	if len(s) == 0 {
		return 0
	}
	s = strings.Trim(s, " \n\r")
	n, err := strconv.Atoi(s)
	if err != nil {
		fmt.Printf("Error converting %s to a number: %v\n", s, err)
		return 0
	}
	return n
}

// This is a short term function. Should be replaced by a multi instanced dropdown selector
func parseCompensation(d *personDetail) {
	ca := strings.Split(d.CompensationStr, ",")
	d.Comps = d.Comps[:0] // clear it
	for i := 0; i < len(ca); i++ {
		d.Comps = append(d.Comps, strToInt(ca[i]))
	}
}

// This is a short term function. Should be replaced by a multi instanced dropdown selector
func parseDeductions(d *personDetail) {
	ca := strings.Split(d.DeductionsStr, ",")
	d.Deductions = d.Deductions[:0] // clear it
	for i := 0; i < len(ca); i++ {
		d.Deductions = append(d.Deductions, strToInt(ca[i]))
	}
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
	d.EligibleForRehire = strToInt(r.FormValue("EligibleForRehire"))
	d.LastReview = r.FormValue("LastReview")
	d.NextReview = r.FormValue("NextReview")
	d.Birthdate = r.FormValue("Birthdate")
	d.CompensationStr = r.FormValue("CompensationStr")
	d.DeductionsStr = r.FormValue("DeductionsStr")
	d.MgrUID = strToInt(r.FormValue("MgrUID"))
	d.Accepted401K = acceptTypeToInt(r.FormValue("Accepted401K"))
	d.AcceptedDentalInsurance = acceptTypeToInt(r.FormValue("AcceptedDentalInsurance"))
	d.AcceptedHealthInsurance = acceptTypeToInt(r.FormValue("AcceptedHealthInsurance"))
	d.Hire = stringToDate(r.FormValue("Hire"))
	d.Termination = stringToDate(r.FormValue("Termination"))

	parseCompensation(&d)
	parseDeductions(&d)

	//----------------------
	// Handle dropdowns...
	//----------------------
	if "none" == d.Salutation {
		d.Salutation = ""
	}

	update, err := Phonebook.db.Prepare("update people set Salutation=?,FirstName=?,MiddleName=?,LastName=?,PreferredName=?,EmergencyContactName=?,EmergencyContactPhone=?,PrimaryEmail=?,SecondaryEmail=?,OfficePhone=?,OfficeFax=?,CellPhone=?,CoCode=?,JobCode=?,PositionControlNumber=?,DeptCode=?,HomeStreetAddress=?,HomeStreetAddress2=?,HomeCity=?,HomeState=?,HomePostalCode=?,HomeCountry=?,status=?,EligibleForRehire=?,Accepted401K=?,AcceptedDentalInsurance=?,AcceptedHealthInsurance=?,Hire=?,Termination=? where people.uid=?")
	errcheck(err)
	_, err = update.Exec(
		d.Salutation, d.FirstName, d.MiddleName, d.LastName, d.PreferredName,
		d.EmergencyContactName, d.EmergencyContactPhone,
		d.PrimaryEmail, d.SecondaryEmail, d.OfficePhone, d.OfficeFax, d.CellPhone, d.CoCode, d.JobCode,
		d.PositionControlNumber, d.DeptCode,
		d.HomeStreetAddress, d.HomeStreetAddress2, d.HomeCity, d.HomeState, d.HomePostalCode, d.HomeCountry,
		d.Status, d.EligibleForRehire, d.Accepted401K, d.AcceptedDentalInsurance, d.AcceptedHealthInsurance,
		dateToDBStr(d.Hire), dateToDBStr(d.Termination),
		uid)

	if nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	//--------------------------------------------------------------------------
	// Remove old compensation type(s) and Insert new compensation type(s)
	//--------------------------------------------------------------------------
	ct, err := Phonebook.db.Prepare("DELETE FROM compensation WHERE uid=?")
	errcheck(err)
	_, err = ct.Exec(d.UID)
	errcheck(err)
	ct, err = Phonebook.db.Prepare("INSERT INTO compensation (uid,type) VALUES(?,?)")
	errcheck(err)
	for i := 0; i < len(d.Comps); i++ {
		_, err := ct.Exec(d.UID, d.Comps[i])
		errcheck(err)
	}

	//--------------------------------------------------------------------------
	// Update deductions...
	//--------------------------------------------------------------------------
	ct, err = Phonebook.db.Prepare("DELETE FROM deductions WHERE uid=?")
	errcheck(err)
	_, err = ct.Exec(d.UID)
	errcheck(err)
	ct, err = Phonebook.db.Prepare("INSERT INTO deductions (uid,deduction) VALUES(?,?)")
	errcheck(err)
	for i := 0; i < len(d.Deductions); i++ {
		_, err := ct.Exec(d.UID, d.Deductions[i])
		errcheck(err)
	}

	http.Redirect(w, r, fmt.Sprintf("/adminView/%d", uid), http.StatusFound)
}
