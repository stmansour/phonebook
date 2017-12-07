#!/bin/bash

#---------------------------------------------------------------
# TOP is the directory where Phonebook begins. It is used
# in base.sh to set other useful directories such as ${BASHDIR}
#---------------------------------------------------------------
TOP=../..

TESTNAME="Web Services"
TESTSUMMARY="Test Web Services"
RRDATERANGE="-j 2016-07-01 -k 2016-08-01"

CREATENEWDB=0

#---------------------------------------------------------------
#  Use the testdb for these tests...
#---------------------------------------------------------------
# echo "Create new database..."
# mysql --no-defaults rentroll < restore.sql

source ../share/base.sh

initDB() {
	pushd ${DBTOOLSDIR} >/dev/null
	./apply.sh > /dev/null
	popd >/dev/null
	${USERSIM} -f -u 300 -c 70 -C 70
}



echo "STARTING PHONEBOOK SERVER"
startPhonebook

# Get version
#--------------------------
doPlainGET "http://localhost:8250/v1/version" "a0" "WebService--Version"

# Test authentication
#--------------------------
echo "%7B%22user%22%3A%22sman%22%2C%22pass%22%3A%22Dai5F0pp%22%7D " > request
doPlainPOST "http://localhost:8250/v1/authenticate" "request" "a"  "WebService--Authenticate"


echo "Shutting down phonebook service..."
stopPhonebook

# logcheck
