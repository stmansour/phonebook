package ws

import (
	"fmt"
	"mojo/util"
	"net/http"
	"phonebook/db"
)

//-------------------------------------------------------------------
//  SEARCH
//-------------------------------------------------------------------

// DirBU is a directory business unit definition
type DirBU struct {
	ClassCode   int64
	CoCode      int64
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
	var (
		g   BUTDResponse
		err error
	)
	util.Console("Entered %s\n", funcname)
	util.Console("handle typedown: GetTransactantsTypeDown( search=%s, limit=%d\n", d.wsTypeDownReq.Search, d.wsTypeDownReq.Max)
	g.Records, err = db.GetBUTypeDown(d.wsTypeDownReq.Search, d.wsTypeDownReq.Max)
	util.Console("GetTransactantTypeDown returned %d matches\n", len(g.Records))
	g.Total = int64(len(g.Records))
	if err != nil {
		e := fmt.Errorf("Error getting typedown matches: %s", err.Error())
		SvcErrorReturn(w, e, funcname)
		return
	}
	g.Status = "success"
	SvcWriteResponse(&g, w)
}
