package ui

import (
	"net/url"
	"phonebook/db"
	"phonebook/lib"
)

// GetImageLocation returns the ImageURL for the specified user.
func GetImageLocation(uid int) string {
	lib.Console("Entered GetImageLocation\n")
	var imagePath string

	defaultImageName := "defaultProfileImage.png"
	u, err := url.Parse(lib.AppConfig.S3BucketHost)
	if err != nil {
		lib.Ulog("Error parsing: %s : %s\n", lib.AppConfig.S3BucketHost, err.Error())
	}
	pic := defaultImageName
	err = db.PrepStmts.GetImagePath.QueryRow(uid).Scan(&imagePath)
	if err == nil && imagePath != "" {
		pic = imagePath
	}
	im := u.String() + "/" + lib.AppConfig.S3BucketName + "/" + pic
	return im
}
