#!/bin/bash

#---------------------------------------------------------------
# TOP is the directory where Phonebook begins. It is used
# in base.sh to set other useful directories such as ${BASHDIR}
#---------------------------------------------------------------
TOP=../..
BINDIR=${TOP}/tmp/phonebook
TESTNAME="Web Services"
TESTSUMMARY="Test Web Services"
RRDATERANGE="-j 2016-07-01 -k 2016-08-01"
USERSIMPATH="../usersim"
USERSIM=${USERSIMPATH}/usersim
ADDUSER=${BINDIR}/pbadduser
DBTOOLSDIR="../../dbtools"
CREATENEWDB=0

source ../share/base.sh

USER=$(grep Tester1Name config.json | awk '{print $2;}' | sed 's/[",]//g')
PASS=$(grep Tester1Pass config.json | awk '{print $2;}' | sed 's/[",]//g')
if [ "${USER}x" = "x" -o "${PASS}x" = "x" ]; then
    echo "Could not establish user and password. Is config.conf correct?"
    exit 2
fi


#------------------------------------------------------------------------------
#  login - will attempt to login to the wreis server. If it is successful
#          it will set two environment variables:
#
#          TOKEN   - will contain the cookie value for AIR login
#          COOKIES - contains the option for CURL to include the AIR cookie
#                    in requests
#
#          dojsonPOST is setup to use ${COOKIES}
#
#  Scenario:
#  Execute the url to ping the server
#
#  Expected Results:
#   1.  It should return the server version
#------------------------------------------------------------------------------

login() {
    if [ "x${COOKIES}" = "x" ]; then
        encodeRequest "{\"user\":\"${USER}\",\"pass\":\"${PASS}\"}"
        OUTFILE="loginrequest"
        dojsonPOST "http://localhost:8250/v1/authenticate/" "request" "${OUTFILE}"  "login"

        #-----------------------------------------------------------------------------
        # Now we need to add the token to the curl command for future calls to
        # the server.  curl -b "air=${TOKEN}"  ...
        # Set the command line for cookies in ${COOKIES} and dojsonPOST will use them.
        #-----------------------------------------------------------------------------
        TOKEN=$(grep Token "${OUTFILE}" | awk '{print $2;}' | sed 's/[",]//g')

		echo "LOGIN SUCCESSFULL:  token = ${TOKEN}"
        COOKIES="-b air=${TOKEN}"

        #-----------------------------------------------------------------------
        # This is needed so that the tests can be entered at any point.
        # login() uses dojsonPOST which updates STEP.  We only want the
        # test steps in the main routine below to update the test counts.
        # login should be written so that it can be called anywhere, anytime
        # and it will not alter the sequencing of the output files.
        #-----------------------------------------------------------------------
        ((STEP--))
    fi
}


#------------------------------------------------------------------------------
#  Create Directory Database - generate a random database...
#  This is here mainly to document how it is done.  We have the
#  db in a local file and it is used for this test.  But if we need
#  to regenerate it, this shows how it is done.
#------------------------------------------------------------------------------
CreateDirDB() {
	if [ ! -f "config.json" ]; then
		cp ../usersim/config.json .
	fi

	pushd ${DBTOOLSDIR};./apply.sh;popd
	pushd ${USERSIMPATH};${USERSIM} -f -u 20 -c 7 -C 7 -N accord;popd
	cmd="${ADDUSER} -f Billy -l Thorton -u billybob -r Viewer -p Testing123"
	echo "cmd = ${cmd}"
	${ADDUSER} -f Billy -l Thorton -r Viewer -p Testing123
}

#---------------------------------------------------------------
#  Use the testdb for these tests...
#---------------------------------------------------------------
rm -f pbconsole.log   # link to server console messages
echo "Create new database..."
mysql --no-defaults accord < dbws.sql

echo "STARTING PHONEBOOK SERVER"
startPhonebook
ln -s ${BINDIR}/pbconsole.log

# Get version
#--------------------------
curl -s "http://localhost:8250/v1/version" >a0 2>>${LOGFILE}
OL=$(cat a0 | wc -c)
if [ ${OL} -ne 10 ]; then
    echo "EXPECTING 10 line for version, got ${OL}"
    exit 2
fi

#------------------------------------------------------------------------------
#  TEST a
#  Authentication and ValidateCookie - authenticate to create a cookie in
#  the session table. Then test the variants of validatecookie on it.
#
#  Scenario:
#      Variant 1:  FLAGS = 0 - is the normal call, it should return with
#          the uid, the user name, imageurl, and more.
#      Variant 2:  FLAGS = 1 - minimal check for existence. Returns
#          status = "success" if the cookie is found. Returns "failure"
#          if the cookie is not found.
#
#  Expected Results:
#    see comments below
#------------------------------------------------------------------------------
TFILES=a
if [ "${SINGLETEST}${TFILES}" = "${TFILES}" -o "${SINGLETEST}${TFILES}" = "${TFILES}${TFILES}" ]; then

	#------------------------------------
	# login kwalsh - should succeed
	#------------------------------------
	encodeRequest '{"user":"kwalsh","pass":"Testing123","useragent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.186 Safari/537.36","remoteaddr":"172.31.63.140:7497"}' > request
	doPlainPOST "http://localhost:8250/v1/authenticate" "request" "${TFILES}0"  "WebService--Authenticate"

	#-----------------------------------------------
	# validate the cookie that was returned...
	#-----------------------------------------------
	C=$(curl -s -X POST http://localhost:8250/v1/authenticate -H "Content-Type: application/json" -d @request | python3 -m json.tool | grep "Token" | awk '{print $2}' | sed 's/[,"]//g')
	echo "%7B%22cookieval%22%3A%22${C}%22%2C%22flags%22%3A0%2C%22useragent%22%3A%22curl%22%2C%22ip%22%3A%221.2.3.4%22%7D" > request
	doPlainPOST "http://localhost:8250/v1/validatecookie" "request" "${TFILES}1"  "WebService--ValidateCookie-FLAGS=0"

	#--------------------------------------------------------------
	# provide an invalid cookie and make sure things fail...
	#--------------------------------------------------------------
	echo "%7B%22cookieval%22%3A%22deadbeefdeadbeef%22%2C%22flags%22%3A1%2C%22useragent%22%3A%22curl%22%2C%22ip%22%3A%221.2.3.4%22%7D" > request
	doPlainPOST "http://localhost:8250/v1/validatecookie" "request" "${TFILES}2"  "WebService--ValidateCookie-FLAGS=1-fail"

	#--------------------------------------------------------------
	# validate that the good cookie still works...
	#--------------------------------------------------------------
	echo "%7B%22cookieval%22%3A%22${C}%22%2C%22flags%22%3A1%2C%22useragent%22%3A%22curl%22%2C%22ip%22%3A%221.2.3.4%22%7D" > request
	doPlainPOST "http://localhost:8250/v1/validatecookie" "request" "${TFILES}3"  "WebService--ValidateCookie-FLAGS=0-succeed"

	#------------------------------------------------------------------------
	# Attempt to login a user who is listed as inactive in Directory.  Even
	# though the username and password are correct, it should fail with an
	# indication that the account is inactive.
	#------------------------------------------------------------------------
	encodeRequest '{"user":"bthorton","pass":"Testing123","useragent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.186 Safari/537.36","remoteaddr":"172.31.63.140:7497"}' > request
	doPlainPOST "http://localhost:8250/v1/authenticate" "request" "${TFILES}4"  "WebService--AuthenticateAttemptOnInactiveAccount"

fi

#------------------------------------------------------------------------------
#  TEST b
#  Validate peopletd (people type down)
#
#  Scenario:
#      Simulate typedown events from w2ui, send replies
#
#  Expected Results:
#	1. We should see results that include a name string and the UID
#------------------------------------------------------------------------------
TFILES=b
if [ "${SINGLETEST}${TFILES}" = "${TFILES}" -o "${SINGLETEST}${TFILES}" = "${TFILES}${TFILES}" ]; then
	dojsonGET "http://localhost:8250/v1/peopletd?request=%7B%22search%22%3A%22m%22%2C%22max%22%3A100%7D" "${TFILES}0" "WebService--PeopleTypeDown"
fi

#------------------------------------------------------------------------------
#  TEST c
#
#  Validate Business Unit designator typedown
#
#  Scenario:
#
#  Simulate a typedown request
#
#  Expected Results:
#   We should get back the matching designator(s), and other info
#
#------------------------------------------------------------------------------
TFILES="c"
STEP=0
if [ "${SINGLETEST}${TFILES}" = "${TFILES}" -o "${SINGLETEST}${TFILES}" = "${TFILES}${TFILES}" ]; then
	dojsonGET "http://localhost:8250/v1/butd?request=%7B%22search%22%3A%22S%22%2C%22max%22%3A10%7D" "${TFILES}0" "WebService--BUTypeDown"
fi

#------------------------------------------------------------------------------
#  TEST d
#
#  Validate simple BUD search
#
#  Scenario:
#
#  Simulate a typedown request
#
#  Expected Results:
#   We should get back the matching designator(s), and other info
#
#------------------------------------------------------------------------------
TFILES="d"
STEP=0
if [ "${SINGLETEST}${TFILES}" = "${TFILES}" -o "${SINGLETEST}${TFILES}" = "${TFILES}${TFILES}" ]; then
	dojsonGET "curl http://localhost:8250/v1/bud?request%3D%7B%22search%22%3A%22CCC%22%7D" "${TFILES}0" "WebService--BUDSearch"
	dojsonGET "curl http://localhost:8250/v1/bud?request%3D%7B%22search%22%3A%22FOG%22%7D" "${TFILES}1" "WebService--BUDSearch"
fi

#------------------------------------------------------------------------------
#  TEST e
#
#  Validate a request for a person
#
#  Scenario:
#
#  Request a known person from the db.  Request an unknown person to make sure
#  we get an error. Request to an unknown business should fail.
#
#  Expected Results:
#   The known person should come back with all the info that is safe
#   The unknown person should generate an error
#
#------------------------------------------------------------------------------
TFILES="e"
STEP=0
if [ "${SINGLETEST}${TFILES}" = "${TFILES}" -o "${SINGLETEST}${TFILES}" = "${TFILES}${TFILES}" ]; then
	mysql --no-defaults accord < accord.sql
	login
	encodeRequest '{"cmd":"get","selected":[],"limit":0,"offset":0}'
    dojsonPOST "http://localhost:8250/v1/people/1/9" "request" "${TFILES}${STEP}"  "get-person-9-biz-1"
	encodeRequest '{"cmd":"get","selected":[],"limit":0,"offset":0}'
    dojsonPOST "http://localhost:8250/v1/people/1/999" "request" "${TFILES}${STEP}"  "get-person-999-biz-1"
	COOKIES=
fi

#------------------------------------------------------------------------------
#  TEST f
#
#  Validate a request for a list of persons
#
#  Scenario:
#
#  Request a list of known persons from the db.  Request an unknown person to make sure
#  we get an error. Request to an unknown business should fail.
#
#  Expected Results:
#   The known person list should come back with all the info that is safe
#   The unknown person should generate an error
#
#------------------------------------------------------------------------------
TFILES="f"
STEP=0
if [ "${SINGLETEST}${TFILES}" = "${TFILES}" -o "${SINGLETEST}${TFILES}" = "${TFILES}${TFILES}" ]; then
	mysql --no-defaults accord < accord.sql
	login
	encodeRequest '{"cmd":"getlist","UIDs":[3,7,9]}'
    dojsonPOST "http://localhost:8250/v1/people/1" "request" "${TFILES}${STEP}"  "get-personlist-biz-1"
	encodeRequest '{"cmd":"getlist","UIDs":[3,7,9,8888]}'
    dojsonPOST "http://localhost:8250/v1/people/1" "request" "${TFILES}${STEP}"  "get-personlist-error-biz-1"
	COOKIES=
fi

#------------------------------------------------------------------------------
#  TEST g
#
#  Validate saving updating and retrieving licenses
#
#  Scenario:
#
#
#  Expected Results:
#   do a save, validate that the save occurred
#   do an update, validate that the update occurred,
#   do a delete, validate that the delete occurred
#
#------------------------------------------------------------------------------
TFILES="g"
STEP=0
if [ "${SINGLETEST}${TFILES}" = "${TFILES}" -o "${SINGLETEST}${TFILES}" = "${TFILES}${TFILES}" ]; then
	mysql --no-defaults accord < accord.sql
	login
    # 0. Add a license for UID 269
	encodeRequest '{"cmd":"save","record":{"LID":0,"UID":269,"State":"MO","LicenseNo":"02069301","FLAGS":0}}'
    dojsonPOST "http://localhost:8250/v1/license/0/0" "request" "${TFILES}${STEP}"  "save-license"
    # 1. Read all UID 269's licenses
	encodeRequest '{"cmd":"get"}'
    dojsonPOST "http://localhost:8250/v1/licenses/0/269" "request" "${TFILES}${STEP}"  "get-licenses"
    # 2. Read license LID=1
	encodeRequest '{"cmd":"get"}'
    dojsonPOST "http://localhost:8250/v1/license/0/1" "request" "${TFILES}${STEP}"  "get-license"
    # 3. Add an additional license for UID 269
	encodeRequest '{"cmd":"save","record":{"LID":0,"UID":269,"State":"CA","LicenseNo":"01234567","FLAGS":0}}'
    dojsonPOST "http://localhost:8250/v1/license/0/0" "request" "${TFILES}${STEP}"  "save-license"
    # 4. Read all UID 269's licenses, there should be 2 now
	encodeRequest '{"cmd":"get"}'
    dojsonPOST "http://localhost:8250/v1/licenses/0/269" "request" "${TFILES}${STEP}"  "get-licenses"
    # 5. delete license 2
	encodeRequest '{"cmd":"delete"}'
    dojsonPOST "http://localhost:8250/v1/license/0/2" "request" "${TFILES}${STEP}"  "delete-licenses"
    # 6. Read all UID 269's licenses, there should only be 1 now
	encodeRequest '{"cmd":"get"}'
    dojsonPOST "http://localhost:8250/v1/licenses/0/269" "request" "${TFILES}${STEP}"  "get-licenses"

    COOKIES=
fi

echo "Shutting down phonebook service..."
stopPhonebook


logcheck
