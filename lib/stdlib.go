package lib

import (
	"crypto/sha512"
	"database/sql"
	"extres"
	"fmt"
	"log"
	"math/rand"
	"runtime/debug"
	"strconv"
	"strings"
)

// JSONDATETIME et. al., are globally available constants
const (
	JSONDATETIME = "2006-01-02T15:04:00Z"
)

// // Errcheck simplifies error handling by putting all the generic
// // code in one place.
// func Errcheck(err error) {
// 	if nil != err {
// 		debug.PrintStack()
// 		fmt.Printf("Error = %v\n", err)
// 		os.Exit(1)
// 	}
// }

// IsSQLNoResultsError returns true if the error provided is a sql err indicating no rows in the solution set.
func IsSQLNoResultsError(err error) bool {
	return err == sql.ErrNoRows
}

// Errcheck - saves a bunch of typing, prints error if it exists
//            and provides a traceback as well
// Note that the error is printed only if the environment is NOT production.
func Errcheck(err error) {
	if err != nil {
		if IsSQLNoResultsError(err) {
			return
		}
		if extres.APPENVPROD != AppConfig.Env {
			fmt.Printf("error = %v\n", err)
		}
		debug.PrintStack()
		log.Fatal(err)
	}
}

// Ulog is Phonebooks's standard logger
func Ulog(format string, a ...interface{}) {
	p := fmt.Sprintf(format, a...)
	log.Print(p)
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
	update, err := db.Prepare("UPDATE people SET passhash=? WHERE UserName=?")
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
