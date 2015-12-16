package main

import (
	"fmt"
	"time"
)

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
	update, err := Phonebook.db.Prepare("update counters set SearchPeople=?,SearchClasses=?,SearchCompanies=?,EditPerson=?,ViewPerson=?,ViewClass=?,ViewCompany=?,AdminEditPerson=?,AdminEditClass=?,AdminEditCompany=?,DeletePerson=?,DeleteClass=?,DeleteCompany=?")
	errcheck(err)
	_, err = update.Exec(Counters.SearchPeople, Counters.SearchClasses, Counters.SearchCompanies,
		Counters.EditPerson, Counters.ViewPerson, Counters.ViewClass, Counters.ViewCompany,
		Counters.AdminEditPerson, Counters.AdminEditClass, Counters.AdminEditCompany,
		Counters.DeletePerson, Counters.DeleteClass, Counters.DeleteCompany)
	if nil != err {
		ulog("Error updating counters table: %v\n", err)
	}
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
			fmt.Printf("UpdateCounters completed. Current counters: %+v\n", Counters)
		}
	}
}
