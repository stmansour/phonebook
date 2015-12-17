package main

import (
	"fmt"
	"math/rand"
	"strings"
)

func randomPhoneNumber() string {
	return fmt.Sprintf("(%d) %3d-%04d", 100+rand.Intn(899), 100+rand.Intn(899), rand.Intn(9999))
}

func randomEmail(lastname string, firstname string) string {
	var providers = []string{"gmail.com", "yahoo.com", "comcast.net", "aol.com", "bdiddy.com", "hotmail.com", "abiz.com"}
	np := len(providers)
	n := rand.Intn(10)
	switch {
	case n < 4:
		return fmt.Sprintf("%s%s%d@%s", firstname[0:1], lastname, rand.Intn(10000), providers[rand.Intn(np)])
	case n > 6:
		return fmt.Sprintf("%s%s%d@%s", firstname, lastname[0:1], rand.Intn(10000), providers[rand.Intn(np)])
	default:
		return fmt.Sprintf("%s%s%d@%s", firstname, lastname, rand.Intn(1000), providers[rand.Intn(np)])
	}
}
func randomCompanyEmail(cn string) string {
	var providers = []string{"gmail.com", "yahoo.com", "comcast.net", "aol.com", "bdiddy.com", "hotmail.com", "abiz.com", "zcorp.com", "belcore.com",
		"netzero.com", "tricore.com", "zephcore.com", "carmelcore.com"}
	np := len(providers)
	n := rand.Intn(10)
	switch {
	case n < 4:
		return fmt.Sprintf("%s%d@%s", cn, rand.Intn(10000), providers[rand.Intn(np)])
	case n > 6:
		return fmt.Sprintf("%s%d@%s", cn[0:1], rand.Intn(10000), providers[rand.Intn(np)])
	default:
		return fmt.Sprintf("%s%d@%s", cn, rand.Intn(1000), providers[rand.Intn(np)])
	}
}

func randomAddress() string {
	return fmt.Sprintf("%d %s", rand.Intn(99999), App.Streets[rand.Intn(len(App.Streets))])
}

func getUsername(firstname, lastname string) string {
	//============================================
	// generate a unique username...
	//============================================
	username := strings.ToLower(firstname[0:1] + lastname)
	if len(username) > 17 {
		username = username[0:17]
	}
	newuser := username
	var xx int
	nUID := 0
	for {
		found := false
		rows, err := App.db.Query("select uid from people where UserName=?", newuser)
		errcheck(err)
		defer rows.Close()
		for rows.Next() {
			errcheck(rows.Scan(&xx))
			nUID++
			found = true
			newuser = fmt.Sprintf("%s%d", username, nUID)
		}
		if !found {
			break
		}
	}
	return strings.ToLower(newuser)
}
