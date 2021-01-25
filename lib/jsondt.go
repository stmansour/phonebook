package lib

import (
	"fmt"
	"time"
)

// TIME0 is the "beginning of time" constant to use when we need
// to set a time far enough in the past so that there won't be a
// date prior issue
var TIME0 = time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

// ENDOFTIME can be used when there is no end time
var ENDOFTIME = time.Date(9999, time.December, 31, 0, 0, 0, 0, time.UTC)

// JSONDate is a wrapper around time.Time. We need it
// in order to be able to control the formatting used
// on the date values sent to the w2ui controls.  Without
// this wrapper, the default time format used by the
// JSON encoder / decoder does not work with the w2ui
// controls
type JSONDate time.Time

var earliestDate = time.Date(1900, time.January, 1, 0, 0, 0, 0, time.UTC)

// MarshalJSON overrides the default time.Time handler and sends
// date strings of the form YYYY-MM-DD. Any date prior to Jan 1, 1900
// is snapped to Jan 1, 1900.
//--------------------------------------------------------------------
func (t *JSONDate) MarshalJSON() ([]byte, error) {
	ts := time.Time(*t)
	if ts.Before(earliestDate) {
		ts = earliestDate
	}
	// val := fmt.Sprintf("\"%s\"", ts.Format("2006-01-02"))
	val := fmt.Sprintf("\"%s\"", ts.Format(RRDATEFMT3))
	return []byte(val), nil
}

// UnmarshalJSON overrides the default time.Time handler and reads in
// date strings of the form YYYY-MM-DD.  Any date prior to Jan 1, 1900
// is snapped to Jan 1, 1900.
//--------------------------------------------------------------------
func (t *JSONDate) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = Stripchars(s, "\"")
	if len(s) == 0 {
		*t = JSONDate(TIME0)
		return nil
	}
	// x, err := time.Parse("2006-01-02", s)
	x, err := StringToDate(s)
	if err != nil {
		return err
	}
	if x.Before(earliestDate) {
		x = earliestDate
	}
	*t = JSONDate(x)
	return nil
}

// JSONDateTime is a wrapper around time.Time. We need it
// in order to be able to control the formatting used
// on the DATETIME values sent to the w2ui controls.  Without
// this wrapper, the default time format used by the
// JSON encoder / decoder does not work with the w2ui
// controls
type JSONDateTime time.Time

// MarshalJSON overrides the default time.Time handler and sends
// date strings of the form YYYY-MM-DD. Any date prior to Jan 1, 1900
// is snapped to Jan 1, 1900.
//--------------------------------------------------------------------
func (t *JSONDateTime) MarshalJSON() ([]byte, error) {
	ts := time.Time(*t)
	if ts.Before(earliestDate) {
		ts = earliestDate
	}
	val := fmt.Sprintf("\"%s\"", ts.Format(RRDATETIMEINPFMT))
	return []byte(val), nil
}

// UnmarshalJSON overrides the default time.Time handler and reads in
// date strings of the form YYYY-MM-DD.  Any date prior to Jan 1, 1900
// is snapped to Jan 1, 1900.
//--------------------------------------------------------------------
func (t *JSONDateTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = Stripchars(s, "\"")
	// x, err := time.Parse("2006-01-02", s)
	if len(s) == 0 {
		*t = JSONDateTime(TIME0)
		return nil
	}
	x, err := StringToDate(s)
	if err != nil {
		return err
	}
	if x.Before(earliestDate) {
		x = earliestDate
	}
	*t = JSONDateTime(x)
	return nil
}

// MarshalJSON deals with XJSONYesNo
//--------------------------------------------------------------------
func (t *XJSONYesNo) MarshalJSON() ([]byte, error) {
	return []byte("\"" + string(*t) + "\""), nil
}

// UnmarshalJSON deals with XJSONYesNo
//--------------------------------------------------------------------
func (t *XJSONYesNo) UnmarshalJSON(b []byte) error {
	*t = XJSONYesNo(Stripchars(string(b), "\""))
	return nil
}
