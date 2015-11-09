package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

func getJobTitle(JobCode int) string {
	if JobCode > 0 {
		var JobTitle string
		rows, err := Phonebook.db.Query("select title from jobtitles where jobcode=?", JobCode)
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
	rows, err := Phonebook.db.Query("select firstname,lastname from people where uid=?", uid)
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
	rows, err := Phonebook.db.Query("select name from departments where deptcode=?", deptcode)
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
	rows, err := Phonebook.db.Query(s)
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		var m person
		errcheck(rows.Scan(&m.UID, &m.LastName, &m.FirstName, &m.JobCode, &m.PrimaryEmail, &m.OfficePhone, &m.CellPhone))
		d.Reports = append(d.Reports, m)
	}
	errcheck(rows.Err())
}

func detailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var sess *session
	var ui uiSupport
	sess = nil
	if 0 < initHandlerSession(sess, &ui, w, r) {
		return
	}

	var d personDetail
	d.Reports = make([]person, 0)
	d.Image = "/images/anon.png"
	path := "/detail/"
	uidstr := r.RequestURI[len(path):]
	if len(uidstr) > 0 {
		uid, _ := strconv.Atoi(uidstr)
		d.UID = uid
		rows, err := Phonebook.db.Query("select lastname,firstname,preferredname,jobcode,primaryemail,"+
			"officephone,cellphone,deptcode,cocode,mgruid,ClassCode,"+
			"HomeStreetAddress,HomeStreetAddress2,HomeCity,HomeState,HomePostalCode,HomeCountry "+
			"from people where uid=?", uid)
		errcheck(err)
		defer rows.Close()
		for rows.Next() {
			errcheck(rows.Scan(&d.LastName, &d.FirstName, &d.PreferredName, &d.JobCode, &d.PrimaryEmail,
				&d.OfficePhone, &d.CellPhone, &d.DeptCode, &d.CoCode, &d.MgrUID, &d.ClassCode,
				&d.HomeStreetAddress, &d.HomeStreetAddress2, &d.HomeCity,
				&d.HomeState, &d.HomePostalCode, &d.HomeCountry))
		}
		errcheck(rows.Err())
		d.MgrName = getNameFromUID(d.MgrUID)
		d.DeptName = getDepartmentFromDeptCode(d.DeptCode)
		d.JobTitle = getJobTitle(d.JobCode)
		getCompanyInfo(d.CoCode, &d.Company)
		getReports(uid, &d)
		d.Class = ui.ClassCodeToName[d.ClassCode]
	}
	t, _ := template.New("detail.html").ParseFiles("detail.html")
	ui.D = &d
	err := t.Execute(w, &ui)
	if nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
