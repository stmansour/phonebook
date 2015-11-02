package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func savePersonDetailsHandler(w http.ResponseWriter, r *http.Request) {
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

	d.PreferredName = r.FormValue("PreferredName")
	d.PrimaryEmail = r.FormValue("PrimaryEmail")
	d.OfficePhone = r.FormValue("OfficePhone")
	d.CellPhone = r.FormValue("CellPhone")
	d.EmergencyContactPhone = r.FormValue("EmergencyContactPhone")
	d.EmergencyContactName = r.FormValue("EmergencyContactName")

	// fmt.Printf("email = %s, officephone = %s, cell = %s", d.PrimaryEmail, d.OfficePhone, d.CellPhone)

	update, err := Phonebook.db.Prepare("update people set PreferredName=?,PrimaryEmail=?, OfficePhone=?, CellPhone=?, EmergencyContactName=?, EmergencyContactPhone=? where people.uid=?")
	errcheck(err)
	_, err = update.Exec(d.PreferredName, d.PrimaryEmail, d.OfficePhone, d.CellPhone, d.EmergencyContactName, d.EmergencyContactPhone, uid)
	if nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, fmt.Sprintf("/detail/%d", uid), http.StatusFound)
}
