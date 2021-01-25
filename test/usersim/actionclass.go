package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func getClassInfo(classcode int64, c *class) {
	// s := fmt.Sprintf("select classcode,Name,Designation,Description from classes where classcode=%d", classcode)
	rows, err := App.prepstmt.classInfo.Query(classcode)
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		errcheck(rows.Scan(&c.ClassCode, &c.CoCode, &c.Name, &c.Designation, &c.Description))
	}
	errcheck(rows.Err())
}

var classHTMLEscaper = strings.NewReplacer(
	`&`, "&amp;",
	`'`, "&#39;", // "&#39;" is shorter than "&apos;" and apos was not in HTML until HTML5.
	`<`, "&lt;",
	`>`, "&gt;",
	`"`, "&#34;", // "&#34;" is shorter than "&quot;".
	`+`, "&#43;",
)

var sessionCookieName = string("air")

func escapeString(s string) string {
	return classHTMLEscaper.Replace(s)
}

// viewClass executes the server command to serve a Class detail page and validates the
// data in the HTML returned.
// RETURNS:
//		true = all data verified correctly
//     false = the test failed for any on of several reasons. If the session is not established
//             after the request, the test fails. If the request succeeds, and one or more of
//             the data fields were not correct then the test fails.
func viewClass(d *personDetail, atr *TestResults) bool {
	// to save confusion, any user doing tests will only read/modify the Class
	// that matches their UID.  That way, when multiple users are editing and modifying
	// companies, no user will stomp on the work being done by another user
	URL := fmt.Sprintf("http://%s:%d/class/%d", App.Host, App.Port, d.UID)
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
	// fmt.Printf("viewClass: Adding session cookie: %#v\n", d.SessionCookie)
	req.AddCookie(d.SessionCookie)
	resp, err := hc.Do(req)
	errcheck(err)
	defer resp.Body.Close()

	// fmt.Printf("viewClass: hc.Do(req) returned err = %v\n", err)

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
	if strings.Contains(m, ProductName+" - Business Unit") && d.SessionCookie != nil {
		pageName := "Business Unit Detail"
		var c class
		getClassInfo(d.UID, &c) // yes, we're getting the Class with cocode == d.UID

		validate := []validationTable{
			/* 00 */ {pageName + ": validate ClassCode", &s, `class="LastName"`, `DESIGNATION`, true, c.Name},
			/* 02 */ {pageName + ": validate Designation", &s, `DESIGNATION`, `DESCRIPTION`, true, c.Designation},
			/* 03 */ {pageName + ": validate Description", &s, `DESCRIPTION`, `value="Done"`, true, c.Description},
		}
		var tc testContext
		tc.d = d
		tc.cl = &c
		executeValSubstrTests(&validate, &tr, &tc)
		if tr.Fail > 0 {
			fmt.Printf("Class uid = %d\n", d.UID)
			fmt.Printf("tr.Fail = %d, len(tr.Failures)=%d\n", tr.Fail, len(tr.Failures))
			dumpTestErrors(&tr)
		}
		aggregateTR(&tr, atr)
		return (tr.Fail == 0)
	}
	return false
}

// adminEditClass executes the server command to serve a Class adminEditCo page and validates the
// data in the HTML returned.
// RETURNS:
//		true = all data verified correctly
//     false = the test failed for any on of several reasons. If the session is not established
//             after the request, the test fails. If the request succeeds, and one or more of
//             the data fields were not correct then the test fails.
func adminEditClass(d *personDetail, atr *TestResults) bool {
	// to save confusion, any user doing tests will only read/modify the Class
	// that matches their UID.  That way, when multiple users are editing and modifying
	// companies, no user will stomp on the work being done by another user
	URL := fmt.Sprintf("http://%s:%d/adminEditClass/%d", App.Host, App.Port, d.UID)
	hc := http.Client{}
	pageName := "AdminEditClass"

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
	// fmt.Printf("viewClass: Adding session cookie: %#v\n", d.SessionCookie)
	req.AddCookie(d.SessionCookie)
	resp, err := hc.Do(req)
	errcheck(err)
	defer resp.Body.Close()

	// fmt.Printf("viewClass: hc.Do(req) returned err = %v\n", err)

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
	if strings.Contains(m, ProductName+" - Edit Class") && d.SessionCookie != nil {
		var c class
		getClassInfo(d.UID, &c) // yes, we're getting the Class with cocode == d.UID

		topline := fmt.Sprintf("Admin Edit - %s (%d)", c.Designation, c.ClassCode)
		validate := []validationTable{
			/* 00 */ {pageName + ": validate top line", &s, `class="AppHeading"`, `<form action="/saveAdminEditClass/`, true, topline},
			/* 01 */ {pageName + ": validate Name", &s, `name="Name"`, `name="Designation"`, false, `value="` + c.Name},
			/* 03 */ {pageName + ": validate Designation", &s, `name="Designation"`, `>DESCRIPTION<`, false, `value="` + c.Designation},
			/* 04 */ {pageName + ": validate Description", &s, `>DESCRIPTION<`, `value="Save"`, true, c.Description},
		}

		var tc testContext
		tc.d = d
		tc.cl = &c
		executeValSubstrTests(&validate, &tr, &tc)
		if tr.Fail > 0 {
			fmt.Printf("Class uid = %d\n", d.UID)
			fmt.Printf("tr.Fail = %d, len(tr.Failures)=%d\n", tr.Fail, len(tr.Failures))
			dumpTestErrors(&tr)
		}
		aggregateTR(&tr, atr)
		return (tr.Fail == 0)
	}
	return false
}
