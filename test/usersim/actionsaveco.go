package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
)

// adminEdit changes the details for the current user, saves them, then validates the change after saving
//    returns true if the save was successful
//            false if save failed
func saveAdminEditCo(d *personDetail, atr *TestResults) bool {
	URL := fmt.Sprintf("http://%s:%d/saveAdminEditCo/%d", App.Host, App.Port, d.UID)
	var c company
	hc := http.Client{}
	getCompanyInfo(d.UID, &c)

	// fmt.Printf("Company %d\n", d.UID)
	// fmt.Printf("Initial:  %#v\n", c)

	form := url.Values{}

	// user d.UID must only change the company names to something that none
	// of the other users will use.
	n := len(App.Companies) / App.TestUsers // this is the range that ea
	nstart := n * (d.UID - 1)               // starting index for this user

	// set new values
	c.LegalName = App.Companies[nstart+rand.Intn(n)]
	c.CommonName = c.LegalName
	c.Designation = genDesignation(c.LegalName, "cocode", "companies")
	c.Email = randomCompanyEmail(c.LegalName)
	c.Phone = randomPhoneNumber()
	c.Fax = randomPhoneNumber()
	c.Active = rand.Intn(2)
	c.EmploysPersonnel = rand.Intn(2)
	c.Address = randomAddress()
	if rand.Intn(10) > 7 {
		c.Address2 = fmt.Sprintf("Suite %d", 1+rand.Intn(10000))
	}
	c.City = App.Cities[rand.Intn(len(App.Cities))]
	c.State = App.States[rand.Intn(len(App.States))]
	c.PostalCode = fmt.Sprintf("%05d", rand.Intn(99999))
	c.Country = "USA"

	// fmt.Printf("Random Updates:  %#v\n", c)

	//===================================================
	// Simulate filling in the fields...
	//===================================================
	form.Add("LegalName", c.LegalName)
	form.Add("CommonName", c.CommonName)
	form.Add("Designation", c.Designation)
	form.Add("Email", c.Email)
	form.Add("Phone", c.Phone)
	form.Add("Fax", c.Fax)
	form.Add("Active", activeToString(c.Active))
	form.Add("EmploysPersonnel", yesnoToString(c.EmploysPersonnel))
	form.Add("Address", c.Address)
	form.Add("Address2", c.Address2)
	form.Add("City", c.City)
	form.Add("State", c.State)
	form.Add("PostalCode", c.PostalCode)
	form.Add("Country", c.Country)

	//===================================================
	// Simulate the button press...
	//===================================================
	form.Add("action", "save")

	req, err := http.NewRequest("POST", URL, bytes.NewBufferString(form.Encode()))
	errcheck(err)

	hdrs := []KeyVal{
		// {"Host:", fmt.Sprintf("%s:%d", App.Host, App.Port)},
		{"Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"},
		{"Accept-Encoding", "gzip, deflate"},
		{"Accept-Language", "en-US,en;q=0.8"},
		{"Cache-Control", "max-age=0"},
		{"Connection", "keep-alive"},
		{"Content-Type", "application/x-www-form-urlencoded"},
		{"User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.80 Safari/537.36"},
	}
	for i := 0; i < len(hdrs); i++ {
		req.Header.Add(hdrs[i].key, hdrs[i].value)
	}
	// fmt.Printf("adminEditSave: Adding session cookie: %#v\n", d.SessionCookie)
	req.AddCookie(d.SessionCookie)

	// fmt.Printf("d before save = %#v\n", d)

	//===================================================
	// SUBMIT THE FORM...
	//===================================================
	resp, err := hc.Do(req)
	if nil != err {
		fmt.Printf("saveAdminEdit:  hc.Do(req) returned error:  %#v\n", err)
		fmt.Printf("err: %s\n", err.Error())
		return false
	}
	defer resp.Body.Close()

	// if 1 > 0 {
	// 	fmt.Printf("DumpResponse:\n")
	// 	dump, err := httputil.DumpResponse(resp, true)
	// 	errcheck(err)
	// 	fmt.Printf("\n\ndumpResponse = %s\n", string(dump))
	// }

	// Verify if the response was ok
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("saveAdminEditCo:  Server return non-200 status: %v\n", resp.Status)
	}

	//==============================================================
	// now read in the updated version
	//==============================================================
	var cnew company
	getCompanyInfo(d.UID, &cnew)

	// fmt.Printf("After db update and readback: %#v\n", cnew)
	res := c.matches(&cnew, atr)

	return res
}
