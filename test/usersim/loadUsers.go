package main

import (
	"fmt"
	"path/filepath"
)

func getCompanyInfo(cocode int, c *company) {
	s := fmt.Sprintf("select cocode,LegalName,CommonName,Address,Address2,City,State,PostalCode,Country,Phone,Fax,Email,Designation,Active,EmploysPersonnel from companies where cocode=%d", cocode)
	rows, err := App.db.Query(s)
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		errcheck(rows.Scan(&c.CoCode, &c.LegalName, &c.CommonName, &c.Address, &c.Address2, &c.City, &c.State, &c.PostalCode, &c.Country, &c.Phone, &c.Fax, &c.Email, &c.Designation, &c.Active, &c.EmploysPersonnel))
	}
	errcheck(rows.Err())
}

func getJobTitle(JobCode int) string {
	if JobCode > 0 {
		var JobTitle string
		rows, err := App.db.Query("select title from jobtitles where jobcode=?", JobCode)
		errcheck(err)
		defer rows.Close()
		for rows.Next() {
			errcheck(rows.Scan(&JobTitle))
			return JobTitle
		}
		errcheck(rows.Err())
	}
	return "Unknown"
}

func getNameFromUID(uid int) string {
	var FirstName string
	var LastName string
	var name string
	rows, err := App.db.Query("select firstname,lastname from people where uid=?", uid)
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		errcheck(rows.Scan(&FirstName, &LastName))
		name = fmt.Sprintf("%s %s", FirstName, LastName)
	}
	errcheck(rows.Err())
	return name
}

func getDepartmentFromDeptCode(deptcode int) string {
	var name string
	rows, err := App.db.Query("select name from departments where deptcode=?", deptcode)
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		errcheck(rows.Scan(&name))
	}
	errcheck(rows.Err())
	return name
}

func getReports(uid int, d *personDetail) {
	s := fmt.Sprintf("select uid,lastname,firstname,jobcode,primaryemail,officephone,cellphone from people where mgruid=%d AND status>0 order by lastname, firstname", uid)
	rows, err := App.db.Query(s)
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		var m person
		errcheck(rows.Scan(&m.UID, &m.LastName, &m.FirstName, &m.JobCode, &m.PrimaryEmail, &m.OfficePhone, &m.CellPhone))
		d.Reports = append(d.Reports, m)
	}
	errcheck(rows.Err())
}

func getImageFilename(uid int) string {
	pat := fmt.Sprintf("pictures/%d.*", uid)
	matches, err := filepath.Glob(pat)
	if err != nil {
		fmt.Printf("filepath.Glob(%s) returned error: %v\n", pat, err)
		return "/images/anon.png"
	}
	if len(matches) > 0 {
		return "/" + matches[0]
	}
	return "/images/anon.png"
}

//===========================================================
//  getPersonDetail can be used by external callers
//  to get detailed person information on a particular user
//  returns 0 if success, err number otherwise
//===========================================================
func getPersonDetail(d *personDetail, uid int) int {
	d.Image = getImageFilename(uid)
	err := App.db.QueryRow("select lastname,firstname,preferredname,jobcode,primaryemail,"+
		"officephone,cellphone,deptcode,cocode,mgruid,ClassCode,"+
		"HomeStreetAddress,HomeStreetAddress2,HomeCity,HomeState,HomePostalCode,HomeCountry "+
		"from people where uid=?", uid).Scan(&d.LastName, &d.FirstName, &d.PreferredName, &d.JobCode, &d.PrimaryEmail,
		&d.OfficePhone, &d.CellPhone, &d.DeptCode, &d.CoCode, &d.MgrUID, &d.ClassCode,
		&d.HomeStreetAddress, &d.HomeStreetAddress2, &d.HomeCity,
		&d.HomeState, &d.HomePostalCode, &d.HomeCountry)
	if nil != err {
		return 1
	}
	d.MgrName = getNameFromUID(d.MgrUID)
	d.DeptName = getDepartmentFromDeptCode(d.DeptCode)
	d.JobTitle = getJobTitle(d.JobCode)
	getCompanyInfo(d.CoCode, &d.Company)
	return 0
}

func getCompensations(d *personDetail) {
	rows, err := App.db.Query("select type from compensation where uid=?", d.UID)
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
	rows, err := App.db.Query("select deduction from deductions where uid=?", d.UID)
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
	rows, err := App.db.Query("select dcode,name from deductionlist")
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
	rows, err := App.db.Query(
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
			"EmergencyContactName,EmergencyContactPhone,RID,UserName "+ // 39
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
