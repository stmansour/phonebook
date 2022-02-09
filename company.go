package main

import (
	"fmt"
	"net/http"
	"phonebook/db"
	"strconv"
)

func companyInit(c *db.Company) {
	c.LegalName = ""
	c.CommonName = ""
	c.Address = ""
	c.Address2 = ""
	c.City = ""
	c.State = ""
	c.PostalCode = ""
	c.Country = ""
	c.Phone = ""
	c.Fax = ""
	c.Email = ""
	c.Designation = ""
	c.Active = 0
}

// func (c *db.Company) filterSecurityRead(ssn *sess.Session, permRequired int) {
// 	filterSecurityRead(c, db.ELEMCOMPANY, ssn, permRequired, 0)
// }

// MapKey is Accord's key for using google maps

// func mapURL(addr, city, state, zip, country string) string {
// 	s := fmt.Sprintf("https://www.google.com/maps/embed/v1/place?key=%s&q=%s,%s+%s+%s+%s",
// 		db.DB.Config.MapKey, addr, city, state, zip, country)
// 	lib.Console("MAP URL:    %s\n", s)
// 	return s
// }

// func (c *db.Company) mapURL() string {
// 	return mapURL(c.Address, c.City, c.State, c.PostalCode, c.Country)
// }

func getCompanyInfo(cocode int64, c *db.Company) {
	Phonebook.ReqCountersMem <- 1    // ask to access the shared mem, blocks until granted
	<-Phonebook.ReqCountersMemAck    // make sure we got it
	Counters.ViewCompany++           // initialize our data
	Phonebook.ReqCountersMemAck <- 1 // tell Dispatcher we're done with the data
	rows, err := Phonebook.prepstmt.companyInfo.Query(cocode)
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		errcheck(rows.Scan(&c.CoCode, &c.LegalName, &c.CommonName, &c.Address, &c.Address2, &c.City, &c.State, &c.PostalCode, &c.Country, &c.Phone, &c.Fax, &c.Email, &c.Designation, &c.Active, &c.EmploysPersonnel))
	}
	errcheck(rows.Err())

	rows, err = Phonebook.prepstmt.CompanyClasses.Query(cocode)
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		var cl db.Class
		errcheck(rows.Scan(&cl.ClassCode, &cl.CoCode, &cl.Name, &cl.Designation, &cl.Description, &cl.LastModTime, &cl.LastModBy))
		c.C = append(c.C, cl)
	}
	errcheck(rows.Err())
}

func companyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var ssn *db.Session
	var ui uiSupport
	ssn = nil
	if 0 < initHandlerSession(ssn, &ui, w, r) {
		return
	}
	ssn = ui.X

	// SECURITY
	if !ssn.ElemPermsAny(db.ELEMCOMPANY, db.PERMVIEW) {
		ulog("Permissions refuse company view page on userid=%d (%s), role=%s\n", ssn.UID, ssn.Firstname, ssn.PMap.Urole.Name)
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}

	var c db.Company
	path := "/company/"
	costr := r.RequestURI[len(path):]
	if len(costr) > 0 {
		cocode, _ := strconv.ParseInt(costr, 10, 64)
		breadcrumbAdd(ssn, "Company", fmt.Sprintf("/company/%d", cocode))
		getCompanyInfo(cocode, &c)
		ui.C = &c
		filterSecurityRead(ui.C, db.ELEMCOMPANY, ssn, db.PERMVIEW, 0)

		err := renderTemplate(w, ui, "company.html")
		if nil != err {
			errmsg := fmt.Sprintf("companyHandler: err = %v\n", err)
			ulog(errmsg)
			fmt.Println(errmsg)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		fmt.Fprintf(w, "cocode = %s\nCould not convert to number\n", costr)
	}
}
