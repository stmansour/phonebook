package main

import (
	"fmt"
	"net/http"
	"text/template"
)

func searchHandler(w http.ResponseWriter, r *http.Request) {
	var d searchResults
	d.Query = r.FormValue("searchstring")
	if len(d.Query) > 0 {
		s := fmt.Sprintf("select uid,lastname,firstname,jobcode,primaryemail,officephone,cellphone,department from people where status>0 and (lastname like \"%%%s%%\" or firstname like \"%%%s%%\" or primaryemail like \"%%%s%%\") order by lastname,firstname", d.Query, d.Query, d.Query)
		rows, err := Phonebook.db.Query(s)
		errcheck(err)
		defer rows.Close()

		for rows.Next() {
			var m person
			errcheck(rows.Scan(&m.UID, &m.LastName, &m.FirstName, &m.JobCode, &m.PrimaryEmail, &m.OfficePhone, &m.CellPhone, &m.Department))
			d.Matches = append(d.Matches, m)
		}
		errcheck(rows.Err())
	}
	t, _ := template.New("search.html").ParseFiles("search.html")
	err := t.Execute(w, &d)

	if nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
