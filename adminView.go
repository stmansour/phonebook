package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

func getComps(uid int, d *personDetail) {
	rows, err := Phonebook.db.Query("select type from compensation where uid=?", uid)
	errcheck(err)
	defer rows.Close()
	var c int
	for rows.Next() {
		errcheck(rows.Scan(&c))
		d.Comps = append(d.Comps, c)
	}
	errcheck(rows.Err())
}

func getCompensationStr(uid int, d *personDetail) {
	getComps(uid, d)
	d.CompensationStr = ""
	for i := 0; i < len(d.Comps); i++ {
		d.CompensationStr += fmt.Sprintf("%d", d.Comps[i])
		if i+1 < len(d.Comps) {
			d.CompensationStr += ", "
		}
	}
}

func getDeductions(d *personDetail) {
	rows, err := Phonebook.db.Query("select deduction from deductions where uid=?", d.UID)
	errcheck(err)
	defer rows.Close()
	var c int
	for rows.Next() {
		errcheck(rows.Scan(&c))
		d.Deductions = append(d.Deductions, c)
	}
	errcheck(rows.Err())
}

func getDeductionsStr(d *personDetail) {
	getDeductions(d)
	d.DeductionsStr = ""
	for i := 0; i < len(d.Deductions); i++ {
		d.DeductionsStr += fmt.Sprintf("%d", d.Deductions[i])
		if i+1 < len(d.Deductions) {
			d.DeductionsStr += ", "
		}
	}
}

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
			"HomeState,HomePostalCode,HomeCountry,"+
			"AcceptedHealthInsurance,AcceptedDentalInsurance,Accepted401K,"+
			"jobcode,"+
			"mgruid,deptcode,cocode,StateOfEmployment,"+
			"CountryOfEmployment,PreferredName,"+
			"EmergencyContactName,EmergencyContactPhone,EligibleForRehire "+
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
			&d.HomeState, &d.HomePostalCode, &d.HomeCountry,
			&d.AcceptedHealthInsurance, &d.AcceptedDentalInsurance, &d.Accepted401K,
			&d.JobCode, /*&d.hire, &d.termination,*/
			&d.MgrUID, &d.DeptCode, &d.CoCode, &d.StateOfEmployment,
			&d.CountryOfEmployment, &d.PreferredName,
			&d.EmergencyContactName, &d.EmergencyContactPhone, &d.EligibleForRehire))
	}
	errcheck(rows.Err())
	d.MgrName = getNameFromUID(d.MgrUID)
	d.Department = getDepartmentFromDeptCode(d.DeptCode)
	d.JobTitle = getJobTitle(d.JobCode)
	getCompanyInfo(d.CoCode, &d.Company)
	getReports(uid, d)
	getCompensationStr(uid, d) // fills the d.Comps array too
	getDeductionsStr(d)
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
	funcMap := template.FuncMap{
		"compToString":      compensationTypeToString,
		"deductionToString": deductionToString,
		"acceptIntToString": acceptIntToString,
	}

	t, _ := template.New("adminView.html").Funcs(funcMap).ParseFiles("adminView.html")
	err := t.Execute(w, &d)
	if nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
