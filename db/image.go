package db

import (
	"net/url"
	"path"
	"phonebook/lib"
)

// GetImageLocation returns the ImageURL for the specified user.
func GetImageLocation(uid int64) string {
	var imagePath string

	defaultImageName := "defaultProfileImage.png"

	im := defaultImageName // If something went wrong  or database doesn't have imagePath than display default image
	err := PrepStmts.GetImagePath.QueryRow(uid).Scan(&imagePath)
	if imagePath != "" && err == nil {
		im = imagePath
	}
	return GenerateImageLocation(im)
}

// GenerateImageLocation return the image URL from the image path
func GenerateImageLocation(imagePath string) string {
	u, err := url.Parse(lib.AppConfig.S3BucketHost)
	if err != nil {
		lib.Ulog("Error parsing: %s : %s\n", lib.AppConfig.S3BucketHost, err.Error())
	}
	u.Path = path.Join(u.Path, lib.AppConfig.S3BucketName, imagePath)
	return u.String()
}
