#!/bin/bash

###############################################################################
#
#                 I M P O R T A N T     N O T I C E  !!!
#
# This test and the usersim application uses the database named "accordtest"
# not the normal directory database named "accord"
#
#
###############################################################################
TOP="../.."
TREPORT="${TOP}/test/testreport.txt"
TESTNAME="Directory Usersim"
TESTSUMMARY="AIR Directory User Simulation"

DBTOOLSDIR="../../dbtools"
USERSIMDIR="."
USERCOUNT=3
USERSIM="${USERSIMDIR}/usersim"

DBNAME="accordtest"
DBUSER="ec2-user"

PBHOST="localhost"
PBPORT="8250"

PHONEBOOKDIR="../.."
STARTPHONEBOOKCMD="./activate.sh -N ${DBNAME} -T start"
STOPPHONEBOOKCMD="./activate.sh stop"

LOGFILE="usersim.log"
FORCEGOOD=0
SKIPCOMPARE=0
GOLD="./gold"
ERRFILE="errors.txt"
TESTCOUNT=0

##########################################################################
# elapsedtime()
# Shows the number of seconds that was needed to run this script
##########################################################################
elapsedtime() {
	duration=$SECONDS
	msg="ElapsedTime: $(($duration / 60)) min $(($duration % 60)) sec"
	echo "${msg}" >>${LOGFILE}
	echo "${msg}"

}

passmsg() {
	if [ ! -f ${TREPORT} ]; then touch ${TREPORT}; fi
	printf "PASSED  %-20.20s  %-40.40s  %6d  \n" "${TESTDIR}" "${TESTNAME}" ${TESTCOUNT} >> ${TREPORT}
}

failmsg() {
	if [ ! -f ${TREPORT} ]; then touch ${TREPORT}; fi
	printf "FAILED  %-20.20s  %-40.40s  %6d  \n" "${TESTDIR}" "${TESTNAME}" ${TESTCOUNT} >> ${TREPORT}
}

forcemsg() {
	if [ ! -f ${TREPORT} ]; then touch ${TREPORT}; fi
	printf "FORCED  %-20.20s  %-40.40s  %6d  \n" "${TESTDIR}" "${TESTNAME}" ${TESTCOUNT} >> ${TREPORT}
}

tdir() {
	local IFS=/
	local p n m
	p=( ${SCRIPTPATH} )
	n=${#p[@]}
	m=$(( n-1 ))
	TESTDIR=${p[$m]}
}

# goldpath simply creates a gold directory if it does not already exist.
goldpath() {
	if [ ! -d "./gold" ]; then
		mkdir gold
	fi
}

##########################################################################
# logcheck()
#   Compares log to log.gold
#   Date related fields are detected with a regular expression and changed
#   to "current time".  More filters may be needed depending on what goes
#   into the logfile.
#	Parameters:
#		none at this time
##########################################################################
logcheck() {
	echo -n "Test completed: " >> ${LOGFILE}
	date >> ${LOGFILE}
	if [ "${FORCEGOOD}" = "1" ]; then
		goldpath
		cp ${LOGFILE} ${GOLD}/${LOGFILE}.gold
		echo "DONE"
	elif [ "${SKIPCOMPARE}" = "0" ]; then
		echo -n "PHASE x: Log file check...  "
		if [ ! -f ${GOLD}/${LOGFILE}.gold  ]; then
			echo "Missing file -- ${GOLD}/${LOGFILE}.gold"
			echo "An empty one is being created to continue"
			touch ${GOLD}/${LOGFILE}.gold
		fi
		if [  ! -f ${LOGFILE} ]; then
			echo "Missing file -- Required files for this check: ${LOGFILE}"
			failmsg
			exit 1
		fi
		declare -a out_filters=(
			's/^Date\/Time:.*/current time/'
			's/^Test completed:.*/current time/'
			's/(20[1-4][0-9]\/[0-1][0-9]\/[0-3][0-9] [0-2][0-9]:[0-5][0-9]:[0-5][0-9] )(.*)/$2/'
			's/(20[1-4][0-9]\/[0-1][0-9]-[0-3][0-9] [0-2][0-9]:[0-5][0-9]:[0-5][0-9] )(.*)/$2/'
			's/(20[1-4][0-9]-[0-1][0-9]-[0-3][0-9] [0-2][0-9]:[0-5][0-9]:[0-5][0-9] )(.*)/$2/'
		)
		cp ${GOLD}/${LOGFILE}.gold ll.g
		cp ${LOGFILE} llog
		for f in "${out_filters[@]}"
		do
			perl -pe "$f" ll.g > llx1; mv llx1 ll.g
			perl -pe "$f" llog > lly1; mv lly1 llog
		done
		UDIFFS=$(diff llog ll.g | wc -l)
		if [ ${UDIFFS} -eq 0 ]; then
			echo "PASSED"
			passmsg
			rm -f ll.g llog
		else
			echo "FAILED:  differences are as follows:" >> ${ERRFILE}
			diff ll.g llog >> ${ERRFILE}
			echo >> ${ERRFILE}
			echo "If the new output is correct:  mv ${LOGFILE} ${GOLD}/${LOGFILE}.gold" >> ${ERRFILE}
			cat ${ERRFILE}
			failmsg
			if [ "${ASKBEFOREEXIT}" = "1" ]; then
				pause ${LOGFILE}
			else
				exit 1
			fi
		fi
	else
		echo "FINISHED...  but did not check output"
	fi
	elapsedtime
}


initDB() {
	pushd ${DBTOOLSDIR} >/dev/null
	./apply.sh > /dev/null
	popd >/dev/null
	${USERSIM} -f -u 300 -c 70 -C 70
}

stopPhonebook() {
	pushd ${PHONEBOOKDIR} >/dev/null
	results=$(${STOPPHONEBOOKCMD})
	# echo "results = ${results}"
	result=$(echo "${results}" | grep OK | wc -l)
	if [ ${result} -ge 0 ]; then
		echo "phonebook stopped"
		sleep 2
	else 
		echo "phonebook did not stop properly.  result = \"${result}\""
		exit 1
	fi
	popd >/dev/null
}

startPhonebook() {
	pushd ${PHONEBOOKDIR} >/dev/null
	results=$(${STARTPHONEBOOKCMD})
	# echo "results = \"${results}\""
	result=$(${STARTPHONEBOOKCMD} | grep OK | wc -l)
	if [ "${result}" -eq 0 ]; then
		echo "phonebook did not start properly.  result = \"${result}\""
		exit 1
	fi
	sleep 1
	popd >/dev/null
}


echo "Stopping any running instance of phonebook..."
stopPhonebook

goldpath
echo    "Test Name:    ${TESTNAME}" > ${LOGFILE}
echo    "Test Purpose: ${TESTSUMMARY}" >> ${LOGFILE}
echo -n "Date/Time:    " >>${LOGFILE}
date >> ${LOGFILE}
echo >>${LOGFILE}

L=$(ps -ef | grep phonebook | grep -v grep | grep -v "ssh phonebook" | wc -l)
if [ ${L} -gt 0 ]; then
	echo "Could not stop running instance of phonebook..."
	ps -ef | grep phonebook | grep -v grep 
	exit 1
fi

echo "Creating new db: ${DBNAME}"
initDB

echo "Starting phonebook service"
startPhonebook

L=$(ps -ef | grep phonebook | grep -v grep | grep -v "ssh phonebook" | wc -l)
if [ ${L} -ne 1 ]; then
	echo "Could not find one and only one running instance of phonebook..."
	ps -ef | grep phonebook | grep -v grep 
	exit 1
fi

usrsimout=$(${USERSIM} -u ${USERCOUNT})
echo "${usrsimout}" > usersim.out

echo "Shutting down phonebook service..."
stopPhonebook
TESTCOUNT=$(cat usersim.out | grep "Total Tests:" | awk '{print $3;}')

logcheck

echo "results: ${usrsimout}"
rs=$(echo "${usrsimout}" | grep "fail: 0" | wc -l)
if [ ${rs} -eq 1 ]; then
	exit 0
else
	exit 1
fi
