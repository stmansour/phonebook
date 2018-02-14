package main

import (
	"crypto/sha512"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"phonebook/authz"
	"phonebook/db"
	"phonebook/sess"
	"strconv"
	"strings"
	"time"
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

// uploadImageFile handles the uploading of a user's picture file and its
// placement in the pictures directory.
//
// Params
// usrfname - name of the file on the user's local system
// usrfile - the open file from the form return in the user's browser
// uid - uid of the user for whom the image applies
//
// Returns:  err from any os file function
func uploadImageFile(usrfname string, usrfile *multipart.File, uid int) error {
	// use the same filetype for the final filename
	ftype := filepath.Ext(usrfname)
	// ulog("user file type: %s\n", ftype)

	// the file name we'll use for this user's picture...
	picturefilename := fmt.Sprintf("pictures/%d%s", uid, ftype)
	// ulog("picturefilename: %s\n", picturefilename)

	//  delete old tmp file if exists:
	tmpFile := fmt.Sprintf("pictures/%d.tmp", uid)
	// ulog("tmpFile to delete if exists: %s\n", tmpFile)

	finfo, err := os.Stat(tmpFile)
	if os.IsNotExist(err) {
		ulog("%s was not found. Nothing to delete\n", tmpFile)
	} else {
		ulog("os.Stat(%s) returns:  err=%v,  finfo=%#v\n", tmpFile, err, finfo)
		err = os.Remove(tmpFile)
		ulog("os.Remove(%s) returns err=%v\n", tmpFile, err)
	}

	// copy the requested file to "<uid>.tmp"
	err = uploadFileCopy(usrfile, tmpFile)
	if nil != err {
		ulog("uploadFileCopy returned error: %v\n", err)
		return err
	}

	// see if there are any files that match the old filename MINUS the filetype...
	m, err := filepath.Glob(fmt.Sprintf("./pictures/%d.*", uid))
	if nil != err {
		ulog("filepath.Glob returned error: %v\n", err)
		return err
	}
	// ulog("filepath.Glob returned the following matches: %v\n", m)
	for i := 0; i < len(m); i++ {
		if filepath.Ext(m[i]) != ".tmp" {
			ulog("removing %s\n", m[i])
			err = os.Remove(m[i])
			if nil != err {
				ulog("error removing file: %s  err = %v\n", m[i], err)
				return err
			}
		}
	}

	// now move our newly uploaded picture into its proper name...
	err = os.Rename(tmpFile, picturefilename)
	if nil != err {
		ulog("os.Rename(%s,%s):  err = %v\n", tmpFile, picturefilename, err)
		return err
	}

	return nil
}

const (
	S3_REGION         = "ap-south-1"                          // This parameter define the region of bucket
	S3_BUCKET         = "upload-images-test"                  // This parameter define the bucket name in S3
	PUBLIC_ACL        = "public-read"                         // This parameter make S3 bucket's object readable
	IMAGE_UPLOAD_PATH = ""                                    // This parameter define in which folder have to upload image
	AWS_PROFILE_NAME  = "akshay"                              // define profile name to get credentials
	S3_HOST           = "https://s3.ap-south-1.amazonaws.com" // define host to access images from the s3
)

func generateFileName(uid int) string {
	// get timestamps
	timestamps := time.Now().UTC()

	// id of user, and timestamps
	s := []string{strconv.Itoa(uid), timestamps.Format("20160102150405")}

	// generate filename to save on s3/db
	filename := strings.Join(s, "_")

	return filename
}

func uploadImageFileToS3(usrfname *multipart.FileHeader, usrfile multipart.File, uid int) (string, string) {

	// generate filename to save on s3/db
	filename := generateFileName(uid)

	// setup credential
	// reading credential from the aws config
	creds := credentials.NewSharedCredentials("", AWS_PROFILE_NAME)
	_, err := creds.Get()
	if err != nil {
		fmt.Printf("Bad credentials: %s", err)
	}

	// Set up configuration
	cfg := aws.NewConfig().WithRegion(S3_REGION).WithCredentials(creds)

	// Set up session
	sess, err := session.NewSession(cfg)
	if err != nil {
		fmt.Println(err)
	}

	// Create new s3 instance
	svc := s3.New(sess)

	// Path after the bucket name
	imagePath := path.Join(IMAGE_UPLOAD_PATH, filename)

	// define parameters to upload image to S3
	params := &s3.PutObjectInput{
		Bucket:               aws.String(S3_BUCKET),
		Key:                  aws.String(imagePath), // it include filename
		Body:                 usrfile,               // data of file
		ServerSideEncryption: aws.String("AES256"),
		ContentType:          aws.String(usrfname.Header["Content-Type"][0]),
		CacheControl:         aws.String("max-age=86400"),
		ACL:                  aws.String(PUBLIC_ACL),
	}

	// Upload image to s3 bucket
	resp, err := svc.PutObject(params)
	if err != nil {
		fmt.Printf("bad response: %s", err)
	}

	fmt.Printf("response %s", awsutil.StringValue(resp))
	fmt.Printf("Image location: %s", path.Join(S3_HOST, S3_BUCKET, imagePath))
	imageLocation := path.Join(S3_HOST, S3_BUCKET, imagePath)

	return imagePath, imageLocation
}

func savePersonDetailsHandler(w http.ResponseWriter, r *http.Request) {
	var ssn *sess.Session
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
	uid, err := strconv.Atoi(uidstr)
	if err != nil {
		fmt.Fprintf(w, "Error converting uid to a number: %v. URI: %s\n", err, r.RequestURI)
		return
	}

	//=================================================================
	//  SECURITY
	//=================================================================
	if !ssn.ElemPermsAny(authz.ELEMPERSON, authz.PERMOWNERMOD) {
		ulog("Permissions refuse savePersonDetails page on userid=%d (%s), role=%s\n", ssn.UID, ssn.Firstname, ssn.PMap.Urole.Name)
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}
	if int64(uid) != ssn.UID {
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
			fmt.Println(imageLocation)
			fmt.Println(imagePath)
			ssn.ImageURL = imageLocation
		} else {
			ulog("err loading picture: %v\n", err)
		}

		//=================================================================
		//  Do the update
		//=================================================================
		_, err = Phonebook.prepstmt.updateMyDetails.Exec(d.PreferredName, d.PrimaryEmail, d.OfficePhone, d.CellPhone,
			d.EmergencyContactName, d.EmergencyContactPhone,
			d.HomeStreetAddress, d.HomeStreetAddress2, d.HomeCity, d.HomeState, d.HomePostalCode, d.HomeCountry, ssn.UID,
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
