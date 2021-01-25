package main

import (
	"fmt"
	"html/template"
	"phonebook/db"
	"phonebook/lib"
)

func breadcrumbBack(ssn *db.Session, n int) string {
	var s string
	m := len(ssn.Breadcrumbs)
	if n <= m {
		s = ssn.Breadcrumbs[m-n].URL
		ssn.Breadcrumbs = ssn.Breadcrumbs[0 : m-n]
	} else {
		s = "/search/"
	}
	return s
}

func breadcrumbToString(ssn *db.Session) string {
	L := len(ssn.Breadcrumbs)
	if L < 1 {
		return ""
	}
	s := ssn.Breadcrumbs[0].Name
	for i := 1; i < len(ssn.Breadcrumbs); i++ {
		s += fmt.Sprintf(" / %s", ssn.Breadcrumbs[i].Name)
	}
	return s
}

func breadcrumbToHTMLString(ssn *db.Session) template.HTML {
	var s string
	L := len(ssn.Breadcrumbs)
	if L < 1 {
		return ""
	}
	if L == 1 {
		s = ssn.Breadcrumbs[0].Name
	} else {
		s = fmt.Sprintf("<a href=\"/pop/%d\">%s</a>", L, ssn.Breadcrumbs[0].Name)
	}
	for i := 1; i < L; i++ {
		if i == L-1 {
			s += " / " + ssn.Breadcrumbs[i].Name
		} else {
			s += fmt.Sprintf(" / <a href=\"/pop/%d\">%s</a>", L-i, ssn.Breadcrumbs[i].Name)
		}
	}
	return template.HTML(s)
}

func breadcrumbAdd(ssn *db.Session, name string, url string) {
	c := lib.Crumb{URL: url, Name: name}
	ssn.Breadcrumbs = append(ssn.Breadcrumbs, c)
}

func breadcrumbReset(ssn *db.Session, name string, url string) {
	ssn.Breadcrumbs = make([]lib.Crumb, 0)
	breadcrumbAdd(ssn, name, url)
}

func getBreadcrumb(token string) string {
	s, ok := db.Sessions[token]
	if !ok {
		fmt.Printf("getBreadcrumb:  Could not find db.Session for %s\n", token)
		return "-/-"
	}
	return breadcrumbToString(s)
}

func getHTMLBreadcrumb(token string) template.HTML {
	s, ok := db.Sessions[token]
	if !ok {
		fmt.Printf("getHTMLBreadcrumb:  Could not find db.Session for %s\n", token)
		return "-/-"
	}
	return breadcrumbToHTMLString(s)
}
