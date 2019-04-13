package ws

import (
	"fmt"
	"net/http"
	"phonebook/db"
	"phonebook/lib"
)

//-------------------------------------------------------------------
//  SEARCH
//-------------------------------------------------------------------

// DirBU is a directory business unit definition
type DirBU struct {
	ClassCode   int
	CoCode      int
	Name        string
	Designation string
	Description string
}

// DirCompany is a directory company definition
type DirCompany struct {
	CoCode           int64
	LegalName        string
	CommonName       string
	Address          string
	Address2         string
	City             string
	State            string
	PostalCode       string
	Country          string
	Phone            string
	Fax              string
	Email            string
	Designation      string
	Active           int
	EmploysPersonnel int
}

// BUTDResponse holds the task list definition list
type BUTDResponse struct {
	Status  string      `json:"status"`
	Total   int64       `json:"total"`
	Records []db.BUInfo `json:"records"`
}

// BUDResponse holds the task list definition list
type BUDResponse struct {
	Status string `json:"status"`
	Record DirBU  `json:"record"`
}

//-------------------------------------------------------------------
//  SAVE
//-------------------------------------------------------------------

// SvcBUTypedown handles typedown requests for Business Units. It returns
// the ClassCode (the uid), the Name, and the Designation.
// wsdoc {
//  @Title  BU Typedown
//	@URL /v1/butd/:BUI?request={"search":"The search string","max":"Maximum number of return items"}
//	@Method GET
//	@Synopsis Fast Search for Transactants matching typed characters
//  @Desc Returns ClassCode, Name, and Designation that
//  @Desc match supplied chars at the beginning of the BU's designation
//  @Input WebTypeDownRequest
//  @Response BUTDResponse
// wsdoc }
//-------------------------------------------------------------------
func SvcBUTypedown(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	const funcname = "SvcBUTypedown"
	var g BUTDResponse
	var err error

	// lib.Console("Entered %s\n", funcname)
	// lib.Console("handle typedown: GetTransactantsTypeDown( search=%s, limit=%d\n", d.wsTypeDownReq.Search, d.wsTypeDownReq.Max)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	g.Records, err = db.GetBUTypeDown(d.wsTypeDownReq.Search, d.wsTypeDownReq.Max)
	// lib.Console("GetTransactantTypeDown returned %d matches\n", len(g.Records))
	g.Total = int64(len(g.Records))
	if err != nil {
		e := fmt.Errorf("Error getting typedown matches: %s", err.Error())
		SvcErrorReturn(w, e, funcname)
		return
	}
	g.Status = "success"
	SvcWriteResponse(&g, w)
}

// SvcBUD dispatch a /v1/bud command
//-------------------------------------------------------------------
func SvcBUD(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	const funcname = "SvcBUD"

	lib.Console("Request: %s\n", d.wsSearchReq.Cmd)

	switch d.wsSearchReq.Cmd {
	case "typedown":
		getBUByBUD(w, r, d)
	default:
		err := fmt.Errorf("Unhandled command: %s", d.wsSearchReq.Cmd)
		SvcErrorReturn(w, err, funcname)
		return
	}
}

// getBUByBUD does a search based on the supplied BUD.  If found, it returns
// the Class information. If not found it returns status = error
// wsdoc {
//  @Title  Get BUD
//	@URL /v1/bud/:BUI?request={"search":"XYZ","max":"10"}
//	@Method GET
//	@Synopsis Find a Business Unit by its BUD
//  @Desc Returns ClassCode, Name, and Designation that
//  @Desc match supplied chars at the beginning of the BU's designation.
//  @Desc Example: (not uri encoded) /v1/butd/10?request={"search":"XYZ"}
//  @Desc searches for Business Units where containing "xy" and return
//  @Desc no more than 10.
//  @Input WebTypeDownRequest
//  @Response BUDResponse
// wsdoc }
//-------------------------------------------------------------------
func getBUByBUD(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	const funcname = "getBUByBUD"
	var g BUDResponse
	var err error
	var a db.Class
	lib.Console("Entered %s\n", funcname)
	lib.Console("handle typedown: GetBUByBUD( search = %q )\n", d.wsTypeDownReq.Search)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	a, err = db.GetBUByBUD(d.wsTypeDownReq.Search)
	// lib.Console("GetTransactantTypeDown returned %d matches\n", len(g.Records))
	if err != nil {
		e := fmt.Errorf("Error getting typedown matches: %s", err.Error())
		SvcErrorReturn(w, e, funcname)
		return
	}
	lib.Console("a.BUD = %s, a.Name = %s\n", a.Designation, a.Name)
	if a.ClassCode == 0 {
		e := fmt.Errorf("No Business Unit with Designation %s exists", d.wsTypeDownReq.Search)
		SvcErrorReturn(w, e, funcname)
		return
	}
	lib.MigrateStructVals(&a, &g.Record)
	g.Status = "success"
	SvcWriteResponse(&g, w)
}
