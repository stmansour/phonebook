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
	testname   string  // name of test
	html       *string // full html text
	reStart    string  // regexp for start of substring
	reStop     string  // regexp for end of substring
	trgtIsHTML bool    // is the target string encoded html?
	target     string  // what we expect to find in the substring
}

// validateSubstring searches for the substring between s1 and s2.
// Returns:
//		true - target string was found in the defined substring
//		false - could mean any of several things:
//					* string s1 was not found
//					* string s2 was not found
//					* target was not found in the string between s1 & s2
func validateSubstring(ps *string, s1 string, s2 string, isHTML bool, target string, reason *string) bool {
	reg1 := regexp.MustCompile(s1)
	reg2 := regexp.MustCompile(s2)
	m1 := reg1.FindStringIndex(*ps)
	m2 := reg2.FindStringIndex(*ps)

	if nil == m1 || nil == m2 {
		*reason = fmt.Sprintf("s1=%s  s2=%s  target=%s\nm1 = %#v\nm2 = %#v\n", s1, s2, target, m1, m2)
		if App.ShowTestMatching {
			fmt.Println(*reason)
		}
		return false
	}
	if m2[0] < m1[1] {
		*reason = fmt.Sprintf("s2 has a bad regexp. A match occurs before s1, at index %d\ns1=%s  s2=%s  target=%s\nm1 = %#v\nm2 = %#v\n", m2[0], s1, s2, target, m1, m2)
		fmt.Println(*reason)
		return false
	}
	m := (*ps)[m1[1]:m2[0]]
	b := strings.Contains(m, target)
	if !b {
		// fmt.Printf("validateSubstring: %s  <==>  %s\n", m, target)
		if isHTML {
			// the HTML that is placed in the templates has been unclear
			// with respect to html encoding. Sometimes it is, sometimes it's not.
			// We'll check it both ways.
			// fmt.Printf("target: %s\nmatch: %s\nencoded target: %s\n", target, m, escapeString(target))
			b = strings.Contains(m, escapeString(target))
		}
		if !b {
			*reason = fmt.Sprintf("s1=%s  s2=%s  target=%s\nsubstring = \"%s\"\ncould not find = \"%s\"\n",
				s1, s2, target, m, target)
			if App.ShowTestMatching {
				fmt.Println(*reason)
			}
		}
	}
	return b
}

// executeValSubstrTests accepts a table of data check tests, executes each check,
// and returns the number of checks that passed and the number that failed
func executeValSubstrTests(m *[]validationTable, tr *TestResults, tc *testContext) {
	for i := 0; i < len(*m); i++ {
		var s string
		if validateSubstring((*m)[i].html, (*m)[i].reStart, (*m)[i].reStop, (*m)[i].trgtIsHTML, (*m)[i].target, &s) {
			tr.Pass++
		} else {
			tr.Fail++
			f := TestFailure{
				/*TestName:*/ (*m)[i].testname,
				/*Context:*/ fmt.Sprintf("user:  %s (%d)", tc.d.UserName, tc.d.UID),
				/*Reason:*/ s,
				/*index:*/ i}
			tr.Failures = append(tr.Failures, f)
		}
	}
}

func aggregateTR(Mytr *TestResults, tr *TestResults) {
	tr.Pass += Mytr.Pass
	tr.Fail += Mytr.Fail
	for i := 0; i < len(Mytr.Failures); i++ {
		tr.Failures = append(tr.Failures, Mytr.Failures[i])
	}
}

// viewPersonDetail executes the server command to serve a person detail page and validates the
// data in the HTML returned.
// RETURNS:
//		true = all data verified correctly
//     false = one or more of the data fields were not correct
func viewPersonDetail(d *personDetail, tr *TestResults) bool {
	URL := fmt.Sprintf("http://%s:%d/detail/%d", App.Host, App.Port, d.UID)
	hc := http.Client{}

	req, err := http.NewRequest("GET", URL, nil)
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
	// fmt.Printf("adding session cookie: Expires = %s, d.SessionCookie = %#v,\n", d.SessionCookie.Expires.Format("2006-01-02 15:04:00 MST"), d.SessionCookie)
	req.AddCookie(d.SessionCookie)
	resp, err := hc.Do(req)
	errcheck(err)
	defer resp.Body.Close()

	// fmt.Printf("viewPersonDetail: hc.Do(req) returned err = %v\n", err)

	// Check that the server actually sent compressed data
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		fmt.Printf("gzip response\n")
		reader, err = gzip.NewReader(resp.Body)
		errcheck(err)
		defer reader.Close()
	default:
		reader = resp.Body
	}

	// Look for the cookie...
	cookies := resp.Cookies()
	// fmt.Printf("viewPersonDetail: Cookies:  %+v\n", cookies)
	d.SessionCookie = nil
	for i := 0; i < len(cookies); i++ {
		if cookies[i].Name == sessionCookieName {
			d.SessionCookie = cookies[i]
			break
		}
	}
	if nil == d.SessionCookie {
		fmt.Printf("d.SessionCookie is nil after executing: %s\n", URL)
	}

	// Verify that we were sent to the Sign In page...
	htmlData, err := ioutil.ReadAll(reader)
	errcheck(err)
	s := string(htmlData)

	// fmt.Printf("Response = %s\n", s)

	m1 := reTitle.FindStringIndex(s)
	m2 := reTitleEnd.FindStringIndex(s)
	m := s[m1[1]:m2[0]]
	// fmt.Printf("Page returned = %s\n", m)
	var Mytr TestResults
	tr.Failures = make([]TestFailure, 0)
	if strings.Contains(m, ProductName) && strings.Contains(m, "Details") && d.SessionCookie != nil {
		myname := d.FirstName + " " + d.MiddleName + " " + d.LastName
		validate := []validationTable{
			{"Validate Full Name", &s, "FULL NAME", "PREFERRED NAME", true, myname},
			{"Validate Preferred Name", &s, "PREFERRED NAME", "EMAIL", true, d.PreferredName},
			{"Validate Primary Email", &s, "EMAIL", "PHONE &", true, d.PrimaryEmail},
			{"Validate Office Phone", &s, "PHONE &", "CELL", true, d.OfficePhone},
			{"Validate Cell Phone", &s, "CELL", "BUSINESS UNIT", true, d.CellPhone},
			{"Validate Class", &s, "BUSINESS UNIT", "DEPARTMENT", true, d.Class},
			{"Validate Department", &s, "DEPARTMENT", "MANAGER", true, d.DeptName},
		}
		var tc testContext
		tc.d = d
		executeValSubstrTests(&validate, &Mytr, &tc)
		if Mytr.Fail > 0 {
			fmt.Printf("Failures for %s (%d)\n", d.UserName, d.UID)
			dumpTestErrors(&Mytr)
		}
		aggregateTR(&Mytr, tr)
		return (Mytr.Fail == 0)
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
func viewAdminPerson(d *personDetail, URL string, pageName string, tr *TestResults) bool {
	hc := http.Client{}

	req, err := http.NewRequest("GET", URL, nil)
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
		if cookies[i].Name == sessionCookieName {
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
		errcheck(err)
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

	var Mytr TestResults
	var tc testContext
	Mytr.Failures = make([]TestFailure, 0)

	// lib.Ulog("Page returned = %s\n", m)
	// lib.Ulog("preparing to check for %s and %s", ProductName, pageName)
	if strings.Contains(m, ProductName) && strings.Contains(m, pageName) && d.SessionCookie != nil {
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
			/* 00 */ {pageName + ": validate Last Name", &s, `name="LastName"`, `name="MiddleName"`, false, `value="` + d.LastName + `"`},
			/* 01 */ {pageName + ": validate Middle Name", &s, `name="MiddleName"`, `name="FirstName"`, false, `value="` + d.MiddleName + `"`},
			/* 02 */ {pageName + ": validate First Name", &s, `name="FirstName"`, `name="PreferredName"`, false, `value="` + d.FirstName + `"`},
			/* 03 */ {pageName + ": validate Preferred Name", &s, `name="PreferredName"`, `OFFICE PHONE`, false, `value="` + d.PreferredName + `"`},
			/* 04 */ {pageName + ": validate Office Phone", &s, `name="OfficePhone"`, `name="OfficeFax"`, false, `value="` + d.OfficePhone + `"`},
			/* 05 */ {pageName + ": validate Office Fax", &s, `name="OfficeFax"`, `name="CellPhone"`, false, `value="` + d.OfficeFax + `"`},
			/* 06 */ {pageName + ": validate Cell Phone", &s, `name="CellPhone"`, `PRIMARY EMAIL`, false, `value="` + d.CellPhone + `"`},
			/* 07 */ {pageName + ": validate Primary Email", &s, `name="PrimaryEmail"`, `name="SecondaryEmail"`, false, `value="` + d.PrimaryEmail + `"`},
			/* 08 */ {pageName + ": validate Secondary Email", &s, `name="SecondaryEmail"`, `HOME STREET ADDRESS`, false, `value="` + d.SecondaryEmail + `"`},
			/* 09 */ {pageName + ": validate Home Street Addr", &s, `name="HomeStreetAddress"`, `name="HomeStreetAddress2"`, false, `value="` + d.HomeStreetAddress + `"`},
			/* 10 */ {pageName + ": validate Home Street Addr2", &s, `name="HomeStreetAddress2"`, `>CITY<`, false, `value="` + d.HomeStreetAddress2 + `"`},
			/* 11 */ {pageName + ": validate Home City", &s, `name="HomeCity"`, `name="HomeState"`, false, `value="` + d.HomeCity + `"`},
			/* 12 */ {pageName + ": validate Home State", &s, `name="HomeState"`, `name="HomePostalCode"`, false, `value="` + d.HomeState + `"`},
			/* 13 */ {pageName + ": validate Home Postal Code", &s, `name="HomePostalCode"`, `name="HomeCountry"`, false, `value="` + d.HomePostalCode + `"`},
			/* 14 */ {pageName + ": validate Home Country", &s, `name="HomeCountry"`, `EMERGENCY CONTACT NAME`, false, `value="` + d.HomeCountry + `"`},
			/* 15 */ {pageName + ": validate EmergencyContactName", &s, `name="EmergencyContactName"`, `name="EmergencyContactPhone"`, false, `value="` + d.EmergencyContactName + `"`},
			/* 16 */ {pageName + ": validate EmergencyContactPhone", &s, `name="EmergencyContactPhone"`, `>COMPANY<`, false, `value="` + d.EmergencyContactPhone + `"`},
			/* 17 */ //{pageName + ": validate Company", &s, `>COMPANY<`, `>JOB TITLE<`, false, fmt.Sprintf("option value=\"%d\"selected>", d.CoCode)},
			/* 18 */ {pageName + ": validate JobTitle", &s, `>JOB TITLE<`, `>MANAGER UID<`, false, fmt.Sprintf("option value=\"%d\"selected>", d.JobCode)},
			/* 19 */ {pageName + ": validate ManagerUID", &s, `>MANAGER UID<`, `>STATE OF EMPLOYMENT<`, false, `value="` + fmt.Sprintf("%d", d.MgrUID) + `"`},
			/* 20 */ {pageName + ": validate StateOfEmployment", &s, `>STATE OF EMPLOYMENT<`, `>COUNTRY OF EMPLOYMENT<`, false, `value="` + d.StateOfEmployment + `"`},
			/* 21 */ {pageName + ": validate CountryOfEmployment", &s, `>COUNTRY OF EMPLOYMENT<`, `>DEPARTMENT<`, false, `value="` + d.CountryOfEmployment + `"`},
			/* 22 */ {pageName + ": validate Department", &s, `>DEPARTMENT<`, `>BUSINESS UNIT<`, false, fmt.Sprintf("option value=\"%d\"selected>", d.DeptCode)},
			/* 23 */ {pageName + ": validate Class", &s, `>BUSINESS UNIT<`, `>POSITION CONTROL NUMBER`, false, fmt.Sprintf("option value=\"%d\"selected>", d.ClassCode)},
			/* 24 */ {pageName + ": validate PositionControlNumber", &s, `>POSITION CONTROL NUMBER`, `>HIRE DATE<`, false, fmt.Sprintf("value=\"%s\"", d.PositionControlNumber)},
			/* 25 */ {pageName + ": validate Hire Date", &s, `>HIRE DATE<`, `>STATUS<`, false, fmt.Sprintf("value=\"%s\"", sHire)},
			/* 26 */ {pageName + ": validate Status", &s, `>STATUS<`, `>ELIGIBLE FOR REHIRE`, false, fmt.Sprintf("option value=\"%s\" selected>", activeToString(d.Status))},
			/* 27 */ {pageName + ": validate EligibleForRehire", &s, `>ELIGIBLE FOR REHIRE`, `>LAST REVIEW<`, false, fmt.Sprintf("option value=\"%s\" selected>", yesnoToString(d.EligibleForRehire))},
			/* 28 */ {pageName + ": validate LastReview", &s, `>LAST REVIEW<`, `>NEXT REVIEW<`, false, fmt.Sprintf("value=\"%s\"", sLastReview)},
			/* 29 */ {pageName + ": validate NextReview", &s, `>NEXT REVIEW<`, `>TERMINATION DATE<`, false, fmt.Sprintf("value=\"%s\"", sNextReview)},
			/* 30 */ {pageName + ": validate Termination", &s, `>TERMINATION DATE<`, `ACCEPTED HEALTH INSURANCE `, false, fmt.Sprintf("value=\"%s\"", sTermination)},
			/* 31 */ {pageName + ": validate AcceptedHealthInsurance", &s, `>ACCEPTED HEALTH INSURANCE `, `>ACCEPTED DENTAL INSURANCE `, false, fmt.Sprintf("value=\"%s\" selected", acceptIntToString(d.AcceptedHealthInsurance))},
			/* 32 */ {pageName + ": validate AcceptedDentalInsurance", &s, `>ACCEPTED DENTAL INSURANCE `, `>ACCEPTED 401K `, false, fmt.Sprintf("value=\"%s\" selected", acceptIntToString(d.AcceptedDentalInsurance))},
			/* 33 */ {pageName + ": validate Accepted401K", &s, `>ACCEPTED 401K `, `>COMPENSATION`, false, fmt.Sprintf("value=\"%s\" selected", acceptIntToString(d.Accepted401K))},
			/* 34 */ {pageName + ": validate birthDOM", &s, `>BIRTHDAY<`, `<input type="submit"`, false, fmt.Sprintf(`name="BirthDOM" value="%d"`, d.BirthDOM)},
		}

		//-----------------------------------------------------------
		// add validation entries to the table for compensation...
		//-----------------------------------------------------------
		for i := 0; i < len(d.MyComps); i++ {
			h := ""
			if d.MyComps[i].HaveIt > 0 {
				h = " checked"
			}
			v := validationTable{pageName + ": validate compensation." + d.MyComps[i].Name, &s, `>COMPENSATION`, `>DEDUCTIONS`, false, fmt.Sprintf(`name="%s" value="%d"%s>`, d.MyComps[i].Name, d.MyComps[i].CompCode, h)}
			validate = append(validate, v)
		}

		//-----------------------------------------------------------
		// add deduction entries to the table for deductions...
		//-----------------------------------------------------------
		for i := 0; i < len(d.MyDeductions); i++ {
			h := ""
			if d.MyDeductions[i].HaveIt > 0 {
				h = " checked"
			}
			v := validationTable{pageName + ": validate deductions." + d.MyDeductions[i].Name, &s, `>DEDUCTIONS`, `>BIRTHDAY<`, false, fmt.Sprintf(`name="%s" value="%d"%s>`, d.MyDeductions[i].Name, d.MyDeductions[i].DCode, h)}
			validate = append(validate, v)
		}

		//-----------------------------------------------------------
		// check birthmonth if present...
		//-----------------------------------------------------------
		if d.BirthMonth > 0 {
			v := validationTable{pageName + ": validate birth month", &s, `>BIRTHDAY<`, `<input type="submit"`, false, fmt.Sprintf(`value="%d" selected>`, d.BirthMonth)}
			validate = append(validate, v)
		}

		tc.d = d
		tc.testtype = ELEMPERSON
		executeValSubstrTests(&validate, &Mytr, &tc)
		if Mytr.Fail > 0 {
			fmt.Printf("Mytr.Fail = %d, len(Mytr.Failures)=%d\n", Mytr.Fail, len(Mytr.Failures))
			dumpTestErrors(&Mytr)
		}
		aggregateTR(&Mytr, tr)
		return (Mytr.Fail == 0)
	}
	return false
}

func adminViewTest(d *personDetail, atr *TestResults) bool {
	URL := fmt.Sprintf("http://%s:%d/adminView/%d", App.Host, App.Port, d.UID)
	return viewAdminPerson(d, URL, ProductName+" - Admin View", atr)
}

func adminEditTest(d *personDetail, atr *TestResults) bool {
	URL := fmt.Sprintf("http://%s:%d/adminEdit/%d", App.Host, App.Port, d.UID)
	return viewAdminPerson(d, URL, ProductName+" - Admin Edit", atr)
}
