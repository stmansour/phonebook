package main

import (
	"fmt"
	"net/http"
	"phonebook/db"
	"strconv"
	"strings"
)

func delCheckError(caller string, ssn *db.Session, err error, s string, w http.ResponseWriter, r *http.Request) bool {
	if nil != err {
		ulog("%s: \"%s\"  err = %v\n", caller, s, err)
		fmt.Printf("%s: \"%s\"  err = %v\n", caller, s, err)
		http.Redirect(w, r, breadcrumbBack(ssn, 2), http.StatusFound)
		return true
	}
	return false
}

func intPersonRefErrHandler(w http.ResponseWriter, r *http.Request, path string) {
	w.Header().Set("Content-Type", "text/html")
	var ssn *db.Session
	var ui uiSupport
	ssn = nil
	if 0 < initHandlerSession(ssn, &ui, w, r) {
		return
	}
	ssn = ui.X

	uidstr := r.RequestURI[len(path):]
	if len(uidstr) > 0 {
		uid, _ := strconv.ParseInt(uidstr, 10, 64)
		var pd db.PersonDetail
		if 0 != getPersonDetail(&pd, uid) {
			ulog("%s: Error retrieving person information for userid=%d\n", path, uid)
			http.Redirect(w, r, "/search/", http.StatusFound)
			return
		}
		ui.D = &pd
		// ui.D.filterSecurityRead(ssn, db.PERMVIEW)
		PDetFilterSecurityRead(ui.D, ssn, db.PERMVIEW)
		breadcrumbAdd(ssn, "Inactivate Person", fmt.Sprintf("/inactivatePerson/%d", uid))

		s := fmt.Sprintf("select uid,lastname,firstname,preferredname,jobcode,primaryemail,officephone,cellphone,deptcode from people where status=1 and mgruid=%d", uid)
		// fmt.Printf("QUERY = %s\n", s)
		rows, err := Phonebook.db.Query(s) // note: the single arg to Query causes the sql impl to NOT create a prepared statement
		errcheck(err)
		defer rows.Close()
		var d searchResults

		for rows.Next() {
			var m db.Person
			errcheck(rows.Scan(&m.UID, &m.LastName, &m.FirstName, &m.PreferredName, &m.JobCode, &m.PrimaryEmail, &m.OfficePhone, &m.CellPhone, &m.DeptCode))
			m.DeptName = getDepartmentFromDeptCode(m.DeptCode)
			pm := &m
			filterSecurityRead(pm, db.ELEMPERSON, ssn, db.PERMVIEW|db.PERMMOD, m.UID)
			d.Matches = append(d.Matches, m)
		}
		errcheck(rows.Err())
		ui.R = &d
		// fmt.Printf("Match Count = %d\n", len(ui.R.Matches))

		err = renderTemplate(w, ui, "delPersonRefErr.html")
		if nil != err {
			errmsg := fmt.Sprintf("intPersonRefErrHandler: err = %v\n", err)
			ulog(errmsg)
			fmt.Println(errmsg)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		fmt.Fprintf(w, "uid = %s\nCould not convert to number\n", uidstr)
	}
}
func delPersonRefErrHandler(w http.ResponseWriter, r *http.Request) {
	intPersonRefErrHandler(w, r, "/delPersonRefErr/")
}

func inactivatePersonHandler(w http.ResponseWriter, r *http.Request) {
	intPersonRefErrHandler(w, r, "/inactivatePerson/")
}

func getDirectReportsCount(uid int64) int64 {
	//===============================================================
	//  Check to see if this person manages anyone before deleting...
	//===============================================================
	s := fmt.Sprintf("select uid from people where status=1 and MgrUID=%d", uid)
	rows, err := Phonebook.db.Query(s)
	errcheck(err)
	defer rows.Close()
	var refuid int
	count := int64(0)
	for rows.Next() {
		errcheck(rows.Scan(&refuid))
		count++
	}
	errcheck(rows.Err())
	return count
}

func delPersonHandler(w http.ResponseWriter, r *http.Request) {
	var ssn *db.Session
	var ui uiSupport
	ssn = nil
	if 0 < initHandlerSession(ssn, &ui, w, r) {
		return
	}
	ssn = ui.X
	Phonebook.ReqCountersMem <- 1    // ask to access the shared mem, blocks until granted
	<-Phonebook.ReqCountersMemAck    // make sure we got it
	Counters.DeletePerson++          // initialize our data
	Phonebook.ReqCountersMemAck <- 1 // tell Dispatcher we're done with the data

	// SECURITY
	if !hasAccess(ssn, db.ELEMPERSON, "ElemEntity", db.PERMDEL) {
		ulog("Permissions refuse delPersonHandler page on userid=%d (%s), role=%s\n", ssn.UID, ssn.Firstname, ssn.PMap.Urole.Name)
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}

	// var d db.PersonDetail
	c := "delPersonHandler"
	m := strings.Split(r.RequestURI, "/")
	uidstr := m[len(m)-1]
	if len(uidstr) == 0 {
		fmt.Fprintf(w, "The RequestURI needs the person's uid. It was not found on the URI:  %s\n", r.RequestURI)
		return
	}
	uid, err := strconv.ParseInt(uidstr, 10, 64)
	if err != nil {
		fmt.Fprintf(w, "Error converting uid to a number: %v. URI: %s\n", err, r.RequestURI)
		return
	}

	count := getDirectReportsCount(uid)
	if count > 0 {
		http.Redirect(w, r, fmt.Sprintf("/delPersonRefErr/%d", uid), http.StatusFound)
		return
	}

	//===============================
	//  ******  BEGIN TRANSACTION  ******
	//===============================
	//------------------------------------------------------------
	// in order to delete a person, we must delete all references
	// to the person in the following database tables:
	//		deductions
	//		compensation
	//------------------------------------------------------------
	s := fmt.Sprintf("DELETE FROM people WHERE UID=%d", uid)
	_, err = Phonebook.prepstmt.delPerson.Exec(uid)
	if delCheckError(c, ssn, err, s, w, r) {
		return
	}

	s = fmt.Sprintf("DELETE FROM deductions WHERE UID=%d", uid)
	_, err = Phonebook.prepstmt.delPersonDeduct.Exec(uid)
	if delCheckError(c, ssn, err, s, w, r) {
		return
	}

	s = fmt.Sprintf("DELETE FROM compensation WHERE UID=%d", uid)
	_, err = Phonebook.prepstmt.delPersonComp.Exec(uid)
	if delCheckError(c, ssn, err, s, w, r) {
		return
	}
	//===============================
	//  ******  END TRANSACTION  ******
	//===============================

	http.Redirect(w, r, "/search/", http.StatusFound)
}

func delClassRefErr(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var ssn *db.Session
	var ui uiSupport
	ssn = nil
	if 0 < initHandlerSession(ssn, &ui, w, r) {
		return
	}
	ssn = ui.X

	path := "/delClassRefErr/"
	costr := r.RequestURI[len(path):]
	if len(costr) > 0 {
		classcode, _ := strconv.ParseInt(costr, 10, 64)
		var c db.Class
		getClassInfo(classcode, &c)
		ui.A = &c
		// ui.A.filterSecurityRead(ssn, db.PERMVIEW)
		filterSecurityRead(ui.A, db.ELEMCLASS, ssn, db.PERMVIEW, 0)

		breadcrumbAdd(ssn, "Delete Class", fmt.Sprintf("/delClassRefErr/%d", classcode))

		s := fmt.Sprintf("select uid,lastname,firstname,preferredname,jobcode,primaryemail,officephone,cellphone,deptcode from people where classcode=%d", classcode)
		rows, err := Phonebook.db.Query(s) // does NOT create a prepared statement
		errcheck(err)
		defer rows.Close()
		var d searchResults

		for rows.Next() {
			var m db.Person
			errcheck(rows.Scan(&m.UID, &m.LastName, &m.FirstName, &m.PreferredName, &m.JobCode, &m.PrimaryEmail, &m.OfficePhone, &m.CellPhone, &m.DeptCode))
			m.DeptName = getDepartmentFromDeptCode(m.DeptCode)
			pm := &m
			// pm.filterSecurityRead(ssn, db.PERMVIEW|db.PERMMOD)
			filterSecurityRead(pm, db.ELEMPERSON, ssn, db.PERMVIEW|db.PERMMOD, m.UID)
			d.Matches = append(d.Matches, m)
		}
		errcheck(rows.Err())
		ui.R = &d

		err = renderTemplate(w, ui, "delClassRefErr.html")
		if nil != err {
			errmsg := fmt.Sprintf("delClassRefErr: err = %v\n", err)
			ulog(errmsg)
			fmt.Println(errmsg)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		fmt.Fprintf(w, "classcode = %s\nCould not convert to number\n", costr)
	}
}

func delClassHandler(w http.ResponseWriter, r *http.Request) {
	var ssn *db.Session
	var ui uiSupport
	ssn = nil
	if 0 < initHandlerSession(ssn, &ui, w, r) {
		return
	}
	ssn = ui.X
	Phonebook.ReqCountersMem <- 1    // ask to access the shared mem, blocks until granted
	<-Phonebook.ReqCountersMemAck    // make sure we got it
	Counters.DeleteClass++           // initialize our data
	Phonebook.ReqCountersMemAck <- 1 // tell Dispatcher we're done with the data

	// SECURITY
	if !hasAccess(ssn, db.ELEMCLASS, "ElemEntity", db.PERMDEL) {
		ulog("Permissions refuse delCoHandler page on userid=%d (%s), role=%s\n", ssn.UID, ssn.Firstname, ssn.PMap.Urole.Name)
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}

	c := "delClassHandler"
	m := strings.Split(r.RequestURI, "/")
	uidstr := m[len(m)-1]
	if len(uidstr) == 0 {
		fmt.Fprintf(w, "The RequestURI needs the class code. It was not found on the URI:  %s\n", r.RequestURI)
		return
	}
	ClassCode, err := strconv.ParseInt(uidstr, 10, 64)
	if err != nil {
		fmt.Fprintf(w, "Error converting class code to a number: %v. URI: %s\n", err, r.RequestURI)
		return
	}

	//===============================================================
	//  Check for references to this db.Class before deleting
	//===============================================================
	s := fmt.Sprintf("select uid from people where classcode=%d", ClassCode)
	rows, err := Phonebook.db.Query(s)
	errcheck(err)
	defer rows.Close()
	count := 0
	for rows.Next() {
		var refuid int
		errcheck(rows.Scan(&refuid))
		count++
	}
	errcheck(rows.Err())

	if count > 0 {
		http.Redirect(w, r, fmt.Sprintf("/delClassRefErr/%d", ClassCode), http.StatusFound)
		return
	}

	//===============================================================
	// in order to delete a person, we must delete all references
	// to the person in the following database tables:
	//		deductions
	//		compensationc,
	//===============================================================
	// s = fmt.Sprintf("DELETE FROM classes WHERE ClassCode=%d", ClassCode)
	// stmt, err := Phonebook.db.Prepare(s)
	if delCheckError(c, ssn, err, s, w, r) {
		return
	}
	_, err = Phonebook.prepstmt.delClass.Exec(ClassCode)
	if delCheckError(c, ssn, err, s, w, r) {
		return
	}
	// we've deleted it, now we need to reload our db.Class list...
	loadClasses()
	http.Redirect(w, r, "/searchcl/", http.StatusFound)
}

func delCoRefErr(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var ssn *db.Session
	var ui uiSupport
	ssn = nil
	if 0 < initHandlerSession(ssn, &ui, w, r) {
		return
	}
	ssn = ui.X

	path := "/delCoRefErr/"
	costr := r.RequestURI[len(path):]
	if len(costr) > 0 {
		cocode, _ := strconv.ParseInt(costr, 10, 64)
		var c db.Company
		getCompanyInfo(cocode, &c)
		ui.C = &c
		// ui.C.filterSecurityRead(ssn, db.PERMVIEW)
		filterSecurityRead(ui.C, db.ELEMCOMPANY, ssn, db.PERMVIEW, 0)

		breadcrumbAdd(ssn, "Delete Company", fmt.Sprintf("/delCoRefErr/%d", cocode))

		// s := fmt.Sprintf("select uid,lastname,firstname,preferredname,jobcode,primaryemail,officephone,cellphone,deptcode from people where cocode=%d", cocode)
		rows, err := Phonebook.prepstmt.delCompany.Query(cocode)
		errcheck(err)
		defer rows.Close()
		var d searchResults

		for rows.Next() {
			var m db.Person
			errcheck(rows.Scan(&m.UID, &m.LastName, &m.FirstName, &m.PreferredName, &m.JobCode, &m.PrimaryEmail, &m.OfficePhone, &m.CellPhone, &m.DeptCode))
			m.DeptName = getDepartmentFromDeptCode(m.DeptCode)
			pm := &m
			// pm.filterSecurityRead(ssn, db.PERMVIEW|db.PERMMOD)
			filterSecurityRead(pm, db.ELEMPERSON, ssn, db.PERMVIEW|db.PERMMOD, m.UID)
			d.Matches = append(d.Matches, m)
		}
		errcheck(rows.Err())
		ui.R = &d

		err = renderTemplate(w, ui, "delCoRefErr.html")
		if nil != err {
			errmsg := fmt.Sprintf("delCoRefErr: err = %v\n", err)
			ulog(errmsg)
			fmt.Println(errmsg)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		fmt.Fprintf(w, "cocode = %s\nCould not convert to number\n", costr)
	}
}

func delCoHandler(w http.ResponseWriter, r *http.Request) {
	var ssn *db.Session
	var ui uiSupport
	ssn = nil
	if 0 < initHandlerSession(ssn, &ui, w, r) {
		return
	}
	ssn = ui.X
	Phonebook.ReqCountersMem <- 1    // ask to access the shared mem, blocks until granted
	<-Phonebook.ReqCountersMemAck    // make sure we got it
	Counters.DeleteCompany++         // initialize our data
	Phonebook.ReqCountersMemAck <- 1 // tell Dispatcher we're done with the data

	// SECURITY
	if !hasAccess(ssn, db.ELEMCOMPANY, "ElemEntity", db.PERMDEL) {
		ulog("Permissions refuse delCoHandler page on userid=%d (%s), role=%s\n", ssn.UID, ssn.Firstname, ssn.PMap.Urole.Name)
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}

	c := "delCoHandler"
	m := strings.Split(r.RequestURI, "/")
	uidstr := m[len(m)-1]
	if len(uidstr) == 0 {
		fmt.Fprintf(w, "The RequestURI needs the company code. It was not found on the URI:  %s\n", r.RequestURI)
		return
	}
	CoCode, err := strconv.Atoi(uidstr)
	if err != nil {
		fmt.Fprintf(w, "Error converting company code to a number: %v. URI: %s\n", err, r.RequestURI)
		return
	}

	//===============================================================
	//  Check for references to this db.Class before deleting
	//===============================================================
	s := fmt.Sprintf("select uid from people where CoCode=%d", CoCode)
	rows, err := Phonebook.db.Query(s)
	errcheck(err)
	defer rows.Close()
	count := 0
	for rows.Next() {
		var refuid int
		errcheck(rows.Scan(&refuid))
		count++
	}
	errcheck(rows.Err())

	if count > 0 {
		http.Redirect(w, r, fmt.Sprintf("/delCoRefErr/%d", CoCode), http.StatusFound)
		return
	}

	//===============================================================
	// in order to delete a person, we must delete all references
	// to the person in the following database tables:
	//		deductions
	//		compensation
	//===============================================================
	s = fmt.Sprintf("DELETE FROM companies WHERE CoCode=%d", CoCode)
	stmt, err := Phonebook.db.Prepare(s)
	if delCheckError(c, ssn, err, s, w, r) {
		return
	}
	_, err = stmt.Exec()
	if delCheckError(c, ssn, err, s, w, r) {
		return
	}
	// we've deleted it, now we need to reload our company list...
	loadCompanies()
	http.Redirect(w, r, "/searchco/", http.StatusFound)
}
