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
	C=$(curl -s -X POST http://localhost:8250/v1/authenticate -H "Content-Type: application/json" -d @request | python -m json.tool | grep "Token" | awk '{print $2}' | sed 's/[,"]//g')
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
	encodeRequest '{"cmd":"get","selected":[],"limit":0,"offset":0}'
    dojsonPOST "http://localhost:8250/v1/people/1/9" "request" "${TFILES}${STEP}"  "get-person-9-biz-1"
	encodeRequest '{"cmd":"get","selected":[],"limit":0,"offset":0}'
    dojsonPOST "http://localhost:8250/v1/people/1/999" "request" "${TFILES}${STEP}"  "get-person-999-biz-1"
fi

echo "Shutting down phonebook service..."
stopPhonebook


logcheck
