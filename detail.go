package main

import (
	"fmt"
	"net/http"
	"phonebook/db"
	"strconv"
	"strings"
)

func getJobTitle(JobCode int64) string {
	if JobCode > 0 {
		var JobTitle string
		rows, err := Phonebook.prepstmt.getJobTitle.Query(JobCode)
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

func getNameFromUID(uid int64) string {
	var FirstName string
	var LastName string
	var name string
	rows, err := Phonebook.prepstmt.nameFromUID.Query(uid)
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		errcheck(rows.Scan(&FirstName, &LastName))
		name = fmt.Sprintf("%s %s", FirstName, LastName)
	}
	errcheck(rows.Err())
	return name
}

func getDepartmentFromDeptCode(deptcode int64) string {
	var name string
	rows, err := Phonebook.prepstmt.deptName.Query(deptcode)
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		errcheck(rows.Scan(&name))
	}
	errcheck(rows.Err())
	return name
}

func getReports(uid int64, d *db.PersonDetail) {
	//s := fmt.Sprintf("select uid,lastname,firstname,jobcode,primaryemail,officephone,cellphone from people where mgruid=%d AND status>0 order by lastname, firstname", uid)
	rows, err := Phonebook.prepstmt.directReports.Query(uid)
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		var m db.Person
		errcheck(rows.Scan(&m.UID, &m.LastName, &m.FirstName, &m.JobCode, &m.PrimaryEmail, &m.OfficePhone, &m.CellPhone))
		d.Reports = append(d.Reports, m)
	}
	errcheck(rows.Err())
}

func detailpopHandler(w http.ResponseWriter, r *http.Request) {
	var sess *db.Session
	var ui uiSupport
	sess = nil
	if 0 < initHandlerSession(sess, &ui, w, r) {
		return
	}
	sess = ui.X
	breadcrumbBack(sess, 1)
	detailHandler(w, r)
}

//===========================================================
//  getPersonDetail can be used by external callers
//  to get detailed person information on a particular user
//  returns 0 if success, err number otherwise
//===========================================================
func getPersonDetail(d *db.PersonDetail, uid int64) int {
	d.Image = db.GetImageLocation(uid)
	err := Phonebook.prepstmt.personDetail.QueryRow(uid).Scan(&d.LastName, &d.MiddleName,
		&d.FirstName, &d.PreferredName, &d.JobCode, &d.PrimaryEmail,
		&d.OfficePhone, &d.CellPhone, &d.DeptCode, &d.CoCode, &d.MgrUID, &d.ClassCode,
		&d.EmergencyContactName, &d.EmergencyContactPhone,
		&d.HomeStreetAddress, &d.HomeStreetAddress2, &d.HomeCity,
		&d.HomeState, &d.HomePostalCode, &d.HomeCountry, &d.OfficeFax)
	if nil != err {
		return 1
	}
	d.MgrName = getNameFromUID(d.MgrUID)
	d.DeptName = getDepartmentFromDeptCode(d.DeptCode)
	d.JobTitle = getJobTitle(d.JobCode)
	getCompanyInfo(d.CoCode, &d.Company)
	return 0
}

func detailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var sess *db.Session
	var uis uiSupport
	sess = nil
	if 0 < initHandlerSession(sess, &uis, w, r) {
		return
	}
	sess = uis.X
	Phonebook.ReqCountersMem <- 1    // ask to access the shared mem, blocks until granted
	<-Phonebook.ReqCountersMemAck    // make sure we got it
	Counters.ViewPerson++            // initialize our data
	Phonebook.ReqCountersMemAck <- 1 // tell Dispatcher we're done with the data

	var d db.PersonDetail
	d.Reports = make([]db.Person, 0)

	var path string
	if strings.Contains(r.RequestURI, "pop") {
		path = "/detailpop/"
	} else {
		path = "/detail/"
	}
	uidstr := r.RequestURI[len(path):]
	uid := int64(0)
	if len(uidstr) > 0 {
		uid, _ = strconv.ParseInt(uidstr, 10, 64)
		d.UID = uid
	}
	breadcrumbAdd(sess, "Person", fmt.Sprintf("/detail/%d", uid))

	//=================================================================
	// SECURITY
	//=================================================================
	if !sess.ElemPermsAny(db.ELEMPERSON, db.PERMVIEW|db.PERMOWNERVIEW) {
		ulog("ViewPersonDetail: Permission refusal on userid=%d (%s), role=%s\n", sess.UID, sess.Firstname, sess.PMap.Urole.Name)
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}

	if uid > 0 {
		d.Image = db.GetImageLocation(uid)
		rows, err := Phonebook.prepstmt.personDetail.Query(uid)
		errcheck(err)
		defer rows.Close()
		for rows.Next() {
			errcheck(rows.Scan(&d.LastName, &d.MiddleName, &d.FirstName, &d.PreferredName,
				&d.JobCode, &d.PrimaryEmail,
				&d.OfficePhone, &d.CellPhone, &d.DeptCode, &d.CoCode, &d.MgrUID,
				&d.ClassCode, &d.EmergencyContactName, &d.EmergencyContactPhone,
				&d.HomeStreetAddress, &d.HomeStreetAddress2, &d.HomeCity,
				&d.HomeState, &d.HomePostalCode, &d.HomeCountry, &d.OfficeFax))
		}
		errcheck(rows.Err())
		d.MgrName = getNameFromUID(d.MgrUID)
		d.DeptName = getDepartmentFromDeptCode(d.DeptCode)
		d.JobTitle = getJobTitle(d.JobCode)
		getCompanyInfo(d.CoCode, &d.Company)
		getReports(uid, &d)
		d.Class = uis.ClassCodeToName[d.ClassCode]
	}

	uis.D = &d

	filterSecurityRead(uis.D, db.ELEMPERSON, sess, db.PERMVIEW, d.UID)

	err := renderTemplate(w, uis, "detail.html")
	if nil != err {
		errmsg := fmt.Sprintf("detailHandler: err = %v\n", err)
		ulog(errmsg)
		fmt.Println(errmsg)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
