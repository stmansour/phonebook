//  This module does the following:
//		a) reads the security roles from disk and stores them for
//         internal access.
//		b) applies security roles to the element structures

package main

import (
	"phonebook/db"
	"reflect"
	"time"
)

// SecRoleAdmin - SecRoleHR
const (
	SecRoleAdmin = 1
	SecRoleHR    = 2
)

func dumpAccessRoles() {
	for i := 0; i < len(db.Authz.Roles); i++ {
		r := db.Authz.Roles[i]
		ulog("Role %d: %s - %s\n", r.RID, r.Name, r.Descr)
		for j := 0; j < len(r.Perms); j++ {
			f := r.Perms[j]
			ulog("Elem:%d, Field:%s Perm:0x%02x\n", f.Elem, f.Field, f.Perm)
		}
	}
}

func readFieldPerms(r *db.Role) {
	rows, err := Phonebook.prepstmt.readFieldPerms.Query(r.RID)
	errcheck(err)
	defer rows.Close()

	for rows.Next() {
		var f db.FieldPerm
		errcheck(rows.Scan(&f.Elem, &f.Field, &f.Perm, &f.Descr))
		r.Perms = append(r.Perms, f)
	}
	errcheck(rows.Err())
	// for i := 0; i < len(r.Perms); i++ {
	// 	fmt.Printf("%d - %s - 0x%02x = %d\n", r.Perms[i].Elem, r.Perms[i].Field, r.Perms[i].Perm, r.Perms[i].Perm)
	// }
}

func readAccessRoles() {
	rows, err := Phonebook.prepstmt.accessRoles.Query()
	errcheck(err)
	defer rows.Close()

	for rows.Next() {
		var r db.Role
		r.Perms = make([]db.FieldPerm, 0)
		errcheck(rows.Scan(&r.RID, &r.Name, &r.Descr))
		readFieldPerms(&r)
		db.Authz.Roles = append(db.Authz.Roles, r)
	}

	errcheck(rows.Err())
}

//=========================================================================================
// SYNOPSIS:
//      filterSecurityRead filters the data in d based on the permissions provided. If the
//		permissions required are not met, the field is zeroed out.  To meet the requirement
//		the ssn.Sessions permission for this field will be logically anded to the supplied
//      perm.  If the result is non-zero, the condition is met.
// ARGS:
//      d            = the struct we want access to
//		el			 = type of element: db.ELEMPERSON, db.ELEMCOMPANY, db.ELEMCLASS
//   	ssn          = session of the logged in user
//   	permRequired = logical or of the required permissions.  Example db.PERMVIEW | db.PERMOWNERVIEW
//		dataUID      = only used if el == PERSON
// RETURNS:
//      ret val = the permissions found logically ANDed with permRequired.  This can be
//				  useful for determining whether or not to check the OWNER uid to that of
//				  the data being accessed. For example, if the data is accessible because of
//				  db.PERMOWNERMOD, the caller can compare the return value to db.PERMOWNERMOD. If
//				  equal, it needs to further check that the session uid matches the uid of the
//				  data being edited before it allows the edit to proceed.
//   	the data elements for this struct filtered based on the permissions associated
//				  with the logged in user's session
//=========================================================================================
func filterSecurityRead(d interface{}, el int, ssn *db.Session, permRequired int64, dataUID int64) int64 {
	var perm int64
	var ok bool

	sulog("filterSecurityRead: d, permRequired=0x%02x, session: %+v\n", permRequired, ssn)
	pcheck := int64(0)
	val := reflect.ValueOf(d).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)         // this is the struct (Foo)
		n := val.Type().Field(i).Name // variable name for field(i)
		t := field.Type().String()    // the variable type

		sulog("%d. %s\n", i, n)
		// Does this field have the required permissions?
		switch el {
		case db.ELEMPERSON:
			perm, ok = ssn.PMap.Pp[n] // here's the permission we have
		case db.ELEMCOMPANY:
			perm, ok = ssn.PMap.Pco[n] // here's the permission we have
		case db.ELEMCLASS:
			perm, ok = ssn.PMap.Pcl[n] // here's the permission we have
		}
		sulog("    permission found: 0x%02x\n", perm)

		if !ok { // this means that the variable was not found in the access list
			sulog("    field not found, will ignore.\n")
			continue // if it's not there, we can ignore it
		}
		sulog("    field found, checking permissions...\n")
		pcheck = permRequired & perm                                 // and it with the required permissions
		ok = 0 != pcheck                                             // if the result is non-zero, the first test passes
		if el == db.ELEMPERSON && ok && pcheck == db.PERMOWNERVIEW { // if this was an ownerView result...
			ok = int64(dataUID) == ssn.UID // the session uid needs to match the data uid
		}
		if ok {
			sulog("    requested permission granted\n")
		} else if field.IsValid() {
			if field.CanSet() {
				sulog("    no permissions for this field - will zero out...\n")
				switch t {
				case "int":
					sulog("No access to %s, type int, setting to 0\n", n)
					field.SetInt(0)
				case "string":
					sulog("No access to %s, type string, setting to 0 length string\n", n)
					field.SetString("")
				case "time.Time":
					sulog("No access to %s, type time.Time, setting to 0\n", n)
					field.Set(reflect.ValueOf(time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)))
				case "[]int":
					sulog("No access to %s, type []int, setting to 0\n", n)
					field.Set(reflect.ValueOf([]int{}))
				default:
					ulog("filterSecurityRead: unhandled variable type. Name = %s, type = %s\n", n, t)
				}
			}
		}
	}
	return pcheck
}

// PDetFilterSecurityRead is a wrapper around filterSecurityRead
func PDetFilterSecurityRead(d *db.PersonDetail, ssn *db.Session, permRequired int64) {
	filterSecurityRead(d, db.ELEMPERSON, ssn, permRequired, d.UID)
	if d.CoCode == 0 {
		companyInit(&d.Company)
	}
	// fmt.Printf("AFTER security filter d = %+v\n\n", d)
}

//=========================================================================================
// SYNOPSIS:
//      filterSecurityMerge merges the data in dNew with that of d based on the permissions
//      in the sess.Session. The fields in d will be updated to the values
//		in dNew provided the field permission allows it. The net result is that the values
//		in d are merged with values of dNew where it is allowed. The resulting d is
//		suitable for writing back to the database.
// ARGS:
//      d            = base structure, should be pulled from db just prior to calling, no mods
// 		ssn          = session of the logged in user
//      el           = type of element: (1 = people, 2 = company, 3 = class, 4 = service)
// 		permRequired = logical or of the required permissions.  Example db.PERMMOD | db.PERMOWNERMOD
// 		dNew         = an updated version of the same structure as d, it should contain the updates we want to apply to d
// RETURNS:
// 		the data elements for this struct filtered based on the supplied perm value
//=========================================================================================
func filterSecurityMerge(d interface{}, ssn *db.Session, el int, permRequired int64, dNew interface{}, UID int64) {
	val := reflect.ValueOf(d).Elem()
	valNew := reflect.ValueOf(dNew).Elem()
	var perm int64
	var ok bool

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)         // the next field in the structure
		fieldNew := valNew.Field(i)   // the corresponding field in the new structure
		n := val.Type().Field(i).Name // variable name for field(i)
		t := field.Type().String()    // the variable type

		// Do we have the required permissions to update this field?
		switch el {
		case db.ELEMPERSON:
			perm, ok = ssn.PMap.Pp[n] // here's the permission we have
		case db.ELEMCOMPANY:
			perm, ok = ssn.PMap.Pco[n] // here's the permission we have
		case db.ELEMCLASS:
			perm, ok = ssn.PMap.Pcl[n] // here's the permission we have
		}

		if !ok { // !ok here means that the variable was not found in the access list
			// lib.Console("filterSecurityMerge: field %s not covered, skipping\n", n)
			continue // if it's not there, we can ignore it
		}

		// lib.Console("Field = %s, perm = %d, permRequired = %d,  ", n, perm, permRequired)

		pcheck := permRequired & perm // AND it with the required permissions
		ok = 0 != pcheck              // if the result is non-zero, the first test passes
		// lib.Console("pcheck = %d, ok = %t\n", pcheck, ok)
		if el == db.ELEMPERSON && ok && pcheck == db.PERMOWNERMOD { // if we passed it still may be an ownerMOD result...
			ok = int64(UID) == ssn.UID // if so, the session uid needs to match the data uid to proceed
		}

		if ok && field.IsValid() {
			if field.CanSet() {
				// lib.Console("type = %v, will attempt to change!\n", t)
				switch t {
				case "int":
					field.Set(reflect.ValueOf(fieldNew.Interface()))
				case "int64":
					field.Set(reflect.ValueOf(fieldNew.Interface()))
				case "string":
					field.SetString(fieldNew.String())
				case "time.Time":
					field.Set(reflect.ValueOf(fieldNew.Interface()))
				case "[]int":
					field.Set(reflect.ValueOf(fieldNew.Interface()))
				default:
					ulog("filterSecurityRead: unhandled variable type: %s\n", t)
				}
			}
		}
	}
}

// func (d *db.PersonDetail) filterSecurityMerge(ssn *db.Session, permRequired int, dNew *db.PersonDetail) {
// 	filterSecurityMerge(d, ssn, db.ELEMPERSON, permRequired, dNew, d.UID)
// }

// func (c *company) filterSecurityMerge(ssn *db.Session, permRequired int, cNew *company) {
// 	filterSecurityMerge(c, ssn, db.ELEMCOMPANY, permRequired, cNew, 0)
// }

// func (c *db.Class) filterSecurityMerge(ssn *db.Session, permRequired int, cNew *db.Class) {
// 	filterSecurityMerge(c, ssn, db.ELEMCLASS, permRequired, cNew, 0)
// }
