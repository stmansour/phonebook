// add a user
//   needs firstname, lastname, username, passwork, role

package main

import (
	"crypto/sha512"
	"database/sql"
	"extres"
	"flag"
	"fmt"
	"os"
	"phonebook/lib"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	_ "github.com/go-sql-driver/mysql"
)

// Role defines a collection of FieldPerms that can be assigned to a person
type Role struct {
	RID  int    // assigned by DB
	Name string // role name
}

// App is the global data structure for this app
var App struct {
	db              *sql.DB
	DBName          string
	DBUser          string
	username        string
	firstname       string
	lastname        string
	passwd          string
	rname           string
	DefaultImgFName string
	RID             int
	DumpRoles       bool
	Roles           []Role // the roles saved in the database
}

func errcheck(err error) {
	if nil != err {
		fmt.Printf("Error = %v\n", err)
		os.Exit(1)
	}
}

func setDefaultImage() error {
	file, err := os.Open(App.DefaultImgFName)
	if err != nil {
		return err
	}
	defer file.Close()
	creds := credentials.NewStaticCredentials(lib.AppConfig.S3BucketKeyID, lib.AppConfig.S3BucketKey, "")
	if _, err = creds.Get(); err != nil {
		return fmt.Errorf("Bad credentials: %s", err)
	}

	cfg := aws.NewConfig().WithRegion(lib.AppConfig.S3Region).WithCredentials(creds)
	sess, err := session.NewSession(cfg)
	if err != nil {
		return fmt.Errorf("Error creating session: %s", err.Error())
	}
	svc := s3.New(sess)
	imagePath := "defaultProfileImage.png"
	params := &s3.PutObjectInput{
		Bucket:               aws.String(lib.AppConfig.S3BucketName),
		Key:                  aws.String(imagePath), // it include filename
		Body:                 file,                  // data of file
		ServerSideEncryption: aws.String("AES256"),
		ContentType:          aws.String("image/png"),
		CacheControl:         aws.String("max-age=86400"),
		ACL:                  aws.String("public-read"),
	}

	// fmt.Printf(`*** PutObject params
	// 	Bucket:               %s
	// 	Key:                  %s
	// 	ServerSideEncryption: %s
	// 	ContentType:          %s
	// 	CacheControl:         %s
	// 	ACL:                  %s\n`, lib.AppConfig.S3BucketName, imagePath, "AES256", "image/png", "max-age=86400", "public-read")

	// Upload image to s3 bucket
	if _, err = svc.PutObject(params); err != nil {
		return fmt.Errorf("Error with PutObject: %s", err.Error())
	}

	return nil
}

func readAccessRoles() {
	rows, err := App.db.Query("select RID,Name from roles")
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		var r Role
		errcheck(rows.Scan(&r.RID, &r.Name))
		App.Roles = append(App.Roles, r)
	}
	errcheck(rows.Err())
}

func readCommandLineArgs() {
	dbuPtr := flag.String("B", "ec2-user", "database user name")
	dbnmPtr := flag.String("N", "accord", "database name (accordtest, accord)")
	uPtr := flag.String("u", "", "username")
	pPtr := flag.String("p", "accord", "password")
	fPtr := flag.String("f", "", "first or given name")
	lPtr := flag.String("l", "", "last or surname name")
	rPtr := flag.String("r", "Viewer", "sets the user's role")
	diPtr := flag.String("defaultimage", "", "filename of image file that is the default image for all users. Must be a .png")
	RPtr := flag.Bool("R", false, "dump roles to stdout")
	flag.Parse()
	App.DefaultImgFName = *diPtr
	App.DBName = *dbnmPtr
	App.DBUser = *dbuPtr
	App.username = *uPtr
	App.firstname = *fPtr
	App.lastname = *lPtr
	App.passwd = *pPtr
	App.rname = *rPtr
	App.DumpRoles = *RPtr
}

func getUsername() {
	//============================================
	// generate a unique username...
	//============================================
	App.username = strings.ToLower(App.firstname[0:1] + App.lastname)
	if len(App.username) > 17 {
		App.username = App.username[0:17]
	}
	UserName := App.username
	var xx int
	nUID := 0
	for {
		found := false
		rows, err := App.db.Query("select uid from people where UserName=?", UserName)
		errcheck(err)
		defer rows.Close()
		for rows.Next() {
			errcheck(rows.Scan(&xx))
			nUID++
			found = true
			UserName = fmt.Sprintf("%s%d", App.username, nUID)
		}
		if !found {
			break
		}
	}
	App.username = UserName
}

func main() {
	readCommandLineArgs()

	var err error
	// s := fmt.Sprintf("%s:@/%s?charset=utf8&parseTime=True", App.DBUser, App.DBName)
	lib.ReadConfig()
	s := extres.GetSQLOpenString(App.DBName, &lib.AppConfig)
	App.db, err = sql.Open("mysql", s)
	if nil != err {
		fmt.Printf("sql.Open: Error = %v\n", err)
	}
	defer App.db.Close()
	err = App.db.Ping()
	if nil != err {
		fmt.Printf("App.db.Ping: Error = %v\n", err)
	}
	readAccessRoles()
	if App.DumpRoles {
		for i := 0; i < len(App.Roles); i++ {
			fmt.Printf("%d - %s\n", App.Roles[i].RID, App.Roles[i].Name)
		}
		os.Exit(0)
	}

	if len(App.DefaultImgFName) > 0 {
		if err = setDefaultImage(); err != nil {
			fmt.Printf("Error reading image file: %s\n", err.Error())
		}
		os.Exit(0)
	}

	App.RID = 0

	for i := 0; i < len(App.Roles); i++ {
		if App.Roles[i].Name == App.rname {
			App.RID = App.Roles[i].RID
		}
	}

	if 0 == App.RID {
		fmt.Printf("Could not find role named: %s\n", App.rname)
		os.Exit(0)
	}

	getUsername()
	sha := sha512.Sum512([]byte(App.passwd))
	passhash := fmt.Sprintf("%x", sha)

	stmt, err := App.db.Prepare("INSERT INTO people (UserName,passhash,FirstName,LastName,RID) VALUES(?,?,?,?,?)")
	if nil != err {
		fmt.Printf("error = %v\n", err)
		os.Exit(1)
	}
	_, err = stmt.Exec(App.username, passhash, App.firstname, App.lastname, App.RID)
	if nil != err {
		fmt.Printf("error = %v\n", err)
	} else {
		fmt.Printf("Added user to database %s:  username: %s, access role: %s\n", App.DBName, App.username, App.rname)
	}
}
