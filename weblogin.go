package main

import (
	"crypto/md5"
	"crypto/sha512"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"phonebook/lib"
	"phonebook/sess"
	"strings"
	"time"

	"gopkg.in/gomail.v2"
)

func handlerInitUIDate(ui *uiSupport) {
	Phonebook.ReqMem <- 1    // ask to access the shared mem, blocks until granted
	<-Phonebook.ReqMemAck    // make sure we got it
	initUIData(ui)           // initialize our data
	Phonebook.ReqMemAck <- 1 // tell Dispatcher we're done with the data
}

// initHandlerSession validates the session cookie and redirects if necessary.
// it also initializes the uiSession variable
// RETURNS:  0 = no problems
//           1 = redirected
func initHandlerSession(ssn *sess.Session, ui *uiSupport, w http.ResponseWriter, r *http.Request) int {
	var ok bool
	cookie, err := r.Cookie("accord")
	if nil != cookie && err == nil {
		ssn, ok = sessionGet(cookie.Value)
		if !ok || ssn == nil {
			http.Redirect(w, r, "/signin/", http.StatusFound)
			return 1
		}
		ssn.Refresh(w, r)
	} else {
		//fmt.Printf("REDIRECT to signin\n")
		http.Redirect(w, r, "/signin/", http.StatusFound)
		return 1
	}
	ui.X = ssn
	handlerInitUIDate(ui)
	return 0
}

func webloginHandler(w http.ResponseWriter, r *http.Request) {

	// debug only
	// dump, err := httputil.DumpRequest(r, false)
	// errcheck(err)
	// fmt.Printf("\n\ndumpRequest = %s\n", string(dump))

	resetpw := r.FormValue("lostpw")
	if resetpw == "resetpw" {
		resetpwHandler(w, r)
		return
	}

	Phonebook.ReqCountersMem <- 1    // ask to access the shared mem, blocks until granted
	<-Phonebook.ReqCountersMemAck    // make sure we got it
	Counters.SignIn++                // initialize our data
	Phonebook.ReqCountersMemAck <- 1 // tell Dispatcher we're done with the data

	n := 0 //error number associated with this login attempt
	loggedIn := false
	myusername := strings.ToLower(r.FormValue("username"))
	password := []byte(r.FormValue("password"))
	sha := sha512.Sum512(password)
	mypasshash := fmt.Sprintf("%x", sha)
	email := ""

	var passhash, firstname, preferredname string
	var uid, RID int
	err := Phonebook.prepstmt.loginInfo.QueryRow(myusername).Scan(&uid, &firstname, &preferredname, &email, &passhash, &RID)
	switch {
	case err == sql.ErrNoRows:
		ulog("No user with username = %s\n", myusername)
		n = 1
	case err != nil:
		ulog("login username: %s,  error = %v\n", myusername, err)
		n = 2
	default:
		// ulog("found username %s in database. UID = %d\n", myusername, uid)
	}

	if passhash == mypasshash {
		loggedIn = true
		ulog("user %s logged in\n", myusername)
		expiration := time.Now().Add(10 * time.Minute)
		//=================================================================================
		// There could be multiple ssn.Sessions from the same user on different browsers.
		// These could be on the same or separate machines. We need the IP and the browser
		// to guarantee uniqueness...
		//=================================================================================
		key := myusername + r.Header.Get("User-Agent") + r.RemoteAddr
		cval := fmt.Sprintf("%x", md5.Sum([]byte(key)))
		name := firstname
		if len(preferredname) > 0 {
			name = preferredname
		}

		s := sess.NewSession(cval, myusername, name, uid, RID)
		cookie := http.Cookie{Name: "accord", Value: s.Token, Expires: expiration}
		cookie.Path = "/"
		http.SetCookie(w, &cookie)
		r.AddCookie(&cookie) // need this so that the redirect to search finds the cookie
	} else {
		ulog("user name or password did not match for: %s\n", myusername)
		n = 1
	}

	if !loggedIn {
		http.Redirect(w, r, fmt.Sprintf("/signin/%d", n), http.StatusFound)
	} else {
		// http.Redirect(w, r, "/search/", http.StatusFound)
		searchHandler(w, r) // redirect loses the cookie, but this seems to work just fine
	}
}

func showResetPwPage(w http.ResponseWriter, r *http.Request, errmsg string) {
	t, _ := template.New("resetpw.html").Funcs(funcMap).ParseFiles("resetpw.html")
	var ui uiSupport
	handlerInitUIDate(&ui)
	var ssn sess.Session
	ssn.Username = r.FormValue("username")
	ui.X = &ssn
	ui.ErrMsg = template.HTML(errmsg)
	err := t.Execute(w, &ui)
	if nil != err {
		errmsg := fmt.Sprintf("signinHandler: err = %v\n", err)
		ulog(errmsg)
		fmt.Println(errmsg)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var supportedDomains = []string{
	"accordinterests.com",
	"l-objet.com",
	"myisolabella.com",
}
var stillNeedHelp = string(`<a href="#" onclick="return alert('For assistance with your username or password please contact Steve Mansour:\nemail: sman@accordinterests.com')">Still need help?</a>`)

// Note that there is no session when this handler is called. The user
// cannot get logged in
func resetpwHandler(w http.ResponseWriter, r *http.Request) {
	var firstname, preferredname, emailAddr, passhash string
	var uid, RID int
	var err error

	pagename := r.FormValue("pagename")

	if pagename != "resetpw" { // if the resetpw page was not the calling page...
		showResetPwPage(w, r, "") // ...then show the resetpw page
		return
	}

	myusername := strings.ToLower(r.FormValue("username"))

	//-------------------------------------
	// validate that myusername exists
	//-------------------------------------
	err = Phonebook.prepstmt.loginInfo.QueryRow(myusername).Scan(&uid, &firstname, &preferredname, &emailAddr, &passhash, &RID)
	switch {
	case err == sql.ErrNoRows:
		errmsg := fmt.Sprintf("Username %s was not found\n", myusername) + stillNeedHelp
		showResetPwPage(w, r, errmsg)
		return
	case err != nil:
		errmsg := fmt.Sprintf("Error: %s", err.Error()) + stillNeedHelp
		showResetPwPage(w, r, errmsg)
		return
	}
	if emailAddr == "" {
		errmsg := fmt.Sprintf("Error: No email address for user: %s", myusername) + stillNeedHelp
		showResetPwPage(w, r, errmsg)
		return
	}

	//-------------------------------------
	// validate domain
	//-------------------------------------
	errmsg := ""
	domain := ""
	k := strings.LastIndex(emailAddr, "@")
	if k > 0 {
		domain = emailAddr[k+1:]
	}
	found := false
	for i := 0; i < len(supportedDomains); i++ {
		if domain == supportedDomains[i] {
			found = true
			break
		}
	}
	if !found {
		errmsg += fmt.Sprintf("Error: %s is not a supported domain for automatic password reset.\n", domain)
		errmsg += stillNeedHelp
		showResetPwPage(w, r, errmsg)
		return
	}

	//-------------------------------------
	// reset the password for myusername
	//-------------------------------------
	password := lib.RandPasswordStringRunes(8)
	err = lib.UpdateUserPassword(myusername, password, Phonebook.db)
	if nil != err {
		errmsg += fmt.Sprintf("Error updating password = %s\n", err.Error())
		showResetPwPage(w, r, errmsg)
		return
	}

	//------------------------------------------------------------------------------
	// send an email to the associated account that the password has been changed.
	//------------------------------------------------------------------------------
	m := gomail.NewMessage()
	m.SetHeader("From", "sman@accordinterests.com")
	ulog("To address is set to: \"%s\"\n", emailAddr)
	m.SetHeader("To", emailAddr)
	msg := fmt.Sprintf("Hello %s,<br><br>Your password has been set to:  %s<br><br>", myusername, password)
	msg += `Please log into <a href="https://directory.airoller.com/">https://directory.airoller.com/</a>`
	m.SetHeader("Subject", "Your password has been updated")
	m.SetBody("text/html", msg)
	if err := lib.SMTPDialAndSend(m); err != nil {
		errmsg += fmt.Sprintf("Error sending emailAddr = %s", err.Error())
	}

	//-------------------------------------
	// notify user
	//-------------------------------------
	t, _ := template.New("pwreset.html").Funcs(funcMap).ParseFiles("pwreset.html")
	var ui uiSupport
	handlerInitUIDate(&ui)
	var ssn sess.Session
	ssn.Username = myusername
	ui.X = &ssn
	ui.ErrMsg = template.HTML(errmsg)
	err = t.Execute(w, &ui)
	if nil != err {
		errmsg := fmt.Sprintf("signinHandler: err = %v\n", err)
		ulog(errmsg)
		fmt.Println(errmsg)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
