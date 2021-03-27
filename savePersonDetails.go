package main

import (
	"crypto/sha512"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"phonebook/db"
	"phonebook/lib"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func uploadFileCopy(from *multipart.File, toname string) error {
	to, err := os.Create(toname)
	if err != nil {
		ulog("uploadImageFile: Error on os.Create(%s) -- err=%v\n", toname, err)
		return err
	}
	defer to.Close()
	_, err = io.Copy(to, *from)
	if err != nil {
		ulog("savePersonDetailsHandler: Error writing picture file: %v\n", err)
	} else {
		ulog("File uploaded successfully to %s\n", toname)
	}
	return err
}

/*
Require constants for S3 configuration.
*/
const (
	S3Region        = "us-east-1" // This parameter define the region of bucket
	ImageUploadPath = ""          // This parameter define in which folder have to upload image
)

// generateFileName generate file name with uid
func generateFileName(uid int64) string {
	s := fmt.Sprintf("%d", uid)
	return strings.Join([]string{s}, "") // return file name with extension e.g., <uid>
}

// uploadImageFileToS3 upload image to AWS S3 bucket
// params
// fileHeader: It contains header information of file. e.g., filename, file type
// usrfile: It contains file data/information
// uid: UID of user
//
// return
// imagePath : 211.jpg
// imageLocation: <host>/<bucket>/211/jpg
func uploadImageFileToS3(fileHeader *multipart.FileHeader, usrfile multipart.File, uid int64) (string, string) {

	// generate filename to save on s3/db
	filename := generateFileName(uid)

	// setup credential
	// reading credential from the aws config
	creds := credentials.NewStaticCredentials(lib.AppConfig.S3BucketKeyID, lib.AppConfig.S3BucketKey, "")
	_, err := creds.Get()
	if err != nil {
		fmt.Printf("Bad credentials: %s", err)
	}

	// Set up configuration
	cfg := aws.NewConfig().WithRegion(S3Region).WithCredentials(creds)

	// Set up session
	sess, err := session.NewSession(cfg)
	if err != nil {
		fmt.Println(err)
	}

	// Create new s3 instance
	svc := s3.New(sess)

	// Path after the bucket name
	imagePath := path.Join(ImageUploadPath, filename)

	// define parameters to upload image to S3
	params := &s3.PutObjectInput{
		Bucket:               aws.String(lib.AppConfig.S3BucketName),
		Key:                  aws.String(imagePath), // it include filename
		Body:                 usrfile,               // data of file
		ServerSideEncryption: aws.String("AES256"),
		ContentType:          aws.String(fileHeader.Header["Content-Type"][0]),
		CacheControl:         aws.String("max-age=86400"),
		ACL:                  aws.String("public-read"),
	}

	fmt.Printf("*** PutObject image path: %s\n", imagePath)

	// Upload image to s3 bucket
	resp, err := svc.PutObject(params)
	if err != nil {
		fmt.Printf("bad response: %s", err)
	}

	// get image location
	//imageLocation := path.Join(lib.AppConfig.S3BucketHost, lib.AppConfig.S3BucketName, imagePath)
	imageLocation := db.GenerateImageLocation(imagePath)

	ulog("Response of Image Uploading: \n%s\n", awsutil.StringValue(resp))
	ulog("Image location: %s", imageLocation)

	return imagePath, imageLocation
}

func savePersonDetailsHandler(w http.ResponseWriter, r *http.Request) {
	var ssn *db.Session
	var uis uiSupport
	ssn = nil
	if 0 < initHandlerSession(ssn, &uis, w, r) {
		return
	}
	ssn = uis.X
	Phonebook.ReqCountersMem <- 1    // ask to access the shared mem, blocks until granted
	<-Phonebook.ReqCountersMemAck    // make sure we got it
	Counters.EditPerson++            // initialize our data
	Phonebook.ReqCountersMemAck <- 1 // tell Dispatcher we're done with the data

	var d db.PersonDetail
	path := "/savePersonDetails/"
	uidstr := r.RequestURI[len(path):]
	if len(uidstr) == 0 {
		fmt.Fprintf(w, "The RequestURI needs the person's uid. It was not found on the URI:  %s\n", r.RequestURI)
		return
	}
	uid, err := strconv.ParseInt(uidstr, 10, 64)
	if err != nil {
		fmt.Fprintf(w, "Error converting uid to a number: %v. URI: %s\n", err, r.RequestURI)
		return
	}

	//=================================================================
	//  SECURITY
	//=================================================================
	if !ssn.ElemPermsAny(db.ELEMPERSON, db.PERMOWNERMOD) {
		ulog("Permissions refuse savePersonDetails page on userid=%d (%s), role=%s\n", ssn.UID, ssn.Firstname, ssn.PMap.Urole.Name)
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}
	if uid != ssn.UID {
		ulog("Permissions refuse savePersonDetails page on userid=%d (%s), role=%s trying to save for UID=%d\n", ssn.UID, ssn.Firstname, ssn.PMap.Urole.Name, uid)
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}
	d.UID = uid
	adminReadDetails(&d) //read current data
	action := strings.ToLower(r.FormValue("action"))
	if "save" == action {
		d.PreferredName = r.FormValue("PreferredName")
		d.PrimaryEmail = r.FormValue("PrimaryEmail")
		d.OfficePhone = r.FormValue("OfficePhone")
		d.CellPhone = r.FormValue("CellPhone")
		d.EmergencyContactPhone = r.FormValue("EmergencyContactPhone")
		d.EmergencyContactName = r.FormValue("EmergencyContactName")
		d.HomeStreetAddress = r.FormValue("HomeStreetAddress")
		d.HomeStreetAddress2 = r.FormValue("HomeStreetAddress2")
		d.HomeCity = r.FormValue("HomeCity")
		d.HomeState = r.FormValue("HomeState")
		d.HomePostalCode = r.FormValue("HomePostalCode")
		d.HomeCountry = r.FormValue("HomeCountry")

		if 0 == len(d.PreferredName) {
			ssn.Firstname = d.FirstName
		} else {
			ssn.Firstname = d.PreferredName
		}

		//=================================================================
		//  handle image
		//=================================================================
		file, header, err := r.FormFile("picturefile")
		// fmt.Printf("file: %v, header: %v, err: %v\n", file, header, err)
		if nil == err {
			defer file.Close()
			//err = uploadImageFile(header.Filename, &file, uid) // Upload image on local disk space
			imagePath, imageLocation := uploadImageFileToS3(header, file, uid) // Upload image to AWS S3
			if nil != err {
				ulog("uploadImageFile returned error: %v\n", err)
			}

			d.ProfileImagePath = imagePath
			d.ProfileImageURL = imageLocation

			// Store imageLocation into session
			ssn.ImageURL = imageLocation
		} else {
			ulog("err loading picture: %v\n", err)
		}

		//=================================================================
		//  Do the update
		//=================================================================
		_, err = Phonebook.prepstmt.updateMyDetails.Exec(d.PreferredName, d.PrimaryEmail, d.OfficePhone, d.CellPhone,
			d.EmergencyContactName, d.EmergencyContactPhone,
			d.HomeStreetAddress, d.HomeStreetAddress2, d.HomeCity, d.HomeState, d.HomePostalCode, d.HomeCountry,
			ssn.UID, d.ProfileImagePath,
			uid)
		if nil != err {
			errmsg := fmt.Sprintf("savePersonDetailsHandler: Phonebook.prepstmt.updateMyDetails.Exec: err = %v\n", err)
			ulog(errmsg)
			fmt.Println(errmsg)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		password := r.FormValue("password")
		if "" != password {
			sha := sha512.Sum512([]byte(password))
			passhash := fmt.Sprintf("%x", sha)
			_, err = Phonebook.prepstmt.updatePasswd.Exec(passhash, uid)
			if nil != err {
				errmsg := fmt.Sprintf("savePersonDetailsHandler: Phonebook.prepstmt.updatePasswd.Exec: err = %v\n", err)
				ulog(errmsg)
				fmt.Println(errmsg)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
	http.Redirect(w, r, breadcrumbBack(ssn, 2), http.StatusFound)
}
