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

#---------------------------------------------------------------
#  Use the testdb for these tests...
#---------------------------------------------------------------
# echo "Create new database..."
# mysql --no-defaults rentroll < restore.sql

source ../share/base.sh

if [ ! -f "config.json" ]; then
	cp ../usersim/config.json .
fi

pushd ${DBTOOLSDIR};./apply.sh;popd
pushd ${USERSIMPATH};${USERSIM} -f -u 200 -c 70 -C 70;popd
cmd="${ADDUSER} -f Billy -l Thorton -u billybob -r Viewer -p Testing123"
echo "cmd = ${cmd}"
${ADDUSER} -f Billy -l Thorton -r Viewer -p Testing123

echo "STARTING PHONEBOOK SERVER"
startPhonebook

# Get version
#--------------------------
curl -s "http://localhost:8250/v1/version" >a0 2>>${LOGFILE}
OL=$(cat a0 | wc -c)
if [ ${OL} -ne 10 ]; then
    echo "EXPECTING 10 line for version, got ${OL}"
    exit 2
fi

#------------------------------------------------------------------------------
#  TESTS a, b, c
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
#	1. The authenticate command should pass (test a)
#   2. The cookie should be stored in ${C} and should be found along with
#          all the other login information. (test b)
#   3. Test c should return "failure" as it will attempt to validate
#          a bogus cookie value.
#   4. Test d should return "success" as it will attempt to validate
#          the cookie value created in test a.
#   5. Ensure that setting FLAGS 1<<1 causes timestamp to update
#------------------------------------------------------------------------------
echo "%7B%22user%22%3A%22bthorton%22%2C%22pass%22%3A%22Testing123%22%2C%22useragent%22%3A%22Mozilla%2F5.0%20(Macintosh%3B%20Intel%20Mac%20OS%20X%2010_12_6)%20AppleWebKit%2F537.36%20(KHTML%2C%20like%20Gecko)%20Chrome%2F64.0.3282.186%20Safari%2F537.36%22%2C%22remoteaddr%22%3A%22172.31.63.140%3A7497%22%7D" > request
doPlainPOST "http://localhost:8250/v1/authenticate" "request" "a"  "WebService--Authenticate"
C=$(curl -s -X POST http://localhost:8250/v1/authenticate -H "Content-Type: application/json" -d @request | python -m json.tool | grep "Token" | awk '{print $2}' | sed 's/[,"]//g')
echo "%7B%22cookieval%22%3A%22${C}%22%2C%22flags%22%3A0%2C%22useragent%22%3A%22curl%22%2C%22ip%22%3A%221.2.3.4%22%7D" > request
doPlainPOST "http://localhost:8250/v1/validatecookie" "request" "b"  "WebService--ValidateCookie-FLAGS=0"
echo "%7B%22cookieval%22%3A%22deadbeefdeadbeef%22%2C%22flags%22%3A1%2C%22useragent%22%3A%22curl%22%2C%22ip%22%3A%221.2.3.4%22%7D" > request
doPlainPOST "http://localhost:8250/v1/validatecookie" "request" "c"  "WebService--ValidateCookie-FLAGS=1-fail"
echo "%7B%22cookieval%22%3A%22${C}%22%2C%22flags%22%3A1%2C%22useragent%22%3A%22curl%22%2C%22ip%22%3A%221.2.3.4%22%7D" > request
doPlainPOST "http://localhost:8250/v1/validatecookie" "request" "d"  "WebService--ValidateCookie-FLAGS=0-succeed"

echo "Shutting down phonebook service..."
stopPhonebook

logcheck
