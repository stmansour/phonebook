package main

import (
	"crypto/sha512"
	"fmt"
	"net/http"
	"strconv"
)

func savePersonDetailsHandler(w http.ResponseWriter, r *http.Request) {
	var sess *session
	var ui uiSupport
	sess = nil
	if 0 < initHandlerSession(sess, &ui, w, r) {
		return
	}
	sess = ui.X

	var d personDetail
	path := "/savePersonDetails/"
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

	//=================================================================
	// SECURITY
	//=================================================================
	if !sess.elemPermsAny(ELEMPERSON, PERMOWNERMOD) {
		ulog("Permissions refuse savePersonDetails page on userid=%d (%s), role=%s\n", sess.UID, sess.Firstname, sess.Urole.Name)
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}
	if uid != sess.UID {
		ulog("Permissions refuse savePersonDetails page on userid=%d (%s), role=%s trying to save for UID=%d\n", sess.UID, sess.Firstname, sess.Urole.Name, uid)
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}

	d.PreferredName = r.FormValue("PreferredName")
	d.PrimaryEmail = r.FormValue("PrimaryEmail")
	d.OfficePhone = r.FormValue("OfficePhone")
	d.CellPhone = r.FormValue("CellPhone")
	d.EmergencyContactPhone = r.FormValue("EmergencyContactPhone")
	d.EmergencyContactName = r.FormValue("EmergencyContactName")
	d.HomeStreetAddress = r.FormValue("HomeStreetAddress")
	d.HomeStreetAddress2 = r.FormValue("HomeStreetAddress2")
	d.HomeCity = r.FormValue("HomeCity")
	d.HomeState = r.FormValue("HomeState")
	d.HomePostalCode = r.FormValue("HomePostalCode")
	d.HomeCountry = r.FormValue("HomeCountry")

	// fmt.Printf("email = %s, officephone = %s, cell = %s", d.PrimaryEmail, d.OfficePhone, d.CellPhone)

	update, err := Phonebook.db.Prepare("update people set PreferredName=?,PrimaryEmail=?,OfficePhone=?,CellPhone=?," +
		"EmergencyContactName=?,EmergencyContactPhone=?," +
		"HomeStreetAddress=?,HomeStreetAddress2=?,HomeCity=?,HomeState=?,HomePostalCode=?,HomeCountry=? " +
		"where people.uid=?")
	errcheck(err)

	_, err = update.Exec(d.PreferredName, d.PrimaryEmail, d.OfficePhone, d.CellPhone,
		d.EmergencyContactName, d.EmergencyContactPhone,
		d.HomeStreetAddress, d.HomeStreetAddress2, d.HomeCity, d.HomeState, d.HomePostalCode, d.HomeCountry,
		uid)
	if nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	password := r.FormValue("password")
	if "" != password {
		sha := sha512.Sum512([]byte(password))
		passhash := fmt.Sprintf("%x", sha)
		update, err = Phonebook.db.Prepare("update people set passhash=? where uid=?")
		errcheck(err)
		_, err = update.Exec(passhash, uid)
		if nil != err {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	http.Redirect(w, r, fmt.Sprintf("/detail/%d", uid), http.StatusFound)
}
