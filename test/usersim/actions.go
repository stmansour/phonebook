package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

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
	if nil == m1 || nil == m2 {
		return false
	}
	// fmt.Printf("m1 = %#v\n", m1)
	// fmt.Printf("m2 = %#v\n", m2)
	m := (*ps)[m1[1]:m2[0]]
	// fmt.Printf("m = %s\n", m)
	return strings.Contains(m, target)
}

// Get a copy of the user's data, modify it, and validate that the changes were saved
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
	// fmt.Printf("Cookies:value: %+v\n", cookies)
	d.SessionCookie = nil
	for i := 0; i < len(cookies); i++ {
		if cookies[i].Name == "accord" {
			d.SessionCookie = cookies[i]
			break
		}
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
	if strings.Contains(m, "Accord") && strings.Contains(m, "Details") && d.SessionCookie != nil {
		myname := d.FirstName + " " + d.MiddleName + " " + d.LastName
		if !validateSubstring(&s, "FULL NAME", "EMAIL", myname) {
			return false
		}
		if !validateSubstring(&s, "PHONE", "CELL", d.OfficePhone) {
			return false
		}
		return true
	}
	return false
}
