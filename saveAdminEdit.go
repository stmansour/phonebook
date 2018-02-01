package main

import (
	"fmt"
	"net/http"
	"phonebook/authz"
	"phonebook/db"
	"phonebook/sess"
	"strconv"
	"strings"
)

func strToInt(s string) int {
	if len(s) == 0 {
		return 0
	}
	s = strings.Trim(s, " \n\r")
	n, err := strconv.Atoi(s)
	if err != nil {
		fmt.Printf("Error converting %s to a number: %v\n", s, err)
		return 0
	}
	return n
}

// This is a short term function. Should be replaced by a multi instanced dropdown selector
func parseDeductions(d *db.PersonDetail) {
	ca := strings.Split(d.DeductionsStr, ",")
	d.Deductions = d.Deductions[:0] // clear it
	for i := 0; i < len(ca); i++ {
		d.Deductions = append(d.Deductions, strToInt(ca[i]))
	}
}

func saveAdminEditHandler(w http.ResponseWriter, r *http.Request) {
	var ssn *sess.Session
	var ui uiSupport
	ssn = nil
	if 0 < initHandlerSession(ssn, &ui, w, r) {
		return
	}
	ssn = ui.X
	Phonebook.ReqCountersMem <- 1    // ask to access the shared mem, blocks until granted
	<-Phonebook.ReqCountersMemAck    // make sure we got it
	Counters.AdminEditPerson++       // initialize our data
	Phonebook.ReqCountersMemAck <- 1 // tell Dispatcher we're done with the data

	// SECURITY
	if !ssn.ElemPermsAny(authz.ELEMPERSON, authz.PERMMOD) {
		ulog("Permissions refuse saveAdminEdit page on userid=%d (%s), role=%s\n", ssn.UID, ssn.Firstname, ssn.PMap.Urole.Name)
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}

	var d db.PersonDetail
	path := "/saveAdminEdit/"
	uidstr := r.RequestURI[len(path):]
	if len(uidstr) == 0 {
		fmt.Fprintf(w, "The RequestURI needs the person's uid. It was not found on the URI:  %s\n", r.RequestURI)
		return
	}
	uid, err := strconv.Atoi(uidstr)
	if err != nil {
		fmt.Fprintf(w, "Error converting uid to a number: %v. URI: %s\n", err, r.RequestURI)
		return
	}

	//====================================================
	//  Determine which button the user pressed and act...
	//====================================================
	action := strings.ToLower(r.FormValue("action"))

	if action == "delete" {
		url := fmt.Sprintf("/delPerson/%d", uid)
		http.Redirect(w, r, url, http.StatusFound)
		return
	}

	if "save" == action {
		d.UID = uid
		d.Salutation = r.FormValue("Salutation")
		d.FirstName = r.FormValue("FirstName")
		d.MiddleName = r.FormValue("MiddleName")
		d.LastName = r.FormValue("LastName")
		d.PreferredName = r.FormValue("PreferredName")
		d.EmergencyContactName = r.FormValue("EmergencyContactName")
		d.EmergencyContactPhone = r.FormValue("EmergencyContactPhone")
		d.CoCode = strToInt(r.FormValue("CoCode"))
		d.JobCode = strToInt(r.FormValue("JobCode"))
		d.DeptCode = strToInt(r.FormValue("DeptCode"))
		d.ClassCode = strToInt(r.FormValue("ClassCode"))
		d.PositionControlNumber = r.FormValue("PositionControlNumber")
		d.HomeStreetAddress = r.FormValue("HomeStreetAddress")
		d.HomeStreetAddress2 = r.FormValue("HomeStreetAddress2")
		d.HomeCity = r.FormValue("HomeCity")
		d.HomeState = r.FormValue("HomeState")
		d.HomePostalCode = r.FormValue("HomePostalCode")
		d.HomeCountry = r.FormValue("HomeCountry")
		d.PrimaryEmail = r.FormValue("PrimaryEmail")
		d.SecondaryEmail = r.FormValue("SecondaryEmail")
		d.OfficePhone = r.FormValue("OfficePhone")
		d.OfficeFax = r.FormValue("OfficeFax")
		d.CellPhone = r.FormValue("CellPhone")
		// d.DeptName = r.FormValue("DeptName")
		d.Status = activeToInt(r.FormValue("Status")) // active or inactive, old values included "not-active"
		d.EligibleForRehire = yesnoToInt(r.FormValue("EligibleForRehire"))
		d.LastReview = stringToDate(r.FormValue("LastReview"))
		d.NextReview = stringToDate(r.FormValue("NextReview"))
		d.BirthDOM = strToInt(r.FormValue("BirthDOM"))
		d.BirthMonth = strToInt(r.FormValue("BirthMonth"))
		d.MgrUID = strToInt(r.FormValue("MgrUID"))
		d.Accepted401K = acceptTypeToInt(r.FormValue("Accepted401K"))
		d.AcceptedDentalInsurance = acceptTypeToInt(r.FormValue("AcceptedDentalInsurance"))
		d.AcceptedHealthInsurance = acceptTypeToInt(r.FormValue("AcceptedHealthInsurance"))
		d.Hire = stringToDate(r.FormValue("Hire"))
		d.Termination = stringToDate(r.FormValue("Termination"))
		d.StateOfEmployment = r.FormValue("StateOfEmployment")
		d.CountryOfEmployment = r.FormValue("CountryOfEmployment")

		if hasAccess(ssn, authz.ELEMPERSON, "Role", authz.PERMMOD) {
			d.RID = strToInt(r.FormValue("Role"))
		}

		initMyComps(&d)
		d.Comps = d.Comps[:0] // clear the compensation types list
		for i := 0; i < len(d.MyComps); i++ {
			if "" != r.FormValue(d.MyComps[i].Name) {
				d.Comps = append(d.Comps, d.MyComps[i].CompCode)
			}
		}

		initMyDeductions(&d)
		for i := 0; i < len(d.MyDeductions); i++ {
			if "" != r.FormValue(d.MyDeductions[i].Name) {
				d.Deductions = append(d.Deductions, d.MyDeductions[i].DCode)
			}
		}

		//----------------------
		// Handle dropdowns...
		//----------------------
		if "none" == d.Salutation {
			d.Salutation = ""
		}

		//-------------------------------
		// SECURITY
		//-------------------------------
		var do db.PersonDetail // container for current info
		do.UID = uid           // init
		adminReadDetails(&do)  //read current data

		//----------------------------------------------------------------------------
		// If we're changing Status to Inactive, then it's like a delete. We'll have
		// to reference checking...
		//----------------------------------------------------------------------------
		if do.Status == ACTIVE && d.Status == INACTIVE {
			count := getDirectReportsCount(uid)
			if count > 0 {
				http.Redirect(w, r, fmt.Sprintf("/inactivatePerson/%d", uid), http.StatusFound)
				return
			}
		}

		// func (d *db.PersonDetail) filterSecurityMerge(ssn *sess.Session, permRequired int, dNew *db.PersonDetail) {
		// 	filterSecurityMerge(d, ssn, authz.ELEMPERSON, permRequired, dNew, d.UID)
		// }
		filterSecurityMerge(&do, ssn, authz.ELEMPERSON, authz.PERMMOD, &d, do.UID) // merge in new data

		if int64(uid) == ssn.UID {
			if 0 == len(do.PreferredName) {
				ssn.Firstname = do.FirstName
			} else {
				ssn.Firstname = do.PreferredName
			}
		}

		if uid == 0 {
			// some first-time setup that needs to be handled
			if do.RID == 0 {
				do.RID = 4 // default security role is Viewer
			}

			//============================================
			// generate a unique username...
			//============================================
			do.UserName = strings.ToLower(do.FirstName[0:1] + do.LastName)
			do.UserName = stripchars(do.UserName, "., -&`~!@#$%^*()_+={}'[]\";:<>/?\\")
			if len(do.UserName) > 17 {
				do.UserName = do.UserName[0:17]
			}
			UserName := do.UserName
			var xx int
			nUID := 0
			for {
				found := false
				rows, err := Phonebook.db.Query("select uid from people where UserName=?", UserName)
				errcheck(err)
				defer rows.Close()
				for rows.Next() {
					errcheck(rows.Scan(&xx))
					nUID++
					found = true
					UserName = fmt.Sprintf("%s%d", do.UserName, nUID)
				}
				if !found {
					break
				}
			}
			do.UserName = UserName

			//============================================
			// OK, now write it to the db...
			//============================================
			_, err = Phonebook.prepstmt.adminInsertPerson.Exec(do.Salutation, do.FirstName, do.MiddleName, do.LastName, do.PreferredName, // 5
				do.EmergencyContactName, do.EmergencyContactPhone, //7
				do.PrimaryEmail, do.SecondaryEmail, do.OfficePhone, do.OfficeFax, do.CellPhone, do.CoCode, do.JobCode, //14
				do.PositionControlNumber, do.DeptCode, //16
				do.HomeStreetAddress, do.HomeStreetAddress2, do.HomeCity, do.HomeState, do.HomePostalCode, do.HomeCountry, // 22
				do.Status, do.EligibleForRehire, do.Accepted401K, do.AcceptedDentalInsurance, do.AcceptedHealthInsurance, // 27
				dateToDBStr(do.Hire), dateToDBStr(do.Termination), do.ClassCode, // 30
				do.BirthMonth, do.BirthDOM, do.MgrUID, do.StateOfEmployment, do.CountryOfEmployment, // 35
				dateToDBStr(do.LastReview), dateToDBStr(do.NextReview), do.RID, ssn.UID, do.UserName) // 37
			errcheck(err)

			// read this record back to get the UID...
			rows, err := Phonebook.prepstmt.adminReadBack.Query(
				do.FirstName, do.LastName, do.PrimaryEmail, do.OfficePhone, do.CoCode, do.JobCode)
			errcheck(err)
			defer rows.Close()
			nUID = 0 // quick way to handle multiple matches... in this case, largest UID wins, it hast to be the latest person added
			for rows.Next() {
				errcheck(rows.Scan(&uid))
				if uid > nUID {
					nUID = uid
				}
			}
			errcheck(rows.Err())
			uid = nUID
			do.UID = uid
		} else {
			//--------------------------
			// update existing record
			//--------------------------
			_, err = Phonebook.prepstmt.adminUpdatePerson.Exec(
				do.Salutation, do.FirstName, do.MiddleName, do.LastName, do.PreferredName,
				do.EmergencyContactName, do.EmergencyContactPhone,
				do.PrimaryEmail, do.SecondaryEmail, do.OfficePhone, do.OfficeFax, do.CellPhone, do.CoCode, do.JobCode,
				do.PositionControlNumber, do.DeptCode,
				do.HomeStreetAddress, do.HomeStreetAddress2, do.HomeCity, do.HomeState, do.HomePostalCode, do.HomeCountry,
				do.Status, do.EligibleForRehire, do.Accepted401K, do.AcceptedDentalInsurance, do.AcceptedHealthInsurance,
				dateToDBStr(do.Hire), dateToDBStr(do.Termination), do.ClassCode,
				do.BirthMonth, do.BirthDOM, do.MgrUID, do.StateOfEmployment, do.CountryOfEmployment,
				dateToDBStr(do.LastReview), dateToDBStr(do.NextReview), ssn.UID, do.RID,
				uid)

			if nil != err {
				errmsg := fmt.Sprintf("saveAdminEditHandler: Phonebook.prepstmt.adminUpdatePerson.Exec: err = %v\n", err)
				ulog(errmsg)
				fmt.Println(errmsg)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		//--------------------------------------------------------------------------
		// Remove old compensation type(s) and Insert new compensation type(s)
		//--------------------------------------------------------------------------
		// ct, err := Phonebook.db.Prepare("DELETE FROM compensation WHERE uid=?")
		// errcheck(err)
		_, err = Phonebook.prepstmt.delPersonComp.Exec(do.UID)
		errcheck(err)
		// ct, err = Phonebook.db.Prepare("INSERT INTO compensation (uid,type) VALUES(?,?)")
		// errcheck(err)
		for i := 0; i < len(do.Comps); i++ {
			_, err := Phonebook.prepstmt.insertComp.Exec(do.UID, do.Comps[i])
			errcheck(err)
		}

		//--------------------------------------------------------------------------
		// Update deductions...
		//--------------------------------------------------------------------------
		// ct, err = Phonebook.db.Prepare("DELETE FROM deductions WHERE uid=?")
		// errcheck(err)
		_, err = Phonebook.prepstmt.delPersonDeduct.Exec(do.UID)
		errcheck(err)
		// ct, err = Phonebook.db.Prepare("INSERT INTO deductions (uid,deduction) VALUES(?,?)")
		// errcheck(err)
		for i := 0; i < len(do.MyDeductions); i++ {
			// fmt.Printf("\"%s\" = %s\n", do.MyDeductions[i].Name, r.FormValue(do.MyDeductions[i].Name))
			if r.FormValue(do.MyDeductions[i].Name) != "" {
				_, err := Phonebook.prepstmt.insertDeduct.Exec(do.UID, do.MyDeductions[i].DCode)
				errcheck(err)
			}
		}
	}

	s := breadcrumbBack(ssn, 2)
	http.Redirect(w, r, s, http.StatusFound)
}
