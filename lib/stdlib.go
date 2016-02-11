package lib

import (
	"fmt"
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
