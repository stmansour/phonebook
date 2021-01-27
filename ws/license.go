package ws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"phonebook/db"
	"phonebook/lib"
)

//-------------------------------------------------------------------
//                        **** SEARCH ****
//-------------------------------------------------------------------

//-------------------------------------------------------------------
//                         **** SAVE ****
//-------------------------------------------------------------------

// SaveLicense is the structure of date send in response to a
// get person request
type SaveLicense struct {
	Cmd    string     `json:"cmd"`
	Record db.License `json:"record"` // set to id of newly inserted record
}

//-------------------------------------------------------------------
//                         **** GET ****
//-------------------------------------------------------------------

// GetLicenseResponse is the structure of date send in response to a
// get person request
type GetLicenseResponse struct {
	Status string     `json:"status"` // typically "success"
	Record db.License `json:"record"` // set to id of newly inserted record
}

// GetLicenseListResponse describes the POST request for getting a list of people
type GetLicenseListResponse struct {
	Status  string       `json:"status"`
	Total   int64        `json:"total"`
	Records []db.License `json:"records"`
}

// SvcHandlerLicense formats a complete data record for an property for use
// with the w2ui Form
//
// The server command can be:
//      get
//      save
//      delete
//------------------------------------------------------------------------------
func SvcHandlerLicense(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "SvcHandlerLicense"
	// lib.Console("Entered SvcHandlerLicense, cmd = %s, d.ID = %d\n", d.wsSearchReq.Cmd, d.ID)

	switch d.wsSearchReq.Cmd {
	case "get":
		if d.ID < 0 {
			SvcErrorReturn(w, fmt.Errorf("LicenseID is required but was not specified"), funcname)
			return
		}
		getLicense(w, r, d)
	case "save":
		saveLicense(w, r, d)
	case "delete":
		deleteLicense(w, r, d)
	default:
		err := fmt.Errorf("Unhandled command: %s", d.wsSearchReq.Cmd)
		SvcErrorReturn(w, err, funcname)
		return
	}
}

// getLicense retrieves a specific license by LID
//
// /v1/licenses/LID
//
// INPUTS
// ctx - db context
// id - LID of the record to read
//
// RETURNS
// Any errors encountered, or nil if no errors
//-----------------------------------------------------------------------------
func getLicense(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	const funcname = "getLicense"
	var g GetLicenseResponse
	var err error

	// var a db.WSLicense
	// a, err = db.GetWSLicense(d.ID)

	var a db.License
	a, err = db.GetLicense(r.Context(), d.ID)
	if err != nil {
		SvcErrorReturn(w, err, funcname)
		return
	}
	if a.LID == 0 {
		err = fmt.Errorf("License with LID %d not found", d.ID)
		SvcErrorReturn(w, err, funcname)
		return
	}
	g.Record = a
	g.Status = "success"
	// lib.Console("g.status = %s, g.record - %#v\n", g.Status, g.Record)
	SvcWriteResponse(&g, w)
}

// deleteLicense retrieves a specific license by LID
//
// /v1/licenses/LID
//
// INPUTS
// ctx - db context
// id - LID of the record to read
//
// RETURNS
// Any errors encountered, or nil if no errors
//-----------------------------------------------------------------------------
func deleteLicense(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	const funcname = "deleteLicense"

	// lib.Console("Entered: %s, delete LID = %d\n", funcname, d.ID)
	err := db.DeleteLicense(r.Context(), d.ID)
	if err != nil {
		SvcErrorReturn(w, err, funcname)
		return
	}
	SvcWriteSuccessResponse(w)
}

// GetLicenses retrieves all licenses associated with UID
//
// /v1/licenses/UID
//
// INPUTS
// ctx - db context
// id - LID of the record to read
//
// RETURNS
// Any errors encountered, or nil if no errors
//-----------------------------------------------------------------------------
func GetLicenses(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	const funcname = "GetLicenses"
	var g GetLicenseListResponse
	var err error
	var m []db.License

	// var a db.WSLicense
	// a, err = db.GetWSLicense(d.ID)

	lib.Console("Entered %s:  search for licenses for UID=%d\n", funcname, d.ID)

	m, err = db.GetLicenses(r.Context(), d.ID)
	if err != nil {
		SvcErrorReturn(w, err, funcname)
		return
	}

	for i := 0; i < len(m); i++ {
		var a db.License
		if err = lib.MigrateStructVals(&m[i], &a); err != nil {
			SvcErrorReturn(w, err, funcname)
			return
		}
		g.Records = append(g.Records, a)
	}
	g.Status = "success"
	g.Total = int64(len(g.Records))
	SvcWriteResponse(&g, w)
}

// SaveLicense returns the requested property
// wsdoc {
//  @Title  Save License
//	@URL /v1/License/LID
//  @Method  GET
//	@Synopsis Update the information on a License with the supplied data, create if necessary.
//  @Description  This service creates a License if LID == 0 or updates a License if LID > 0 with
//  @Description  the information supplied. All fields must be supplied.
//	@Input SaveLicense
//  @Response SvcStatusResponse
// wsdoc }
//-----------------------------------------------------------------------------
func saveLicense(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "saveLicense"
	lib.Console("Entered %s\n", funcname)
	// lib.Console("record data = %s\n", d.data)
	// lib.Console("LID = %d\n", d.ID)

	var p SaveLicense

	err := json.Unmarshal([]byte(d.data), &p)
	if err != nil {
		e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
		SvcErrorReturn(w, e, funcname)
		return
	}

	// util.Console("license info: p.Record.LicenseNo = %s, whole struct: %#v\n", p.Record.LicenseNo, p)

	if p.Record.LID < 1 {
		if _, err = db.InsertLicense(r.Context(), &p.Record); err != nil {
			e := fmt.Errorf("%s: Error with db.InsertLicense:  %s", funcname, err.Error())
			SvcErrorReturn(w, e, funcname)
			return
		}
	} else {
		if err = db.UpdateLicense(r.Context(), &p.Record); err != nil {
			e := fmt.Errorf("%s: Error with db.UpdateLicense:  %s", funcname, err.Error())
			SvcErrorReturn(w, e, funcname)
			return
		}
	}
	// lib.Console("UpdateLicense completed successfully\n")
	SvcWriteSuccessResponseWithID(w, p.Record.LID)
}
