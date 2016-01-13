package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// viewCompany executes the server command to serve a company detail page and validates the
// data in the HTML returned.
// RETURNS:
//		true = all data verified correctly
//     false = the test failed for any on of several reasons. If the session is not established
//             after the request, the test fails. If the request succeeds, and one or more of
//             the data fields were not correct then the test fails.
func viewCompany(d *personDetail, atr *TestResults) bool {
	// to save confusion, any user doing tests will only read/modify the company
	// that matches their UID.  That way, when multiple users are editing and modifying
	// companies, no user will stomp on the work being done by another user
	URL := fmt.Sprintf("http://%s:%d/company/%d", App.Host, App.Port, d.UID)
	hc := http.Client{}
	pageName := "Company Detail"

	req, err := http.NewRequest("GET", URL, nil)
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
	// fmt.Printf("viewCompany: Adding session cookie: %#v\n", d.SessionCookie)
	req.AddCookie(d.SessionCookie)
	resp, err := hc.Do(req)
	errcheck(err)
	defer resp.Body.Close()

	// fmt.Printf("viewCompany: hc.Do(req) returned err = %v\n", err)

	cookies := resp.Cookies()
	// fmt.Printf("Cookies:value: %+v\n", cookies)

	d.SessionCookie = nil
	for i := 0; i < len(cookies); i++ {
		if cookies[i].Name == "accord" {
			d.SessionCookie = cookies[i]
			break
		}
	}

	if d.SessionCookie == nil {
		fmt.Printf("Session cookie is nil after hc.Do(req)\n")
		return false
	}

	// Check that the server actually sent compressed data
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		fmt.Printf("gzip response\n")
		reader, err = gzip.NewReader(resp.Body)
		defer reader.Close()
	default:
		reader = resp.Body
	}

	//==================================================
	// Verify that we were sent
	//==================================================
	htmlData, err := ioutil.ReadAll(reader)
	errcheck(err)
	s := string(htmlData)
	m1 := reTitle.FindStringIndex(s)
	m2 := reTitleEnd.FindStringIndex(s)
	m := s[m1[1]:m2[0]]

	var tr TestResults
	tr.Failures = make([]TestFailure, 0)

	// fmt.Printf("Page returned = %s\n", m)
	if strings.Contains(m, "Accord") && strings.Contains(m, "Company") && d.SessionCookie != nil {
		var c company
		getCompanyInfo(d.UID, &c) // yes, we're getting the company with cocode == d.UID

		addrCityState := fmt.Sprintf("%s, %s&nbsp;&nbsp;%s, %s", c.City, c.State, c.PostalCode, c.Country)

		validate := []validationTable{
			/* 00 */ {pageName + ": validate Legal Name", &s, `class="LastName"`, `class="FirstName"`, true, c.LegalName},
			/* 01 */ {pageName + ": validate Common Name", &s, `class="FirstName"`, `ADDRESS`, true, c.CommonName},
			/* 02 */ {pageName + ": validate Address", &s, `ADDRESS`, `DESIGNATION`, true, c.Address},
			/* 03 */ {pageName + ": validate CityState", &s, `ADDRESS`, `DESIGNATION`, true, addrCityState},
			/* 04 */ {pageName + ": validate Designation", &s, `>DESIGNATION<`, `>COMPANY CODE<`, true, c.Designation},
			/* 05 */ {pageName + ": validate Company Code", &s, `>COMPANY CODE<`, `>PHONE<`, true, fmt.Sprintf("%d", c.CoCode)},
			/* 06 */ {pageName + ": validate Phone", &s, `>PHONE<`, `>FAX<`, true, c.Phone},
			/* 07 */ {pageName + ": validate Fax", &s, `>FAX<`, `>EMAIL<`, true, c.Fax},
			/* 08 */ {pageName + ": validate email", &s, `>EMAIL<`, `>ACTIVE<`, true, c.Email},
			/* 09 */ {pageName + ": validate Status", &s, `>ACTIVE<`, `>EMPLOYS PEOPLE<`, true, fmt.Sprintf("%s", yesnoToString(c.Active))},
			/* 10 */ {pageName + ": validate EmploysPersonnel", &s, `>EMPLOYS PEOPLE<`, `/adminViewBtn/`, true, fmt.Sprintf("%s", yesnoToString(c.EmploysPersonnel))},
		}

		var tc testContext
		tc.d = d
		tc.co = &c
		executeValSubstrTests(&validate, &tr, &tc)
		if tr.Fail > 0 {
			fmt.Printf("Company uid = %d\n", d.UID)
			fmt.Printf("tr.Fail = %d, len(tr.Failures)=%d\n", tr.Fail, len(tr.Failures))
			dumpTestErrors(&tr)
		}
		aggregateTR(&tr, atr)
		return (tr.Fail == 0)
	}
	return false
}

// adminEditCompany executes the server command to serve a company adminEditCo page and validates the
// data in the HTML returned.
// RETURNS:
//		true = all data verified correctly
//     false = the test failed for any on of several reasons. If the session is not established
//             after the request, the test fails. If the request succeeds, and one or more of
//             the data fields were not correct then the test fails.
func adminEditCompany(d *personDetail, atr *TestResults) bool {
	// to save confusion, any user doing tests will only read/modify the company
	// that matches their UID.  That way, when multiple users are editing and modifying
	// companies, no user will stomp on the work being done by another user
	URL := fmt.Sprintf("http://%s:%d/adminEditCo/%d", App.Host, App.Port, d.UID)
	hc := http.Client{}
	pageName := "AdminEditCo"

	req, err := http.NewRequest("GET", URL, nil)
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
	// fmt.Printf("viewCompany: Adding session cookie: %#v\n", d.SessionCookie)
	req.AddCookie(d.SessionCookie)
	resp, err := hc.Do(req)
	errcheck(err)
	defer resp.Body.Close()

	// fmt.Printf("viewCompany: hc.Do(req) returned err = %v\n", err)

	cookies := resp.Cookies()
	// fmt.Printf("Cookies:value: %+v\n", cookies)

	d.SessionCookie = nil
	for i := 0; i < len(cookies); i++ {
		if cookies[i].Name == "accord" {
			d.SessionCookie = cookies[i]
			break
		}
	}

	if d.SessionCookie == nil {
		fmt.Printf("Session cookie is nil after hc.Do(req)\n")
		return false
	}

	// Check that the server actually sent compressed data
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		fmt.Printf("gzip response\n")
		reader, err = gzip.NewReader(resp.Body)
		defer reader.Close()
	default:
		reader = resp.Body
	}

	//==================================================
	// Verify that we were sent
	//==================================================
	htmlData, err := ioutil.ReadAll(reader)
	errcheck(err)
	s := string(htmlData)
	m1 := reTitle.FindStringIndex(s)
	m2 := reTitleEnd.FindStringIndex(s)
	m := s[m1[1]:m2[0]]

	var tr TestResults
	tr.Failures = make([]TestFailure, 0)

	// fmt.Printf("Page returned = %s\n", m)
	if strings.Contains(m, "Accord") && strings.Contains(m, "Admin Edit") && d.SessionCookie != nil {
		var c company
		getCompanyInfo(d.UID, &c) // yes, we're getting the company with cocode == d.UID

		var sAct = string(`value="Inactive" selected>`)
		if c.Active > 0 {
			sAct = `value="Active" selected>`
		}
		var sEmp = string(`value="No" selected>`)
		if c.EmploysPersonnel > 0 {
			sEmp = `value="Yes" selected>`
		}

		validate := []validationTable{
			/* 00 */ {pageName + ": validate top line", &s, `class="AppHeading"`, `<form action="/saveAdminEditCo/`, true, fmt.Sprintf("%s (%d)", c.LegalName, c.CoCode)},
			/* 01 */ {pageName + ": validate Legal Name", &s, `name="LegalName"`, `name="CommonName"`, false, `value="` + c.LegalName},
			/* 02 */ {pageName + ": validate Common Name", &s, `name="CommonName"`, `name="Designation"`, false, `value="` + c.CommonName},
			/* 03 */ {pageName + ": validate Designation", &s, `name="Designation"`, `>PHONE<`, false, `value="` + c.Designation},
			/* 04 */ {pageName + ": validate Address", &s, `name="Address"`, `name="Address2"`, false, `value="` + c.Address},
			/* 05 */ {pageName + ": validate City", &s, `name="City"`, `name="State"`, false, `value="` + c.City},
			/* 06 */ {pageName + ": validate State", &s, `name="State"`, `name="PostalCode"`, false, `value="` + c.State},
			/* 07 */ {pageName + ": validate Country", &s, `name="Country"`, `>STATUS<`, false, `value="` + c.Country},
			/* 08 */ {pageName + ": validate PostalCode", &s, `name="PostalCode"`, `name="Country"`, false, `value="` + c.PostalCode},
			/* 09 */ {pageName + ": validate Phone", &s, `name="Phone"`, `name="Fax"`, false, `value="` + c.Phone},
			/* 10 */ {pageName + ": validate Fax", &s, `name="Fax"`, `name="Email"`, false, `value="` + c.Fax},
			/* 11 */ {pageName + ": validate email", &s, `name="Email"`, `>ADDRESS<`, false, `value="` + c.Email},
			/* 09 */ {pageName + ": validate Status", &s, `>STATUS<`, `>EMPLOYS PERSONNEL`, false, sAct},
			/* 10 */ {pageName + ": validate EmploysPersonnel", &s, `>EMPLOYS PERSONNEL`, `type="submit"`, false, sEmp},
		}

		var tc testContext
		tc.d = d
		tc.co = &c
		executeValSubstrTests(&validate, &tr, &tc)
		if tr.Fail > 0 {
			fmt.Printf("Company uid = %d\n", d.UID)
			fmt.Printf("tr.Fail = %d, len(tr.Failures)=%d\n", tr.Fail, len(tr.Failures))
			dumpTestErrors(&tr)
		}
		aggregateTR(&tr, atr)

		return (tr.Fail == 0)
	}
	return false
}
