package db

import (
	"database/sql"
	"phonebook/lib"
	"strings"
)

var mySQLRpl = string("?")
var myRpl = mySQLRpl

// PrepStmts defines the structure of prepared sql statemetns
var PrepStmts struct {
	DeleteSessionCookie       *sql.Stmt
	DeleteExpiredCookies      *sql.Stmt
	GetAllSessionCookies      *sql.Stmt
	GetSessionCookie          *sql.Stmt
	FindMatchingSessionCookie *sql.Stmt
	InsertSessionCookie       *sql.Stmt
	UpdateSessionCookie       *sql.Stmt
	LoginInfo                 *sql.Stmt
	GetImagePath              *sql.Stmt
	GetPeopleTypeDown         *sql.Stmt
	GetPerson                 *sql.Stmt
	GetByUsername             *sql.Stmt
	GetWSPerson               *sql.Stmt
	GetBUTypeDown             *sql.Stmt
	GetBUByBUD                *sql.Stmt
	GetPeople                 *sql.Stmt
	InsertPeople              *sql.Stmt
	UpdatePeople              *sql.Stmt
	DeletePeople              *sql.Stmt
	GetLicense                *sql.Stmt
	GetLicenses               *sql.Stmt
	InsertLicense             *sql.Stmt
	UpdateLicense             *sql.Stmt
	DeleteLicense             *sql.Stmt
}

// GenSQLInsertAndUpdateStrings generates a string suitable for SQL INSERT and UPDATE statements given the fields as used in SELECT statements.
//
//  example:
//	given this string:      "LID,BID,RAID,GLNumber,Status,Type,Name,AcctType,LastModTime,LastModBy"
//  we return these five strings:
//  1)  "BID,RAID,GLNumber,Status,Type,Name,AcctType,LastModBy"                 -- use for SELECT
//  2)  "?,?,?,?,?,?,?,?"  														-- use for INSERT
//  3)  "BID=?RAID=?,GLNumber=?,Status=?,Type=?,Name=?,AcctType=?,LastModBy=?"  -- use for UPDATE
//  4)  "LID,BID,RAID,GLNumber,Status,Type,Name,AcctType,LastModBy", 			-- use for INSERT (no PRIMARYKEY), add "WHERE LID=?"
//  5)  "?,?,?,?,?,?,?,?,?"  													-- use for INSERT (no PRIMARYKEY)
//
// Note that in this convention, we remove LastModTime from insert and update statements (the db is set up to update them by default) and
// we remove the initial ID as that number is AUTOINCREMENT on INSERTs and is not updated on UPDATE.
func GenSQLInsertAndUpdateStrings(s string) (string, string, string, string, string) {
	fields := strings.Split(s, ",")

	// mostly 0th element is ID, but it is not necessary
	s0 := fields[0]
	s2 := fields[1:] // skip the ID

	insertFields := []string{} // fields which are allowed while INSERT
	updateFields := []string{} // fields which are allowed while while UPDATE

	// remove fields which value automatically handled by database while insert and update op.
	for _, fld := range s2 {
		fld = strings.TrimSpace(fld)
		if fld == "" { // if nothing then continue
			continue
		}
		// INSERT FIELDS Inclusion
		if fld != "LastModTime" && fld != "CreateTime" { // remove these fields for INSERT
			insertFields = append(insertFields, fld)
		}
		// UPDATE FIELDS Inclusion
		if fld != "LastModTime" && fld != "CreateTime" && fld != "CreateBy" { // remove these fields for UPDATE
			updateFields = append(updateFields, fld)
		}
	}

	var s3, s4 string
	for i := range insertFields {
		if i == len(insertFields)-1 {
			s3 += myRpl
		} else {
			s3 += myRpl + ","
		}
	}

	for i, uFld := range updateFields {
		if i == len(updateFields)-1 {
			s4 += uFld + "=" + myRpl
		} else {
			s4 += uFld + "=" + myRpl + ","
		}
	}

	// list down insert fields with comma separation
	s = strings.Join(insertFields, ",")

	s5 := s0 + "," + s     // for INSERT where first val is not AUTOINCREMENT
	s6 := s3 + "," + myRpl // for INSERT where first val is not AUTOINCREMENT
	return s, s3, s4, s5, s6
}

// CreatePreparedStmts creates prepared sql statements
func CreatePreparedStmts() {
	var err error
	var s1, s2, s3, flds string

	DB.DBFields = make(map[string]string, 0)

	//==========================================
	// License
	//==========================================
	flds = "LID,UID,State,LicenseNo,FLAGS,LastModTime,LastModBy,CreateTime,CreateBy"

	DB.DBFields["license"] = flds
	PrepStmts.GetLicense, err = DB.DirDB.Prepare("SELECT " + flds + " FROM license WHERE LID=?")
	Errcheck(err)
	PrepStmts.GetLicenses, err = DB.DirDB.Prepare("SELECT " + flds + " FROM license WHERE UID=?")
	Errcheck(err)
	s1, s2, s3, _, _ = GenSQLInsertAndUpdateStrings(flds)
	PrepStmts.InsertLicense, err = DB.DirDB.Prepare("INSERT INTO license (" + s1 + ") VALUES(" + s2 + ")")
	Errcheck(err)
	PrepStmts.UpdateLicense, err = DB.DirDB.Prepare("UPDATE license SET " + s3 + " WHERE UID=?")
	Errcheck(err)
	PrepStmts.DeleteLicense, err = DB.DirDB.Prepare("DELETE FROM license WHERE LID=?")
	Errcheck(err)

	//==========================================
	// People
	//==========================================
	flds = "UID,UserName,LastName,MiddleName,FirstName,PreferredName,Salutation,PositionControlNumber,OfficePhone,OfficeFax,CellPhone,PrimaryEmail,SecondaryEmail,BirthMonth,BirthDoM,HomeStreetAddress,HomeStreetAddress2,HomeCity,HomeState,HomePostalCode,HomeCountry,JobCode,Hire,Termination,MgrUID,DeptCode,CoCode,ClassCode,StateOfEmployment,CountryOfEmployment,EmergencyContactName,EmergencyContactPhone,Status,EligibleForRehire,AcceptedHealthInsurance,AcceptedDentalInsurance,Accepted401K,LastReview,NextReview,passhash,RID,ImagePath,LastModTime,LastModBy,CreateTime,CreateBy"

	DB.DBFields["people"] = flds
	PrepStmts.GetPeople, err = DB.DirDB.Prepare("SELECT " + flds + " FROM people WHERE UID=?")
	Errcheck(err)
	PrepStmts.GetByUsername, err = DB.DirDB.Prepare("SELECT " + flds + " FROM people WHERE Username=?")
	Errcheck(err)
	s1, s2, s3, _, _ = GenSQLInsertAndUpdateStrings(flds)
	PrepStmts.InsertPeople, err = DB.DirDB.Prepare("INSERT INTO people (" + s1 + ") VALUES(" + s2 + ")")
	Errcheck(err)
	PrepStmts.UpdatePeople, err = DB.DirDB.Prepare("UPDATE people SET " + s3 + " WHERE UID=?")
	Errcheck(err)
	PrepStmts.DeletePeople, err = DB.DirDB.Prepare("DELETE FROM people WHERE UID=?")
	Errcheck(err)

	PrepStmts.GetWSPerson, err = DB.DirDB.Prepare("SELECT UID,FirstName,MiddleName,LastName,PreferredName FROM people WHERE UID=?")
	Errcheck(err)

	flds = "UID,UserName,Cookie,DtExpire,UserAgent,IP"
	DB.DBFields["sessions"] = flds
	PrepStmts.InsertSessionCookie, err = DB.DirDB.Prepare("INSERT INTO sessions (" + flds + ") VALUES(?,?,?,?,?,?)")
	lib.Errcheck(err)
	PrepStmts.GetAllSessionCookies, err = DB.DirDB.Prepare("SELECT " + flds + " FROM sessions ORDER BY DtExpire ASC")
	lib.Errcheck(err)
	PrepStmts.GetSessionCookie, err = DB.DirDB.Prepare("SELECT " + flds + " FROM sessions WHERE Cookie=?")
	lib.Errcheck(err)
	PrepStmts.FindMatchingSessionCookie, err = DB.DirDB.Prepare("SELECT " + flds + " FROM sessions WHERE UserName=? AND IP=? AND UserAgent=?")
	lib.Errcheck(err)
	PrepStmts.UpdateSessionCookie, err = DB.DirDB.Prepare("UPDATE sessions SET DtExpire=? WHERE Cookie=?")
	lib.Errcheck(err)
	PrepStmts.DeleteSessionCookie, err = DB.DirDB.Prepare("DELETE FROM sessions WHERE Cookie=?")
	lib.Errcheck(err)
	PrepStmts.DeleteExpiredCookies, err = DB.DirDB.Prepare("DELETE FROM sessions WHERE DtExpire <= ?")
	lib.Errcheck(err)

	PrepStmts.LoginInfo, err = DB.DirDB.Prepare("SELECT uid,firstname,preferredname,PrimaryEmail,passhash,rid FROM people WHERE UserName=?")
	lib.Errcheck(err)

	// get image path from the people table
	PrepStmts.GetImagePath, err = DB.DirDB.Prepare("SELECT ImagePath from people WHERE UID=?")
	lib.Errcheck(err)

	//-----------------------
	// People
	//-----------------------
	PrepStmts.GetPeopleTypeDown, err = DB.DirDB.Prepare("SELECT UID,FirstName,MiddleName,LastName,PreferredName FROM people WHERE FirstName LIKE ? OR MiddleName LIKE ? OR LastName LIKE ? or PreferredName LIKE ? LIMIT ?")
	lib.Errcheck(err)
	PrepStmts.GetPerson, err = DB.DirDB.Prepare("SELECT UID,FirstName,MiddleName,LastName,PreferredName FROM people WHERE UID=?")
	lib.Errcheck(err)

	//--------------------
	// Business Unit...
	//--------------------
	PrepStmts.GetBUTypeDown, err = DB.DirDB.Prepare("SELECT ClassCode,CoCode,Name,Designation FROM classes WHERE Designation LIKE ? ORDER BY Designation ASC LIMIT ?")
	lib.Errcheck(err)
	PrepStmts.GetBUByBUD, err = DB.DirDB.Prepare("SELECT ClassCode,CoCode,Name,Designation,Description FROM classes WHERE Designation=?")
	lib.Errcheck(err)

}
