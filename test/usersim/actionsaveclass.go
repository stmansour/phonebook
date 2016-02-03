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
func saveAdminEditClass(d *personDetail, atr *TestResults) bool {
	URL := fmt.Sprintf("http://%s:%d/saveAdminEditClass/%d", App.Host, App.Port, d.UID)
	var c class
	hc := http.Client{}
	getClassInfo(d.UID, &c)

	form := url.Values{}
	// set new values
	c.Name = App.RandClasses[rand.Intn(len(App.RandClasses))]
	c.Designation = genDesignation(c.Name, "classcode", "classes")
	c.Description = RandomString(14)

	//===================================================
	// Simulate filling in the fields...
	//===================================================
	form.Add("Name", c.Name)
	form.Add("Designation", c.Designation)
	form.Add("Description", c.Description)

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
		fmt.Printf("saveAdminEditClass:  Server return non-200 status: %v\n", resp.Status)
	}

	//==============================================================
	// now read in the updated version
	//==============================================================
	var cnew class
	getClassInfo(d.UID, &cnew)
	res := c.matches(&cnew, atr)

	return res
}
