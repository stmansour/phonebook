package main

import (
	"fmt"
	"net/http"
	"phonebook/authz"
	"phonebook/sess"
	"strconv"
	"text/template"
)

func getCompensations(d *personDetail) {
	rows, err := Phonebook.prepstmt.getComps.Query(d.UID)
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
	rows, err := Phonebook.prepstmt.deductList.Query(d.UID)
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
	rows, err := Phonebook.prepstmt.myDeductions.Query()
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
	rows, err := Phonebook.prepstmt.adminPersonDetails.Query(d.UID)
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
			&d.EmergencyContactName, &d.EmergencyContactPhone, &d.RID, &d.UserName))
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
	var ssn *sess.Session
	var ui uiSupport
	ssn = nil
	if 0 < initHandlerSession(ssn, &ui, w, r) {
		return
	}
	ssn = ui.X
	Phonebook.ReqCountersMem <- 1    // ask to access the shared mem, blocks until granted
	<-Phonebook.ReqCountersMemAck    // make sure we got it
	Counters.ViewPerson++            // initialize our data
	Phonebook.ReqCountersMemAck <- 1 // tell Dispatcher we're done with the data

	var d personDetail
	d.Reports = make([]person, 0)
	d.Image = "/images/anon.png"
	path := "/adminView/"
	uidstr := r.RequestURI[len(path):]
	if len(uidstr) > 0 {
		uid, _ := strconv.Atoi(uidstr)
		d.UID = uid
		breadcrumbAdd(ssn, "AdminView Person", fmt.Sprintf("/adminView/%d", uid))
		adminReadDetails(&d)
	}

	//============================================================
	// SECURITY
	//============================================================
	if !ssn.ElemPermsAll(authz.ELEMPERSON, authz.PERMVIEW|authz.PERMMOD) {
		fmt.Printf("ssn.ElemPermsAny(authz.ELEMPERSON, authz.PERMVIEW|authz.PERMMOD) returned 0\n")
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}
	// Ensure that the user has permissions to view everything we're about
	// to display.
	d.filterSecurityRead(ssn, authz.PERMVIEW|authz.PERMMOD)

	t, _ := template.New("adminView.html").Funcs(funcMap).ParseFiles("adminView.html")
	ui.D = &d
	// fmt.Printf("ui.D = %#v\n", ui.D)
	err := t.Execute(w, &ui)
	if nil != err {
		errmsg := fmt.Sprintf("adminViewHandler: err = %v\n", err)
		ulog(errmsg)
		fmt.Println(errmsg)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
