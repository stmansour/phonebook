package lib

import (
	"fmt"
	"strings"
	"time"
)

// DateToString rounds the supplied amount to the nearest cent.
func DateToString(t time.Time) string {
	return t.Format("01/02/2006")
}

// RRDATEFMT, et al, are possible date formats
const (
	RRDATEFMT         = "01/02/06"
	RRDATEFMT2        = "1/2/06"
	RRDATEFMT3        = "1/2/2006"
	RRDATEFMT4        = "01/02/2006"
	RRDATEINPFMT      = "2006-01-02"
	RRDATEFMTSQL      = RRDATEINPFMT
	RRDATETIMESQL     = "2006-01-02 15:04:05"
	RRJSUTCDATETIME   = "Mon, 02 Jan 2006 15:04:05 MST"
	RRDATETIMEINPFMT  = "2006-01-02 15:04:05 MST"
	RRDATETIMEFMT     = "2006-01-02T15:04:05Z"
	RRDATETIMEW2UIFMT = "1/2/2006 3:04 pm"
	RRDATEREPORTFMT   = "Jan 2, 2006"
	RRDATETIMERPTFMT  = "Jan 2, 2006 3:04pm MST"
	RRDATERECEIPTFMT  = "January 2, 2006"
)

// AcceptedDateFmts is the array of string formats that StringToDate accepts
var AcceptedDateFmts = []string{
	RRDATEFMT,
	RRDATEFMT2,
	RRDATEFMT3,
	RRDATEFMT4,
	RRDATEINPFMT,
	RRDATEFMTSQL,
	RRDATETIMESQL,
	RRJSUTCDATETIME,
	RRDATETIMEINPFMT,
	RRDATETIMEFMT,
	RRDATETIMEW2UIFMT,
	RRDATEREPORTFMT,
	RRDATETIMERPTFMT,
	RRDATERECEIPTFMT,
}

// StringToDate tries to convert the supplied string to a time.Time value. It will use the
// formats called out in dbtypes.go:  RRDATEFMT, RRDATEINPFMT, RRDATEINPFMT2, ...
//
// for further experimentation, try: https://play.golang.org/p/JNUnA5zbMoz
//----------------------------------------------------------------------------------
func StringToDate(s string) (time.Time, error) {
	// try the ansi std date format first
	var Dt time.Time
	var err error
	s = strings.TrimSpace(s)
	for i := 0; i < len(AcceptedDateFmts); i++ {
		Dt, err = time.Parse(AcceptedDateFmts[i], s)
		if nil == err {
			return Dt, nil
		}
	}
	return Dt, fmt.Errorf("date could not be decoded: %s", s)
}

// EqualDtToJSONDate compares a time.time to a JSONDate
//----------------------------------------------------------------------------
func EqualDtToJSONDate(dt *time.Time, jdt *JSONDate) bool {
	d := time.Time(*jdt) // check for Date change
	return d.Equal(d)
}

// EqualDtToJSONDateTime compares a time.time to a JSONDate
//----------------------------------------------------------------------------
func EqualDtToJSONDateTime(dt *time.Time, jdt *JSONDateTime) bool {
	d := time.Time(*jdt) // check for Date change
	return d.Equal(d)
}
