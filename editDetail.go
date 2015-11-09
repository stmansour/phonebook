package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

func editDetailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var sess *session
	var ui uiSupport
	sess = nil
	if 0 < initHandlerSession(sess, &ui, w, r) {
		return
	}

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
	rows, err := Phonebook.db.Query("select lastname,firstname,preferredname,jobcode,primaryemail,"+
		"officephone,cellphone,deptcode,cocode,mgruid,ClassCode,EmergencyContactName,EmergencyContactPhone,"+
		"HomeStreetAddress,HomeStreetAddress2,HomeCity,"+
		"HomeState,HomePostalCode,HomeCountry "+
		"from people where uid=?", uid)
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		errcheck(rows.Scan(&d.LastName, &d.FirstName, &d.PreferredName, &d.JobCode, &d.PrimaryEmail,
			&d.OfficePhone, &d.CellPhone, &d.DeptCode, &d.CoCode, &d.MgrUID,
			&d.ClassCode, &d.EmergencyContactName, &d.EmergencyContactPhone,
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
	ui.D = &d

	t, _ := template.New("editDetail.html").ParseFiles("editDetail.html")
	err = t.Execute(w, &ui)
	if nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
