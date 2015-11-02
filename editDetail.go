package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

func editDetailHandler(w http.ResponseWriter, r *http.Request) {
	var d personDetail
	d.Reports = make([]person, 0)
	uidstr := r.RequestURI[12:]
	if len(uidstr) == 0 {
		fmt.Fprintf(w, "the RequestURI needs to know the person's uid. It was not found on the URI:  %s\n", r.RequestURI)
		return
	}
	uid, err := strconv.Atoi(uidstr)
	if err != nil {
		fmt.Fprintf(w, "Error converting uid to a number: %v. URI: %s\n", err, r.RequestURI)
		return
	}

	d.UID = uid
	rows, err := Phonebook.db.Query("select lastname,firstname,preferredname,jobcode,primaryemail,officephone,cellphone,deptcode,cocode,mgruid,Class,EmergencyContactName,EmergencyContactPhone from people where uid=?", uid)
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		errcheck(rows.Scan(&d.LastName, &d.FirstName, &d.PreferredName, &d.JobCode, &d.PrimaryEmail,
			&d.OfficePhone, &d.CellPhone, &d.DeptCode, &d.CoCode, &d.MgrUID,
			&d.Class, &d.EmergencyContactName, &d.EmergencyContactPhone))
	}
	errcheck(rows.Err())
	d.MgrName = getNameFromUID(d.MgrUID)
	d.DeptName = getDepartmentFromDeptCode(d.DeptCode)
	d.JobTitle = getJobTitle(d.JobCode)
	getCompanyInfo(d.CoCode, &d.Company)
	getReports(uid, &d)
	t, _ := template.New("editDetail.html").ParseFiles("editDetail.html")
	err = t.Execute(w, &d)
	if nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
