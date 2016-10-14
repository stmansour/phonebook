package lib

import (
	"crypto/sha512"
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"runtime/debug"
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
