package db

import (
	"context"
	"database/sql"
	"phonebook/lib"
	"time"
)

// WSPerson is the People information we pass over the web service call.
// We can add information as we need it after reviewing the security
// implications.
//-------------------------------------------------------------------------
type WSPerson struct {
	UID           int64
	FirstName     string
	MiddleName    string
	LastName      string
	PreferredName string
}

// People defines a date and a rent amount for a People. A People record
// is part of a group or list. The group is defined by the RSLID
//-----------------------------------------------------------------------------
type People struct {
	UID                     int64
	UserName                string
	LastName                string
	MiddleName              string
	FirstName               string
	PreferredName           string
	Salutation              string
	PositionControlNumber   string
	OfficePhone             string
	OfficeFax               string
	CellPhone               string
	PrimaryEmail            string
	SecondaryEmail          string
	BirthMonth              int64
	BirthDoM                int64
	HomeStreetAddress       string
	HomeStreetAddress2      string
	HomeCity                string
	HomeState               string
	HomePostalCode          string
	HomeCountry             string
	JobCode                 int64
	Hire                    time.Time
	Termination             time.Time
	MgrUID                  int64
	DeptCode                int64
	CoCode                  int64
	ClassCode               int64
	StateOfEmployment       string
	CountryOfEmployment     string
	EmergencyContactName    string
	EmergencyContactPhone   string
	Status                  int64
	EligibleForRehire       int64
	AcceptedHealthInsurance int64
	AcceptedDentalInsurance int64
	Accepted401K            int64
	LastReview              time.Time
	NextReview              time.Time
	passhash                string
	RID                     int64
	ImagePath               string
	LastModTime             time.Time
	LastModBy               int64
	CreateTime              time.Time
	CreateBy                int64
}

// GetPeopleTypeDown returns a slice of session cookies
//
// INPUTS
//
// RETURNS
//  []SessionCookie - a slice with all the rows in the sessions table.
//
//  err      Any errors encountered
//-----------------------------------------------------------------------------
func GetPeopleTypeDown(s1 string, limit int64) ([]PeopleTypeDown, error) {
	funcname := "GetPeopleTypeDown"
	var m []PeopleTypeDown
	s := "%" + s1 + "%"
	lib.Console("s = %q\n", s)

	rows, err := PrepStmts.GetPeopleTypeDown.Query(s, s, s, s, limit)
	if err != nil {
		lib.Ulog("%s: error getting rows: %s\n", funcname, err.Error())
		return m, err
	}
	defer rows.Close()

	for rows.Next() {
		var p PeopleTypeDown
		var first, middle, last, preferred string
		// SELECT UID,FirstName,MiddleName,LastName,PreferredName
		err := rows.Scan(&p.UID, &first, &middle, &last, &preferred)
		if err != nil {
			lib.Ulog("%s: error getting row:  %v\n", funcname, err)
			return m, err
		}
		fn := first
		if len(preferred) > 0 {
			fn = preferred
		}
		p.Name = fn + " "
		if len(middle) > 0 {
			p.Name += middle + " "
		}
		p.Name += last
		m = append(m, p)
		// lib.Console("UID=%d, FirstName=%q, MiddleName=%q, LastName=%q PreferredName=%q\n", p.UID, first, middle, last, preferred)
	}

	return m, nil
}

// GetWSPerson returns the requested People
// wsdoc {
//  @Title  Get People
//	@URL /v1/task/:BUI/TID
//  @Method  GET
//	@Synopsis Get information on a People
//  @Description  Return all fields for assessment :TID
//	@Input WebGridSearchRequest
//  @Response People struct, any error encountered
// wsdoc }
//-----------------------------------------------------------------------------
func GetWSPerson(id int64) (WSPerson, error) {
	var a WSPerson
	fields := []interface{}{id}
	row := PrepStmts.GetWSPerson.QueryRow(fields...)
	err := row.Scan(&a.UID, &a.FirstName, &a.MiddleName, &a.LastName, &a.PreferredName)
	SkipSQLNoRowsError(&err)
	return a, err
}

// DeletePeople deletes the People with the specified id from the database
//
// INPUTS
// ctx - db context
// id - PRID of the record to read
//
// RETURNS
// Any errors encountered, or nil if no errors
//-----------------------------------------------------------------------------
func DeletePeople(ctx context.Context, id int64) error {
	return genericDelete(ctx, "People", PrepStmts.DeletePeople, id)
}

// GetPeople reads and returns a People structure
//
// INPUTS
// ctx - db context
// id - PRID of the record to read
//
// RETURNS
// ErrSessionRequired if the session is invalid
// nil if the session is valid
//-----------------------------------------------------------------------------
func GetPeople(ctx context.Context, id int64) (People, error) {
	var a People
	if !ValidateSession(ctx) {
		return a, ErrSessionRequired
	}

	fields := []interface{}{id}
	stmt, row := getRowFromDB(ctx, PrepStmts.GetPeople, fields)
	if stmt != nil {
		defer stmt.Close()
	}
	return a, ReadPeople(row, &a)
}

// GetByUsername reads and returns a People structure
//
// INPUTS
// ctx - db context
// id - username of the record to read
//
// RETURNS
// ErrSessionRequired if the session is invalid
// nil if the session is valid
//-----------------------------------------------------------------------------
func GetByUsername(ctx context.Context, id string) (People, error) {
	var a People
	if !ValidateSession(ctx) {
		return a, ErrSessionRequired
	}

	fields := []interface{}{id}
	stmt, row := getRowFromDB(ctx, PrepStmts.GetByUsername, fields)
	if stmt != nil {
		defer stmt.Close()
	}
	return a, ReadPeople(row, &a)
}

// InsertPeople writes a new People record to the database
//
// INPUTS
// ctx - db context
// a   - pointer to struct to fill
//
// RETURNS
// id of the record just inserted
// any error encountered or nil if no error
//-----------------------------------------------------------------------------
func InsertPeople(ctx context.Context, a *People) (int64, error) {
	sess, ok := GetSessionFromContext(ctx)
	if !ok {
		return a.UID, ErrSessionRequired
	}
	fields := []interface{}{
		a.UID,
		a.UserName,
		a.LastName,
		a.MiddleName,
		a.FirstName,
		a.PreferredName,
		a.Salutation,
		a.PositionControlNumber,
		a.OfficePhone,
		a.OfficeFax,
		a.CellPhone,
		a.PrimaryEmail,
		a.SecondaryEmail,
		a.BirthMonth,
		a.BirthDoM,
		a.HomeStreetAddress,
		a.HomeStreetAddress2,
		a.HomeCity,
		a.HomeState,
		a.HomePostalCode,
		a.HomeCountry,
		a.JobCode,
		a.Hire,
		a.Termination,
		a.MgrUID,
		a.DeptCode,
		a.CoCode,
		a.ClassCode,
		a.StateOfEmployment,
		a.CountryOfEmployment,
		a.EmergencyContactName,
		a.EmergencyContactPhone,
		a.Status,
		a.EligibleForRehire,
		a.AcceptedHealthInsurance,
		a.AcceptedDentalInsurance,
		a.Accepted401K,
		a.LastReview,
		a.NextReview,
		a.passhash,
		a.RID,
		a.ImagePath,
		sess.UID,
		sess.UID,
	}

	var err error
	a.CreateBy, a.LastModBy, a.UID, err = genericInsert(ctx, "People", PrepStmts.InsertPeople, fields, a)
	return a.UID, err
}

// ReadPeople reads a full People structure of data from the database based
// on the supplied Rows pointer.
//
// INPUTS
// row - db Row pointer
// a   - pointer to struct to fill
//
// RETURNS
//
// ErrSessionRequired if the session is invalid
// nil if the session is valid
//-----------------------------------------------------------------------------
func ReadPeople(row *sql.Row, a *People) error {
	err := row.Scan(
		&a.UID,
		&a.UserName,
		&a.LastName,
		&a.MiddleName,
		&a.FirstName,
		&a.PreferredName,
		&a.Salutation,
		&a.PositionControlNumber,
		&a.OfficePhone,
		&a.OfficeFax,
		&a.CellPhone,
		&a.PrimaryEmail,
		&a.SecondaryEmail,
		&a.BirthMonth,
		&a.BirthDoM,
		&a.HomeStreetAddress,
		&a.HomeStreetAddress2,
		&a.HomeCity,
		&a.HomeState,
		&a.HomePostalCode,
		&a.HomeCountry,
		&a.JobCode,
		&a.Hire,
		&a.Termination,
		&a.MgrUID,
		&a.DeptCode,
		&a.CoCode,
		&a.ClassCode,
		&a.StateOfEmployment,
		&a.CountryOfEmployment,
		&a.EmergencyContactName,
		&a.EmergencyContactPhone,
		&a.Status,
		&a.EligibleForRehire,
		&a.AcceptedHealthInsurance,
		&a.AcceptedDentalInsurance,
		&a.Accepted401K,
		&a.LastReview,
		&a.NextReview,
		&a.passhash,
		&a.RID,
		&a.ImagePath,
		&a.LastModTime,
		&a.LastModBy,
		&a.CreateTime,
		&a.CreateBy,
	)
	SkipSQLNoRowsError(&err)
	return err
}

// ReadPeopleList reads a full People structure of data from the database based
// on the supplied Rows pointer.
//
// INPUTS
// row - db Row pointer
// a   - pointer to struct to fill
//
// RETURNS
//
// ErrSessionRequired if the session is invalid
// nil if the session is valid
//-----------------------------------------------------------------------------
func ReadPeopleList(rows *sql.Rows, a *People) error {
	err := rows.Scan(
		&a.UID,
		&a.UserName,
		&a.LastName,
		&a.MiddleName,
		&a.FirstName,
		&a.PreferredName,
		&a.Salutation,
		&a.PositionControlNumber,
		&a.OfficePhone,
		&a.OfficeFax,
		&a.CellPhone,
		&a.PrimaryEmail,
		&a.SecondaryEmail,
		&a.BirthMonth,
		&a.BirthDoM,
		&a.HomeStreetAddress,
		&a.HomeStreetAddress2,
		&a.HomeCity,
		&a.HomeState,
		&a.HomePostalCode,
		&a.HomeCountry,
		&a.JobCode,
		&a.Hire,
		&a.Termination,
		&a.MgrUID,
		&a.DeptCode,
		&a.CoCode,
		&a.ClassCode,
		&a.StateOfEmployment,
		&a.CountryOfEmployment,
		&a.EmergencyContactName,
		&a.EmergencyContactPhone,
		&a.Status,
		&a.EligibleForRehire,
		&a.AcceptedHealthInsurance,
		&a.AcceptedDentalInsurance,
		&a.Accepted401K,
		&a.LastReview,
		&a.NextReview,
		&a.passhash,
		&a.RID,
		&a.ImagePath,
		&a.LastModTime,
		&a.LastModBy,
		&a.CreateTime,
		&a.CreateBy,
	)
	SkipSQLNoRowsError(&err)
	return err
}

// UpdatePeople updates an existing People record in the database
//
// INPUTS
// ctx - db context
// a   - pointer to struct to fill
//
// RETURNS
// id of the record just inserted
// any error encountered or nil if no error
//-----------------------------------------------------------------------------
func UpdatePeople(ctx context.Context, a *People) error {
	sess, ok := GetSessionFromContext(ctx)
	if !ok {
		return ErrSessionRequired
	}
	fields := []interface{}{
		a.UserName,
		a.LastName,
		a.MiddleName,
		a.FirstName,
		a.PreferredName,
		a.Salutation,
		a.PositionControlNumber,
		a.OfficePhone,
		a.OfficeFax,
		a.CellPhone,
		a.PrimaryEmail,
		a.SecondaryEmail,
		a.BirthMonth,
		a.BirthDoM,
		a.HomeStreetAddress,
		a.HomeStreetAddress2,
		a.HomeCity,
		a.HomeState,
		a.HomePostalCode,
		a.HomeCountry,
		a.JobCode,
		a.Hire,
		a.Termination,
		a.MgrUID,
		a.DeptCode,
		a.CoCode,
		a.ClassCode,
		a.StateOfEmployment,
		a.CountryOfEmployment,
		a.EmergencyContactName,
		a.EmergencyContactPhone,
		a.Status,
		a.EligibleForRehire,
		a.AcceptedHealthInsurance,
		a.AcceptedDentalInsurance,
		a.Accepted401K,
		a.LastReview,
		a.NextReview,
		a.passhash,
		a.RID,
		a.LastModTime,
		a.LastModBy,
		a.ImagePath,
		sess.UID,
		a.UID,
	}
	var err error
	a.LastModBy, err = genericUpdate(ctx, PrepStmts.UpdatePeople, fields)
	return updateError(err, "People", *a)
}
