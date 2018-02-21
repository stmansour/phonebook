package ui

import (
	"path"
	"phonebook/db"
	"phonebook/lib"
)

// GetImageLocation returns the ImageURL for the specified user.
func GetImageLocation(uid int) string {
	var imagePath string

	defaultImageName := "defaultProfileImage.png"

	err := db.PrepStmts.GetImagePath.QueryRow(uid).Scan(&imagePath)
	if err != nil {
		lib.Ulog("Error while getting profile imagePath: %s\n", err)
		return path.Join(lib.AppConfig.S3BucketHost, lib.AppConfig.S3BucketName, defaultImageName) // If something went wrong than display default image
	}

	if imagePath != "" {
		return path.Join(lib.AppConfig.S3BucketHost, lib.AppConfig.S3BucketName, imagePath)
	} else {
		return path.Join(lib.AppConfig.S3BucketHost, lib.AppConfig.S3BucketName, defaultImageName) // If database doesn't have imagePath than assign default image
	}

}
