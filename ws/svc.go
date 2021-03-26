package ws

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"phonebook/db"
	"phonebook/lib"
	"rentroll/rlib"
	"strings"
	"time"
)

// SvcCtx holds global data needed by the service routines
var SvcCtx struct {
	db *sql.DB
}

// GenSearch describes a search condition
type GenSearch struct {
	Field    string `json:"field"`
	Type     string `json:"type"`
	Value    string `json:"value"`
	Operator string `json:"operator"`
}

// ColSort is what the UI uses to indicate how the return values should be sorted
type ColSort struct {
	Field     string `json:"field"`
	Direction string `json:"direction"`
}

// WebGridSearchRequest is a struct suitable for describing a webservice operation.
type WebGridSearchRequest struct {
	Cmd           string      `json:"cmd"`           // get, save, delete
	Limit         int         `json:"limit"`         // max number to return
	Offset        int         `json:"offset"`        // solution set offset
	Selected      []int       `json:"selected"`      // selected rows
	SearchLogic   string      `json:"searchLogic"`   // OR | AND
	Search        []GenSearch `json:"search"`        // what fields and what values
	Sort          []ColSort   `json:"sort"`          // sort criteria
	SearchDtStart time.Time   `json:"searchDtStart"` // for time-sensitive searches
	SearchDtStop  time.Time   `json:"searchDtStop"`  // for time-sensitive searches
	Bool1         bool        `json:"Bool1"`         // a general purpose bool flag for postData from client
	Client        string      `json:"client"`        // name of requesting client.  ex: "roller", "receipts"
	RentableName  string      `json:"RentableName"`  // RECEIPT-ONLY CLIENT EXTENSION - to be removed when Receipt-Only client goes away
}

// WebGridSearchRequestJSON is a struct suitable for describing a webservice operation.
// It is the wire format data. It will be merged into another object where JSONDate values
// are converted to time.Time
type WebGridSearchRequestJSON struct {
	Cmd           string        `json:"cmd"`           // get, save, delete
	Limit         int           `json:"limit"`         // max number to return
	Offset        int           `json:"offset"`        // solution set offset
	Selected      []int         `json:"selected"`      // selected rows
	SearchLogic   string        `json:"searchLogic"`   // OR | AND
	Search        []GenSearch   `json:"search"`        // what fields and what values
	Sort          []ColSort     `json:"sort"`          // sort criteria
	SearchDtStart rlib.JSONDate `json:"searchDtStart"` // for time-sensitive searches
	SearchDtStop  rlib.JSONDate `json:"searchDtStop"`  // for time-sensitive searches
	Bool1         bool          `json:"Bool1"`         // a general purpose bool flag for postData from client
	Client        string        `json:"client"`        // name of requesting client
	RentableName  string        `json:"RentableName"`  // RECEIPT-ONLY CLIENT EXTENSION - to be removed when Receipt-Only client goes away
}

// WebTypeDownRequest is a search call made by a client while the user is
// typing in something to search for and the expecation is that the solution
// set will be sent back in realtime to aid the user.  Search is a string
// to search for -- it's what the user types in.  Max is the maximum number
// of matches to return.
type WebTypeDownRequest struct {
	Search string `json:"search"`
	Max    int    `json:"max"`
}

// ServiceData is the generalized data gatherer for svcHandler. It allows all
// the common data to be centrally parsed and passed to a handler, which may
// need to parse further to get its unique data.  It includes fields for
// common data elements in web svc requests
type ServiceData struct {
	Service       string               // the service requested (position 1)
	DetVal        string               // value of 3rd path element if present (it is not always a number)
	UID           int64                // user id of requester
	ID            int64                // ID associated with the request -- example  /v1/people/123 means get info about person with UID 123
	BID           int64                // business id
	pathElements  []string             // the parts of the uri
	data          string               // the raw unparsed data
	wsSearchReq   WebGridSearchRequest // what did the search requester ask for
	wsTypeDownReq WebTypeDownRequest   // fast for typedown
	QueryParams   map[string][]string  // parameters when HTTP GET is used
	Files         map[string][]*multipart.FileHeader
	MFValues      map[string][]string
	sess          *db.Session
}

// ServiceHandler describes the handler for all services
type ServiceHandler struct {
	Cmd           string
	Handler       func(http.ResponseWriter, *http.Request, *ServiceData)
	AuthNRequired bool
}

// SvcError is the generalized error structure to return errors to the grid widget
type SvcError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// SvcSuccess is the general success return value when no data is required
type SvcSuccess struct {
	Status string `json:"status"`
}

// SvcStatusResponse is the response to return status when no other data
// needs to be returned
type SvcStatusResponse struct {
	Status string `json:"status"` // typically "success"
	Recid  int64  `json:"recid"`  // set to id of newly inserted record
}

// Svcs is the table of all service handlers
var Svcs = []ServiceHandler{
	{"authenticate", SvcAuthenticate, false},
	{"bud", SvcBUD, false},
	{"butd", SvcBUTypedown, false},
	{"discon", SvcDisableConsole, false},
	{"encon", SvcEnableConsole, false},
	{"license", SvcHandlerLicense, true},
	{"licenses", GetLicenses, true},
	{"logoff", SvcLogoff, false},
	{"people", SvcPeople, true},
	{"peopletd", SvcPeopleTypeDown, false},
	{"resetpw", SvcResetPWHandler, true},
	{"validatecookie", SvcValidateCookie, false},
	{"version", SvcHandlerVersion, false},
}

// InitServices initializes the context data needed by service routines
//-----------------------------------------------------------------------------
func InitServices(db *sql.DB) {
	SvcCtx.db = db
}

// initServiceData reads the fields that are common in all uris and loads them
// into the approptiate slots in d (ServiceData)
//
// pathElements:  0   1            2     3
//               /v1/{subservice}/{BUI}/{ID} into an array of strings
// BID is common to nearly all commands
//-----------------------------------------------------------------------------
func initServiceData(w http.ResponseWriter, r *http.Request, d *ServiceData) error {
	var err error
	ss := strings.Split(r.RequestURI[1:], "?") // it could be GET command
	d.pathElements = strings.Split(ss[0], "/")
	d.Service = d.pathElements[1]
	if len(d.pathElements) >= 3 {
		d.BID, err = getBIDfromBUI(d.pathElements[2])
		if err != nil {
			return fmt.Errorf("Could not determine business from %s", d.pathElements[2])
		}
		if d.BID < 0 {
			return fmt.Errorf("Invalid business id: %s", d.pathElements[2])
		}
	}
	if len(d.pathElements) >= 4 {
		d.DetVal = d.pathElements[3]
		d.ID, err = lib.IntFromString(d.DetVal, "bad request integer value") // assume it's a BID
		if err != nil {
			d.ID = 0
		}
	}

	showRequestHeaders(r)
	return nil
}

// V1ServiceHandler is the main dispatch point for WEB SERVICE requests
//
// The expected input is of the form:
//		request=%7B%22cmd%22%3A%22get%22%2C%22selected%22%3A%5B%5D%2C%22limit%22%3A100%2C%22offset%22%3A0%7D
// This is exactly what the w2ui grid sends as a request.
//
// Decoded, this message looks something like this:
//		request={"cmd":"get","selected":[],"limit":100,"offset":0}
//
// The leading "request=" is optional. This routine parses the basic information, then contacts an appropriate
// handler for more detailed processing.  It will set the Cmd member variable.
//
// W2UI sometimes sends requests that look like this: request=%7B%22search%22%3A%22s%22%2C%22max%22%3A250%7D
// using HTTP GET (rather than its more typical POST).  The command decodes to this:
// request={"search":"s","max":250}
//-----------------------------------------------------------------------------------------------------------
func V1ServiceHandler(w http.ResponseWriter, r *http.Request) {
	funcname := "V1ServiceHandler"
	lib.Console("YOU HAVE ENTERED V1-SERVICE-HANDLER\n")

	svcDebugTxn(funcname, r)
	var d ServiceData
	var err error

	//-----------------------------------------------------------------------
	// pathElements:  0   1
	//               /v1/{command}/bid/id
	//
	// ex:           /v1/authenticate/bid
	//               /v1/people/1/201
	//-----------------------------------------------------------------------
	lib.Console("RequestURI = %s\n", r.RequestURI)
	err = initServiceData(w, r, &d)

	lib.Console("%s: 001. r.Method = %s\n", funcname, r.Method)
	switch r.Method {
	case "POST":
		lib.Console("%s: 001-a   calling getPOSTdata\n", funcname)
		if nil != getPOSTdata(w, r, &d) {
			lib.Console("%s: 001-b   returning to client\n", funcname)
			return
		}
	case "GET":
		if nil != getGETdata(w, r, &d) {
			return
		}
	}

	lib.Console("%s: 002 \n", funcname)

	//-----------------------------------------------------------------------
	//  Now call the appropriate handler to do the rest
	//-----------------------------------------------------------------------
	sid := -1
	for i := 0; i < len(Svcs); i++ {
		if Svcs[i].Cmd == d.Service {
			sid = i
			break
		}
	}
	lib.Console("%s: 003   sid = %d\n", funcname, sid)
	if sid < 0 {
		lib.Console("**** YIPES! **** %s - Handler not found\n", r.RequestURI)
		err = fmt.Errorf("Service not recognized: %s", d.Service)
		lib.Console("***ERROR IN URL***  %s", err.Error())
		SvcErrorReturn(w, err, funcname)
		return
	}

	//-----------------------------------------------------------------------
	// Is authentication required for this command?  If so validate that we
	// have a cookie.
	//-----------------------------------------------------------------------
	lib.Console("debug> SVC: A\n")
	if Svcs[sid].AuthNRequired {
		lib.Console("debug> SVC: B\n")
		c, err := db.ValidateSessionCookie(r, false) // this updates the expire time
		if err != nil {
			SvcErrorReturn(w, err, funcname)
			return
		}
		lib.Console("debug> SVC: C\n")
		if !c {
			SvcErrorReturn(w, db.ErrSessionRequired, funcname)
			return
		}

		lib.Console("debug> SVC: D\n")
		//----------------------------------------------------------------------
		// The air cookie is valid.  Create (or get) the internal session. This
		// is needed to identify the person associated with the request. All
		// database updates performed by this user will be updated captured
		// in the database writes/updates. We maintain info about this user
		// so we don't have to look it up every time they make a db change.
		//----------------------------------------------------------------------
		if d.sess, err = db.GetSession(r.Context(), w, r); err != nil {
			SvcErrorReturn(w, err, funcname)
			return
		}
		//----------------------------------------------------------------------
		// If we make it here, we have a good session.  Add it to the request
		// context
		//----------------------------------------------------------------------
		lib.Console("debug> SVC: E,   d.sess.{UID, Username, Firstname} = %d, %s, %s\n", d.sess.UID, d.sess.Username, d.sess.Firstname)
		ctx := db.SetSessionContextKey(r.Context(), d.sess)
		r = r.WithContext(ctx)

		lib.Console("debug> SVC: F\n")
	}

	Svcs[sid].Handler(w, r, &d)
	svcDebugTxnEnd()
}

// getBIDfromBUI reads the business field from the uri and converts it as needed
//------------------------------------------------------------------------------
func getBIDfromBUI(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return int64(0), nil
	}
	d, err := rlib.IntFromString(s, "bad request integer value") // assume it's a BID
	if err != nil {
		err = nil // clear the slate

		// need to do a db lookup for the BUD to determine the BID
		// var ok bool // OK, let's see if it's a BUD
		// d, ok = rlib.RRdb.BUDlist[s]
		// if !ok {
		// 	d = 0
		// 	err = fmt.Errorf("Could not find Business for %q", s)
		// }
		d = 0
		err = fmt.Errorf("Could not find Business for %q", s)
	}
	return d, err
}

func getPOSTdata(w http.ResponseWriter, r *http.Request, d *ServiceData) error {
	funcname := "getPOSTdata"
	var err error

	const _1MB = (1 << 20) * 1024
	lib.Console("Entered %s\n", funcname)

	// if content type is form data then
	ct := r.Header.Get("Content-Type")
	ct, _, err = mime.ParseMediaType(ct)
	if err != nil {
		e := fmt.Errorf("%s: Error while parsing content type: %s", funcname, err.Error())
		SvcErrorReturn(w, e, funcname)
		return e
	}
	if ct == "multipart/form-data" {
		// parse multipart form first
		err = r.ParseMultipartForm(_1MB)
		if err != nil {
			e := fmt.Errorf("%s: Error while parsing multipart form: %s", funcname, err.Error())
			SvcErrorReturn(w, e, funcname)
			return e
		}

		// check for headers
		for _, fheaders := range r.MultipartForm.File {
			for _, fh := range fheaders {
				cd := "Content-Disposition"
				if _, ok := fh.Header["Content-Disposition"]; !ok {
					e := fmt.Errorf("%s: Header missing (%s)", funcname, cd)
					SvcErrorReturn(w, e, funcname)
					return e
				}
				ct := "Content-Type"
				if _, ok := fh.Header["Content-Type"]; !ok {
					e := fmt.Errorf("%s: Header missing (%s)", funcname, ct)
					SvcErrorReturn(w, e, funcname)
					return e
				}
			}
		}

		d.Files = r.MultipartForm.File
		d.MFValues = r.MultipartForm.Value
	}

	htmlData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		e := fmt.Errorf("%s: Error reading message Body: %s", funcname, err.Error())
		SvcErrorReturn(w, e, funcname)
		return e
	}
	lib.Console("\t- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -\n")
	lib.Console("\thtmlData = %s\n", htmlData)
	u, err := url.QueryUnescape(string(htmlData))
	if err != nil {
		e := fmt.Errorf("%s: Error with QueryUnescape: %s", funcname, err.Error())
		SvcErrorReturn(w, e, funcname)
		return e
	}
	lib.Console("\tUnescaped htmlData = %s\n", u)

	u = strings.TrimPrefix(u, "request=") // strip off "request=" if it is present
	d.data = u
	var wjs WebGridSearchRequestJSON
	err = json.Unmarshal([]byte(u), &wjs)
	if err != nil {
		e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
		SvcErrorReturn(w, e, funcname)
		return e
	}
	rlib.MigrateStructVals(&wjs, &d.wsSearchReq)
	rlib.Console("Client = %s\n", d.wsSearchReq.Client)

	return err
}

func getGETdata(w http.ResponseWriter, r *http.Request, d *ServiceData) error {
	funcname := "getGETdata"
	lib.Console("Entered %s\n", funcname)
	s, err := url.QueryUnescape(strings.TrimSpace(r.URL.String()))
	if err != nil {
		e := fmt.Errorf("%s: Error with url.QueryUnescape:  %s", funcname, err.Error())
		SvcErrorReturn(w, e, funcname)
		return e
	}
	lib.Console("Unescaped query = %s\n", s)
	d.QueryParams = r.URL.Query()
	lib.Console("Query Parameters: %v\n", d.QueryParams)

	w2uiPrefix := "request="
	n := strings.Index(s, w2uiPrefix)
	lib.Console("n = %d\n", n)
	if n > 0 {
		lib.Console("Will process as Typedown\n")
		d.data = s[n+len(w2uiPrefix):]
		lib.Console("%s: will unmarshal: %s\n", funcname, d.data)
		if err = json.Unmarshal([]byte(d.data), &d.wsTypeDownReq); err != nil {
			e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
			SvcErrorReturn(w, e, funcname)
			return e
		}
		d.wsSearchReq.Cmd = "typedown"
	} else {
		lib.Console("Will process as web search command\n")
		d.wsSearchReq.Cmd = r.URL.Query().Get("cmd")
	}

	return nil
}
