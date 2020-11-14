package db

import (
	"phonebook/lib"
)

// WSPerson is the person information we pass over the web service call.
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

// GetPeopleTypeDown returns a slice of session cookies
//
// INPUTS
//
// RETURNS
//  []SessionCookie - a slice with all the rows in the sessions table.
//
//  err      Any errors encountered
//-----------------------------------------------------------------------------
func GetPeopleTypeDown(s1 string, limit int) ([]PeopleTypeDown, error) {
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
		if len(first) > 0 {
			fn = preferred
		}
		p.Name = fn + " "
		if len(middle) > 0 {
			p.Name += middle + " "
		}
		p.Name += last
		m = append(m, p)
	}

	return m, nil
}

// GetWSPerson returns the requested Person
// wsdoc {
//  @Title  Get Person
//	@URL /v1/task/:BUI/TID
//  @Method  GET
//	@Synopsis Get information on a Person
//  @Description  Return all fields for assessment :TID
//	@Input WebGridSearchRequest
//  @Response Person struct, any error encountered
// wsdoc }
//-----------------------------------------------------------------------------
func GetWSPerson(id int64) (WSPerson, error) {
	var a WSPerson
	fields := []interface{}{id}
	row := PrepStmts.GetPerson.QueryRow(fields...)
	err := row.Scan(&a.UID, &a.FirstName, &a.MiddleName, &a.LastName, &a.PreferredName)
	SkipSQLNoRowsError(&err)
	return a, err
}
