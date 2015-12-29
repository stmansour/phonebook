package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type validationTable struct {
	testname string  // name of test
	html     *string // full html text
	reStart  string  // regexp for start of substring
	reStop   string  // regexp for end of substring
	target   string  // what we expect to find in the substring
}

// validateSubstring searches for the substring between s1 and s2.
// Returns:
//		true - target string was found in the defined substring
//		false - could mean any of several things:
//					* string s1 was not found
//					* string s2 was not found
//					* target was not found in the string between s1 & s2
func validateSubstring(ps *string, s1 string, s2 string, target string) bool {
	reg1 := regexp.MustCompile(s1)
	reg2 := regexp.MustCompile(s2)
	m1 := reg1.FindStringIndex(*ps)
	m2 := reg2.FindStringIndex(*ps)
	if App.ShowTestMatching {
		fmt.Printf("s1=%s  s2=%s  target=%s\n", s1, s2, target)
		fmt.Printf("m1 = %#v\n", m1)
		fmt.Printf("m2 = %#v\n", m2)
	}
	if nil == m1 || nil == m2 {
		return false
	}
	if m2[0] < m1[1] {
		fmt.Printf("s2 has a bad regexp. A match occurs before s1, at index %d\n", m2[0])
	}
	m := (*ps)[m1[1]:m2[0]]
	b := strings.Contains(m, target)
	if App.ShowTestMatching {
		fmt.Printf("m = %s\n", m)
		fmt.Printf("found target = %v\n", b)
	}
	return b
}

// executeValSubstrTests accepts a table of data check tests, executes each check,
// and returns the number of checks that passed and the number that failed
func executeValSubstrTests(m *[]validationTable, tr *TestResults) {
	for i := 0; i < len(*m); i++ {
		if validateSubstring((*m)[i].html, (*m)[i].reStart, (*m)[i].reStop, (*m)[i].target) {
			tr.Pass++
		} else {
			tr.Fail++
			f := TestFailure{(*m)[i].testname, i}
			tr.Failures = append(tr.Failures, f)
		}
	}
}

// viewPersonDetail executes the server command to serve a person detail page and validates the
// data in the HTML returned.
// RETURNS:
//		true = all data verified correctly
//     false = one or more of the data fields were not correct
func viewPersonDetail(d *personDetail) bool {
	URL := fmt.Sprintf("http://%s:%d/detail/%d", App.Host, App.Port, d.UID)
	hc := http.Client{}

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
	req.AddCookie(d.SessionCookie)
	resp, err := hc.Do(req)
	errcheck(err)
	defer resp.Body.Close()

	// fmt.Printf("viewPersonDetail: hc.Do(req) returned err = %v\n", err)

	cookies := resp.Cookies()
	// fmt.Printf("viewPersonDetail: Cookies:  %+v\n", cookies)
	d.SessionCookie = nil
	for i := 0; i < len(cookies); i++ {
		if cookies[i].Name == "accord" {
			d.SessionCookie = cookies[i]
			break
		}
	}
	if nil == d.SessionCookie {
		fmt.Printf("d.SessionCookie is nil after executing: %s\n", URL)
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

	// Verify that we were sent to the Sign In page...
	htmlData, err := ioutil.ReadAll(reader)
	errcheck(err)
	s := string(htmlData)
	m1 := reTitle.FindStringIndex(s)
	m2 := reTitleEnd.FindStringIndex(s)
	m := s[m1[1]:m2[0]]
	// fmt.Printf("Page returned = %s\n", m)
	var tr TestResults
	tr.Failures = make([]TestFailure, 0)
	if strings.Contains(m, "Accord") && strings.Contains(m, "Details") && d.SessionCookie != nil {
		myname := d.FirstName + " " + d.MiddleName + " " + d.LastName
		validate := []validationTable{
			{"Validate Full Name", &s, "FULL NAME", "PREFERRED NAME", myname},
			{"Validate Preferred Name", &s, "PREFERRED NAME", "EMAIL", d.PreferredName},
			{"Validate Primary Email", &s, "EMAIL", "PHONE &", d.PrimaryEmail},
			{"Validate Office Phone", &s, "PHONE &", "CELL", d.OfficePhone},
			{"Validate Cell Phone", &s, "CELL", "CLASS", d.CellPhone},
			{"Validate Class", &s, "CLASS", "DEPARTMENT", d.Class},
			{"Validate Department", &s, "DEPARTMENT", "MANAGER", d.DeptName},
		}
		executeValSubstrTests(&validate, &tr)
		if tr.Fail > 0 {
			dumpTestErrors(&tr)
		}
		return (tr.Fail == 0)
	}
	return false
}

// viewAdminPersonDetail executes the server command to serve a person detail page and validates the
// data in the HTML returned.
// RETURNS:
//		true = all data verified correctly
//     false = the test failed for any on of several reasons. If the session is not established
//             after the request, the test fails. If the request succeeds, and one or more of
//             the data fields were not correct then the test fails.
func viewAdminPersonDetail(d *personDetail) bool {
	URL := fmt.Sprintf("http://%s:%d/adminView/%d", App.Host, App.Port, d.UID)
	hc := http.Client{}

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
	// fmt.Printf("viewAdminPersonDetail: Adding session cookie: %#v\n", d.SessionCookie)
	req.AddCookie(d.SessionCookie)
	resp, err := hc.Do(req)
	errcheck(err)
	defer resp.Body.Close()

	// fmt.Printf("viewPersonDetail: hc.Do(req) returned err = %v\n", err)

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

	// Verify that we were sent to the Sign In page...
	htmlData, err := ioutil.ReadAll(reader)
	errcheck(err)
	s := string(htmlData)
	m1 := reTitle.FindStringIndex(s)
	m2 := reTitleEnd.FindStringIndex(s)
	m := s[m1[1]:m2[0]]

	var tr TestResults
	tr.Failures = make([]TestFailure, 0)

	// fmt.Printf("Page returned = %s\n", m)
	if strings.Contains(m, "Accord") && strings.Contains(m, "Admin - View") && d.SessionCookie != nil {
		t := time.Date(2000, time.December, 31, 23, 59, 59, 0, time.UTC)
		sHire := ""
		if d.Hire.After(t) {
			sHire = dateToString(d.Hire)
		}
		sLastReview := ""
		if d.LastReview.After(t) {
			sLastReview = dateToString(d.LastReview)
		}
		sNextReview := ""
		if d.NextReview.After(t) {
			sNextReview = dateToString(d.NextReview)
		}
		sTermination := ""
		if d.Termination.After(t) {
			sTermination = dateToString(d.Termination)
		}
		validate := []validationTable{
			/* 00 */ {"adminView validate Last Name", &s, `name="LastName"`, `name="MiddleName"`, `value="` + d.LastName + `"`},
			/* 01 */ {"adminView validate Middle Name", &s, `name="MiddleName"`, `name="FirstName"`, `value="` + d.MiddleName + `"`},
			/* 02 */ {"adminView validate First Name", &s, `name="FirstName"`, `name="PreferredName"`, `value="` + d.FirstName + `"`},
			/* 03 */ {"adminView validate Preferred Name", &s, `name="PreferredName"`, `OFFICE PHONE`, `value="` + d.PreferredName + `"`},
			/* 04 */ {"adminView validate Office Phone", &s, `name="OfficePhone"`, `name="OfficeFax"`, `value="` + d.OfficePhone + `"`},
			/* 05 */ {"adminView validate Office Fax", &s, `name="OfficeFax"`, `name="CellPhone"`, `value="` + d.OfficeFax + `"`},
			/* 06 */ {"adminView validate Cell Phone", &s, `name="CellPhone"`, `PRIMARY EMAIL`, `value="` + d.CellPhone + `"`},
			/* 07 */ {"adminView validate Primary Email", &s, `name="PrimaryEmail"`, `name="SecondaryEmail"`, `value="` + d.PrimaryEmail + `"`},
			/* 08 */ {"adminView validate Secondary Email", &s, `name="SecondaryEmail"`, `HOME STREET ADDRESS`, `value="` + d.SecondaryEmail + `"`},
			/* 09 */ {"adminView validate Home Street Addr", &s, `name="HomeStreetAddress"`, `name="HomeStreetAddress2"`, `value="` + d.HomeStreetAddress + `"`},
			/* 10 */ {"adminView validate Home Street Addr2", &s, `name="HomeStreetAddress2"`, `>CITY<`, `value="` + d.HomeStreetAddress2 + `"`},
			/* 11 */ {"adminView validate Home City", &s, `name="HomeCity"`, `name="HomeState"`, `value="` + d.HomeCity + `"`},
			/* 12 */ {"adminView validate Home State", &s, `name="HomeState"`, `name="HomePostalCode"`, `value="` + d.HomeState + `"`},
			/* 13 */ {"adminView validate Home Postal Code", &s, `name="HomePostalCode"`, `name="HomeCountry"`, `value="` + d.HomePostalCode + `"`},
			/* 14 */ {"adminView validate Home Country", &s, `name="HomeCountry"`, `EMERGENCY CONTACT NAME`, `value="` + d.HomeCountry + `"`},
			/* 15 */ {"adminView validate EmergencyContactName", &s, `name="EmergencyContactName"`, `name="EmergencyContactPhone"`, `value="` + d.EmergencyContactName + `"`},
			/* 16 */ {"adminView validate EmergencyContactPhone", &s, `name="EmergencyContactPhone"`, `>COMPANY<`, `value="` + d.EmergencyContactPhone + `"`},
			/* 17 */ {"adminView validate Company", &s, `>COMPANY<`, `>JOB TITLE<`, fmt.Sprintf("option value=\"%d\"selected>", d.CoCode)},
			/* 18 */ {"adminView validate JobTitle", &s, `>JOB TITLE<`, `>MANAGER UID<`, fmt.Sprintf("option value=\"%d\"selected>", d.JobCode)},
			/* 19 */ {"adminView validate ManagerUID", &s, `>MANAGER UID<`, `>STATE OF EMPLOYMENT<`, `value="` + fmt.Sprintf("%d", d.MgrUID) + `"`},
			/* 20 */ {"adminView validate StateOfEmployment", &s, `>STATE OF EMPLOYMENT<`, `>COUNTRY OF EMPLOYMENT<`, `value="` + d.StateOfEmployment + `"`},
			/* 21 */ {"adminView validate CountryOfEmployment", &s, `>COUNTRY OF EMPLOYMENT<`, `>DEPARTMENT<`, `value="` + d.CountryOfEmployment + `"`},
			/* 22 */ {"adminView validate Department", &s, `>DEPARTMENT<`, `>CLASS<`, fmt.Sprintf("option value=\"%d\"selected>", d.DeptCode)},
			/* 23 */ {"adminView validate Class", &s, `>CLASS<`, `>POSITION CONTROL NUMBER`, fmt.Sprintf("option value=\"%d\"selected>", d.ClassCode)},
			/* 24 */ {"adminView validate PositionControlNumber", &s, `>POSITION CONTROL NUMBER`, `>HIRE DATE<`, fmt.Sprintf("value=\"%s\"", d.PositionControlNumber)},
			/* 25 */ {"adminView validate Hire Date", &s, `>HIRE DATE<`, `>STATUS<`, fmt.Sprintf("value=\"%s\"", sHire)},
			/* 26 */ {"adminView validate Status", &s, `>STATUS<`, `>ELIGIBLE FOR REHIRE`, fmt.Sprintf("option value=\"%s\" selected>", activeToString(d.Status))},
			/* 27 */ {"adminView validate EligibleForRehire", &s, `>ELIGIBLE FOR REHIRE`, `>LAST REVIEW<`, fmt.Sprintf("option value=\"%s\" selected>", yesnoToString(d.EligibleForRehire))},
			/* 28 */ {"adminView validate LastReview", &s, `>LAST REVIEW<`, `>NEXT REVIEW<`, fmt.Sprintf("value=\"%s\"", sLastReview)},
			/* 29 */ {"adminView validate NextReview", &s, `>NEXT REVIEW<`, `>TERMINATION DATE<`, fmt.Sprintf("value=\"%s\"", sNextReview)},
			/* 30 */ {"adminView validate Termination", &s, `>TERMINATION DATE<`, `ACCEPTED HEALTH INSURANCE `, fmt.Sprintf("value=\"%s\"", sTermination)},
			/* 31 */ {"adminView validate AcceptedHealthInsurance", &s, `>ACCEPTED HEALTH INSURANCE `, `>ACCEPTED DENTAL INSURANCE `, fmt.Sprintf("value=\"%s\" selected", acceptIntToString(d.AcceptedHealthInsurance))},
			/* 32 */ {"adminView validate AcceptedDentalInsurance", &s, `>ACCEPTED DENTAL INSURANCE `, `>ACCEPTED 401K `, fmt.Sprintf("value=\"%s\" selected", acceptIntToString(d.AcceptedDentalInsurance))},
			/* 33 */ {"adminView validate Accepted401K", &s, `>ACCEPTED 401K `, `>COMPENSATION`, fmt.Sprintf("value=\"%s\" selected", acceptIntToString(d.Accepted401K))},
			/* 34 */ {"adminView validate birthDOM", &s, `>BIRTHDAY<`, `action="/adminViewBtn/`, fmt.Sprintf(`name="BirthDOM" value="%d"`, d.BirthDOM)},
		}

		//-----------------------------------------------------------
		// add validation entries to the table for compensation...
		//-----------------------------------------------------------
		for i := 0; i < len(d.MyComps); i++ {
			h := ""
			if d.MyComps[i].HaveIt > 0 {
				h = " checked"
			}
			v := validationTable{"adminView validate compensation." + d.MyComps[i].Name, &s, `>COMPENSATION`, `>DEDUCTIONS`, fmt.Sprintf(`name="%s" value="%d"%s>`, d.MyComps[i].Name, d.MyComps[i].CompCode, h)}
			validate = append(validate, v)
		}

		//-----------------------------------------------------------
		// add deduction entries to the table for compensation...
		//-----------------------------------------------------------
		for i := 0; i < len(d.MyDeductions); i++ {
			h := ""
			if d.MyDeductions[i].HaveIt > 0 {
				h = " checked"
			}
			v := validationTable{"adminView validate deductions." + d.MyDeductions[i].Name, &s, `>DEDUCTIONS`, `>BIRTHDAY<`, fmt.Sprintf(`name="%s" value="%d"%s>`, d.MyDeductions[i].Name, d.MyDeductions[i].DCode, h)}
			validate = append(validate, v)
		}

		//-----------------------------------------------------------
		// check birthmonth if present...
		//-----------------------------------------------------------
		if d.BirthMonth > 0 {
			v := validationTable{"adminView validate birth month", &s, `>BIRTHDAY<`, `action="/adminViewBtn/`, fmt.Sprintf(`value="%d" selected>`, d.BirthMonth)}
			validate = append(validate, v)
		}

		executeValSubstrTests(&validate, &tr)
		if tr.Fail > 0 {
			dumpTestErrors(&tr)
		}

		return (tr.Fail == 0)
	}
	return false
}
