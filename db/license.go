package db

import (
	"context"
	"database/sql"
	"phonebook/lib"
	"time"
)

// License defines the structure associated with a realtor's seller / broker license
type License struct {
	LID         int64
	UID         int64
	State       string
	LicenseNo   string
	FLAGS       int64 // bit 0:  0 = seller license, 1 = broker license
	LastModTime time.Time
	LastModBy   int64
	CreateTime  time.Time
	CreateBy    int64
}

// DeleteLicense deletes the License with the specified id from the database
//
// INPUTS
// ctx - db context
// id - PRID of the record to read
//
// RETURNS
// Any errors encountered, or nil if no errors
//-----------------------------------------------------------------------------
func DeleteLicense(ctx context.Context, id int64) error {
	return genericDelete(ctx, "license", PrepStmts.DeleteLicense, id)
}

// GetLicense reads and returns a License structure
//
// INPUTS
// ctx - db context
// id - LID of the record to read
//
// RETURNS
// ErrSessionRequired if the session is invalid
// nil if the session is valid
//-----------------------------------------------------------------------------
func GetLicense(ctx context.Context, id int64) (License, error) {
	var a License
	if !ValidateSession(ctx) {
		return a, ErrSessionRequired
	}

	fields := []interface{}{id}
	stmt, row := getRowFromDB(ctx, PrepStmts.GetLicense, fields)
	if stmt != nil {
		defer stmt.Close()
	}
	return a, ReadLicense(row, &a)
}

// GetLicenses reads and returns a slice of License structures where the UID
// matches id
//
// INPUTS
// ctx - db context
// id - UID of the records to read
//
// RETURNS
// ErrSessionRequired if the session is invalid
// nil if the session is valid
//-----------------------------------------------------------------------------
func GetLicenses(ctx context.Context, id int64) ([]License, error) {
	var a []License
	var err error
	if !ValidateSession(ctx) {
		return a, ErrSessionRequired
	}

	fields := []interface{}{id}
	stmt2, rows, err := getRowsFromDB(ctx, PrepStmts.GetLicenses, fields)
	if err != nil {
		return a, err
	}
	if stmt2 != nil {
		defer stmt2.Close()
	}
	for i := 0; rows.Next(); i++ {
		var x License
		if err = ReadLicenseList(rows, &x); err != nil {
			return a, err
		}
		a = append(a, x)
	}
	if err = rows.Err(); err != nil {
		return a, err
	}
	return a, nil
}

// InsertLicense writes a new License record to the database
//
// INPUTS
// ctx - db context
// a   - pointer to struct to fill
//
// RETURNS
// id of the record just inserted
// any error encountered or nil if no error
//-----------------------------------------------------------------------------
func InsertLicense(ctx context.Context, a *License) (int64, error) {
	sess, ok := GetSessionFromContext(ctx)
	lib.Console("session info from context:  ok = %t\n", ok)
	lib.Console("sess.{UID, Username, Firstname} = %d, %s, %s", sess.UID, sess.Username, sess.Firstname)
	if !ok {
		return a.LID, ErrSessionRequired
	}
	fields := []interface{}{
		a.UID,
		a.State,
		a.LicenseNo,
		a.FLAGS,
		sess.UID,
		sess.UID,
	}
	var err error
	a.CreateBy, a.LastModBy, a.LID, err = genericInsert(ctx, "License", PrepStmts.InsertLicense, fields, a)
	return a.LID, err
}

// ReadLicense reads a full License structure of data from the database based
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
func ReadLicense(row *sql.Row, a *License) error {
	err := row.Scan(
		&a.LID,
		&a.UID,
		&a.State,
		&a.LicenseNo,
		&a.FLAGS,
		&a.LastModTime,
		&a.LastModBy,
		&a.CreateTime,
		&a.CreateBy,
	)
	SkipSQLNoRowsError(&err)
	return err
}

// ReadLicenseList reads a full License structure of data from the database based
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
func ReadLicenseList(rows *sql.Rows, a *License) error {
	err := rows.Scan(
		&a.LID,
		&a.UID,
		&a.State,
		&a.LicenseNo,
		&a.FLAGS,
		&a.LastModTime,
		&a.LastModBy,
		&a.CreateTime,
		&a.CreateBy,
	)
	SkipSQLNoRowsError(&err)
	return err
}

// UpdateLicense updates an existing License record in the database
//
// INPUTS
// ctx - db context
// a   - pointer to struct to fill
//
// RETURNS
// id of the record just inserted
// any error encountered or nil if no error
//-----------------------------------------------------------------------------
func UpdateLicense(ctx context.Context, a *License) error {
	sess, ok := GetSessionFromContext(ctx)
	if !ok {
		return ErrSessionRequired
	}
	fields := []interface{}{
		a.UID,
		a.State,
		a.LicenseNo,
		a.FLAGS,
		sess.UID,
		a.LID,
	}
	var err error
	a.LastModBy, err = genericUpdate(ctx, PrepStmts.UpdateLicense, fields)
	return updateError(err, "License", *a)
}
