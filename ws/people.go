package ws

import (
	"fmt"
	"net/http"
	"phonebook/db"
)

// PeopleTypedownResponse is the data structure for the response to a search for people
type PeopleTypedownResponse struct {
	Status  string              `json:"status"`
	Total   int64               `json:"total"`
	Records []db.PeopleTypeDown `json:"records"`
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

	// util.Console("Entered %s\n", funcname)
	// util.Console("handle typedown: GetPeopleTypeDown( search=%s, limit=%d )\n", d.wsTypeDownReq.Search, d.wsTypeDownReq.Max)
	g.Records, err = db.GetPeopleTypeDown(d.wsTypeDownReq.Search, d.wsTypeDownReq.Max)
	if err != nil {
		e := fmt.Errorf("Error getting typedown matches: %s", err.Error())
		SvcErrorReturn(w, e, funcname)
		return
	}
	// util.Console("GetPeopleTypeDown returned %d matches\n", len(g.Records))
	g.Total = int64(len(g.Records))
	g.Status = "success"
	SvcWriteResponse(&g, w)
}
