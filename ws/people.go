package ws

import (
	"fmt"
	"net/http"
	"phonebook/db"
	lib "phonebook/lib"
)

//-------------------------------------------------------------------
//                        **** SEARCH ****
//-------------------------------------------------------------------

// PeopleTypedownResponse is the data structure for the response to a search for people
type PeopleTypedownResponse struct {
	Status  string              `json:"status"`
	Total   int64               `json:"total"`
	Records []db.PeopleTypeDown `json:"records"`
}

//-------------------------------------------------------------------
//                         **** SAVE ****
//-------------------------------------------------------------------

//-------------------------------------------------------------------
//                         **** GET ****
//-------------------------------------------------------------------

// GetPersonResponse is the structure of date send in response to a
// get person request
type GetPersonResponse struct {
	Status string      `json:"status"` // typically "success"
	Record db.WSPerson `json:"recid"`  // set to id of newly inserted record
}

// SvcPeopleTypeDown handles typedown requests for People.  It returns
// Name, and TLID
// wsdoc {
//  @Title  Get People Typedown
//	@URL /v1/Peopletd/:BUI?request={"search":"The search string","max":"Maximum number of return items"}
//	@Method GET
//	@Synopsis Fast Search for Peoples matching typed characters
//  @Desc Returns TLID, FirstName, Middlename, and LastName of Peoples that
//  @Desc match supplied chars at the beginning of FirstName or LastName
//  @Input WebTypeDownRequest
//  @Response PeoplesTypedownResponse
// wsdoc }
//----------------------------------------------------------------------------
func SvcPeopleTypeDown(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	const funcname = "SvcPeopleTypeDown"
	var g PeopleTypedownResponse
	var err error

	// lib.Console("Entered %s\n", funcname)
	// lib.Console("handle typedown: GetPeopleTypeDown( search=%s, limit=%d )\n", d.wsTypeDownReq.Search, d.wsTypeDownReq.Max)
	g.Records, err = db.GetPeopleTypeDown(d.wsTypeDownReq.Search, d.wsTypeDownReq.Max)
	if err != nil {
		e := fmt.Errorf("Error getting typedown matches: %s", err.Error())
		SvcErrorReturn(w, e, funcname)
		return
	}
	// lib.Console("GetPeopleTypeDown returned %d matches\n", len(g.Records))
	g.Total = int64(len(g.Records))
	g.Status = "success"
	SvcWriteResponse(&g, w)
}

// SvcPeople handles requests for persons.  It returns the fields that we have
// vetted as being safe for web service calls.
//
// For this call, we expect the URI to contain the BID and ID, in this case the
// ID is the UID of the person we're interested in.
//
//           0  1    2   3
// uri 		/v1/asm/:BUI/ID
// The server command can be:
//      get
//      save
//      delete
//----------------------------------------------------------------------------
func SvcPeople(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	const funcname = "SvcPeople"
	var err error
	lib.Console("Entered %s\n", funcname)

	if d.ID, err = SvcExtractIDFromURI(r.RequestURI, "ID", 3, w); err != nil {
		SvcErrorReturn(w, err, funcname)
		return
	}

	lib.Console("Request: %s:  BID = %d,  ID = %d\n", d.wsSearchReq.Cmd, d.BID, d.ID)

	switch d.wsSearchReq.Cmd {
	case "get":
		getPerson(w, r, d)
		break

	// case "save":
	// 	savePerson(w, r, d)
	// 	break
	// case "delete":
	// 	deletePerson(w, r, d)
	// 	break

	default:
		err = fmt.Errorf("Unhandled command: %s", d.wsSearchReq.Cmd)
		SvcErrorReturn(w, err, funcname)
		return
	}
}

// getPerson returns the requested Person
// wsdoc {
//  @Title  Get Person
//	@URL /v1/people/:BUI/ID
//  @Method  GET
//	@Synopsis Get information on a Person
//  @Description  Return all fields for Person :UID
//	@Input WebGridSearchRequest
//  @Response GetPersonResponse
// wsdoc }
//-----------------------------------------------------------------------------
func getPerson(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	const funcname = "getPerson"
	var g GetPersonResponse
	var a db.WSPerson
	var err error

	lib.Console("entered %s, getting UID = %d\n", funcname, d.ID)
	a, err = db.GetWSPerson(d.ID)
	if err != nil {
		SvcErrorReturn(w, err, funcname)
		return
	}
	g.Record = a
	g.Status = "success"
	lib.Console("g.status = %s, g.record - %#v\n", g.Status, g.Record)
	SvcWriteResponse(&g, w)
}
