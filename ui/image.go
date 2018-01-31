package ui

import (
	"fmt"
	"path/filepath"
)

// GetImageFilename returns the path to an image for the specified user.
// TBD:  this needs to be rewritten to return a blobstore url path
//-----------------------------------------------------------------------
func GetImageFilename(uid int) string {
	pat := fmt.Sprintf("pictures/%d.*", uid)
	matches, err := filepath.Glob(pat)
	if err != nil {
		fmt.Printf("filepath.Glob(%s) returned error: %v\n", pat, err)
		return "/images/anon.png"
	}
	if len(matches) > 0 {
		return "/" + matches[0]
	}
	return "/images/anon.png"
}
