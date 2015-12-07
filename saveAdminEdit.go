package main

import (
	"fmt"
	"net/http"
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
func parseDeductions(d *personDetail) {
	ca := strings.Split(d.DeductionsStr, ",")
	d.Deductions = d.Deductions[:0] // clear it
	for i := 0; i < len(ca); i++ {
		d.Deductions = append(d.Deductions, strToInt(ca[i]))
	}
}

func saveAdminEditHandler(w http.ResponseWriter, r *http.Request) {
	var sess *session
	var ui uiSupport
	sess = nil
	if 0 < initHandlerSession(sess, &ui, w, r) {
		return
	}
	sess = ui.X

	// SECURITY
	if !sess.elemPermsAny(ELEMPERSON, PERMMOD) {
		ulog("Permissions refuse saveAdminEdit page on userid=%d (%s), role=%s\n", sess.UID, sess.Firstname, sess.Urole.Name)
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}

	var d personDetail
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
		d.PrimaryEmail = r.FormValue("PrimaryEmail")
		d.OfficePhone = r.FormValue("OfficePhone")
		d.CellPhone = r.FormValue("CellPhone")
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
		d.DeptName = r.FormValue("DeptName")
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

		//fmt.Printf("r.FormValue(BirthMonth) = %s,  convert to num -> %d\n", r.FormValue("BirthMonth"), d.BirthMonth)

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
		var do personDetail                       // container for current info
		do.UID = uid                              // init
		adminReadDetails(&do)                     //read current data
		do.filterSecurityMerge(sess, PERMMOD, &d) // merge in new data

		if uid == 0 {
			// some first-time setup that needs to be handled
			if do.RID == 0 {
				do.RID = 4 // default security role is Viewer
			}
			// TODO: set username

			insert, err := Phonebook.db.Prepare("INSERT INTO people (Salutation,FirstName,MiddleName,LastName,PreferredName," +
				"EmergencyContactName,EmergencyContactPhone," +
				"PrimaryEmail,SecondaryEmail,OfficePhone,OfficeFax,CellPhone,CoCode,JobCode," +
				"PositionControlNumber,DeptCode," +
				"HomeStreetAddress,HomeStreetAddress2,HomeCity,HomeState,HomePostalCode,HomeCountry," +
				"status,EligibleForRehire,Accepted401K,AcceptedDentalInsurance,AcceptedHealthInsurance," +
				"Hire,Termination,ClassCode," +
				"BirthMonth,BirthDOM,mgruid,StateOfEmployment,CountryOfEmployment," +
				"LastReview,NextReview,RID) " +
				//      1                 10                  20                  30
				"VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
			errcheck(err)
			_, err = insert.Exec(do.Salutation, do.FirstName, do.MiddleName, do.LastName, do.PreferredName, // 5
				do.EmergencyContactName, do.EmergencyContactPhone, //7
				do.PrimaryEmail, do.SecondaryEmail, do.OfficePhone, do.OfficeFax, do.CellPhone, do.CoCode, do.JobCode, //14
				do.PositionControlNumber, do.DeptCode, //16
				do.HomeStreetAddress, do.HomeStreetAddress2, do.HomeCity, do.HomeState, do.HomePostalCode, do.HomeCountry, // 22
				do.Status, do.EligibleForRehire, do.Accepted401K, do.AcceptedDentalInsurance, do.AcceptedHealthInsurance, // 27
				dateToDBStr(do.Hire), dateToDBStr(do.Termination), do.ClassCode, // 30
				do.BirthMonth, do.BirthDOM, do.MgrUID, do.StateOfEmployment, do.CountryOfEmployment, // 35
				dateToDBStr(do.LastReview), dateToDBStr(do.NextReview), d.RID) // 37
			errcheck(err)

			// read this record back to get the UID...
			rows, err := Phonebook.db.Query("select uid from people where FirstName=? and LastName=? and PrimaryEmail=? and OfficePhone=? and CoCode=? and JobCode=?",
				do.FirstName, do.LastName, do.PrimaryEmail, do.OfficePhone, do.CoCode, do.JobCode)
			errcheck(err)
			defer rows.Close()
			nUID := 0 // quick way to handle multiple matches... in this case, largest UID wins, it hast to be the latest person added
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
			update, err := Phonebook.db.Prepare("update people set Salutation=?,FirstName=?,MiddleName=?,LastName=?,PreferredName=?," + // 5
				"EmergencyContactName=?,EmergencyContactPhone=?," + // 7
				"PrimaryEmail=?,SecondaryEmail=?,OfficePhone=?,OfficeFax=?,CellPhone=?,CoCode=?,JobCode=?," + // 14
				"PositionControlNumber=?,DeptCode=?," + // 16
				"HomeStreetAddress=?,HomeStreetAddress2=?,HomeCity=?,HomeState=?,HomePostalCode=?,HomeCountry=?," + // 22
				"status=?,EligibleForRehire=?,Accepted401K=?,AcceptedDentalInsurance=?,AcceptedHealthInsurance=?," + // 27
				"Hire=?,Termination=?,ClassCode=?," + // 30
				"BirthMonth=?,BirthDOM=?,mgruid=?,StateOfEmployment=?,CountryOfEmployment=?," + // 35
				"LastReview=?,NextReview=? " + // 37
				"where people.uid=?")
			errcheck(err)
			_, err = update.Exec(
				do.Salutation, do.FirstName, do.MiddleName, do.LastName, do.PreferredName,
				do.EmergencyContactName, do.EmergencyContactPhone,
				do.PrimaryEmail, do.SecondaryEmail, do.OfficePhone, do.OfficeFax, do.CellPhone, do.CoCode, do.JobCode,
				do.PositionControlNumber, do.DeptCode,
				do.HomeStreetAddress, do.HomeStreetAddress2, do.HomeCity, do.HomeState, do.HomePostalCode, do.HomeCountry,
				do.Status, do.EligibleForRehire, do.Accepted401K, do.AcceptedDentalInsurance, do.AcceptedHealthInsurance,
				dateToDBStr(do.Hire), dateToDBStr(do.Termination), do.ClassCode,
				do.BirthMonth, do.BirthDOM, do.MgrUID, do.StateOfEmployment, do.CountryOfEmployment,
				dateToDBStr(do.LastReview), dateToDBStr(do.NextReview),
				uid)

			if nil != err {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
		//--------------------------------------------------------------------------
		// Remove old compensation type(s) and Insert new compensation type(s)
		//--------------------------------------------------------------------------
		ct, err := Phonebook.db.Prepare("DELETE FROM compensation WHERE uid=?")
		errcheck(err)
		_, err = ct.Exec(do.UID)
		errcheck(err)
		ct, err = Phonebook.db.Prepare("INSERT INTO compensation (uid,type) VALUES(?,?)")
		errcheck(err)
		for i := 0; i < len(do.Comps); i++ {
			_, err := ct.Exec(do.UID, do.Comps[i])
			errcheck(err)
		}

		//--------------------------------------------------------------------------
		// Update deductions...
		//--------------------------------------------------------------------------
		ct, err = Phonebook.db.Prepare("DELETE FROM deductions WHERE uid=?")
		errcheck(err)
		_, err = ct.Exec(do.UID)
		errcheck(err)
		ct, err = Phonebook.db.Prepare("INSERT INTO deductions (uid,deduction) VALUES(?,?)")
		errcheck(err)
		for i := 0; i < len(do.MyDeductions); i++ {
			// fmt.Printf("\"%s\" = %s\n", do.MyDeductions[i].Name, r.FormValue(do.MyDeductions[i].Name))
			if r.FormValue(do.MyDeductions[i].Name) != "" {
				_, err := ct.Exec(do.UID, do.MyDeductions[i].DCode)
				errcheck(err)
			}
		}
	}

	s := breadcrumbBack(sess, 2)
	http.Redirect(w, r, s, http.StatusFound)
}
