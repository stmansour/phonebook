package main

import (
	"fmt"
	"phonebook/lib"
)

func dumpPD(d *personDetail) {
	lib.Ulog("\tUser:  %s, %d\n", d.UserName, d.UID)
	lib.Ulog("\tF,M,L,P:  %s, %s, %s, %s\n", d.FirstName, d.MiddleName, d.LastName, d.PreferredName)
	lib.Ulog("\tJobCode: %d, DeptCode: %d, ClassCode: %d, CoCode: %d\n", d.JobCode, d.DeptCode, d.ClassCode, d.CoCode)
}
func dumpPersonDetails(n int, d, d2 *personDetail, desc string) {
	lib.Ulog("Context= %s (TotalSaveAdminEditCalls = %d): %d PERSON differences:\n", desc, TotalSaveAdminEditCalls, n)
	dumpPD(d)
	dumpPD(d2)
}

func (d *personDetail) matches(d2 *personDetail, tr *TestResults, desc string) bool {
	n := 0 // number of miscompares
	if d.FirstName != d2.FirstName {
		var f TestFailure
		f.Index = 0
		f.TestName = fmt.Sprintf("Miscompare on FirstName:  d(%s) : d2(%s)\n", d.FirstName, d2.FirstName)
		tr.Failures = append(tr.Failures, f)
		n++
	}
	if d.LastName != d2.LastName {
		var f TestFailure
		f.Index = 0
		f.TestName = fmt.Sprintf("Miscompare on LastName:  d(%s) : d2(%s)\n", d.LastName, d2.LastName)
		tr.Failures = append(tr.Failures, f)
		n++
	}
	if d.UserName != d2.UserName {
		var f TestFailure
		f.Index = 0
		f.TestName = fmt.Sprintf("Miscompare on UserName:  d(%s) : d2(%s)\n", d.UserName, d2.UserName)
		tr.Failures = append(tr.Failures, f)
		n++
	}
	if d.PreferredName != d2.PreferredName {
		var f TestFailure
		f.Index = 0
		f.TestName = fmt.Sprintf("Miscompare on PreferredName:  d(%s) : d2(%s)\n", d.PreferredName, d2.PreferredName)
		tr.Failures = append(tr.Failures, f)
		n++
	}
	if d.Status != d2.Status {
		var f TestFailure
		f.Index = 0
		f.TestName = fmt.Sprintf("Miscompare on Status:  d(%d) : d2(%d)\n", d.Status, d2.Status)
		tr.Failures = append(tr.Failures, f)
		n++
	}
	if d.OfficePhone != d2.OfficePhone {
		var f TestFailure
		f.Index = 0
		f.TestName = fmt.Sprintf("Miscompare on OfficePhone:  d(%s) : d2(%s)\n", d.OfficePhone, d2.OfficePhone)
		tr.Failures = append(tr.Failures, f)
		n++
	}
	if d.CellPhone != d2.CellPhone {
		var f TestFailure
		f.Index = 0
		f.TestName = fmt.Sprintf("Miscompare on CellPhone:  d(%s) : d2(%s)\n", d.CellPhone, d2.CellPhone)
		tr.Failures = append(tr.Failures, f)
		n++
	}
	if d.OfficeFax != d2.OfficeFax {
		var f TestFailure
		f.Index = 0
		f.TestName = fmt.Sprintf("Miscompare on OfficeFax:  d(%s) : d2(%s)\n", d.OfficeFax, d2.OfficeFax)
		tr.Failures = append(tr.Failures, f)
		n++
	}
	if d.EmergencyContactName != d2.EmergencyContactName {
		var f TestFailure
		f.Index = 0
		f.TestName = fmt.Sprintf("Miscompare on EmergencyContactName:  d(%s) : d2(%s)\n", d.EmergencyContactName, d2.EmergencyContactName)
		tr.Failures = append(tr.Failures, f)
		n++
	}
	if d.HomeStreetAddress != d2.HomeStreetAddress {
		var f TestFailure
		f.Index = 0
		f.TestName = fmt.Sprintf("Miscompare on HomeStreetAddress:  d(%s) : d2(%s)\n", d.HomeStreetAddress, d2.HomeStreetAddress)
		tr.Failures = append(tr.Failures, f)
		n++
	}
	if d.HomeCity != d2.HomeCity {
		var f TestFailure
		f.Index = 0
		f.TestName = fmt.Sprintf("Miscompare on HomeCity:  d(%s) : d2(%s)\n", d.HomeCity, d2.HomeCity)
		tr.Failures = append(tr.Failures, f)
		n++
	}
	if d.HomeState != d2.HomeState {
		var f TestFailure
		f.Index = 0
		f.TestName = fmt.Sprintf("Miscompare on HomeState:  d(%s) : d2(%s)\n", d.HomeState, d2.HomeState)
		tr.Failures = append(tr.Failures, f)
		n++
	}
	if d.HomePostalCode != d2.HomePostalCode {
		var f TestFailure
		f.Index = 0
		f.TestName = fmt.Sprintf("Miscompare on HomePostalCode:  d(%s) : d2(%s)\n", d.HomePostalCode, d2.HomePostalCode)
		tr.Failures = append(tr.Failures, f)
		n++
	}
	if d.HomeCountry != d2.HomeCountry {
		var f TestFailure
		f.Index = 0
		f.TestName = fmt.Sprintf("Miscompare on HomeCountry:  d(%s) : d2(%s)\n", d.HomeCountry, d2.HomeCountry)
		tr.Failures = append(tr.Failures, f)
		n++
	}
	if d.DeptCode != d2.DeptCode {
		var f TestFailure
		f.Index = 0
		f.TestName = fmt.Sprintf("Miscompare on DeptCode:  d(%d) : d2(%d)\n", d.DeptCode, d2.DeptCode)
		tr.Failures = append(tr.Failures, f)
		n++
	}
	if d.JobCode != d2.JobCode {
		var f TestFailure
		f.Index = 0
		f.TestName = fmt.Sprintf("Miscompare on JobCode:  d(%d) : d2(%d)\n", d.JobCode, d2.JobCode)
		tr.Failures = append(tr.Failures, f)
		n++
	}
	if d.EmergencyContactPhone != d2.EmergencyContactPhone {
		var f TestFailure
		f.Index = 0
		f.TestName = fmt.Sprintf("Miscompare on EmergencyContactPhone:  d(%s) : d2(%s)\n", d.EmergencyContactPhone, d2.EmergencyContactPhone)
		tr.Failures = append(tr.Failures, f)
		n++
	}
	if d.PrimaryEmail != d2.PrimaryEmail {
		var f TestFailure
		f.Index = 0
		f.TestName = fmt.Sprintf("Miscompare on PrimaryEmail:  d(%s) : d2(%s)\n", d.PrimaryEmail, d2.PrimaryEmail)
		tr.Failures = append(tr.Failures, f)
		n++
	}
	if d.SecondaryEmail != d2.SecondaryEmail {
		var f TestFailure
		f.Index = 0
		f.TestName = fmt.Sprintf("Miscompare on SecondaryEmail:  d(%s) : d2(%s)\n", d.SecondaryEmail, d2.SecondaryEmail)
		tr.Failures = append(tr.Failures, f)
		n++
	}
	if d.ClassCode != d2.ClassCode {
		var f TestFailure
		f.Index = 0
		f.TestName = fmt.Sprintf("Miscompare on ClassCode:  d(%d) : d2(%d)\n", d.ClassCode, d2.ClassCode)
		tr.Failures = append(tr.Failures, f)
		n++
	}
	if d.CoCode != d2.CoCode {
		var f TestFailure
		f.Index = 0
		f.TestName = fmt.Sprintf("Miscompare on CoCode:  d(%d) : d2(%d)\n", d.CoCode, d2.CoCode)
		tr.Failures = append(tr.Failures, f)
		n++
	}
	if n > 0 {
		dumpPersonDetails(n, d, d2, desc)
	}
	return (n == 0)
}

func (c *company) matches(c2 *company, tr *TestResults) bool {
	n := 0
	if c.LegalName != c2.LegalName {
		var f TestFailure
		f.Index = 0
		f.TestName = fmt.Sprintf("Miscompare on LegalName:  c(%s) : c2(%s)\n", c.LegalName, c2.LegalName)
		tr.Failures = append(tr.Failures, f)
		n++
	}
	if c.CommonName != c2.CommonName {
		var f TestFailure
		f.Index = 0
		f.TestName = fmt.Sprintf("Miscompare on CommonName:  c(%s) : c2(%s)\n", c.CommonName, c2.CommonName)
		tr.Failures = append(tr.Failures, f)
		n++
	}
	if c.Address != c2.Address {
		var f TestFailure
		f.Index = 0
		f.TestName = fmt.Sprintf("Miscompare on Address:  c(%s) : c2(%s)\n", c.Address, c2.Address)
		tr.Failures = append(tr.Failures, f)
		n++
	}
	if c.Address2 != c2.Address2 {
		var f TestFailure
		f.Index = 0
		f.TestName = fmt.Sprintf("Miscompare on Address2:  c(%s) : c2(%s)\n", c.Address2, c2.Address2)
		tr.Failures = append(tr.Failures, f)
		n++
	}
	if c.City != c2.City {
		var f TestFailure
		f.Index = 0
		f.TestName = fmt.Sprintf("Miscompare on City:  c(%s) : c2(%s)\n", c.City, c2.City)
		tr.Failures = append(tr.Failures, f)
		n++
	}
	if c.State != c2.State {
		var f TestFailure
		f.Index = 0
		f.TestName = fmt.Sprintf("Miscompare on State:  c(%s) : c2(%s)\n", c.State, c2.State)
		tr.Failures = append(tr.Failures, f)
		n++
	}
	if c.PostalCode != c2.PostalCode {
		var f TestFailure
		f.Index = 0
		f.TestName = fmt.Sprintf("Miscompare on PostalCode:  c(%s) : c2(%s)\n", c.PostalCode, c2.PostalCode)
		tr.Failures = append(tr.Failures, f)
		n++
	}
	if c.Country != c2.Country {
		var f TestFailure
		f.Index = 0
		f.TestName = fmt.Sprintf("Miscompare on Country:  c(%s) : c2(%s)\n", c.Country, c2.Country)
		tr.Failures = append(tr.Failures, f)
		n++
	}
	if c.Phone != c2.Phone {
		var f TestFailure
		f.Index = 0
		f.TestName = fmt.Sprintf("Miscompare on Phone:  c(%s) : c2(%s)\n", c.Phone, c2.Phone)
		tr.Failures = append(tr.Failures, f)
		n++
	}
	if c.Fax != c2.Fax {
		var f TestFailure
		f.Index = 0
		f.TestName = fmt.Sprintf("Miscompare on Fax:  c(%s) : c2(%s)\n", c.Fax, c2.Fax)
		tr.Failures = append(tr.Failures, f)
		n++
	}
	if c.Email != c2.Email {
		var f TestFailure
		f.Index = 0
		f.TestName = fmt.Sprintf("Miscompare on Email:  c(%s) : c2(%s)\n", c.Email, c2.Email)
		tr.Failures = append(tr.Failures, f)
		n++
	}
	if c.Designation != c2.Designation {
		var f TestFailure
		f.Index = 0
		f.TestName = fmt.Sprintf("Miscompare on Designation:  c(%s) : c2(%s)\n", c.Designation, c2.Designation)
		tr.Failures = append(tr.Failures, f)
		n++
	}
	if c.Active != c2.Active {
		var f TestFailure
		f.Index = 0
		f.TestName = fmt.Sprintf("Miscompare on Active:  c(%d) : c2(%d)\n", c.Active, c2.Active)
		tr.Failures = append(tr.Failures, f)
		n++
	}
	return (n == 0)
}

func (c *class) matches(c2 *class, tr *TestResults) bool {
	n := 0

	if c.Name != c2.Name {
		var f TestFailure
		f.Index = 0
		f.TestName = fmt.Sprintf("Miscompare on Name: cl(%s) : cl2(%s)\n", c.Name, c2.Name)
		tr.Failures = append(tr.Failures, f)
		n++
	}
	if c.Designation != c2.Designation {
		var f TestFailure
		f.Index = 0
		f.TestName = fmt.Sprintf("Miscompare on Designation: cl(%s) : cl2(%s)\n", c.Designation, c2.Designation)
		tr.Failures = append(tr.Failures, f)
		n++
	}
	if c.Description != c2.Description {
		var f TestFailure
		f.Index = 0
		f.TestName = fmt.Sprintf("Miscompare on Description: cl(%s) : cl2(%s)\n", c.Description, c2.Description)
		tr.Failures = append(tr.Failures, f)
		n++
	}

	return (n == 0)
}
