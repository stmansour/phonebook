package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

// Behavior describes a "behavior" of the virtual user.
// That is, a "habit" defining something the user does
type Behavior struct {
	Op     string // the operation performed
	Chance int    // number from 1 to 100 representing the percentage odds of doing this operation
	delay  int    // random number of seconds before executing the next command
}

// Profile is a named collection of Behaviors.  A profile describes
// how a virtual user will utilize the system.
type Profile struct {
	Name      string
	Behaviors []Behavior
}

// TestFailure provides a bit of detail about any test that fails...
// its name and table index as appropriate
type TestFailure struct {
	Name  string
	Index int
}

// TestResults is a container for the number of passed and failed tests
type TestResults struct {
	SimUserID int           // the simulation uses this user id
	Pass      int           // number of tests that passed
	Fail      int           // number of tests that failed
	Failures  []TestFailure // more info about failures
}

// Tester profile does everything that Phonebook can do
var Tester Profile

// Regular Expressions for parsing replies
var reTitle = regexp.MustCompile("<title>")
var reTitleEnd = regexp.MustCompile("</title>")

func initProfiles() {
	Tester.Name = "Tester"
	Tester.Behaviors = []Behavior{{"search", 80, 5},
		{"detail", 10, 10},
		{"searchco", 2, 4},
		{"company", 1, 4},
		{"searchcl", 2, 10},
		{"class", 1, 5},
		{"weblogin", 2, 2},
		{"logoff", 2, 2},
	}
}

// logoff the supplied personDetail
//    returns true if login was successful
//            false if login failed
func logoff(d *personDetail) bool {
	URL := fmt.Sprintf("http://%s:%d/logoff/", App.Host, App.Port)
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
	if strings.Contains(m, "Accord") && strings.Contains(m, "Sign In") && d.SessionCookie == nil {
		// fmt.Printf("Logoff successful\n")
		return true
	}
	return false
}

// login the supplied personDetail
//    returns true if login was successful
//            false if login failed
func login(d *personDetail) bool {
	URL := fmt.Sprintf("http://%s:%d/weblogin/", App.Host, App.Port)
	hc := http.Client{}

	form := url.Values{}
	form.Add("username", d.UserName)
	form.Add("password", "accord")
	// req.PostForm = form
	// req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	req, err := http.NewRequest("POST", URL, bytes.NewBufferString(form.Encode()))
	errcheck(err)

	hdrs := []KeyVal{
		{"Host:", "localhost:8250"},
		{"Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"},
		{"Accept-Encoding", "gzip, deflate"},
		{"Accept-Language", "en-US,en;q=0.8"},
		{"Cache-Control", "max-age=0"},
		{"Connection", "keep-alive"},
		{"Content-Type", "application/x-www-form-urlencoded"},
		{"Origin", "http://localhost:8250"},
		{"Referer", "http://localhost:8250/signin/"},
		{"Upgrade-Insecure-Requests", "1"},
		{"User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.80 Safari/537.36"},
	}
	for i := 0; i < len(hdrs); i++ {
		req.Header.Add(hdrs[i].key, hdrs[i].value)
	}
	// if 1 > 0 {
	// 	fmt.Printf("DumpRequest:\n")
	// 	dump, err := httputil.DumpRequest(req, false)
	// 	errcheck(err)
	// 	fmt.Printf("\n\ndumpRequest = %s\n", string(dump))
	// }

	resp, err := hc.Do(req)
	if nil != err {
		fmt.Printf("login:  hc.Do(req) returned error:  %#v\n", err)
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

	// dump headers...
	// fmt.Printf("Headers:\n")
	// for k, v := range resp.Header {
	// 	fmt.Println("key:", k, "value:", v)
	// }

	// cookies:
	cookies := resp.Cookies()
	// fmt.Printf("Cookies:value: %+v\n", cookies)
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

	// Verify that we were sent to the search page...
	htmlData, err := ioutil.ReadAll(reader)
	errcheck(err)
	s := string(htmlData)
	m1 := reTitle.FindStringIndex(s)
	m2 := reTitleEnd.FindStringIndex(s)
	m := s[m1[1]:m2[0]]
	// fmt.Printf("Page returned = %s\n", m)
	if strings.Contains(m, "Phonebook") && strings.Contains(m, "Search") && d.SessionCookie.Name == "accord" {
		// fmt.Printf("Login successful\n")
		return true
	}
	return false
}

func testResult(testname string, success bool, tr *TestResults) bool {
	if success {
		tr.Pass++
	} else {
		tr.Fail++
	}
	return success
}

func usersim(userindex, iterations, duration int, TestResChan chan TestResults, TestResChanAck chan int) {
	v := App.Peeps[userindex]
	tr := TestResults{v.UID, 0, 0, nil}

	if duration == 0 {
		for i := 0; i < iterations; i++ {
			if v.SessionCookie == nil {
				testResult("login", login(v), &tr)
			}

			if nil == v.SessionCookie {
				fmt.Printf("usersim: could not find accord cookie after login!\n")
				break
			}

			testResult("detail", viewPersonDetail(v), &tr)

			if nil == v.SessionCookie {
				fmt.Printf("usersim: could not find accord cookie after viewPersonDetail!\n")
				break
			}
			testResult("adminView", adminViewTest(v), &tr)

			if nil == v.SessionCookie {
				fmt.Printf("usersim: could not find accord cookie after adminViewTest!\n")
				break
			}

			testResult("adminEdit", adminEditTest(v), &tr)

			if nil == v.SessionCookie {
				fmt.Printf("usersim: could not find accord cookie after adminEditTest!\n")
				break
			}

			// testResult("saveAdminEdit", saveAdminEdit(v), &tr)

			// if nil == v.SessionCookie {
			// 	fmt.Printf("usersim: could not find accord cookie after saveAdminEdit!\n")
			// 	break
			// }

			if v.SessionCookie != nil {
				testResult("logoff", logoff(v), &tr)
			} else {
				fmt.Printf("v.SessionCookie was nil\n")
			}
		}
	}

	TestResChan <- tr // push our results to the simulation executor
	<-TestResChanAck  // wait for receipt before continuing
}

func executeSimulation() {
	StartTime := time.Now()
	TestResChan := make(chan TestResults) // usersim reports results via this struct
	TestResChanAck := make(chan int)      // ack receipt

	if App.TestDuration == 0 {
		for j := 0; j < App.TestUsers; j++ {
			go usersim(j, App.TestIterations, App.TestDuration, TestResChan, TestResChanAck)
		}
	}

	var totTR TestResults                // net results
	for i := 0; i < App.TestUsers; i++ { // i is the number of usersims completed
		select {
		case tr := <-TestResChan: // get the data the usersim collected
			totTR.Fail += tr.Fail // update cumulative totals
			totTR.Pass += tr.Pass // update cumulative totals
			TestResChanAck <- 1   // acknowledge receipt
		}
	}

	fmt.Printf("Total Tests: %d   pass: %d   fail: %d\n", totTR.Fail+totTR.Pass, totTR.Pass, totTR.Fail)
	if len(totTR.Failures) > 0 {
		for i := 0; i < len(totTR.Failures); i++ {
			fmt.Printf("%d. %s[%d] \n", i, totTR.Failures[i].Name, totTR.Failures[i].Index)
		}
	}
	Elapsed := time.Since(StartTime)
	fmt.Printf("Simulation Time: %s\n", Elapsed /*Round(Elapsed, 0.5e9)*/)
}
