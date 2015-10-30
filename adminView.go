package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

func getCompensations(d *personDetail) {
	rows, err := Phonebook.db.Query("select type from compensation where uid=?", d.UID)
	errcheck(err)
	defer rows.Close()
	var c int
	for rows.Next() {
		errcheck(rows.Scan(&c))
		d.Comps = append(d.Comps, c)
	}
	errcheck(rows.Err())
}

func getCompensationStr(d *personDetail) {
	getCompensations(d)
	d.CompensationStr = ""
	for i := 0; i < len(d.Comps); i++ {
		d.CompensationStr += fmt.Sprintf("%d", d.Comps[i])
		if i+1 < len(d.Comps) {
			d.CompensationStr += ", "
		}
	}
}

func initMyComps(d *personDetail) {
	d.MyComps = make([]myComp, 0)
	for i := CTUNSET + 1; i < CTEND; i++ {
		var c myComp
		c.CompCode = i
		c.Name = compensationTypeToString(i)
		c.HaveIt = 0
		d.MyComps = append(d.MyComps, c)
	}
}

func buildMyCompsMap(d *personDetail) {
	getCompensations(d)
	initMyComps(d)
	for i := 0; i < len(d.MyComps); i++ {
		for j := 0; j < len(d.Comps); j++ {
			if d.Comps[j] == d.MyComps[i].CompCode {
				d.MyComps[i].HaveIt = 1
			}
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
	//getDeductions(d)
	d.DeductionsStr = ""
	for i := 0; i < len(d.Deductions); i++ {
		d.DeductionsStr += fmt.Sprintf("%d", d.Deductions[i])
		if i+1 < len(d.Deductions) {
			d.DeductionsStr += ", "
		}
	}
}

func initMyDeductions(d *personDetail) {
	rows, err := Phonebook.db.Query("select dcode,name from DeductionList")
	errcheck(err)
	defer rows.Close()
	d.MyDeductions = make([]aDeduction, 0)
	for rows.Next() {
		var b aDeduction
		errcheck(rows.Scan(&b.DCode, &b.Name))
		if b.DCode == DDUNKNOWN {
			continue
		}
		d.MyDeductions = append(d.MyDeductions, b)
	}
	errcheck(rows.Err())
}

func loadDeductionList(d *personDetail) {
	getDeductions(d)
	initMyDeductions(d)
	for i := 0; i < len(d.MyDeductions); i++ {
		for j := 0; j < len(d.Deductions); j++ {
			if d.Deductions[j] == d.MyDeductions[i].DCode {
				d.MyDeductions[i].HaveIt = 1
			}
		}
	}
}

func adminReadDetails(d *personDetail) {

	//-----------------------------------------------------------
	// query for all the fields in table People
	//-----------------------------------------------------------
	rows, err := Phonebook.db.Query(
		"select LastName,FirstName,MiddleName,Salutation,"+
			"Class,Status,PositionControlNumber,"+
			"OfficePhone,OfficeFax,CellPhone,PrimaryEmail,"+
			"SecondaryEmail,EligibleForRehire,LastReview,NextReview,"+
			"Birthdate,BirthMonth,BirthDOM,HomeStreetAddress,HomeStreetAddress2,HomeCity,"+
			"HomeState,HomePostalCode,HomeCountry,"+
			"AcceptedHealthInsurance,AcceptedDentalInsurance,Accepted401K,"+
			"jobcode,hire,termination,"+
			"mgruid,deptcode,cocode,StateOfEmployment,"+
			"CountryOfEmployment,PreferredName,"+
			"EmergencyContactName,EmergencyContactPhone,EligibleForRehire "+
			"from people where uid=?",
		d.UID)
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		errcheck(rows.Scan(
			&d.LastName, &d.FirstName, &d.MiddleName, &d.Salutation,
			&d.Class, &d.Status, &d.PositionControlNumber,
			&d.OfficePhone, &d.OfficeFax, &d.CellPhone, &d.PrimaryEmail,
			&d.SecondaryEmail, &d.EligibleForRehire, &d.LastReview, &d.NextReview,
			&d.Birthdate, &d.BirthMonth, &d.BirthDOM, &d.HomeStreetAddress, &d.HomeStreetAddress2, &d.HomeCity,
			&d.HomeState, &d.HomePostalCode, &d.HomeCountry,
			&d.AcceptedHealthInsurance, &d.AcceptedDentalInsurance, &d.Accepted401K,
			&d.JobCode, &d.Hire, &d.Termination,
			&d.MgrUID, &d.DeptCode, &d.CoCode, &d.StateOfEmployment,
			&d.CountryOfEmployment, &d.PreferredName,
			&d.EmergencyContactName, &d.EmergencyContactPhone, &d.EligibleForRehire))
	}
	errcheck(rows.Err())
	d.MgrName = getNameFromUID(d.MgrUID)
	d.DeptName = getDepartmentFromDeptCode(d.DeptCode)
	d.JobTitle = getJobTitle(d.JobCode)
	getCompanyInfo(d.CoCode, &d.Company)
	getReports(d.UID, d)
	buildMyCompsMap(d) // fills the d.MyCompsMap and d.Comps array too
	loadDeductionList(d)
	getDeductionsStr(d)
}

func adminViewHandler(w http.ResponseWriter, r *http.Request) {
	var d personDetail
	d.Reports = make([]person, 0)
	d.Image = "/images/anon.png"
	path := "/adminView/"
	uidstr := r.RequestURI[len(path):]
	if len(uidstr) > 0 {
		uid, _ := strconv.Atoi(uidstr)
		d.UID = uid
		adminReadDetails(&d)
	}
	funcMap := template.FuncMap{
		"compToString":      compensationTypeToString,
		"deductionToString": deductionIntToString,
		"acceptIntToString": acceptIntToString,
		"dateToString":      dateToString,
		"activeToString":    activeToInt,
		"yesnoToString":     yesnoToInt,
		"monthStringToInt":  monthStringToInt,
	}
	t, _ := template.New("adminView.html").Funcs(funcMap).ParseFiles("adminView.html")
	PhonebookUI.D = &d
	//fmt.Printf("PhonebookUI = %#v\n", PhonebookUI)
	err := t.Execute(w, &PhonebookUI)
	if nil != err {
		fmt.Printf("Error executing template: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
