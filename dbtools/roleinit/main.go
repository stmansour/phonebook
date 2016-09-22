package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"phonebook/lib"
	"time"
)
import _ "github.com/go-sql-driver/mysql"

type company struct {
	CoCode           int
	LegalName        string
	CommonName       string
	Address          string
	Address2         string
	City             string
	State            string
	PostalCode       string
	Country          string
	Phone            string
	Fax              string
	Email            string
	Designation      string
	Active           int
	EmploysPersonnel int
}

type person struct {
	UID           int
	LastName      string
	FirstName     string
	PreferredName string
	PrimaryEmail  string
	JobCode       int
	OfficePhone   string
	CellPhone     string
	DeptCode      int
	DeptName      string
	Employer      string
}

type myComp struct {
	CompCode int    // code for this comp type
	Name     string // name for this code
	HaveIt   int    // 0 = does not have it, 1 = has it
}

type aDeduction struct {
	DCode  int    // code for this deduction
	Name   string // name for this deduction
	HaveIt int    // 0 = does not have it, 1 = has it
}

type personDetail struct {
	UID                     int
	LastName                string
	FirstName               string
	PrimaryEmail            string
	JobCode                 int
	OfficePhone             string
	CellPhone               string
	DeptName                string
	MiddleName              string
	Salutation              string
	Status                  int
	PositionControlNumber   string
	OfficeFax               string
	SecondaryEmail          string
	EligibleForRehire       int
	LastReview              time.Time
	NextReview              time.Time
	Birthdate               string
	BirthMonth              int
	BirthDOM                int
	HomeStreetAddress       string
	HomeStreetAddress2      string
	HomeCity                string
	HomeState               string
	HomePostalCode          string
	HomeCountry             string
	StateOfEmployment       string
	CountryOfEmployment     string
	PreferredName           string
	Comps                   []int // an array of CompensationType values (ints)
	SecList                 []int
	CompensationStr         string //used in the admin edit interface
	DeptCode                int
	Company                 company
	CoCode                  int
	MgrUID                  int
	JobTitle                string
	Class                   string
	ClassCode               int
	MgrName                 string
	Image                   string // ptr to image -- URI
	Reports                 []person
	Deductions              []int
	DeductionsStr           string
	EmergencyContactName    string
	EmergencyContactPhone   string
	AcceptedHealthInsurance int
	AcceptedDentalInsurance int
	Accepted401K            int
	Hire                    time.Time
	Termination             time.Time
	NameToCoCode            map[string]int
	NameToJobCode           map[string]int
	AcceptCodeToName        map[int]string
	NameToDeptCode          map[string]int // department name to dept code
	MyComps                 []myComp
	MyDeductions            []aDeduction
}

type class struct {
	ClassCode   int
	Name        string
	Designation string
	Description string
}

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

type FieldPerm struct {
	Elem  int    // Element: Person, Company, or Class
	Field string // field within the Element
	Perm  int    // 'logical or' of all permissions on this field
	Descr string // description of the field
}

// Role
type Role struct {
	RID   int         // assigned by DB
	Name  string      // role name
	Descr string      // role description
	Perms []FieldPerm // permissions for all fields, all entities
}

var App struct {
	db          *sql.DB
	DBName      string
	DBUser      string
	presetRoles bool
}

func errcheck(err error) {
	if nil != err {
		fmt.Printf("Error = %v\n", err)
		os.Exit(1)
	}
}

func createRoleTables(db *sql.DB) {
	fmt.Printf("Creating new tables for roles, fieldperms\n")
	ps, err := db.Prepare("DROP TABLE IF EXISTS roles,fieldperms,role,fieldperm")
	errcheck(err)
	_, err = ps.Exec()
	errcheck(err)
	ps, err = db.Prepare("CREATE TABLE roles (RID MEDIUMINT NOT NULL AUTO_INCREMENT,Name VARCHAR(25),Descr VARCHAR(512), PRIMARY KEY (RID))")
	errcheck(err)
	_, err = ps.Exec()
	errcheck(err)
	ps, err = db.Prepare("CREATE TABLE fieldperms (RID MEDIUMINT NOT NULL,Elem MEDIUMINT NOT NULL,Field VARCHAR(25) NOT NULL,Perm MEDIUMINT NOT NULL,Descr VARCHAR(256))")
	errcheck(err)
	_, err = ps.Exec()
	errcheck(err)
}

func makeNewRole(db *sql.DB, r *Role) {
	insert, err := db.Prepare("INSERT INTO roles (Name,Descr) VALUES(?,?)")
	errcheck(err)
	_, err = insert.Exec(r.Name, r.Descr)
	errcheck(err)

	// TODO: ensure the name is unique

	// get the RID
	errcheck(db.QueryRow("select RID from roles where Name=?", r.Name).Scan(&r.RID))

	insert, err = db.Prepare("INSERT INTO fieldperms (RID,Elem,Field,Perm,Descr) VALUES(?,?,?,?,?)")
	errcheck(err)

	// fmt.Println(r.Name)

	// add the FieldPerm...
	for i := 0; i < len(r.Perms); i++ {
		f := r.Perms[i]
		_, err = insert.Exec(r.RID, f.Elem, f.Field, f.Perm, f.Descr)
		errcheck(err)
		// fmt.Printf("RID=%d   Elem=%d   Field=%s   Perm=0x%02x = %d\n", r.RID, f.Elem, f.Field, f.Perm, f.Perm)
	}
}

func makeDefaultRoles(db *sql.DB) {
	AdministratorPerms := []FieldPerm{
		{ELEMPERSON, "Status", PERMVIEW | PERMCREATE | PERMMOD | PERMPRINT, "Indicates whether the person is an active employee."},
		{ELEMPERSON, "EligibleForRehire", PERMVIEW | PERMCREATE | PERMMOD | PERMPRINT, "Indicates whether a past employee can be rehired."},
		{ELEMPERSON, "UID", PERMVIEW | PERMCREATE | PERMPRINT, "A unique identifier associated with the employee. Once created, it never changes."},
		{ELEMPERSON, "Salutation", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "'Mr.', 'Mrs.', 'Ms.', etc."},
		{ELEMPERSON, "FirstName", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "The person's common name."},
		{ELEMPERSON, "MiddleName", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "The person's middle name."},
		{ELEMPERSON, "LastName", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "The person's surname or last name."},
		{ELEMPERSON, "PreferredName", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "Less formal name but more commonly used, for example 'Mike' rather than 'Michael'."},
		{ELEMPERSON, "PrimaryEmail", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "The primary email address to use for this person."},
		{ELEMPERSON, "OfficePhone", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "This person's office telephone number."},
		{ELEMPERSON, "CellPhone", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "This person's cellphone number."},
		{ELEMPERSON, "EmergencyContactName", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "Name of someone to contact in the event of an emergency."},
		{ELEMPERSON, "EmergencyContactPhone", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "Phone number for the emergency contact."},
		{ELEMPERSON, "HomeStreetAddress", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "HomeStreetAddress2", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "HomeCity", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "HomeState", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "HomePostalCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "HomeCountry", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "PrimaryEmail", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "SecondaryEmail", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "OfficePhone", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "OfficeFax", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "CellPhone", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "BirthDOM", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "BirthMonth", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "CoCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "The company code associated with this user."},
		{ELEMPERSON, "JobCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "ClassCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "DeptCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "PositionControlNumber", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "MgrUID", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "Accepted401K", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "AcceptedDentalInsurance", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "AcceptedHealthInsurance", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "Hire", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "Termination", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "LastReview", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "NextReview", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "StateOfEmployment", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "CountryOfEmployment", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "Comps", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "Deductions", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "MyDeductions", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "RID", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "Role", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "Permissions role"},
		{ELEMPERSON, "ElemEntity", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "Permissions to delete the entity"},
		{ELEMCOMPANY, "CoCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "LegalName", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "CommonName", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "Address", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "Address2", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "City", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "State", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "PostalCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "Country", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "Phone", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "Fax", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "Email", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "Designation", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "Active", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "EmploysPersonnel", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "ElemEntity", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCLASS, "ClassCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCLASS, "CoCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "The parent company for this business unit"},
		{ELEMCLASS, "Name", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCLASS, "Designation", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCLASS, "Description", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCLASS, "ElemEntity", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPBSVC, "Shutdown", PERMEXEC, "Permission to shutdown the service"},
		{ELEMPBSVC, "Restart", PERMEXEC, "Permission to restart the service"},
	}
	r := Role{1, "Administrator", "This role has permission to do everything", AdministratorPerms}
	makeNewRole(db, &r)

	HRPerms := []FieldPerm{
		{ELEMPERSON, "Status", PERMVIEW | PERMCREATE | PERMMOD | PERMPRINT, "Indicates whether the person is an active employee."},
		{ELEMPERSON, "EligibleForRehire", PERMVIEW | PERMCREATE | PERMMOD | PERMPRINT, "Indicates whether a past employee can be rehired."},
		{ELEMPERSON, "UID", PERMVIEW | PERMCREATE | PERMPRINT, "A unique identifier associated with the employee. Once created, it never changes."},
		{ELEMPERSON, "Salutation", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "'Mr.', 'Mrs.', 'Ms.', etc."},
		{ELEMPERSON, "FirstName", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "The person's common name."},
		{ELEMPERSON, "MiddleName", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "The person's middle name."},
		{ELEMPERSON, "LastName", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "The person's surname or last name."},
		{ELEMPERSON, "PreferredName", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "Less formal name but more commonly used, for example 'Mike' rather than 'Michael'."},
		{ELEMPERSON, "PrimaryEmail", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "The primary email address to use for this person."},
		{ELEMPERSON, "OfficePhone", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "This person's office telephone number."},
		{ELEMPERSON, "CellPhone", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "This person's cellphone number."},
		{ELEMPERSON, "EmergencyContactName", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "Name of someone to contact in the event of an emergency."},
		{ELEMPERSON, "EmergencyContactPhone", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "Phone number for the emergency contact."},
		{ELEMPERSON, "HomeStreetAddress", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "HomeStreetAddress2", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "HomeCity", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "HomeState", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "HomePostalCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "HomeCountry", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "PrimaryEmail", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "SecondaryEmail", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "OfficePhone", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "OfficeFax", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "CellPhone", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "BirthDOM", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "BirthMonth", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "CoCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "The company code associated with this user."},
		{ELEMPERSON, "JobCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "DeptCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "ClassCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "PositionControlNumber", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "MgrUID", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "Accepted401K", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "AcceptedDentalInsurance", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "AcceptedHealthInsurance", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "Hire", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "Termination", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "LastReview", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "NextReview", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "StateOfEmployment", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "CountryOfEmployment", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "Comps", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "Deductions", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "MyDeductions", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "Role", PERMVIEW, "Permissions Role"},
		{ELEMPERSON, "RID", PERMVIEW | PERMPRINT, "def"},
		{ELEMPERSON, "ElemEntity", PERMNONE, "Permissions to delete the entity"},
		{ELEMCOMPANY, "CoCode", PERMVIEW | PERMPRINT, "def"},
		{ELEMCOMPANY, "LegalName", PERMVIEW | PERMPRINT, "def"},
		{ELEMCOMPANY, "CommonName", PERMVIEW | PERMPRINT, "def"},
		{ELEMCOMPANY, "Address", PERMVIEW | PERMPRINT, "def"},
		{ELEMCOMPANY, "Address2", PERMVIEW | PERMPRINT, "def"},
		{ELEMCOMPANY, "City", PERMVIEW | PERMPRINT, "def"},
		{ELEMCOMPANY, "State", PERMVIEW | PERMPRINT, "def"},
		{ELEMCOMPANY, "PostalCode", PERMVIEW | PERMPRINT, "def"},
		{ELEMCOMPANY, "Country", PERMVIEW | PERMPRINT, "def"},
		{ELEMCOMPANY, "Phone", PERMVIEW | PERMPRINT, "def"},
		{ELEMCOMPANY, "Fax", PERMVIEW | PERMPRINT, "def"},
		{ELEMCOMPANY, "Email", PERMVIEW | PERMPRINT, "def"},
		{ELEMCOMPANY, "Designation", PERMVIEW | PERMPRINT, "def"},
		{ELEMCOMPANY, "Active", PERMVIEW | PERMPRINT, "def"},
		{ELEMCOMPANY, "EmploysPersonnel", PERMVIEW | PERMPRINT, "def"},
		{ELEMCOMPANY, "ElemEntity", PERMNONE, "def"},
		{ELEMCLASS, "ClassCode", PERMVIEW | PERMPRINT, "def"},
		{ELEMCLASS, "CoCode", PERMVIEW | PERMPRINT, "def"},
		{ELEMCLASS, "Name", PERMVIEW | PERMPRINT, "The parent company for this business unit"},
		{ELEMCLASS, "Designation", PERMVIEW | PERMPRINT, "def"},
		{ELEMCLASS, "Description", PERMVIEW | PERMPRINT, "def"},
		{ELEMCLASS, "ElemEntity", PERMNONE, "def"},
		{ELEMPBSVC, "Shutdown", PERMNONE, "Permission to shutdown the service"},
		{ELEMPBSVC, "Restart", PERMNONE, "Permission to restart the service"},
	}
	r = Role{2, "Human Resources", "This role has full permissions on people, read and print permissions for Companies and Classes.", HRPerms}
	makeNewRole(db, &r)

	FinancePerms := []FieldPerm{
		{ELEMPERSON, "Status", PERMVIEW | PERMPRINT, "Indicates whether the person is an active employee."},
		{ELEMPERSON, "EligibleForRehire", PERMVIEW | PERMCREATE | PERMMOD | PERMPRINT, "Indicates whether a past employee can be rehired."},
		{ELEMPERSON, "UID", PERMVIEW | PERMCREATE | PERMPRINT, "A unique identifier associated with the employee. Once created, it never changes."},
		{ELEMPERSON, "Salutation", PERMVIEW | PERMPRINT, "'Mr.', 'Mrs.', 'Ms.', etc."},
		{ELEMPERSON, "FirstName", PERMVIEW | PERMPRINT, "The person's common name."},
		{ELEMPERSON, "MiddleName", PERMVIEW | PERMPRINT, "The person's middle name."},
		{ELEMPERSON, "LastName", PERMVIEW | PERMPRINT, "The person's surname or last name."},
		{ELEMPERSON, "PreferredName", PERMVIEW | PERMPRINT | PERMOWNERMOD, "Less formal name but more commonly used, for example 'Mike' rather than 'Michael'."},
		{ELEMPERSON, "PrimaryEmail", PERMVIEW | PERMPRINT | PERMOWNERMOD, "The primary email address to use for this person."},
		{ELEMPERSON, "OfficePhone", PERMVIEW | PERMPRINT | PERMOWNERMOD, "This person's office telephone number."},
		{ELEMPERSON, "CellPhone", PERMVIEW | PERMPRINT | PERMOWNERMOD, "This person's cellphone number."},
		{ELEMPERSON, "EmergencyContactName", PERMOWNERVIEW | PERMPRINT | PERMOWNERMOD, "Name of someone to contact in the event of an emergency."},
		{ELEMPERSON, "EmergencyContactPhone", PERMOWNERVIEW | PERMPRINT | PERMOWNERMOD, "Phone number for the emergency contact."},
		{ELEMPERSON, "HomeStreetAddress", PERMOWNERVIEW | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "HomeStreetAddress2", PERMOWNERVIEW | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "HomeCity", PERMOWNERVIEW | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "HomeState", PERMOWNERVIEW | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "HomePostalCode", PERMOWNERVIEW | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "HomeCountry", PERMOWNERVIEW | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "PrimaryEmail", PERMVIEW | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "SecondaryEmail", PERMVIEW | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "OfficePhone", PERMVIEW | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "OfficeFax", PERMVIEW | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "CellPhone", PERMVIEW | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "BirthDOM", PERMOWNERVIEW | PERMPRINT, "def"},
		{ELEMPERSON, "BirthMonth", PERMOWNERVIEW | PERMPRINT, "def"},
		{ELEMPERSON, "CoCode", PERMVIEW | PERMPRINT, "The company code associated with this user."},
		{ELEMPERSON, "JobCode", PERMVIEW | PERMPRINT, "def"},
		{ELEMPERSON, "DeptCode", PERMVIEW | PERMPRINT, "def"},
		{ELEMPERSON, "ClassCode", PERMVIEW | PERMPRINT, "def"},
		{ELEMPERSON, "MgrUID", PERMVIEW | PERMPRINT, "def"},
		{ELEMPERSON, "Accepted401K", PERMVIEW | PERMPRINT, "def"},
		{ELEMPERSON, "AcceptedDentalInsurance", PERMVIEW | PERMPRINT, "def"},
		{ELEMPERSON, "AcceptedHealthInsurance", PERMVIEW | PERMPRINT, "def"},
		{ELEMPERSON, "PositionControlNumber", PERMVIEW | PERMPRINT, "def"},
		{ELEMPERSON, "Hire", PERMOWNERVIEW | PERMPRINT, "def"},
		{ELEMPERSON, "Termination", PERMVIEW | PERMPRINT, "def"},
		{ELEMPERSON, "LastReview", PERMNONE, "def"},
		{ELEMPERSON, "NextReview", PERMNONE, "def"},
		{ELEMPERSON, "StateOfEmployment", PERMVIEW | PERMPRINT, "def"},
		{ELEMPERSON, "CountryOfEmployment", PERMVIEW | PERMPRINT, "def"},
		{ELEMPERSON, "Comps", PERMVIEW | PERMPRINT, "def"},
		{ELEMPERSON, "Deductions", PERMVIEW | PERMPRINT, "def"},
		{ELEMPERSON, "MyDeductions", PERMVIEW | PERMPRINT, "def"},
		{ELEMPERSON, "RID", PERMNONE, "def"},
		{ELEMPERSON, "Role", PERMNONE, "Permissions Role"},
		{ELEMPERSON, "ElemEntity", PERMNONE, "Permissions to delete the entity"},
		{ELEMCOMPANY, "CoCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "LegalName", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "CommonName", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "Address", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "Address2", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "City", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "State", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "PostalCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "Country", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "Phone", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "Fax", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "Email", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "Designation", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "Active", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "EmploysPersonnel", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "ElemEntity", PERMNONE, "def"},
		{ELEMCLASS, "ClassCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCLASS, "CoCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "The parent company for this business unit"},
		{ELEMCLASS, "Name", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCLASS, "Designation", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCLASS, "Description", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCLASS, "ElemEntity", PERMNONE, "def"},
		{ELEMPBSVC, "Shutdown", PERMNONE, "Permission to shutdown the service"},
		{ELEMPBSVC, "Restart", PERMNONE, "Permission to restart the service"},
	}
	r = Role{3, "Finance", "This role has full permissions on Companies and Classes, read and print permissions on People.", FinancePerms}
	makeNewRole(db, &r)

	ROPerms := []FieldPerm{
		{ELEMPERSON, "Status", PERMVIEW, "Indicates whether the person is an active employee."},
		{ELEMPERSON, "EligibleForRehire", PERMVIEW, "Indicates whether a past employee can be rehired."},
		{ELEMPERSON, "UID", PERMVIEW | PERMPRINT, "A unique identifier associated with the employee. Once created, it never changes."},
		{ELEMPERSON, "Salutation", PERMVIEW, "'Mr.', 'Mrs.', 'Ms.', etc."},
		{ELEMPERSON, "FirstName", PERMVIEW, "The person's common name."},
		{ELEMPERSON, "MiddleName", PERMVIEW, "The person's middle name."},
		{ELEMPERSON, "LastName", PERMVIEW, "The person's surname or last name."},
		{ELEMPERSON, "PreferredName", PERMVIEW | PERMOWNERMOD, "Less formal name but more commonly used, for example 'Mike' rather than 'Michael'."},
		{ELEMPERSON, "PrimaryEmail", PERMVIEW | PERMOWNERMOD, "The primary email address to use for this person."},
		{ELEMPERSON, "OfficePhone", PERMVIEW | PERMOWNERMOD, "This person's office telephone number."},
		{ELEMPERSON, "CellPhone", PERMVIEW | PERMOWNERMOD, "This person's cellphone number."},
		{ELEMPERSON, "EmergencyContactName", PERMVIEW | PERMOWNERMOD | PERMOWNERPRINT, "Name of someone to contact in the event of an emergency."},
		{ELEMPERSON, "EmergencyContactPhone", PERMVIEW | PERMOWNERMOD | PERMOWNERPRINT, "Phone number for the emergency contact."},
		{ELEMPERSON, "HomeStreetAddress", PERMVIEW | PERMOWNERMOD | PERMOWNERPRINT, "def"},
		{ELEMPERSON, "HomeStreetAddress2", PERMVIEW | PERMOWNERMOD | PERMOWNERPRINT, "def"},
		{ELEMPERSON, "HomeCity", PERMVIEW | PERMOWNERMOD | PERMOWNERPRINT, "def"},
		{ELEMPERSON, "HomeState", PERMVIEW | PERMOWNERMOD | PERMOWNERPRINT, "def"},
		{ELEMPERSON, "HomePostalCode", PERMVIEW | PERMOWNERMOD | PERMOWNERPRINT, "def"},
		{ELEMPERSON, "HomeCountry", PERMVIEW | PERMOWNERMOD | PERMPRINT, "def"},
		{ELEMPERSON, "PrimaryEmail", PERMVIEW | PERMOWNERMOD | PERMPRINT, "def"},
		{ELEMPERSON, "SecondaryEmail", PERMVIEW | PERMOWNERMOD | PERMPRINT, "def"},
		{ELEMPERSON, "OfficePhone", PERMVIEW | PERMOWNERMOD | PERMPRINT, "def"},
		{ELEMPERSON, "OfficeFax", PERMVIEW | PERMOWNERMOD | PERMPRINT, "def"},
		{ELEMPERSON, "CellPhone", PERMVIEW | PERMOWNERMOD | PERMPRINT, "def"},
		{ELEMPERSON, "BirthDOM", PERMOWNERVIEW | PERMOWNERPRINT, "def"},
		{ELEMPERSON, "BirthMonth", PERMOWNERVIEW | PERMOWNERPRINT, "def"},
		{ELEMPERSON, "CoCode", PERMOWNERVIEW | PERMOWNERPRINT, "The company code associated with this user."},
		{ELEMPERSON, "JobCode", PERMOWNERVIEW | PERMOWNERPRINT, "def"},
		{ELEMPERSON, "DeptCode", PERMOWNERVIEW | PERMOWNERPRINT, "def"},
		{ELEMPERSON, "ClassCode", PERMVIEW | PERMPRINT, "def"},
		{ELEMPERSON, "PositionControlNumber", PERMOWNERVIEW | PERMOWNERPRINT, "def"},
		{ELEMPERSON, "MgrUID", PERMVIEW | PERMPRINT, "def"},
		{ELEMPERSON, "Accepted401K", PERMOWNERVIEW | PERMOWNERPRINT, "def"},
		{ELEMPERSON, "AcceptedDentalInsurance", PERMOWNERVIEW | PERMOWNERPRINT, "def"},
		{ELEMPERSON, "AcceptedHealthInsurance", PERMOWNERVIEW | PERMOWNERPRINT, "def"},
		{ELEMPERSON, "Hire", PERMOWNERVIEW | PERMOWNERPRINT, "def"},
		{ELEMPERSON, "Termination", PERMOWNERVIEW, "def"},
		{ELEMPERSON, "LastReview", PERMOWNERVIEW, "def"},
		{ELEMPERSON, "NextReview", PERMOWNERVIEW, "def"},
		{ELEMPERSON, "StateOfEmployment", PERMOWNERVIEW | PERMOWNERPRINT, "def"},
		{ELEMPERSON, "CountryOfEmployment", PERMOWNERVIEW | PERMOWNERPRINT, "def"},
		{ELEMPERSON, "Comps", PERMOWNERVIEW | PERMOWNERPRINT, "Compensation type(s) for this person."},
		{ELEMPERSON, "Deductions", PERMOWNERVIEW | PERMOWNERPRINT, "The deductions for this person."},
		{ELEMPERSON, "MyDeductions", PERMOWNERVIEW | PERMOWNERPRINT, "The deductions for this person."},
		{ELEMPERSON, "RID", PERMVIEW | PERMPRINT, "def"},
		{ELEMPERSON, "Role", PERMNONE, "Permissions Rol"},
		{ELEMPERSON, "ElemEntity", PERMNONE, "Permissions to delete the entity"},
		{ELEMCOMPANY, "CoCode", PERMVIEW, "def"},
		{ELEMCOMPANY, "LegalName", PERMVIEW, "def"},
		{ELEMCOMPANY, "CommonName", PERMVIEW, "def"},
		{ELEMCOMPANY, "Address", PERMVIEW, "def"},
		{ELEMCOMPANY, "Address2", PERMVIEW, "def"},
		{ELEMCOMPANY, "City", PERMVIEW, "def"},
		{ELEMCOMPANY, "State", PERMVIEW, "def"},
		{ELEMCOMPANY, "PostalCode", PERMVIEW, "def"},
		{ELEMCOMPANY, "Country", PERMVIEW, "def"},
		{ELEMCOMPANY, "Phone", PERMVIEW, "def"},
		{ELEMCOMPANY, "Fax", PERMVIEW, "def"},
		{ELEMCOMPANY, "Email", PERMVIEW, "def"},
		{ELEMCOMPANY, "Designation", PERMVIEW, "def"},
		{ELEMCOMPANY, "Active", PERMVIEW, "def"},
		{ELEMCOMPANY, "EmploysPersonnel", PERMVIEW, "def"},
		{ELEMCLASS, "ClassCode", PERMVIEW, "def"},
		{ELEMCLASS, "CoCode", PERMVIEW, "The parent company for this business unit"},
		{ELEMCLASS, "Name", PERMVIEW, "def"},
		{ELEMCLASS, "Designation", PERMVIEW, "def"},
		{ELEMCLASS, "Description", PERMVIEW, "def"},
		{ELEMCLASS, "ElemEntity", PERMNONE, "def"},
		{ELEMPBSVC, "Shutdown", PERMNONE, "Permission to shutdown the service"},
		{ELEMPBSVC, "Restart", PERMNONE, "Permission to restart the service"},
	}
	r = Role{4, "Viewer", "This role has read-only permissions on everything. Viewers can modify their own information.", ROPerms}
	makeNewRole(db, &r)

	TesterPerms := []FieldPerm{
		{ELEMPERSON, "Status", PERMVIEW | PERMCREATE | PERMMOD | PERMPRINT, "Indicates whether the person is an active employee."},
		{ELEMPERSON, "EligibleForRehire", PERMVIEW, "Indicates whether a past employee can be rehired."},
		{ELEMPERSON, "UID", PERMVIEW | PERMCREATE, "A unique identifier associated with the employee. Once created, it never changes."},
		{ELEMPERSON, "Salutation", PERMMOD, "'Mr.', 'Mrs.', 'Ms.', etc."},
		{ELEMPERSON, "FirstName", PERMDEL, "The person's common name."},
		{ELEMPERSON, "MiddleName", PERMPRINT, "The person's middle name."},
		{ELEMPERSON, "LastName", PERMNONE, "The person's surname or last name."},
		{ELEMPERSON, "PreferredName", PERMVIEW | PERMPRINT, "Less formal name but more commonly used, for example 'Mike' rather than 'Michael'."},
		{ELEMPERSON, "PrimaryEmail", PERMVIEW, "The primary email address to use for this person."},
		{ELEMPERSON, "OfficePhone", PERMNONE, "This person's office telephone number."},
		{ELEMPERSON, "CellPhone", PERMVIEW | PERMCREATE | PERMMOD, "This person's cellphone number."},
		{ELEMPERSON, "EmergencyContactName", PERMNONE, "Name of someone to contact in the event of an emergency."},
		{ELEMPERSON, "EmergencyContactPhone", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "Phone number for the emergency contact."},
		{ELEMPERSON, "HomeStreetAddress", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "HomeStreetAddress2", PERMVIEW, "def"},
		{ELEMPERSON, "HomeCity", PERMNONE, "def"},
		{ELEMPERSON, "HomeState", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "HomePostalCode", PERMNONE, "def"},
		{ELEMPERSON, "HomeCountry", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "PrimaryEmail", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "SecondaryEmail", PERMNONE, "def"},
		{ELEMPERSON, "OfficePhone", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "OfficeFax", PERMNONE, "def"},
		{ELEMPERSON, "CellPhone", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "BirthDOM", PERMNONE, "def"},
		{ELEMPERSON, "BirthMonth", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "CoCode", PERMNONE, "The company code associated with this user."},
		{ELEMPERSON, "JobCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "ClassCode", PERMNONE, "def"},
		{ELEMPERSON, "DeptCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "PositionControlNumber", PERMNONE, "def"},
		{ELEMPERSON, "MgrUID", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "Accepted401K", PERMNONE, "def"},
		{ELEMPERSON, "AcceptedDentalInsurance", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "AcceptedHealthInsurance", PERMNONE, "def"},
		{ELEMPERSON, "Hire", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "Termination", PERMNONE, "def"},
		{ELEMPERSON, "LastReview", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "NextReview", PERMNONE, "def"},
		{ELEMPERSON, "StateOfEmployment", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "CountryOfEmployment", PERMNONE, "def"},
		{ELEMPERSON, "Comps", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "Deductions", PERMVIEW | PERMPRINT, "def"},
		{ELEMPERSON, "MyDeductions", PERMVIEW | PERMPRINT, "def"},
		{ELEMPERSON, "RID", PERMVIEW | PERMPRINT, "def"},
		{ELEMPERSON, "Role", PERMNONE, "Permissions Rol"},
		{ELEMPERSON, "ElemEntity", PERMNONE, "Permissions to delete the entity"},
		{ELEMCOMPANY, "CoCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "LegalName", PERMNONE, "def"},
		{ELEMCOMPANY, "CommonName", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "Address", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "Address2", PERMNONE, "def"},
		{ELEMCOMPANY, "City", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "State", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "PostalCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "Country", PERMNONE, "def"},
		{ELEMCOMPANY, "Phone", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "Fax", PERMNONE, "def"},
		{ELEMCOMPANY, "Email", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "Designation", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "Active", PERMNONE, "def"},
		{ELEMCOMPANY, "EmploysPersonnel", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "ElemEntity", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCLASS, "ClassCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCLASS, "CoCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "The parent company for this business unit"},
		{ELEMCLASS, "Name", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCLASS, "Designation", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCLASS, "Description", PERMNONE, "def"},
		{ELEMCLASS, "ElemEntity", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPBSVC, "Shutdown", PERMEXEC, "Permission to shutdown the service"},
		{ELEMPBSVC, "Restart", PERMEXEC, "Permission to restart the service"},
	}
	r = Role{5, "Tester", "This role is for testing", TesterPerms}
	makeNewRole(db, &r)

	OfficeAdminPerms := []FieldPerm{
		{ELEMPERSON, "Status", PERMVIEW | PERMCREATE | PERMMOD | PERMPRINT, "Indicates whether the person is an active employee."},
		{ELEMPERSON, "EligibleForRehire", PERMVIEW | PERMCREATE | PERMMOD | PERMPRINT, "Indicates whether a past employee can be rehired."},
		{ELEMPERSON, "UID", PERMVIEW | PERMCREATE | PERMPRINT, "A unique identifier associated with the employee. Once created, it never changes."},
		{ELEMPERSON, "Salutation", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "'Mr.', 'Mrs.', 'Ms.', etc."},
		{ELEMPERSON, "FirstName", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "The person's common name."},
		{ELEMPERSON, "MiddleName", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "The person's middle name."},
		{ELEMPERSON, "LastName", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "The person's surname or last name."},
		{ELEMPERSON, "PreferredName", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "Less formal name but more commonly used, for example 'Mike' rather than 'Michael'."},
		{ELEMPERSON, "PrimaryEmail", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "The primary email address to use for this person."},
		{ELEMPERSON, "OfficePhone", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "This person's office telephone number."},
		{ELEMPERSON, "CellPhone", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "This person's cellphone number."},
		{ELEMPERSON, "EmergencyContactName", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "Name of someone to contact in the event of an emergency."},
		{ELEMPERSON, "EmergencyContactPhone", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "Phone number for the emergency contact."},
		{ELEMPERSON, "HomeStreetAddress", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "HomeStreetAddress2", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "HomeCity", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "HomeState", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "HomePostalCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "HomeCountry", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "PrimaryEmail", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "SecondaryEmail", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "OfficePhone", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "OfficeFax", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "CellPhone", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "BirthDOM", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "BirthMonth", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "CoCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "The company code associated with this user."},
		{ELEMPERSON, "JobCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "DeptCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "ClassCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "PositionControlNumber", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "MgrUID", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "Accepted401K", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "AcceptedDentalInsurance", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "AcceptedHealthInsurance", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "Hire", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "Termination", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "LastReview", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "NextReview", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "StateOfEmployment", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "CountryOfEmployment", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "Comps", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "Deductions", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "MyDeductions", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "Role", PERMVIEW, "Permissions Rol"},
		{ELEMPERSON, "RID", PERMVIEW | PERMPRINT, "def"},
		{ELEMPERSON, "ElemEntity", PERMNONE, "Permissions to delete the entity"},
		{ELEMCOMPANY, "CoCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "LegalName", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "CommonName", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "Address", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "Address2", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "City", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "State", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "PostalCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "Country", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "Phone", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "Fax", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "Email", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "Designation", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "Active", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "EmploysPersonnel", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "ElemEntity", PERMNONE, "def"},
		{ELEMCLASS, "ClassCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCLASS, "CoCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "The parent company for this business unit"},
		{ELEMCLASS, "Name", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCLASS, "Designation", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCLASS, "Description", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCLASS, "ElemEntity", PERMNONE, "def"},
		{ELEMPBSVC, "Shutdown", PERMNONE, "Permission to shutdown the service"},
		{ELEMPBSVC, "Restart", PERMNONE, "Permission to restart the service"},
	}
	r = Role{6, "OfficeAdministrator", "This role is both HR and Finance.", OfficeAdminPerms}
	makeNewRole(db, &r)

	OfficeInfoAdminPerms := []FieldPerm{
		{ELEMPERSON, "Status", PERMVIEW | PERMCREATE | PERMMOD | PERMPRINT, "Indicates whether the person is an active employee."},
		{ELEMPERSON, "EligibleForRehire", PERMVIEW | PERMCREATE | PERMMOD | PERMPRINT, "Indicates whether a past employee can be rehired."},
		{ELEMPERSON, "UID", PERMVIEW | PERMCREATE | PERMPRINT, "A unique identifier associated with the employee. Once created, it never changes."},
		{ELEMPERSON, "Salutation", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "'Mr.', 'Mrs.', 'Ms.', etc."},
		{ELEMPERSON, "FirstName", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "The person's common name."},
		{ELEMPERSON, "MiddleName", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "The person's middle name."},
		{ELEMPERSON, "LastName", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "The person's surname or last name."},
		{ELEMPERSON, "PreferredName", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "Less formal name but more commonly used, for example 'Mike' rather than 'Michael'."},
		{ELEMPERSON, "PrimaryEmail", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "The primary email address to use for this person."},
		{ELEMPERSON, "OfficePhone", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "This person's office telephone number."},
		{ELEMPERSON, "CellPhone", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "This person's cellphone number."},
		{ELEMPERSON, "EmergencyContactName", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "Name of someone to contact in the event of an emergency."},
		{ELEMPERSON, "EmergencyContactPhone", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "Phone number for the emergency contact."},
		{ELEMPERSON, "HomeStreetAddress", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "HomeStreetAddress2", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "HomeCity", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "HomeState", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "HomePostalCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "HomeCountry", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "PrimaryEmail", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "SecondaryEmail", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "OfficePhone", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "OfficeFax", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "CellPhone", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT | PERMOWNERMOD, "def"},
		{ELEMPERSON, "BirthDOM", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "BirthMonth", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "CoCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "The company code associated with this user."},
		{ELEMPERSON, "JobCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "DeptCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "ClassCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "PositionControlNumber", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "MgrUID", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "Accepted401K", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "AcceptedDentalInsurance", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "AcceptedHealthInsurance", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "Hire", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "Termination", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "LastReview", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "NextReview", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "StateOfEmployment", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "CountryOfEmployment", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "Comps", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "Deductions", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "MyDeductions", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPERSON, "Role", PERMVIEW, "Permissions Rol"},
		{ELEMPERSON, "RID", PERMVIEW | PERMPRINT, "def"},
		{ELEMPERSON, "ElemEntity", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "Permissions to create/delete the entity"},
		{ELEMCOMPANY, "CoCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "LegalName", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "CommonName", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "Address", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "Address2", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "City", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "State", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "PostalCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "Country", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "Phone", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "Fax", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "Email", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "Designation", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "Active", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "EmploysPersonnel", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCOMPANY, "ElemEntity", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCLASS, "ClassCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCLASS, "CoCode", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "The parent company of this business unit."},
		{ELEMCLASS, "Name", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCLASS, "Designation", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCLASS, "Description", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMCLASS, "ElemEntity", PERMVIEW | PERMCREATE | PERMMOD | PERMDEL | PERMPRINT, "def"},
		{ELEMPBSVC, "Shutdown", PERMNONE, "Permission to shutdown the service"},
		{ELEMPBSVC, "Restart", PERMNONE, "Permission to restart the service"},
	}
	r = Role{7, "OfficeInfoAdministrator", "This role is like Office Administrator but also enables delete.", OfficeInfoAdminPerms}
	makeNewRole(db, &r)
}

// func addRoleToPeople(db *sql.DB) {

// 	// alter, err := db.Prepare("ALTER TABLE people drop column RID")
// 	// _, err = alter.Exec()
// 	// if err != nil {
// 	// 	fmt.Printf("Note: could not add column 'RID'. It may already exist.\n")
// 	// }

// 	// alter, err = db.Prepare("ALTER TABLE people add column RID MEDIUMINT")
// 	// _, err = alter.Exec()
// 	// if err != nil {
// 	// 	fmt.Printf("Note: could not add column 'RID'. It may already exist.\n")
// 	// }

// 	update, err := db.Prepare("Update people set RID=4") // everyone starts with ReadOnly
// 	errcheck(err)
// 	_, err = update.Exec()
// 	errcheck(err)
// 	update, err = db.Prepare("Update people set RID=1 where UID=211") // Steve gets Admin
// 	errcheck(err)
// 	_, err = update.Exec()
// 	errcheck(err)
// 	update, err = db.Prepare("Update people set RID=1 where UID=198") // Joe gets Admin
// 	errcheck(err)
// 	_, err = update.Exec()
// 	errcheck(err)
// 	update, err = db.Prepare("Update people set RID=6 where UID=202") // Stacey gets Office Administrator
// 	errcheck(err)
// 	_, err = update.Exec()
// 	errcheck(err)
// 	update, err = db.Prepare("Update people set RID=7 where UID=200") // Darla gets Finance
// 	errcheck(err)
// 	_, err = update.Exec()
// 	errcheck(err)
// 	update, err = db.Prepare("Update people set RID=5 where UID=3") //  vagers gets Tester
// 	errcheck(err)
// 	_, err = update.Exec()
// 	errcheck(err)
// 	update, err = db.Prepare("Update people set RID=2 where UID=207") //  mwheeler gets HR
// 	errcheck(err)
// 	_, err = update.Exec()
// 	errcheck(err)
// }

func readFieldPerms(db *sql.DB, r *Role) {
	rows, err := db.Query("select Elem,Field,Perm,Descr from fieldperms where RID=?", r.RID)
	errcheck(err)
	defer rows.Close()

	for rows.Next() {
		var f FieldPerm
		errcheck(rows.Scan(&f.Elem, &f.Field, &f.Perm, &f.Descr))
		// fmt.Printf("%d - %s - 0x%02x = %d\n", f.Elem, f.Field, f.Perm, f.Perm)
		// r.Perms = append(r.Perms, f)
	}
	errcheck(rows.Err())
}

func readAccessRoles(db *sql.DB) {
	rows, err := db.Query("select RID,Name,Descr from roles")
	errcheck(err)
	defer rows.Close()

	for rows.Next() {
		var r Role
		r.Perms = make([]FieldPerm, 0)
		errcheck(rows.Scan(&r.RID, &r.Name, &r.Descr))
		// fmt.Println(r.Name)
		readFieldPerms(db, &r)
		// Phonebook.Roles = append(Phonebook.Roles, r)
	}

	errcheck(rows.Err())
}

func readCommandLineArgs() {
	dbuPtr := flag.String("B", "ec2-user", "database user name")
	dbnmPtr := flag.String("N", "accord", "database name (accordtest, accord)")
	rolePtr := flag.Bool("r", false, "List current database roles")

	flag.Parse()

	App.DBUser = *dbuPtr
	App.DBName = *dbnmPtr
	App.presetRoles = *rolePtr
}

func main() {
	readCommandLineArgs()

	var err error
	// s := fmt.Sprintf("%s:@/%s?charset=utf8&parseTime=True", App.DBUser, App.DBName)
	lib.ReadConfig()
	s := lib.GetSQLOpenString(App.DBUser, App.DBName)
	// fmt.Printf("DBOPEN:  %s\n", s)
	App.db, err = sql.Open("mysql", s)
	if nil != err {
		fmt.Printf("sql.Open: Error = %v\n", err)
	}
	defer App.db.Close()

	err = App.db.Ping()
	if nil != err {
		fmt.Printf("App.db.Ping: Error = %v\n", err)
	}
	fmt.Printf("Successfully opened database %s as user %s\n", App.DBName, App.DBUser)

	createRoleTables(App.db)
	makeDefaultRoles(App.db)
	fmt.Printf("Added roles and fieldperms\n")
	//addRoleToPeople(App.db)

	if App.presetRoles {
		readAccessRoles(App.db)
	}
}
