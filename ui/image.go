package ui

import (
	"fmt"
	"path"
	"path/filepath"
	"phonebook/db"
	"phonebook/lib"
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

// GetImageLocation returns the ImageURL for the specified user.
func GetImageLocation(uid int) string {
	var imagePath string

	defaultImageName := "usericon.png"

	err := db.PrepStmts.GetImagePath.QueryRow(uid).Scan(&imagePath)
	if err != nil {
		//ulog(err) // TODO(Akshay): Add proper log statment
		return path.Join(lib.AppConfig.S3BucketHost, lib.AppConfig.S3BucketName, defaultImageName) // If something went wrong than display default image
	}

	if imagePath != "" {
		return path.Join(lib.AppConfig.S3BucketHost, lib.AppConfig.S3BucketName, imagePath)
	} else {
		return path.Join(lib.AppConfig.S3BucketHost, lib.AppConfig.S3BucketName, defaultImageName) // If database doesn't have imagePath than assign default image
	}

}
