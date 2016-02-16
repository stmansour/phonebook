package main

import "time"

// ReadTotalCounters hits the database to update the total counter
// values across all running instances
func ReadTotalCounters() {
	errcheck(Phonebook.db.QueryRow("select SearchPeople,SearchClasses,SearchCompanies,"+
		"EditPerson,ViewPerson,ViewClass,ViewCompany,"+
		"AdminEditPerson,AdminEditClass,AdminEditCompany,"+
		"DeletePerson,DeleteClass,DeleteCompany,SignIn,Logoff from counters").Scan(
		&TotCounters.SearchPeople, &TotCounters.SearchClasses, &TotCounters.SearchCompanies,
		&TotCounters.EditPerson, &TotCounters.ViewPerson, &TotCounters.ViewClass, &TotCounters.ViewCompany,
		&TotCounters.AdminEditPerson, &TotCounters.AdminEditClass, &TotCounters.AdminEditCompany,
		&TotCounters.DeletePerson, &TotCounters.DeleteClass, &TotCounters.DeleteCompany,
		&TotCounters.SignIn, &TotCounters.Logoff))
}

// CounterDispatcher controls access to shared resources. Any routine that
// needs to work with the counters must get the channel first
func CounterDispatcher() {
	for {
		select {
		case <-Phonebook.ReqCountersMem:
			Phonebook.ReqCountersMemAck <- 1 // tell caller go ahead
			<-Phonebook.ReqCountersMemAck    // block until caller is done with mem
		}
	}
}

// UpdateCountersTable writes the current values in the Counters struct to the database
func UpdateCountersTable() {
	// update, err := Phonebook.db.Prepare("update counters set SearchPeople=?,SearchClasses=?,SearchCompanies=?,EditPerson=?,ViewPerson=?,ViewClass=?,ViewCompany=?,AdminEditPerson=?,AdminEditClass=?,AdminEditCompany=?,DeletePerson=?,DeleteClass=?,DeleteCompany=?,SignIn=?,Logoff=?")
	// errcheck(err)

	// fmt.Printf("Updating Counters.  Current vals: %#v\n", Counters)
	_, err := Phonebook.prepstmt.countersUpdate.Exec(Counters.SearchPeople, Counters.SearchClasses, Counters.SearchCompanies,
		Counters.EditPerson, Counters.ViewPerson, Counters.ViewClass, Counters.ViewCompany,
		Counters.AdminEditPerson, Counters.AdminEditClass, Counters.AdminEditCompany,
		Counters.DeletePerson, Counters.DeleteClass, Counters.DeleteCompany, Counters.SignIn, Counters.Logoff)

	Counters.SearchPeople = 0
	Counters.SearchClasses = 0
	Counters.SearchCompanies = 0
	Counters.EditPerson = 0
	Counters.ViewPerson = 0
	Counters.ViewClass = 0
	Counters.ViewCompany = 0
	Counters.AdminEditPerson = 0
	Counters.AdminEditClass = 0
	Counters.AdminEditCompany = 0
	Counters.DeletePerson = 0
	Counters.DeleteClass = 0
	Counters.DeleteCompany = 0
	Counters.SignIn = 0
	Counters.Logoff = 0

	if nil != err {
		ulog("Error updating counters table: %v\n", err)
	}

	ReadTotalCounters()
	// fmt.Printf("Done: Current vals: %#v\nTotal vals: %#v\n", Counters, TotCounters)
}

// UpdateCounters periodically saves the value of all counters to the database.
func UpdateCounters() {
	for {
		select {
		case <-time.After(time.Duration(Phonebook.CountersUpdateTime) * time.Minute):
			Phonebook.ReqCountersMem <- 1    // ask to access the counters mem, blocks until granted
			<-Phonebook.ReqCountersMemAck    // make sure we got it
			UpdateCountersTable()            // do the db update
			Phonebook.ReqCountersMemAck <- 1 // tell CountersDispatcher we're done with the data
			// fmt.Printf("UpdateCounters completed. Current counters: %+v\n", Counters)
		}
	}
}
