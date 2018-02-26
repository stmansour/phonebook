package main

import (
	"fmt"
	"html/template"
	"net/http"
)

// renderTemplate
// It renders provide template with filenames
func renderTemplate(w http.ResponseWriter, ui uiSupport, filenames string) error {
	t, err := template.New("phoneBookBaseTemplate").Funcs(funcMap).ParseFiles("base.html", "header.html", filenames)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = t.ExecuteTemplate(w, "base", &ui)
	return err
}
