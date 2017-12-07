package ws

import (
	"net/http"
	"phonebook/lib"
	"strings"
)

func showRequestHeaders(r *http.Request) {
	lib.Console("Headers:\n")
	for k, v := range r.Header {
		lib.Console("\t%s: ", k)
		for i := 0; i < len(v); i++ {
			lib.Console("%q  ", v[i])
		}
		lib.Console("\n")
	}
}
func svcDebugTxn(funcname string, r *http.Request) {
	lib.Console("\n%s\n", lib.Mkstr(80, '-'))
	lib.Console("URL:      %s\n", r.URL.String())
	lib.Console("METHOD:   %s\n", r.Method)
	lib.Console("Handler:  %s\n", funcname)
}

func svcDebugURL(r *http.Request, d *ServiceData) {
	//-----------------------------------------------------------------------
	// pathElements: 0         1     2
	// Break up {subservice}/{BUI}/{ID} into an array of strings
	// BID is common to nearly all commands
	//-----------------------------------------------------------------------
	ss := strings.Split(r.RequestURI[1:], "?") // it could be GET command
	pathElements := strings.Split(ss[0], "/")
	lib.Console("\t%s\n", r.URL.String()) // print before we strip it off
	for i := 0; i < len(pathElements); i++ {
		lib.Console("\t\t%d. %s\n", i, pathElements[i])
	}
}

func svcDebugTxnEnd() {
	lib.Console("END\n")
}
