package lib

import (
	"crypto/sha512"
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
)

// Errcheck simplifies error handling by putting all the generic
// code in one place.
func Errcheck(err error) {
	if nil != err {
		debug.PrintStack()
		fmt.Printf("Error = %v\n", err)
		os.Exit(1)
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.,?()#@!~|")

// RandPasswordStringRunes returns a random password with n characters
func RandPasswordStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// UpdateUserPassword sets the supplied user's password to the supplied password.
func UpdateUserPassword(user, password string, db *sql.DB) error {
	sha := sha512.Sum512([]byte(password))
	passhash := fmt.Sprintf("%x", sha)
	update, err := db.Prepare("update people set passhash=? where username=?")
	if nil != err {
		return err
	}
	_, err = update.Exec(passhash, user)
	return err
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
			if "" != errmsg {
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
