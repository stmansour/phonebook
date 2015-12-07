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
		"select LastName,FirstName,MiddleName,Salutation,"+ // 4
			"ClassCode,Status,PositionControlNumber,"+ // 7
			"OfficePhone,OfficeFax,CellPhone,PrimaryEmail,"+ // 11
			"SecondaryEmail,EligibleForRehire,LastReview,NextReview,"+ // 15
			"BirthMonth,BirthDOM,HomeStreetAddress,HomeStreetAddress2,HomeCity,"+ // 20
			"HomeState,HomePostalCode,HomeCountry,"+ // 23
			"AcceptedHealthInsurance,AcceptedDentalInsurance,Accepted401K,"+ // 26
			"jobcode,hire,termination,"+ // 29
			"mgruid,deptcode,cocode,StateOfEmployment,"+ // 33
			"CountryOfEmployment,PreferredName,"+ // 35
			"EmergencyContactName,EmergencyContactPhone,RID "+ // 38
			"from people where uid=?", d.UID)
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		errcheck(rows.Scan(
			&d.LastName, &d.FirstName, &d.MiddleName, &d.Salutation,
			&d.ClassCode, &d.Status, &d.PositionControlNumber,
			&d.OfficePhone, &d.OfficeFax, &d.CellPhone, &d.PrimaryEmail,
			&d.SecondaryEmail, &d.EligibleForRehire, &d.LastReview, &d.NextReview,
			&d.BirthMonth, &d.BirthDOM, &d.HomeStreetAddress, &d.HomeStreetAddress2, &d.HomeCity,
			&d.HomeState, &d.HomePostalCode, &d.HomeCountry,
			&d.AcceptedHealthInsurance, &d.AcceptedDentalInsurance, &d.Accepted401K,
			&d.JobCode, &d.Hire, &d.Termination,
			&d.MgrUID, &d.DeptCode, &d.CoCode, &d.StateOfEmployment,
			&d.CountryOfEmployment, &d.PreferredName,
			&d.EmergencyContactName, &d.EmergencyContactPhone, &d.RID))
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
	w.Header().Set("Content-Type", "text/html")
	var sess *session
	var ui uiSupport
	sess = nil
	if 0 < initHandlerSession(sess, &ui, w, r) {
		return
	}
	sess = ui.X

	var d personDetail
	d.Reports = make([]person, 0)
	d.Image = "/images/anon.png"
	path := "/adminView/"
	uidstr := r.RequestURI[len(path):]
	if len(uidstr) > 0 {
		uid, _ := strconv.Atoi(uidstr)
		d.UID = uid
		breadcrumbAdd(sess, "AdminView Person", fmt.Sprintf("/adminView/%d", uid))
		adminReadDetails(&d)
	}

	//============================================================
	// SECURITY
	//============================================================
	if !sess.elemPermsAll(ELEMPERSON, PERMVIEW|PERMMOD) {
		fmt.Printf("sess.elemPermsAny(ELEMPERSON, PERMVIEW|PERMMOD) returned 0\n")
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}
	// Ensure that the user has permissions to view everything we're about
	// to display.
	d.filterSecurityRead(sess, PERMVIEW|PERMMOD)

	t, _ := template.New("adminView.html").Funcs(funcMap).ParseFiles("adminView.html")
	ui.D = &d
	//fmt.Printf("ui.D = %#v\n", ui.D)
	err := t.Execute(w, &ui)
	if nil != err {
		ulog("Error executing template: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
