package authz

import "phonebook/lib"

//--------------------------------------------------------------------
//  ROLE SECURITY
//--------------------------------------------------------------------
const (
	PERMNONE       = 0      // no permissions to see, view, modify, delete, print, or anything to this field
	PERMVIEW       = 1 << 0 // OK to view   this field for any element (Person, Company, Class)
	PERMCREATE     = 1 << 1 // OK to create   "
	PERMMOD        = 1 << 2 // OK to modify   "
	PERMDEL        = 1 << 3 // OK to delete   "
	PERMPRINT      = 1 << 4 // OK to print    "
	PERMOWNERVIEW  = 1 << 5 // OK for the owner to view this field  (applies to Person elements)
	PERMOWNERMOD   = 1 << 6 // OK for the owner to modify this field
	PERMOWNERPRINT = 1 << 7 // OK for the owner to modify this field
	PERMEXEC       = 1 << 8 // OK to execute

	ELEMPERSON  = 1 // people
	ELEMCOMPANY = 2 // companies
	ELEMCLASS   = 3 // classes
	ELEMPBSVC   = 4 // the executable service
)

// FieldPerm defines how a specific element field can be accessed
type FieldPerm struct {
	Elem  int    // Element: Person, Company, or Class
	Field string // field within the Element
	Perm  int    // 'logical or' of all permissions on this field
	Descr string // description of the field
}

// Role defines a collection of FieldPerms that can be assigned to a person
type Role struct {
	RID   int         // assigned by DB
	Name  string      // role name
	Descr string      // role description
	Perms []FieldPerm // permissions for all fields, all entities
}

// PermMaps provides maps for quick access to a field.
// This is a handy structure for a session.
//-----------------------------------------------------------------------------
type PermMaps struct {
	Urole Role           // user's role
	Pp    map[string]int // quick way to reference person permissions based on field name
	Pco   map[string]int // quick way to reference company permissions based on field name
	Pcl   map[string]int // quick way to reference db.Class permissions based on field name
	Ppr   map[string]int
}

// Authz is the context structure for the authorization framework
//-----------------------------------------------------------------------------
var Authz struct {
	Roles         []Role
	SecurityDebug bool // push security debug messages to the logfile
}

// Init initializes the authorization framework
//-----------------------------------------------------------------------------
func Init(debug bool) {
	Authz.Roles = make([]Role, 0)
	Authz.SecurityDebug = debug
}

// GetRoleInfo populates the PermMaps
//-----------------------------------------------------------------------------
func GetRoleInfo(rid int, s *PermMaps) {
	found := -1
	idx := -1

	// try to find the requested index
	// lib.Ulog("len(Authz.Roles)=%d\n", len(Authz.Roles))
	// lib.Ulog("GetRoleInfo - looking for rid=%d\n", rid)
	for i := 0; i < len(Authz.Roles); i++ {
		// lib.Ulog("Authz.Roles[%d] = %+v\n", i, Authz.Roles[i])
		if rid == Authz.Roles[i].RID {
			found = i
			idx = i
			s.Urole.Name = Authz.Roles[i].Name
			s.Urole.RID = rid
			break
		}
	}

	if found < 0 {
		idx = 0
		lib.Ulog("Did not find rid == %d, all permissions set to read-only\n", rid)
	}

	r := Authz.Roles[idx]
	s.Pp = make(map[string]int)
	s.Pco = make(map[string]int)
	s.Pcl = make(map[string]int)
	s.Ppr = make(map[string]int)

	for i := 0; i < len(r.Perms); i++ {
		var f FieldPerm
		f.Elem = r.Perms[i].Elem
		f.Field = r.Perms[i].Field
		if found < 0 {
			f.Perm = PERMVIEW
		} else {
			f.Perm = r.Perms[i].Perm
		}
		s.Urole.Perms = append(s.Urole.Perms, f)

		// fast access maps:
		switch f.Elem {
		case ELEMPERSON:
			s.Pp[f.Field] = f.Perm
		case ELEMCOMPANY:
			s.Pco[f.Field] = f.Perm
		case ELEMCLASS:
			s.Pcl[f.Field] = f.Perm
		case ELEMPBSVC:
			s.Ppr[f.Field] = f.Perm
		}
	}
}
