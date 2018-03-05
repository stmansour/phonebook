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
doPlainGET "http://localhost:8250/v1/version" "a0" "WebService--Version"

# Test authentication
#--------------------------
echo "%7B%22user%22%3A%22bthorton%22%2C%22pass%22%3A%22Testing123%22%2C%22useragent%22%3A%22Mozilla%2F5.0%20(Macintosh%3B%20Intel%20Mac%20OS%20X%2010_12_6)%20AppleWebKit%2F537.36%20(KHTML%2C%20like%20Gecko)%20Chrome%2F64.0.3282.186%20Safari%2F537.36%22%2C%22remoteaddr%22%3A%22172.31.63.140%3A7497%22%7D" > request
doPlainPOST "http://localhost:8250/v1/authenticate" "request" "a"  "WebService--Authenticate"


echo "Shutting down phonebook service..."
stopPhonebook

logcheck
