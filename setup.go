package main

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"phonebook/db"
	"strings"
)

func setupHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Printf("entered setup handler\n")
	w.Header().Set("Content-Type", "text/html")
	var sess *db.Session
	var ui uiSupport
	sess = nil
	if 0 < initHandlerSession(sess, &ui, w, r) {
		return
	}
	sess = ui.X
	breadcrumbReset(sess, "Setup", "/setup/")

	// SECURITY
	if !(sess.ElemPermsAny(db.ELEMPBSVC, db.PERMEXEC)) {
		ulog("Permissions refuse setup page on userid=%d (%s), role=%s\n", sess.UID, sess.Firstname, sess.PMap.Urole.Name)
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}

	err := renderTemplate(w, ui, "setup.html")
	if nil != err {
		errmsg := fmt.Sprintf("setupHandler: err = %v\n", err)
		ulog(errmsg)
		fmt.Println(errmsg)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func fileCopy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	cerr := out.Close()
	if err != nil {
		return err
	}
	return cerr
}

func rmFilesWithBaseName(base, except string) {
	//----------------------------------------------------------
	// Remove any files that match the old filename
	// (that is, minus the path and minus the file extension)
	//----------------------------------------------------------
	m, err := filepath.Glob(fmt.Sprintf("./images/%s.*", base))
	if nil != err {
		fmt.Printf("filepath.Glob returned error: %v\n", err)
	}
	fmt.Printf("filepath.Glob returned the following matches: %v\n", m)
	for i := 0; i < len(m); i++ {
		if filepath.Ext(m[i]) != except {
			fmt.Printf("removing %s\n", m[i])
			err = os.Remove(m[i])
			if nil != err {
				fmt.Printf("error removing file: %s  err = %v\n", m[i], err)
			}
		}
	}
}

// uploadSetupImage handles uploads for branding the interface.
//
// Params
//
// Returns:  err from any os file function
func uploadSetupImage(dest string, userfname string, usrfile *multipart.File) (string, error) {
	//-----------------------------------------------
	// use the same filetype for the final filename
	//-----------------------------------------------
	ftype := filepath.Ext(userfname)
	fmt.Printf("user file type: %s\n", ftype)

	//-----------------------------------------------
	// the final image name
	//-----------------------------------------------
	// fmt.Printf("dest = %s\n", dest)
	destfilename := fmt.Sprintf("images/%s%s", dest, ftype)
	// fmt.Printf("destfilename: %s\n", destfilename)

	//-----------------------------------------------
	//  delete old tmp file if exists:
	//-----------------------------------------------
	tmpFile := fmt.Sprintf("images/%s.tmp", dest)
	fmt.Printf("tmpFile to delete if exists: %s\n", tmpFile)
	finfo, err := os.Stat(tmpFile)
	if os.IsNotExist(err) {
		fmt.Printf("%s was not found. Nothing to delete\n", tmpFile)
	} else {
		fmt.Printf("os.Stat(%s) returns:  err=%v,  finfo=%#v\n", tmpFile, err, finfo)
		err = os.Remove(tmpFile)
		fmt.Printf("os.Remove(%s) returns err=%v\n", tmpFile, err)
	}

	//-----------------------------------------------
	// copy the requested new file to "<dest>.tmp"
	//-----------------------------------------------
	err = uploadFileCopy(usrfile, tmpFile)
	if nil != err {
		fmt.Printf("uploadFileCopy returned error: %v\n", err)
		return "", err
	}

	//----------------------------------------------------------
	// Remove any files that match the old filename
	// (that is, minus the path and minus the file extension)
	//----------------------------------------------------------
	rmFilesWithBaseName(dest, ".tmp")

	//-------------------------------------------------------------
	// now move our newly uploaded picture into its final name...
	//-------------------------------------------------------------
	err = os.Rename(tmpFile, destfilename)
	if nil != err {
		fmt.Printf("os.Rename(%s,%s):  err = %v\n", tmpFile, destfilename, err)
		return tmpFile, err
	}

	return destfilename, nil
}

func resetImage(img string, r *http.Request) {
	file, header, err := r.FormFile("imgfile")
	// fmt.Printf("file: %v, header: %v, err: %v\n", file, header, err)
	if nil == err {
		defer file.Close()
		fname, err := uploadSetupImage(img, header.Filename, &file)
		if nil != err {
			ulog("uploadImageFile returned error: %v\n", err)
		}
		if len(fname) > 0 {
			PhonebookUI.Images[img] = fname
		}
	} else if err != nil {
		fmt.Printf("err = %v\n", err)
	}
}

func saveSetupHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Printf("Entered saveSetupHandler\n")
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Cache-Control", "no-store")
	var sess *db.Session
	var ui uiSupport
	sess = nil
	if 0 < initHandlerSession(sess, &ui, w, r) {
		return
	}
	sess = ui.X
	breadcrumbReset(sess, "Setup", "/setup/")

	// SECURITY
	if !(sess.ElemPermsAny(db.ELEMPBSVC, db.PERMEXEC)) {
		ulog("Permissions refuse setup page on userid=%d (%s), role=%s\n", sess.UID, sess.Firstname, sess.PMap.Urole.Name)
		http.Redirect(w, r, "/search/", http.StatusFound)
		return
	}

	action := strings.ToLower(r.FormValue("action"))
	section := r.FormValue("section")
	// fmt.Printf("action = %s,  section = %s\n", action, section)
	if action == "save" {
		resetImage(section, r)
	} else if action == "reset all images" {
		for i := 0; i < len(uiDflt); i++ {
			rmFilesWithBaseName(uiDflt[i], "")
			PhonebookUI.Images[uiDflt[i]] = "./images/" + uiDflt[i] + ".png"
			fileCopy("./images/default/"+uiDflt[i]+".png", PhonebookUI.Images[uiDflt[i]])
		}
	}

	err := renderTemplate(w, ui, "setup.html")
	if nil != err {
		errmsg := fmt.Sprintf("saveSetupHandler: err = %v\n", err)
		ulog(errmsg)
		fmt.Println(errmsg)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
