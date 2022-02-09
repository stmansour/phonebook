package lib

import (
	"fmt"
	"strconv"
	"strings"
)

// Stripchars returns a string with the characters from chars removed
func Stripchars(str, chars string) string {
	return strings.Map(func(r rune) rune {
		if !strings.ContainsRune(chars, r) {
			return r
		}
		return -1
	}, str)
}

// Tline returns a string of dashes that is the specified length
func Tline(n int) string {
	p := make([]byte, n)
	for i := 0; i < n; i++ {
		p[i] = '-'
	}
	return string(p)
}

// Mkstr returns a string of n of the supplied character that is the specified length
func Mkstr(n int, c byte) string {
	p := make([]byte, n)
	for i := 0; i < n; i++ {
		p[i] = c
	}
	return string(p)
}

// IntFromString converts the supplied string to an int64 value. If there
// is a problem in the conversion, it generates an error message. To suppress
// the error message, pass in "" for errmsg.
func IntFromString(sa string, errmsg string) (int64, error) {
	var n = int64(0)
	s := strings.TrimSpace(sa)
	if len(s) > 0 {
		i, err := strconv.Atoi(s)
		if err != nil {
			if errmsg != "" {
				return 0, fmt.Errorf("IntFromString: %s: %s", errmsg, s)
			}
			return n, err
		}
		n = int64(i)
	}
	return n, nil
}

// FloatFromString converts the supplied string to an int64 value. If there
// is a problem in the conversion, it generates an error message.  If the string
// contains a '%' at the end, it treats the number as a percentage (divides by 100)
func FloatFromString(sa string, errmsg string) (float64, string) {
	var f = float64(0)
	s := strings.TrimSpace(sa)
	i := strings.Index(s, "%")
	if i > 0 {
		s = s[:i]
	}
	if len(s) > 0 {
		x, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return f, fmt.Sprintf("FloatFromString: %s: %s\n", errmsg, sa)
		}
		f = x
	}
	if i > 0 {
		f /= 100.0
	}
	return f, ""
}
