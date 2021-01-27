package ws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"phonebook/lib"
	"rentroll/rlib"
	"strings"
)

// SvcExtractIDFromURI extracts an int64 id value from position pos of the supplied uri.
// The URI is of the form returned by http.Request.RequestURI .  In particular:
//
//	pos:     0    1      2  3
//  uri:    /v1/rentable/34/421
//
// So, in the example uri above, a call where pos = 3 would return int64(421). errmsg
// is a string that will be used in the error message if the requested position had an
// error during conversion to int64. So in the example above, pos 3 is the RID, so
// errmsg would probably be set to "RID"
func SvcExtractIDFromURI(uri, errmsg string, pos int, w http.ResponseWriter) (int64, error) {
	var ID = int64(0)
	var err error
	var funcname = "SvcExtractIDFromURI"

	sa := strings.Split(uri[1:], "/")
	// rlib.Console("uri parts:  %v\n", sa)
	if len(sa) < pos+1 {
		err = fmt.Errorf("Expecting at least %d elements in URI: %s, but found only %d", pos+1, uri, len(sa))
		// rlib.Console("err = %s\n", err)
		SvcErrorReturn(w, err, funcname)
		return ID, err
	}
	// rlib.Console("sa[pos] = %s\n", sa[pos])
	ID, err = SvcGetInt64(sa[pos], errmsg, w)
	return ID, err
}

// SvcGetInt64 tries to read an int64 value from the supplied string.
// If it fails for any reason, it sends writes an error message back
// to the caller and returns the error.  Otherwise, it returns an
// int64 and returns nil
//-----------------------------------------------------------------------------
func SvcGetInt64(s, errmsg string, w http.ResponseWriter) (int64, error) {
	i, err := rlib.IntFromString(s, "not an integer number")
	if err != nil {
		err = fmt.Errorf("%s: %s", errmsg, err.Error())
		SvcErrorReturn(w, err, "SvcGetInt64")
		return i, err
	}
	return i, nil
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
	// lib.Console("\nResponse Data:  %s\n\n", string(b))
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

// SvcWriteSuccessResponseWithID is used to complete a successful write
// operation on w2ui form save requests.
func SvcWriteSuccessResponseWithID(w http.ResponseWriter, id int64) {
	var g = SvcStatusResponse{Status: "success", Recid: id}
	w.Header().Set("Content-Type", "application/json")
	SvcWriteResponse(&g, w)
}
