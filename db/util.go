package db

import (
	"context"
	"database/sql"
	"errors"
	"extres"
	"fmt"
	"phonebook/lib"
	"runtime/debug"
	"wreis/session"
)

// ErrSessionRequired session required error
var ErrSessionRequired = errors.New("Session Required, Please Login")

// Errcheck - saves a bunch of typing, prints error if it exists
//            and provides a traceback as well
// Note that the error is printed only if the environment is NOT production.
//------------------------------------------------------------------------------
func Errcheck(err error) {
	if err != nil {
		if IsSQLNoResultsError(err) {
			return
		}
		if extres.APPENVPROD != lib.AppConfig.Env {
			fmt.Printf("error = %v\n", err)
		}
		debug.PrintStack()
		// log.Fatal(err)
	}
}

// IsSQLNoResultsError returns true if the error provided is a sql err
// indicating no rows in the solution set.
//------------------------------------------------------------------------------
func IsSQLNoResultsError(err error) bool {
	return err == sql.ErrNoRows
}

// SkipSQLNoRowsError assing nil to original err variable
// if its kind of no rows in result error from sql package
//------------------------------------------------------------------------------
func SkipSQLNoRowsError(err *error) {
	if IsSQLNoResultsError(*err) {
		*err = nil
	}
}

// ValidateSession validates that we either have a valid session or that a
// session is not required.  If there's a problem it returns an error.
//
// INPUTS
// context - session Context
//
// RETURNS
// true if the session is valid
// false otherwise
//------------------------------------------------------------------------------
func ValidateSession(ctx context.Context) bool {
	if !(lib.AppConfig.Env != extres.APPENVPROD) {
		_, ok := session.GetSessionFromContext(ctx)
		if !ok {
			return false
		}
	}
	return true
}

// ValidateSessionForDBDelete validates that we either have a valid session or that a
// session is not required.  If there's a problem it returns an error.
//
// INPUTS
// context - session Context
//
// RETURNS
// true if the session is valid
// false otherwise
//------------------------------------------------------------------------------
func ValidateSessionForDBDelete(ctx context.Context) error {
	if !(DB.noAuth && lib.AppConfig.Env != extres.APPENVPROD) {
		_, ok := session.GetSessionFromContext(ctx)
		if !ok {
			return ErrSessionRequired
		}
	}
	return nil
}

// ValidateSessionForDBInsert is a convenience function that replaces 8 lines
// of code with about 3. Since these lines are needed for every insert call
// it saves a lot of lines.  Added this routine at the time Task,TaskList,
// TaskDescriptor and  TaskListDefinition were added.
//
// INPUTS
// ctx - session Context
// id1 - pointer to creator UID, will be returned
// id2 - pointer to lastmod UID, will be returned
//
// RETURNS
// any error encountered, nil otherwise
//-----------------------------------------------------------------------------
func ValidateSessionForDBInsert(ctx context.Context, id1, id2 *int64) error {
	if !(DB.noAuth && lib.AppConfig.Env != extres.APPENVPROD) {
		sess, ok := session.GetSessionFromContext(ctx)
		if !ok {
			return ErrSessionRequired
		}
		(*id1) = sess.UID
		(*id2) = sess.UID
		return nil
	}
	return nil
}

// ValidateSessionForDBUpdate is a convenience function that replaces 8 lines
// of code with about 4. Since these lines are needed for every update call
// it saves a lot of lines.  Added this routine at the time Task,TaskList,
// TaskDescriptor and  TaskListDefinition were added.
//
// INPUTS
// ctx - session Context
// id2 - pointer to lastmod UID, will be returned
//
// RETURNS
// any error encountered, nil otherwise
//-----------------------------------------------------------------------------
func ValidateSessionForDBUpdate(ctx context.Context, id2 *int64) error {
	if !(DB.noAuth && lib.AppConfig.Env != extres.APPENVPROD) {
		sess, ok := session.GetSessionFromContext(ctx)
		if !ok {
			return ErrSessionRequired
		}
		(*id2) = sess.UID
		return nil
	}
	return nil
}

// GetRowCountRaw returns the number of database rows in the supplied table with
// the supplied where clause. The where clause can be empty.
//
// INPUTS
//    table - table are we querying
//    joins - any join info, can be nil or an empty string
//    where - the where clause, can be nil or an empty string
//-----------------------------------------------------------------------------
func GetRowCountRaw(table, joins, where string) (int64, error) {
	count := int64(0)
	var err error
	s := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
	if len(joins) > 0 {
		s += " " + joins
	}
	if len(where) > 0 {
		s += " " + where
	}
	lib.Console("\n\nGetRowCountRaw: QUERY = %s\n", s)
	de := DB.DirDB.QueryRow(s).Scan(&count)
	if de != nil {
		err = fmt.Errorf("GetRowCountRaw: query=\"%s\"    err = %s", s, de.Error())
	}
	return count, err
}

func updateError(err error, n string, a interface{}) error {
	if nil != err {
		lib.Ulog("Update%s: error updating %s:  %v\n", n, n, err)
		lib.Ulog("%s = %#v\n", n, a)
	}
	return err
}

func insertError(err error, n string, a interface{}) error {
	if nil != err {
		lib.Ulog("Insert%s: error inserting %s:  %v\n", n, n, err)
		lib.Ulog("%s = %#v\n", n, a)
	}
	return err
}

// deleteDBRow encapsulates about 6 or 7 lines of code to handle deleting a
// database row either in the context of a transaction or not.  These lines
// appear in every Get call, so it simplifies the code a lot.
//
// INPUTS
// ctx        - db context
// s          - name of TABLE, used only for error logging
// stmt       - the prepared statement to execute, caller must de
// fields     - the fields to supply to the prepared statement
//
// RETURNS
// *sql.Stmt  - the statement -- nil if not in transaction, non-nil otherwise
// 				note that the caller needs to defer stmt.Close() when stmt is
//              not nil.
// error      - any error encountered
//-----------------------------------------------------------------------------
func deleteDBRow(ctx context.Context, s string, stmt *sql.Stmt, fields []interface{}) error {
	var err error

	if tx, ok := TxFromContext(ctx); ok { // if transaction is supplied
		stmt := tx.Stmt(stmt)
		defer stmt.Close()
		_, err = stmt.Exec(fields...)
	} else {
		_, err = stmt.Exec(fields...)
	}
	if err != nil {
		lib.Ulog("Error deleting %s fields[]=%#v, error: %v\n", s, fields, err)
	}

	return err
}

// genericDelete encapsulates the most generic form of the database delete
// functionality supplied in this package.
//
// INPUTS
// ctx        - db context
// s          - name of TABLE, used only for error logging
// stmt       - the prepared statement to execute, caller must de
// id         - id of the table record to delete
//
// RETURNS
// error      - any error encountered
//-----------------------------------------------------------------------------
func genericDelete(ctx context.Context, s string, g *sql.Stmt, id int64) error {
	var err error

	if err = ValidateSessionForDBDelete(ctx); err != nil {
		return err
	}

	fields := []interface{}{id}
	if err = deleteDBRow(ctx, s, g, fields); err != nil {
		return err
	}
	return err
}

// getRowFromDB encapsulates about 6 or 7 lines of code to handle getting a
// database row either in the context of a transaction or not.  These lines
// appear in every Get call, so it simplifies the code a lot.
//
// INPUTS
// ctx        - db context
// stmt       - the prepared statement to execute
// fields     - the fields to supply to the prepared statement
//
// RETURNS
// *sql.Stmt  - the statement -- THIS NEEDS TO BE NIL UNLESS IT'S A TXN
//				nil if not in transaction, non-nil otherwise
// 				note that the caller needs to defer stmt.Close() when stmt is
//              not nil.
// *sql.Row   - the database row to read
//-----------------------------------------------------------------------------
func getRowFromDB(ctx context.Context, stmt *sql.Stmt, fields []interface{}) (*sql.Stmt, *sql.Row) {
	var row *sql.Row
	if tx, ok := TxFromContext(ctx); ok { // if transaction is supplied
		stmt := tx.Stmt(stmt)
		row = stmt.QueryRow(fields...)
		return stmt, row
	}
	row = stmt.QueryRow(fields...)
	return nil, row
}

// getRowsFromDB encapsulates about 6 or 7 lines of code to handle getting a
// database row either in the context of a transaction or not.  These lines
// appear in every Get call, so it simplifies the code a lot.
//
// INPUTS
// ctx        - db context
// stmt       - the prepared statement to execute
// fields     - the fields to supply to the prepared statement
//
// RETURNS
// *sql.Stmt  - the statement -- nil if not in transaction, non-nil otherwise
// 				Note that the caller needs to defer stmt.Close() when stmt is
//              not nil.
// *sql.Row   - the database row to read
// error      - any error encountered
//-----------------------------------------------------------------------------
func getRowsFromDB(ctx context.Context, stmt *sql.Stmt, fields []interface{}) (*sql.Stmt, *sql.Rows, error) {
	var rows *sql.Rows
	var err error
	if tx, ok := TxFromContext(ctx); ok { // if transaction is supplied
		stmt := tx.Stmt(stmt)
		rows, err = stmt.Query(fields...)
		return stmt, rows, err
	}
	rows, err = stmt.Query(fields...)
	return nil, rows, err
}

// insertRowToDB encapsulates about 6 or 7 lines of code to handle inserting a
// database row either in the context of a transaction or not.  These lines
// appear in every Insert call, so it simplifies the code a lot.
//
// INPUTS
// ctx        - db context
// stmt       - the prepared statement to execute
// fields     - the fields to supply to the prepared statement
//
// RETURNS
// *sql.Stmt  - the statement -- nil if not in transaction, non-nil otherwise
// 				note that the caller needs to defer stmt.Close() when stmt is
//              not nil.
// sql.Result - the database Result
// error      - any error encountered
//-----------------------------------------------------------------------------
func insertRowToDB(ctx context.Context, stmt *sql.Stmt, fields []interface{}) (*sql.Stmt, sql.Result, error) {
	var res sql.Result
	var err error
	if tx, ok := TxFromContext(ctx); ok { // if transaction is supplied
		stmt := tx.Stmt(stmt)
		res, err = stmt.Exec(fields...)
		return stmt, res, err
	}
	res, err = stmt.Exec(fields...)
	return nil, res, err
}

// genericInsert encapsulates the code to insert a new record.
//
// INPUTS
// ctx        - db context
// stmt       - the prepared statement to execute
// fields     - the fields to supply to the prepared statement
//
// RETURNS
// crid       - creator id
// upid       - modifier id
// id         - id of the record inserted
// error      - any error encountered
//-----------------------------------------------------------------------------
func genericInsert(ctx context.Context, s string, g *sql.Stmt, fields []interface{}, a interface{}) (int64, int64, int64, error) {
	var crid, upid, id int64
	var err error
	if err = ValidateSessionForDBInsert(ctx, &crid, &upid); err != nil {
		return crid, upid, id, err
	}

	stmt, res, err := insertRowToDB(ctx, g, fields)
	if stmt != nil {
		defer stmt.Close()
	}

	id, err = getIDFromResult(s, res, a, err)
	return crid, upid, id, err
}

// updateDBRow encapsulates about 6 or 7 lines of code to handle updating a
// database row either in the context of a transaction or not.  These lines
// appear in every Update call, so it simplifies the code a lot.
//
// INPUTS
// ctx        - db context
// stmt       - the prepared statement to execute
// fields     - the fields to supply to the prepared statement
//
// RETURNS
// error      - any error encountered
//-----------------------------------------------------------------------------
func updateDBRow(ctx context.Context, stmt *sql.Stmt, fields []interface{}) error {
	var err error

	if tx, ok := TxFromContext(ctx); ok { // if transaction is supplied
		stmt := tx.Stmt(stmt)
		//defer stmt.Close()
		_, err = stmt.Exec(fields...)
	} else {
		_, err = stmt.Exec(fields...)
	}
	return err
}

// genericUpdate encapsulates the code to update an existing record.
//
// INPUTS
// ctx        - db context
// stmt       - the prepared statement to execute
// fields     - the fields to supply to the prepared statement
//
// RETURNS
// id         - LastModBy id
// error      - any error encountered
//-----------------------------------------------------------------------------
func genericUpdate(ctx context.Context, g *sql.Stmt, fields []interface{}) (int64, error) {
	var err error
	var id int64
	if err = ValidateSessionForDBUpdate(ctx, &id); err != nil {
		return id, err
	}
	return id, updateDBRow(ctx, g, fields)
}

func getIDFromResult(s string, res sql.Result, a interface{}, err error) (int64, error) {
	if nil == err {
		x, err := res.LastInsertId()
		if err == nil {
			id := int64(x)
			return id, nil
		}
	}

	return 0, insertError(err, s, a)
}
