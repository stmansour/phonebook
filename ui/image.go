package ui

import (
	"fmt"
	"path"
	"path/filepath"
	"phonebook/db"
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

	const (
		S3_BUCKET = "upload-images-test"                  // This parameter define the bucket name in S3 TODO(Akshay): Move this parameter to config file
		S3_HOST   = "https://s3.ap-south-1.amazonaws.com" // define host to access images from the s3. TODO(Akshay): Move this Host to config file
	)

	err := db.PrepStmts.GetImagePath.QueryRow(uid).Scan(&imagePath)
	if err != nil {
		// Something went wrong
		fmt.Println(err) // TODO(Akshay): Proper log statement
	}

	return path.Join(S3_HOST, S3_BUCKET, imagePath)
}
