package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Constant values for employment status as well as any yes/no value
const (
	INACTIVE = 0
	ACTIVE   = 1
	NO       = 0
	YES      = 1
)

func dumpTestErrors(tr *TestResults) {
	for i := 0; i < len(tr.Failures); i++ {
		fmt.Printf("%2d. %s\n", tr.Failures[i].Index, tr.Failures[i].TestName)
		fmt.Printf("    Context: %s\n", tr.Failures[i].Context)
		fmt.Printf("    Reason: %s\n", tr.Failures[i].Reason)
	}
}

func stripchars(str string, chr string) string {
	return strings.Map(func(r rune) rune {
		if strings.IndexRune(chr, r) < 0 {
			return r
		}
		return -1
	}, str)
}

func errcheck(err error) {
	if nil != err {
		fmt.Printf("Error = %#v\n", err)
		// PrintStack()
		os.Exit(1)
	}
}

func yesnoToInt(s string) int64 {
	s = strings.ToUpper(s)
	switch {
	case s == "Y" || s == "YES":
		return YES
	case s == "N" || s == "NO":
		return NO
	default:
		fmt.Printf("Unrecognized yes/no response: %s. Returning default = No\n", s)
		return NO
	}
}

func yesnoToString(i int64) string {
	switch {
	case i == NO:
		return "No"
	case i == YES:
		return "Yes"
	default:
		fmt.Printf("Value for yes/no out of range: %d. Returning default = No\n", i)
		return "No"
	}
}

func activeToInt(s string) int64 {
	s = strings.ToUpper(s)
	switch {
	case s == "ACTIVE":
		return ACTIVE
	case s == "INACTIVE" || s == "IN-ACTIVE" || s == "NOTACTIVE" || s == "NOT-ACTIVE":
		return INACTIVE
	default:
		fmt.Printf("Unrecognized yes/no response: %s. Returning default = Inactive\n", s)
		return NO
	}
}

func activeToString(i int64) string {
	switch {
	case i == INACTIVE:
		return "Inactive"
	case i == ACTIVE:
		return "Active"
	default:
		fmt.Printf("Value for yes/no out of range: %d. Returning default = Inactive\n", i)
		return "Inactive"
	}
}

// PBDateFmt specifies the format of dates for all user facing dates
var PBDateFmt = string("2006-01-02")

// PBDateSaveFmt specifies the format of dates for database write operations
var PBDateSaveFmt = string("2006-01-02")

func dateToString(d time.Time) string {
	if d.Year() < 1970 {
		return "N/A"
	}
	return d.Format(PBDateFmt)
}

// Returns 1(jan) thru 12() if the string matches
// if not, it returns 0
var fmtMonths = []string{
	"January", "February", "March", "April",
	"May", "June", "July", "August",
	"September", "October", "November", "December",
}

func monthStringToInt(s string) int64 {
	for i := int64(0); i < int64(len(fmtMonths)); i++ {
		if fmtMonths[i][0:3] == s[0:3] {
			return i + 1
		}
	}
	return 0
}

func dateToDBStr(d time.Time) string {
	if d.Year() < 1970 {
		return "0000-00-00"
	}
	return d.Format(PBDateSaveFmt)
}

func dateYear(d time.Time) int64 {
	return int64(d.Year())
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

// ACPTUNKNOWN - x are general values for Yes, No, NotApplicable, Unknown
const (
	ACPTUNKNOWN = 0           // no selection
	ACPTYES     = 1           // Yes
	ACPTNO      = 2           // No
	ACPTNOTAPPL = 3           // N/A
	ACPTLAST    = ACPTNOTAPPL // loops go from ACPTUNKNOWN to ACPTLAST
)

func acceptTypeToInt(s string) int64 {
	var i int64
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

func acceptIntToString(i int64) string {
	var s string
	switch {
	case i == ACPTUNKNOWN:
		s = "Unknown"
	case i == ACPTYES:
		s = "Yes"
	case i == ACPTNO:
		s = "No"
	case i == ACPTNOTAPPL:
		s = "N/A"
	default:
		fmt.Printf("Unknown acceptance value: %d\n", i)
		s = "Unknown"
	}
	return s
}

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

func deductionStringToInt(s string) int64 {
	var i int64
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

func deductionIntToString(i int64) string {
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

// CTUNSET through CTBYPRODUCTION are constants that
// represent the compensation type. A person will have one or more
// of these records in the compensation table
const (
	CTUNSET        = iota // compensation type is unknown
	CTSALARY              // Salary
	CTHOURLY              // Paid hourly
	CTCOMMISSION          // Paid by commission
	CTBYPRODUCTION        // Paid by production or by piecework
	CTEND                 // all compensation ids are less than this
)

func compensationTypeToInt(s string) int64 {
	var i int64
	s = strings.ToUpper(s)
	s = strings.Replace(s, " ", "", -1)
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

func compensationTypeToString(i int64) string {
	var s string
	switch {
	case i == CTUNSET:
		s = "Unset"
	case i == CTSALARY:
		s = "Salary"
	case i == CTHOURLY:
		s = "Hourly"
	case i == CTCOMMISSION:
		s = "Commission"
	case i == CTBYPRODUCTION:
		s = "By Production"
	default:
		s = "Uknown compensation type"
	}
	return s
}

/****************************************************************************
***  EXAMPLE USAGE OF THE ROUND() FUNCTION
	import (
		"fmt"
		"time"
	)

	func main() {
		samples := []time.Duration{9.63e6, 1.23456789e9, 1.5e9, 1.4e9, -1.4e9, -1.5e9, 8.91234e9, 34.56789e9, 12345.6789e9}
		format := "% 13s % 13s % 13s % 13s % 13s % 13s % 13s\n"
		fmt.Printf(format, "duration", "ms", "0.5s", "s", "10s", "m", "h")
		for _, d := range samples {
			fmt.Printf(
				format,
				d,
				Round(d, time.Millisecond),
				Round(d, 0.5e9),
				Round(d, time.Second),
				Round(d, 10*time.Second),
				Round(d, time.Minute),
				Round(d, time.Hour),
			)
		}
	}

OUTPUT:
     duration            ms          0.5s             s           10s             m             h
       9.63ms          10ms             0             0             0             0             0
  1.23456789s        1.235s            1s            1s             0             0             0
         1.5s          1.5s          1.5s            2s             0             0             0
         1.4s          1.4s          1.5s            1s             0             0             0
        -1.4s         -1.4s         -1.5s           -1s             0             0             0
        -1.5s         -1.5s         -1.5s           -2s             0             0             0
     8.91234s        8.912s            9s            9s           10s             0             0
    34.56789s       34.568s         34.5s           35s           30s          1m0s             0
3h25m45.6789s  3h25m45.679s    3h25m45.5s      3h25m46s      3h25m50s       3h26m0s        3h0m0s
***
****************************************************************************/

// Round is used to reduce the number of digits in a duration.
func Round(d, r time.Duration) time.Duration {
	if r <= 0 {
		return d
	}
	neg := d < 0
	if neg {
		d = -d
	}
	if m := d % r; m+m < r {
		d = d - m
	} else {
		d = d + r - m
	}
	if neg {
		return -d
	}
	return d
}
