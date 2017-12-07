package ws

import (
	"net/http"
	"phonebook/lib"
)

// SvcDisableConsole disables console messages from printing out
func SvcDisableConsole(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	lib.DisableConsole()
	SvcWriteSuccessResponse(w)
}

// SvcEnableConsole enables console messages to print out
func SvcEnableConsole(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	lib.EnableConsole()
	SvcWriteSuccessResponse(w)
}
