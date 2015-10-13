package main

import (
	"net/http"
	"strconv"
	"text/template"
)

func adminReadDetails(uid int, d *personDetail) {
	d.UID = uid
	//-----------------------------------------------------------
	// query for all the fields in table People
	//-----------------------------------------------------------
	rows, err := Phonebook.db.Query(
		"select LastName,FirstName,MiddleName,Salutation,"+
			"CostCenter,Status,Department,PositionControlNumber,"+
			"OfficePhone,OfficeFax,CellPhone,PrimaryEmail,"+
			"SecondaryEmail,EligibleForRehire,LastReview,NextReview,"+
			"Birthdate,HomeStreetAddress,HomeStreetAddress2,HomeCity,"+
			"HomeState,HomePostalCode,HomeCountry,CompensationType,"+
			"jobcode,"+
			"mgruid,deptcode,cocode,StateOfEmployment,"+
			"CountryOfEmployment,PreferredName,"+
			"EmergencyContactName,EmergencyContactPhone "+
			"from people where uid=?",
		uid)
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		errcheck(rows.Scan(
			&d.LastName, &d.FirstName, &d.MiddleName, &d.Salutation,
			&d.CostCenter, &d.Status, &d.Department, &d.PositionControlNumber,
			&d.OfficePhone, &d.OfficeFax, &d.CellPhone, &d.PrimaryEmail,
			&d.SecondaryEmail, &d.EligibleForRehire, &d.LastReview, &d.NextReview,
			&d.Birthdate, &d.HomeStreetAddress, &d.HomeStreetAddress2, &d.HomeCity,
			&d.HomeState, &d.HomePostalCode, &d.HomeCountry, &d.CompensationType,
			//&d.HealthInsuranceAccepted, &d.DentalInsuranceAccepted, &d.Accepted401K,
			&d.JobCode, /*&d.hire, &d.termination,*/
			&d.MgrUID, &d.DeptCode, &d.CoCode, &d.StateOfEmployment,
			&d.CountryOfEmployment, &d.PreferredName,
			&d.EmergencyContactName, &d.EmergencyContactPhone))
	}
	errcheck(rows.Err())
	d.MgrName = getNameFromUID(d.MgrUID)
	d.Department = getDepartmentFromDeptCode(d.DeptCode)
	d.JobTitle = getJobTitle(d.JobCode)
	getCompanyInfo(d.CoCode, &d.Company)
	getReports(uid, d)
}

func adminViewHandler(w http.ResponseWriter, r *http.Request) {
	var d personDetail
	d.Reports = make([]person, 0)
	path := "/adminView/"
	uidstr := r.RequestURI[len(path):]
	if len(uidstr) > 0 {
		uid, _ := strconv.Atoi(uidstr)
		adminReadDetails(uid, &d)
	}
	t, _ := template.New("adminView.html").ParseFiles("adminView.html")
	err := t.Execute(w, &d)
	if nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
