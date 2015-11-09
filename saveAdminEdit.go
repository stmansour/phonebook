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

	//----------------------
	// Handle dropdowns...
	//----------------------
	if "none" == d.Salutation {
		d.Salutation = ""
	}

	if uid == 0 {
		insert, err := Phonebook.db.Prepare("INSERT INTO people (Salutation,FirstName,MiddleName,LastName,PreferredName," +
			"EmergencyContactName,EmergencyContactPhone," +
			"PrimaryEmail,SecondaryEmail,OfficePhone,OfficeFax,CellPhone,CoCode,JobCode," +
			"PositionControlNumber,DeptCode," +
			"HomeStreetAddress,HomeStreetAddress2,HomeCity,HomeState,HomePostalCode,HomeCountry," +
			"status,EligibleForRehire,Accepted401K,AcceptedDentalInsurance,AcceptedHealthInsurance," +
			"Hire,Termination,ClassCode," +
			"BirthMonth,BirthDOM,mgruid,StateOfEmployment,CountryOfEmployment," +
			"LastReview,NextReview) " +
			//      1                 10                  20                  30
			"VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
		errcheck(err)
		_, err = insert.Exec(d.Salutation, d.FirstName, d.MiddleName, d.LastName, d.PreferredName, // 5
			d.EmergencyContactName, d.EmergencyContactPhone, //7
			d.PrimaryEmail, d.SecondaryEmail, d.OfficePhone, d.OfficeFax, d.CellPhone, d.CoCode, d.JobCode, //14
			d.PositionControlNumber, d.DeptCode, //16
			d.HomeStreetAddress, d.HomeStreetAddress2, d.HomeCity, d.HomeState, d.HomePostalCode, d.HomeCountry, // 22
			d.Status, d.EligibleForRehire, d.Accepted401K, d.AcceptedDentalInsurance, d.AcceptedHealthInsurance, // 27
			dateToDBStr(d.Hire), dateToDBStr(d.Termination), d.ClassCode, // 30
			d.BirthMonth, d.BirthDOM, d.MgrUID, d.StateOfEmployment, d.CountryOfEmployment, // 35
			dateToDBStr(d.LastReview), dateToDBStr(d.NextReview)) // 37
		errcheck(err)

		// read this record back to get the UID...
		rows, err := Phonebook.db.Query("select uid from people where FirstName=? and LastName=? and PrimaryEmail=? and OfficePhone=? and CoCode=? and JobCode=?",
			d.FirstName, d.LastName, d.PrimaryEmail, d.OfficePhone, d.CoCode, d.JobCode)
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
		d.UID = uid
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
			d.Salutation, d.FirstName, d.MiddleName, d.LastName, d.PreferredName,
			d.EmergencyContactName, d.EmergencyContactPhone,
			d.PrimaryEmail, d.SecondaryEmail, d.OfficePhone, d.OfficeFax, d.CellPhone, d.CoCode, d.JobCode,
			d.PositionControlNumber, d.DeptCode,
			d.HomeStreetAddress, d.HomeStreetAddress2, d.HomeCity, d.HomeState, d.HomePostalCode, d.HomeCountry,
			d.Status, d.EligibleForRehire, d.Accepted401K, d.AcceptedDentalInsurance, d.AcceptedHealthInsurance,
			dateToDBStr(d.Hire), dateToDBStr(d.Termination), d.ClassCode,
			d.BirthMonth, d.BirthDOM, d.MgrUID, d.StateOfEmployment, d.CountryOfEmployment,
			dateToDBStr(d.LastReview), dateToDBStr(d.NextReview),
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
	_, err = ct.Exec(d.UID)
	errcheck(err)
	ct, err = Phonebook.db.Prepare("INSERT INTO compensation (uid,type) VALUES(?,?)")
	errcheck(err)
	for i := 0; i < len(d.Comps); i++ {
		_, err := ct.Exec(d.UID, d.Comps[i])
		errcheck(err)
	}

	//--------------------------------------------------------------------------
	// Update deductions...
	//--------------------------------------------------------------------------
	ct, err = Phonebook.db.Prepare("DELETE FROM deductions WHERE uid=?")
	errcheck(err)
	_, err = ct.Exec(d.UID)
	errcheck(err)
	ct, err = Phonebook.db.Prepare("INSERT INTO deductions (uid,deduction) VALUES(?,?)")
	errcheck(err)
	for i := 0; i < len(d.MyDeductions); i++ {
		//fmt.Printf("%s = %s\n", d.MyDeductions[i].Name, r.FormValue(d.MyDeductions[i].Name))
		if r.FormValue(d.MyDeductions[i].Name) != "" {
			_, err := ct.Exec(d.UID, d.MyDeductions[i].DCode)
			errcheck(err)
		}
	}

	http.Redirect(w, r, fmt.Sprintf("/adminView/%d", uid), http.StatusFound)
}
