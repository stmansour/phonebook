package main

import "fmt"

// Crumb is a screen the user is visiting. A list of Crumbs provides a path indicating
// where the user is in the screen hierarchy
type Crumb struct {
	URL  string
	Name string
}

func breadcrumbBack(sess *session, n int) string {
	var s string
	m := len(sess.Breadcrumbs)
	if n <= m {
		s = sess.Breadcrumbs[m-n].URL
		sess.Breadcrumbs = sess.Breadcrumbs[0 : m-n]
	} else {
		s = "/search/"
	}
	return s
}

func breadcrumbToString(sess *session) string {
	L := len(sess.Breadcrumbs)
	if L < 1 {
		return ""
	}
	s := sess.Breadcrumbs[0].Name
	for i := 1; i < len(sess.Breadcrumbs); i++ {
		s += fmt.Sprintf(" / %s", sess.Breadcrumbs[i].Name)
	}
	return s
}

func breadcrumbToHTMLString(sess *session) string {
	var s string
	L := len(sess.Breadcrumbs)
	if L < 1 {
		return ""
	}
	if L == 1 {
		s = sess.Breadcrumbs[0].Name
	} else {
		s = fmt.Sprintf("<a href=\"/pop/%d\">%s</a>", L, sess.Breadcrumbs[0].Name)
	}
	for i := 1; i < L; i++ {
		if i == L-1 {
			s += " / " + sess.Breadcrumbs[i].Name
		} else {
			s += fmt.Sprintf(" / <a href=\"/pop/%d\">%s</a>", L-i, sess.Breadcrumbs[i].Name)
		}
	}
	return s
}

func breadcrumbAdd(sess *session, name string, url string) {
	c := Crumb{url, name}
	sess.Breadcrumbs = append(sess.Breadcrumbs, c)
}

func breadcrumbReset(sess *session, name string, url string) {
	sess.Breadcrumbs = make([]Crumb, 0)
	breadcrumbAdd(sess, name, url)
}

func getBreadcrumb(token string) string {
	s, ok := sessions[token]
	if !ok {
		fmt.Printf("getBreadcrumb:  Could not find session for %s\n", token)
		return "-/-"
	}
	return breadcrumbToString(s)
}

func getHTMLBreadcrumb(token string) string {
	s, ok := sessions[token]
	if !ok {
		fmt.Printf("getBreadcrumb:  Could not find session for %s\n", token)
		return "-/-"
	}
	return breadcrumbToHTMLString(s)
}
