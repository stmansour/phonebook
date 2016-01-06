package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var salutations = []string{
	"Mr",
	"Mrs",
	"Ms",
	"Dr",
}

// RandomString returns a random string of length n containing alpha numeric characters
func RandomString(n int) string {
	const c = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	s := make([]byte, n)
	for i := 0; i < n; i++ {
		s[i] = c[rand.Intn(len(c))]
	}
	return string(s)
}

// adminEdit changes the details for the current user, saves them, then validates the change after saving
//    returns true if the save was successful
//            false if save failed
func saveAdminEdit(d *personDetail) bool {
	URL := fmt.Sprintf("http://%s:%d/saveAdminEdit/%d", App.Host, App.Port, d.UID)
	hc := http.Client{}
	Nlast := len(App.LastNames)
	Nfirst := len(App.FirstNames)

	form := url.Values{}
	d.BirthDOM = 1 + rand.Intn(29)   // choose a random date
	d.BirthMonth = 1 + rand.Intn(12) // choose a random birthmonth
	d.Salutation = salutations[rand.Intn(len(salutations))]
	d.FirstName = strings.ToLower(App.FirstNames[rand.Intn(Nfirst)])
	d.MiddleName = strings.ToLower(App.FirstNames[rand.Intn(Nfirst)])
	d.LastName = strings.ToLower(App.LastNames[rand.Intn(Nlast)])
	d.PreferredName = strings.ToLower(App.FirstNames[rand.Intn(Nfirst)])
	d.PrimaryEmail = randomEmail(d.LastName, d.FirstName)
	d.SecondaryEmail = randomEmail(d.LastName, d.FirstName)
	d.OfficePhone = randomPhoneNumber()
	d.CellPhone = randomPhoneNumber()
	d.OfficeFax = randomPhoneNumber()
	d.EmergencyContactName = strings.ToLower(App.FirstNames[rand.Intn(Nfirst)] + " " + App.LastNames[rand.Intn(Nlast)])
	d.EmergencyContactPhone = randomPhoneNumber()
	d.ClassCode = 1 + rand.Intn(len(App.NameToClassCode))
	d.CoCode = 1 + rand.Intn(len(App.NameToCoCode))
	d.DeptCode = 1 + rand.Intn(1+rand.Intn(App.DeptHi-App.DeptLo-1))
	d.JobCode = 1 + rand.Intn(1+rand.Intn(App.JCHi-App.JCLo-1))
	d.Status = 1
	d.HomeStreetAddress = randomAddress()
	d.HomeCity = App.Cities[rand.Intn(len(App.Cities))]
	d.HomeState = App.States[rand.Intn(len(App.States))]
	d.HomePostalCode = fmt.Sprintf("%05d", rand.Intn(99999))
	d.HomeCountry = "USA"
	d.StateOfEmployment = App.States[rand.Intn(len(App.States))]
	d.CountryOfEmployment = "Tazmania"
	d.Accepted401K = rand.Intn(1 + ACPTLAST)
	d.AcceptedDentalInsurance = rand.Intn(1 + ACPTLAST)
	d.AcceptedHealthInsurance = rand.Intn(1 + ACPTLAST)
	d.EligibleForRehire = rand.Intn(2)
	d.LastReview = stringToDate(fmt.Sprintf("%04d-%02d-%02d", 1990+rand.Intn(27), 1+rand.Intn(12), 1+rand.Intn(28)))
	d.NextReview = stringToDate(fmt.Sprintf("%04d-%02d-%02d", 1990+rand.Intn(27), 1+rand.Intn(12), 1+rand.Intn(28)))
	d.Hire = stringToDate(fmt.Sprintf("%04d-%02d-%02d", 1990+rand.Intn(27), 1+rand.Intn(12), 1+rand.Intn(28)))
	d.Termination = stringToDate(fmt.Sprintf("%04d-%02d-%02d", 1990+rand.Intn(27), 1+rand.Intn(12), 1+rand.Intn(28)))
	d.MgrUID = rand.Intn(len(App.Peeps))
	d.PositionControlNumber = RandomString(10)

	//===================================================
	// Simulate filling in the fields...
	//===================================================
	form.Add("BirthDOM", fmt.Sprintf("%d", d.BirthDOM))
	form.Add("BirthMonth", fmt.Sprintf("%d", d.BirthMonth))
	form.Add("Salutation", d.Salutation)
	form.Add("FirstName", d.FirstName)
	form.Add("LastName", d.LastName)
	form.Add("MiddleName", d.MiddleName)
	form.Add("PreferredName", d.PreferredName)
	form.Add("PrimaryEmail", d.PrimaryEmail)
	form.Add("SecondaryEmail", d.SecondaryEmail)
	form.Add("OfficePhone", d.OfficePhone)
	form.Add("OfficeFax", d.OfficeFax)
	form.Add("CellPhone", d.CellPhone)
	form.Add("EmergencyContactName", d.EmergencyContactName)
	form.Add("EmergencyContactPhone", d.EmergencyContactPhone)
	form.Add("CoCode", fmt.Sprintf("%d", d.CoCode))
	form.Add("ClassCode", fmt.Sprintf("%d", d.ClassCode))
	form.Add("DeptCode", fmt.Sprintf("%d", d.DeptCode))
	form.Add("JobCode", fmt.Sprintf("%d", d.JobCode))
	form.Add("Status", activeToString(d.Status))
	form.Add("HomeStreetAddress", d.HomeStreetAddress)
	form.Add("HomeStreetAddress2", d.HomeStreetAddress2)
	form.Add("HomeCity", d.HomeCity)
	form.Add("HomeState", d.HomeState)
	form.Add("HomeCountry", d.HomeCountry)
	form.Add("HomePostalCode", d.HomePostalCode)
	form.Add("StateOfEmployment", d.StateOfEmployment)
	form.Add("CountryOfEmployment", d.CountryOfEmployment)
	form.Add("AcceptedDentalInsurance", acceptIntToString(d.AcceptedDentalInsurance))
	form.Add("AcceptedHealthInsurance", acceptIntToString(d.AcceptedHealthInsurance))
	form.Add("Accepted401K", acceptIntToString(d.Accepted401K))
	form.Add("EligibleForRehire", yesnoToString(d.EligibleForRehire))
	form.Add("LastReview", dateToString(d.LastReview))
	form.Add("NextReview", dateToString(d.NextReview))
	form.Add("Hire", dateToString(d.Hire))
	form.Add("Termination", dateToString(d.Termination))
	form.Add("MgrUID", fmt.Sprintf("%d", d.MgrUID))
	form.Add("PositionControlNumber", d.PositionControlNumber)
	form.Add("Role", fmt.Sprintf("%d", d.RID))

	//---------------------------------------------------
	// handle the compensation types
	//---------------------------------------------------
	d.Comps = make([]int, 0)
	for i := 0; i < len(d.MyComps); i++ {
		d.MyComps[i].HaveIt = rand.Intn(2)
		h := ""
		if d.MyComps[i].HaveIt > 0 {
			d.Comps = append(d.Comps, d.MyComps[i].CompCode)
			h = " checked"
		}
		form.Add(d.MyComps[i].Name, h)
	}

	//---------------------------------------------------
	// handle deductions...
	//---------------------------------------------------
	d.Deductions = make([]int, 0)
	for i := 0; i < len(d.MyDeductions); i++ {
		d.MyDeductions[i].HaveIt = rand.Intn(2)
		h := ""
		if d.MyDeductions[i].HaveIt > 0 {
			d.Deductions = append(d.Deductions, d.MyDeductions[i].DCode)
			h = " checked"
		}
		form.Add(d.MyDeductions[i].Name, h)
	}

	//===================================================
	// Simulate the button press...
	//===================================================
	form.Add("action", "save")

	req, err := http.NewRequest("POST", URL, bytes.NewBufferString(form.Encode()))
	errcheck(err)

	hdrs := []KeyVal{
		{"Host:", fmt.Sprintf("%s:%d", App.Host, App.Port)},
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
		os.Exit(1)
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
		fmt.Printf("Server return non-200 status: %v\n", resp.Status)
	}

	//==============================================================
	//  Sleep for 2 seconds and let the server make the update...
	//==============================================================
	//time.Sleep(500 * time.Millisecond)

	//==============================================================
	// now read in the updated version
	//==============================================================
	var dnew personDetail
	dnew.UID = d.UID
	adminReadDetails(&dnew)
	res := d.matches(&dnew)
	adminReadDetails(d) // ensure that all values are correct for remaining tests...
	return res
}
