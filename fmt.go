package main

import (
	"fmt"
	"strings"
	"time"
)

// PBDateFmt specifies the format of dates for all user facing dates
var PBDateFmt = string("2006-01-02")

// PBDateSaveFmt specifies the format of dates for database write operations
var PBDateSaveFmt = string("2006-01-02")

func dateToString(d time.Time) string {
	if d.Year() < 1900 {
		return "N/A"
	}
	return d.Format(PBDateFmt)
}

func dateToDBStr(d time.Time) string {
	if d.Year() < 1900 {
		return "0000-00-00"
	}
	return d.Format(PBDateSaveFmt)
}

func stringToDate(s string) time.Time {
	var d time.Time
	var e error
	s = strings.ToUpper(s)
	if len(s) == 0 || s == "N/A" || s == "NA" {
		d = time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	} else {
		d, e = time.Parse(PBDateFmt, s)
		if e != nil {
			fmt.Printf("input: %s  -- Date parse error: %v\n", s, e)
			d = time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
		}
	}
	return d
}

// ACPTUNKNONW - x are general values for Yes, No, NotApplicable, Unknown
const (
	ACPTUNKNOWN = iota // no selection
	ACPTYES            // Yes
	ACPTNO             // No
	ACPTNOTAPPL        // N/A
)

func acceptTypeToInt(s string) int {
	var i int
	s = strings.ToUpper(s)
	s = strings.Replace(s, " ", "", -1)
	switch {
	case s == "UNKNOWN":
		i = ACPTUNKNOWN
	case s == "YES" || s == "Y":
		i = ACPTYES
	case s == "NO" || s == "N":
		i = ACPTNO
	case s == "N/A" || s == "NA" || s == "NOTAPPLICABLE":
		i = ACPTNOTAPPL
	default:
		fmt.Printf("Unknown acceptance type: %s\n", s)
		i = ACPTUNKNOWN
	}
	return i
}

func acceptIntToString(i int) string {
	var s string
	switch {
	case i == ACPTUNKNOWN:
		s = "UNKNOWN"
	case i == ACPTYES:
		s = "YES"
	case i == ACPTNO:
		s = "NO"
	case i == ACPTNOTAPPL:
		s = "N/A"
	default:
		fmt.Printf("Unknown acceptance value: %d\n", i)
		s = "UNKNOWN"
	}
	return s
}

// CTUNSET through CTBYPRODUCTION are constants that
// represent the compensation type. A person will have one or more
// of these records in the compensation table
const (
	CTUNSET        = iota // compensation type is unknown
	CTSALARY              // Salary
	CTHOURLY              // Paid hourly
	CTCOMMISSION          // Paid by commission
	CTBYPRODUCTION        // Paid by production or by piecework
)

// DDUNKNOWN through DDTAXES are constants to represent
// the enumerations for Deductions
const (
	DDUNKNOWN      = iota // an unknown deduction
	DD401K                // 401K deduction
	DD401KLOAN            // 401K loan deduction
	DDCHILDSUPPORT        // Child Support deduction
	DDDENTAL              // dental coverage deduction
	DDFSA                 // FSA
	DDGARN                // garnished wages
	DDGROUPLIFE           // group life insurance
	DDHOUSING             // housing deduction
	DDMEDICAL             // medical insurance deducrtion
	DDMISCDED             // misc
	DDTAXES               // taxes
)

func deductionTypeToInt(s string) int {
	var i int
	s = strings.ToUpper(s)
	s = strings.Replace(s, " ", "", -1)
	switch {
	case s == "401K":
		i = DD401K
	case s == "401KLOAN":
		i = DD401KLOAN
	case s == "CHILDSUPPORT":
		i = DDCHILDSUPPORT
	case s == "DENTAL":
		i = DDDENTAL
	case s == "FSA":
		i = DDFSA
	case s == "GARN":
		i = DDGARN
	case s == "GROUPLIFE":
		i = DDGROUPLIFE
	case s == "HOUSING":
		i = DDHOUSING
	case s == "MEDICAL":
		i = DDMEDICAL
	case s == "MISCDED":
		i = DDMISCDED
	case s == "TAXES":
		i = DDTAXES
	default:
		fmt.Printf("Unknown compensation type: %s\n", s)
		i = DDUNKNOWN
	}
	return i
}

func deductionToString(i int) string {
	var s string
	switch {
	case i == DD401K:
		s = "40"
	case i == DD401KLOAN:
		s = "401KLOAN"
	case i == DDCHILDSUPPORT:
		s = "CHILDSUPPORT"
	case i == DDDENTAL:
		s = "DENTAL"
	case i == DDFSA:
		s = "FSA"
	case i == DDGARN:
		s = "GARN"
	case i == DDGROUPLIFE:
		s = "GROUPLIFE"
	case i == DDHOUSING:
		s = "HOUSING"
	case i == DDDENTAL:
		s = "DENTAL"
	case i == DDMEDICAL:
		s = "MEDICAL"
	case i == DDMISCDED:
		s = "MISCDED"
	case i == DDTAXES:
		s = "TAXES"
	default:
		s = "UKNOWN COMPENSATION TYPE"
	}
	return s
}

func compensationTypeToInt(s string) int {
	var i int
	s = strings.ToUpper(s)
	switch {
	case s == "UNSET":
		i = CTUNSET
	case s == "SALARY":
		i = CTSALARY
	case s == "HOURLY":
		i = CTHOURLY
	case s == "COMMISSION":
		i = CTCOMMISSION
	case s == "BYPRODUCTION" || s == "PIECEWORK":
		i = CTBYPRODUCTION
	default:
		fmt.Printf("Unknown compensation type: %s\n", s)
		i = CTUNSET
	}
	return i
}

func compensationTypeToString(i int) string {
	var s string
	switch {
	case i == CTUNSET:
		s = "UNSET"
	case i == CTSALARY:
		s = "SALARY"
	case i == CTHOURLY:
		s = "HOURLY"
	case i == CTCOMMISSION:
		s = "COMMISSION"
	case i == CTBYPRODUCTION:
		s = "BYPRODUCTION"
	default:
		s = "UKNOWN COMPENSATION TYPE"
	}
	return s
}
