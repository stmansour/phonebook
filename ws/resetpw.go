package ws

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"phonebook/lib"
	"strings"

	"gopkg.in/gomail.v2"
)

// ResetPWDomains is the list of domains supported by the
// password reset comma nd
// var supportedDomains = []string{
// 	"accordinterests.com",
// 	"l-objet.com",
// 	"myisolabella.com",
// }

// ResetPWData is the struct with the username and password
// used for authentication
type ResetPWData struct {
	Username string `json:"username"`
}

// SvcResetPWHandler resets the password for the supplied user
// cannot get logged in
func SvcResetPWHandler(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	var err error
	funcname := "SvcResetPWHandler"
	lib.Console("Entered %s\n", funcname)

	var foo ResetPWData

	data := []byte(d.data)
	if err = json.Unmarshal(data, &foo); err != nil {
		e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
		SvcErrorReturn(w, e, funcname)
		return
	}

	lib.Console("Username = %s\n", foo.Username)
	myusername := strings.ToLower(foo.Username)

	//-------------------------------------
	// validate that myusername exists
	//-------------------------------------
	var PrimaryEmail string
	q := fmt.Sprintf("SELECT PrimaryEmail FROM people WHERE UserName=%q", myusername)
	err = SvcCtx.db.QueryRow(q).Scan(&PrimaryEmail)

	switch {
	case err == sql.ErrNoRows:
		err = fmt.Errorf("Username %s was not found", myusername)
		SvcErrorReturn(w, err, funcname)
		return
	case err != nil:
		err = fmt.Errorf("Error retrieving information for user %s: %s", myusername, err.Error())
		SvcErrorReturn(w, err, funcname)
		return
	}
	if PrimaryEmail == "" {
		err = fmt.Errorf("Error: No email address for user: %s", myusername)
		SvcErrorReturn(w, err, funcname)
		return
	}

	//-------------------------------------
	// validate domain
	//-------------------------------------
	errmsg := ""
	domain := ""
	k := strings.LastIndex(PrimaryEmail, "@")
	if k > 0 {
		domain = PrimaryEmail[k+1:]
	}
	found := false
	for i := 0; i < len(lib.AppConfig.ResetPWList); i++ {
		if domain == lib.AppConfig.ResetPWList[i] {
			found = true
			break
		}
	}
	if !found {
		e := fmt.Errorf("Error: %s is not a supported domain for automatic password reset", domain)
		SvcErrorReturn(w, e, funcname)
		return
	}

	//-------------------------------------
	// reset the password for myusername
	//-------------------------------------
	password := lib.RandPasswordStringRunes(8)
	err = lib.UpdateUserPassword(myusername, password, SvcCtx.db)
	if nil != err {
		e := fmt.Errorf("Error updating password: %s", err.Error())
		SvcErrorReturn(w, e, funcname)
		return
	}

	//------------------------------------------------------------------------------
	// send an email to the associated account that the password has been changed.
	//------------------------------------------------------------------------------
	m := gomail.NewMessage()
	m.SetHeader("From", "sman@accordinterests.com")
	m.SetHeader("To", PrimaryEmail)
	msg := fmt.Sprintf("Hello %s,<br><br>Your Accord password has been reset to:  %s<br><br>", myusername, password)
	m.SetHeader("Subject", "Your Accord password has been updated")
	m.SetBody("text/html", msg)
	if err := lib.SMTPDialAndSend(m); err != nil {
		errmsg += fmt.Sprintf("Error sending PrimaryEmail = %s", err.Error())
		e := fmt.Errorf("Error  sending email to %s: %s", PrimaryEmail, err.Error())
		SvcErrorReturn(w, e, funcname)
		return
	}

	SvcWriteSuccessResponse(w)
}
