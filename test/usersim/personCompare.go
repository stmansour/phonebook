package main

import "fmt"

func (d *personDetail) matches(d2 *personDetail) bool {
	n := 0 // number of miscompares
	if d.FirstName != d2.FirstName {
		if App.ShowTestMatching {
			fmt.Printf("Miscompare on FirstName:  d(%s) : d2(%s)\n", d.FirstName, d2.FirstName)
		}
		n++
	}
	if d.LastName != d2.LastName {
		if App.ShowTestMatching {
			fmt.Printf("Miscompare on LastName:  d(%s) : d2(%s)\n", d.LastName, d2.LastName)
		}
		n++
	}
	if d.UserName != d2.UserName {
		if App.ShowTestMatching {
			fmt.Printf("Miscompare on UserName:  d(%s) : d2(%s)\n", d.UserName, d2.UserName)
		}
		n++
	}
	if d.PreferredName != d2.PreferredName {
		if App.ShowTestMatching {
			fmt.Printf("Miscompare on PreferredName:  d(%s) : d2(%s)\n", d.PreferredName, d2.PreferredName)
		}
		n++
	}
	if d.Status != d2.Status {
		if App.ShowTestMatching {
			fmt.Printf("Miscompare on Status:  d(%d) : d2(%d)\n", d.Status, d2.Status)
		}
		n++
	}
	if d.OfficePhone != d2.OfficePhone {
		if App.ShowTestMatching {
			fmt.Printf("Miscompare on OfficePhone:  d(%s) : d2(%s)\n", d.OfficePhone, d2.OfficePhone)
		}
		n++
	}
	if d.CellPhone != d2.CellPhone {
		if App.ShowTestMatching {
			fmt.Printf("Miscompare on CellPhone:  d(%s) : d2(%s)\n", d.CellPhone, d2.CellPhone)
		}
		n++
	}
	if d.OfficeFax != d2.OfficeFax {
		if App.ShowTestMatching {
			fmt.Printf("Miscompare on OfficeFax:  d(%s) : d2(%s)\n", d.OfficeFax, d2.OfficeFax)
		}
		n++
	}
	if d.EmergencyContactName != d2.EmergencyContactName {
		if App.ShowTestMatching {
			fmt.Printf("Miscompare on EmergencyContactName:  d(%s) : d2(%s)\n", d.EmergencyContactName, d2.EmergencyContactName)
		}
		n++
	}
	if d.HomeStreetAddress != d2.HomeStreetAddress {
		if App.ShowTestMatching {
			fmt.Printf("Miscompare on HomeStreetAddress:  d(%s) : d2(%s)\n", d.HomeStreetAddress, d2.HomeStreetAddress)
		}
		n++
	}
	if d.HomeCity != d2.HomeCity {
		if App.ShowTestMatching {
			fmt.Printf("Miscompare on HomeCity:  d(%s) : d2(%s)\n", d.HomeCity, d2.HomeCity)
		}
		n++
	}
	if d.HomeState != d2.HomeState {
		if App.ShowTestMatching {
			fmt.Printf("Miscompare on HomeState:  d(%s) : d2(%s)\n", d.HomeState, d2.HomeState)
		}
		n++
	}
	if d.HomePostalCode != d2.HomePostalCode {
		if App.ShowTestMatching {
			fmt.Printf("Miscompare on HomePostalCode:  d(%s) : d2(%s)\n", d.HomePostalCode, d2.HomePostalCode)
		}
		n++
	}
	if d.HomeCountry != d2.HomeCountry {
		if App.ShowTestMatching {
			fmt.Printf("Miscompare on HomeCountry:  d(%s) : d2(%s)\n", d.HomeCountry, d2.HomeCountry)
		}
		n++
	}
	if d.DeptCode != d2.DeptCode {
		if App.ShowTestMatching {
			fmt.Printf("Miscompare on DeptCode:  d(%d) : d2(%d)\n", d.DeptCode, d2.DeptCode)
		}
		n++
	}
	if d.JobCode != d2.JobCode {
		if App.ShowTestMatching {
			fmt.Printf("Miscompare on JobCode:  d(%d) : d2(%d)\n", d.JobCode, d2.JobCode)
		}
		n++
	}
	if d.EmergencyContactPhone != d2.EmergencyContactPhone {
		if App.ShowTestMatching {
			fmt.Printf("Miscompare on EmergencyContactPhone:  d(%s) : d2(%s)\n", d.EmergencyContactPhone, d2.EmergencyContactPhone)
		}
		n++
	}
	if d.PrimaryEmail != d2.PrimaryEmail {
		if App.ShowTestMatching {
			fmt.Printf("Miscompare on PrimaryEmail:  d(%s) : d2(%s)\n", d.PrimaryEmail, d2.PrimaryEmail)
		}
		n++
	}
	if d.SecondaryEmail != d2.SecondaryEmail {
		if App.ShowTestMatching {
			fmt.Printf("Miscompare on SecondaryEmail:  d(%s) : d2(%s)\n", d.SecondaryEmail, d2.SecondaryEmail)
		}
		n++
	}
	if d.ClassCode != d2.ClassCode {
		if App.ShowTestMatching {
			fmt.Printf("Miscompare on ClassCode:  d(%d) : d2(%d)\n", d.ClassCode, d2.ClassCode)
		}
		n++
	}
	if d.CoCode != d2.CoCode {
		if App.ShowTestMatching {
			fmt.Printf("Miscompare on CoCode:  d(%d) : d2(%d)\n", d.CoCode, d2.CoCode)
		}
		n++
	}
	return (n == 0)
}
