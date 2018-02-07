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
	"phonebook/lib"
	"strings"
)

// SvcCtx holds global data needed by the service routines
var SvcCtx struct {
	db *sql.DB
}

// ServiceData is the generalized data gatherer for svcHandler. It allows all
// the common data to be centrally parsed and passed to a handler, which may
// need to parse further to get its unique data.  It includes fields for
// common data elements in web svc requests
type ServiceData struct {
	Service      string              // the service requested (position 1)
	DetVal       string              // value of 3rd path element if present (it is not always a number)
	UID          int64               // user id of requester
	pathElements []string            // the parts of the uri
	data         string              // the raw unparsed data
	QueryParams  map[string][]string // parameters when HTTP GET is used
	Files        map[string][]*multipart.FileHeader
	MFValues     map[string][]string
}

// ServiceHandler describes the handler for all services
type ServiceHandler struct {
	Cmd     string
	Handler func(http.ResponseWriter, *http.Request, *ServiceData)
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
	{"authenticate", SvcAuthenticate},
	{"discon", SvcDisableConsole},
	{"encon", SvcEnableConsole},
	{"logoff", SvcLogoff},
	{"resetpw", SvcResetPWHandler},
	{"validatecookie", SvcValidateCookie},
	{"version", SvcHandlerVersion},
}

// InitServices initializes the context data needed by service routines
//-----------------------------------------------------------------------------
func InitServices(db *sql.DB) {
	SvcCtx.db = db
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
	svcDebugTxn(funcname, r)
	var d ServiceData

	//-----------------------------------------------------------------------
	// pathElements:  0   1
	//               /v1/{subservice}/
	// ex:           /v1/authenticate/
	//-----------------------------------------------------------------------
	lib.Console("RequestURI = %s\n", r.RequestURI)
	ss := strings.Split(r.RequestURI[1:], "?") // it could be GET command
	d.pathElements = strings.Split(ss[0], "/")
	for i := 0; i < len(d.pathElements); i++ {
		lib.Console("%d. %s\n", i, d.pathElements[i])
	}
	d.Service = d.pathElements[1]

	svcDebugURL(r, &d)
	showRequestHeaders(r)

	switch r.Method {
	case "POST":
		if nil != getPOSTdata(w, r, &d) {
			return
		}
	case "GET":
		if nil != getGETdata(w, r, &d) {
			return
		}
	}

	//-----------------------------------------------------------------------
	//  Now call the appropriate handler to do the rest
	//-----------------------------------------------------------------------
	found := false
	for i := 0; i < len(Svcs); i++ {
		if Svcs[i].Cmd == d.Service {
			Svcs[i].Handler(w, r, &d)
			found = true
			break
		}
	}
	if !found {
		lib.Console("**** YIPES! **** %s - Handler not found\n", r.RequestURI)
		e := fmt.Errorf("Service not recognized: %s", d.Service)
		lib.Console("***ERROR IN URL***  %s", e.Error())
		SvcErrorReturn(w, e, funcname)
		return
	}
	svcDebugTxnEnd()
}

func getPOSTdata(w http.ResponseWriter, r *http.Request, d *ServiceData) error {
	funcname := "getPOSTdata"
	var err error

	const _1MB = (1 << 20) * 1024

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
	return err
}

func getGETdata(w http.ResponseWriter, r *http.Request, d *ServiceData) error {
	funcname := "getGETdata"
	s, err := url.QueryUnescape(strings.TrimSpace(r.URL.String()))
	if err != nil {
		e := fmt.Errorf("%s: Error with url.QueryUnescape:  %s", funcname, err.Error())
		SvcErrorReturn(w, e, funcname)
		return e
	}
	lib.Console("Unescaped query = %s\n", s)
	d.QueryParams = r.URL.Query()
	lib.Console("Query Parameters: %v\n", d.QueryParams)
	return nil
}

// SvcSuccessReturn sends a success message to the requester
func SvcSuccessReturn(w http.ResponseWriter) {
	g := SvcSuccess{Status: "success"}
	SvcWriteResponse(&g, w)
}

// SvcErrorReturn formats an error return to the grid widget and sends it
func SvcErrorReturn(w http.ResponseWriter, err error, funcname string) {
	// lib.Console("<Function>: %s | <Error>: %s\n", funcname, err.Error())
	lib.Console("%s: %s\n", funcname, err.Error())
	var e SvcError
	e.Status = "error"
	e.Message = fmt.Sprintf("Error: %s", err.Error())
	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(e)
	SvcWrite(w, b)
}

// SvcWriteResponse finishes the transaction with the W2UI client
func SvcWriteResponse(g interface{}, w http.ResponseWriter) {
	funcname := "SvcWriteResponse"
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(g)
	if err != nil {
		e := fmt.Errorf("Error marshaling json data: %s", err.Error())
		lib.Ulog("%s: %s\n", funcname, err.Error())
		SvcErrorReturn(w, e, funcname)
		return
	}
	SvcWrite(w, b)
}

// SvcWrite is a general write routine for service calls... it is a bottleneck
// where we can place debug statements as needed.
func SvcWrite(w http.ResponseWriter, b []byte) {
	lib.Console("first 300 chars of response: %-300.300s\n", string(b))
	// util.Console("\nResponse Data:  %s\n\n", string(b))
	w.Write(b)
}

// SvcWriteSuccessResponse is used to complete a successful write operation on w2ui form save requests.
func SvcWriteSuccessResponse(w http.ResponseWriter) {
	var g = SvcStatusResponse{Status: "success"}
	w.Header().Set("Content-Type", "application/json")
	SvcWriteResponse(&g, w)
}

// SvcHandlerVersion returns the server version number
//  @Title Verrsion
//  @URL /v1/version
//  @Method  POST or GET
//  @Synopsis Get the current server version
//  @Description Returns the server build number appended to the major/minor
//  @Description version number.
//  @Input
//  @Response version number
// wsdoc }
func SvcHandlerVersion(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	lib.Ulog("Entered SvcHandlerVersion\n")
	lib.Ulog("lib.GetVersionNo() returns %s\n", lib.GetVersionNo())
	fmt.Fprintf(w, "%s", lib.GetVersionNo())
}
