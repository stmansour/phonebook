#############################################################################
# RentRoll test infrastructure base
#
# Include this script in your test script to get the base testing capabilities.
#
# Your script can override the following values:
#
# 	BUD="REX"    - name of the business to use for your tests
# 	RRDATERANGE  - default RentRoll date range. You can override this
#                  value with a different range.
#############################################################################

TOOLSDIR=${TOP}/tools
BASHDIR=${TOOLSDIR}/bashtools
PHONEBOOKDIR="../../tmp/phonebook"
STARTPHONEBOOKCMD="./phonebook"
STOPPHONEBOOKCMD="./activate.sh stop"
LOGFILE="log"
PYTHON="python3"

#############################################################################
# Set default values
#############################################################################
SECONDS=0
ERRFILE="err.txt"
UNAME=$(uname)
MYSQLOPTS=""
MYSQL=$(which mysql)
TESTCOUNT=0			## this is an internal counter, your external script should not touch it
SHOWCOMMAND=0
SCRIPTPATH=$(pwd -P)
CASPERTEST="casperjs test"

if [ "x${NOCONSOLE}" = "x" ]; then
	NOCONSOLE="-nocon"
else
	NOCONSOLE=""
fi

if [ "x${MANAGESERVER}" = "x" ]; then
	MANAGESERVER=1
fi
if [ "x${RRPORT}" = "x" ]; then
	RRPORT="8270"
fi
if [ "x${RRBIN}" = "x" ]; then
	RRBIN="../../tmp/rentroll"
else
	echo "RBIN was pre-set to:  \"${RRBIN}\""
fi
TREPORT="${TOP}/test/testreport.txt"

RENTROLL="${RRBIN}/rentroll -A ${NOCONSOLE}"
CSVLOAD="${RRBIN}/rrloadcsv  ${NOCONSOLE}"
GOLD="./gold"

PAUSE=0
SKIPCOMPARE=0
FORCEGOOD=0
ASKBEFOREEXIT=0
TESTCOUNT=0

#############################################################################
#  This code ensures that mysql does not touch production databases.
#  The way identity is kept, default usage of mysql or mysqldump often
#  goes straight to the production databases.
#############################################################################
if [ "${UNAME}" == "Darwin" -o "${IAMJENKINS}" == "jenkins" ]; then
	MYSQLOPTS="--no-defaults"
fi

#############################################################################
# pause()
#   Description:
#		Ask the user how to proceed.
#
#   Params:
#       ${1} = name of the file to move to gold/${1}.gold if 'm' is pressed
#############################################################################
pause() {
	echo
	echo
	read -p "Press [Enter] to continue, M to move ${2} to gold/${2}.gold, Q or X to quit..." x
	x=$(echo "${x}" | tr "[:upper:]" "[:lower:]")
	if [ "${x}" == "q" -o ${x} == "x" ]; then
		exit 0
	elif [[ ${x} == "m" ]]; then
		echo "********************************************"
		echo "********************************************"
		echo "********************************************"
		echo "mv ${1} gold/${1}.gold"
		mv ${1} gold/${1}.gold
	fi

}


#############################################################################
# checkPause()
#   Description:
#		if PAUSE = 1 then call pause() after executing a command. This is
#       handy when checking the database or other output as each command
#       executes to see when something happens or when something in the
#       database goes bad.
#
#   Params:
#
#############################################################################
checkPause() {
	if [ "${PAUSE}" = "1" ]; then
		pause
	fi
}

#############################################################################
# incStep()
#   Description:
#		Increment the STEP variable.  It is encapsulated here because
#       there may be additional steps to perform in the future.
#
#   Params:
#       none
#############################################################################
incStep() {
	((STEP++))
}

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

#############################################################################
# docsvtest()
#    The purpose of this routine is to call rrloadcsv with the
#     parameters supplied in $2 and send its output to a file
#     named $1. After trrloadcsv completes, the output in $1 will
#     be compared with the output in gold/$1.gold.  If there are
#     no diffs, then the test passes.  If there are diffs, then
#     it terminates execution of the script after doing
#     the following:
#
#        (a) Displays the diffs
#        (b) Displays the mv command to use if the newly generated
#            output is now correct and the gold/$1.gold file needs
#            to be updated.  You can just copy the command and paste
#            it into your command line.  Very handy
#        (c) Displays the full command it used to generate the output
#            in $1. This is very handy for reproducing a problem.
#
#     Additionally, there are some Environment Variables that cause
#     it to perform several functions that are very handy:
#
#        SKIPCOMPARE - ${SKIPCOMPARE} defaults to 0. As long as its
#            value is 0 the output in $1 will be compared to
#            gold/$1.gold .  However, there may be times where
#            you want the script to run to completion even if the
#            output miscompares with what is in gold/*  By convention,
#            all of my "functest.sh" scripts use the -f option to
#            set this value.
#
#        FORCEGOOD - ${FORCEGOOD} is set to 0 by default. If it is set
#            set to 1 it means that the output generated and stored in
#            $1 during this run is known to be "correct", even though
#            it may be different than what is in gold/$1.gold.  It will
#            automatically copy $1 to gold/$1.go. This is
#            extremely handy if a change was made to the table output
#            generator, or if any new fields were added to the database
#            and you've validated in some other way that everything is
#            working after such a change.  By convention, all of my
#            "function.sh" scripts use the -o option to set FORCEGOOD
#            to 1.
#
#	Parameters:
# 		$1 = base file name
#		$2 = app options to reproduce
# 		$3 = title for reporting.  No spaces
#############################################################################
docsvtest () {
	TESTCOUNT=$((TESTCOUNT + 1))
	printf "PHASE %2s  %3s  %s... " ${TESTCOUNT} $1 $3

	if [ "x${2}" != "x" ]; then
		${CSVLOAD} ${2} >${1} 2>&1
	fi

	checkPause

	if [ "${FORCEGOOD}" = "1" ]; then
		goldpath
		cp ${1} ${GOLD}/${1}.gold
		echo "DONE"
	elif [ "${SKIPCOMPARE}" = "0" ]; then
		if [ ! -f ${GOLD}/${1}.gold ]; then
			echo "UNSET CONTENT" > ${GOLD}/${1}.gold
			echo "Created a default ${GOLD}/$1.gold for you. Update this file with known-good output."
		fi
		UDIFFS=$(diff ${1} ${GOLD}/${1}.gold | wc -l)
		if [ ${UDIFFS} -eq 0 ]; then
			if [ ${SHOWCOMMAND} -eq 1 ]; then
				echo "PASSED	cmd: ${CSVLOAD} ${2}"
			else
				echo "PASSED"
			fi
		else
			echo "FAILED..." >> ${ERRFILE}
			echo "Differences in ${1} are as follows:" >> ${ERRFILE}
			diff ${GOLD}/${1}.gold ${1} >> ${ERRFILE}
			echo "If correct:  mv ${1} ${GOLD}/${1}.gold" >> ${ERRFILE}
			echo "Command to reproduce:  ${CSVLOAD} ${2}" >> ${ERRFILE}
			cat ${ERRFILE}
			failmsg
			if [ "${ASKBEFOREEXIT}" = "1" ]; then
				pause ${1}
			else
				exit 1
			fi
		fi
	else
		echo
	fi
}


########################################
# dorrtest()
#	Parameters:
# 		$1 = base file name
#		$2 = app options to reproduce
# 		$3 = title
########################################
dorrtest () {
	TESTCOUNT=$((TESTCOUNT + 1))
	printf "PHASE %2s  %3s  %s... " ${TESTCOUNT} $1 $3
	${RENTROLL} ${2} >${1} 2>&1

	checkPause

	if [ "${FORCEGOOD}" = "1" ]; then
		goldpath
		cp ${1} ${GOLD}/${1}.gold
		echo "DONE"
	elif [ "${SKIPCOMPARE}" = "0" ]; then
		if [ ! -f ${GOLD}/${1}.gold ]; then
			echo "UNSET CONTENT" > ${GOLD}/${1}.gold
			echo "Created a default ${GOLD}/$1.gold for you. Update this file with known-good output."
		fi
		UDIFFS=$(diff ${1} ${GOLD}/${1}.gold | wc -l)
		if [ ${UDIFFS} -eq 0 ]; then
			if [ ${SHOWCOMMAND} -eq 1 ]; then
				echo "PASSED	cmd: ${RENTROLL} ${2}"
			else
				echo "PASSED"
			fi
		else
			echo "FAILED..." >> ${ERRFILE}
			echo "Differences in ${1} are as follows:" >> ${ERRFILE}
			diff ${GOLD}/${1}.gold ${1} >> ${ERRFILE}
			echo "If correct:  mv ${1} ${GOLD}/${1}.gold" >> ${ERRFILE}
			echo "Command to reproduce:  ${RENTROLL} ${2}" >> ${ERRFILE}
			cat ${ERRFILE}
			failmsg
			if [ "${ASKBEFOREEXIT}" = "1" ]; then
				pause ${1}
			else
				exit 1
			fi
		fi
	else
		echo
	fi
}

#############################################################################
# doCompareIgnoreDates()
#	just like docsvtest only for Onesite
#
#	Parameters:
# 		$1 = base file name
#		$2 = app options to reproduce
# 		$3 = title for reporting.  No spaces
#############################################################################
doCompareIgnoreDates() {
	declare -a out_filters=(
		's/^Time:.*/current time/'
		's/^Date:.*/current time/'
		's/^Import File:.*/import file/'
		's/\s+[0-1]?[0-9]\/[0-3]?[0-9]\/[0-9][0-9][^-]*/date/g'
		's/\s+[0-1]?[0-9]\/[0-3]?[0-9]\/20[0-9][0-9][^-]*/date/g'
	)
	cp ${GOLD}/${1}.gold ${GOLD}/${1}.g
	cp ${1} ${1}.g
	for f in "${out_filters[@]}"
	do
		perl -pe "$f" ${GOLD}/${1}.g > ${GOLD}/${1}.t; mv ${GOLD}/${1}.t ${GOLD}/${1}.g
		perl -pe "$f" ${1}.g > ${1}.t; mv ${1}.t ${1}.g
	done
	UDIFFS=$(diff ${1}.g ${GOLD}/${1}.g | wc -l)
	if [ ${UDIFFS} -eq 0 ]; then
		if [ ${SHOWCOMMAND} -eq 1 ]; then
			echo "PASSED	cmd: ${CSVLOAD} ${2}"
		else
			echo "PASSED"
		fi
		rm -f ${1}.g ${GOLD}/${1}.g
	else
		echo "FAILED...   if correct:  mv ${1} ${GOLD}/${1}.gold" >> ${ERRFILE}
		echo "Command to reproduce:  ${CSVLOAD} ${2}" >> ${ERRFILE}
		echo "Differences in ${1} are as follows:" >> ${ERRFILE}
		diff ${GOLD}/${1}.g ${1}.g >> ${ERRFILE}
		cat ${ERRFILE}
		failmsg
		if [ "${ASKBEFOREEXIT}" = "1" ]; then
			pause ${1}
		else
			exit 1
		fi
	fi
}

#############################################################################
# doOnesiteTest()
#	just like docsvtest only for Onesite
#
#	Parameters:
# 		$1 = base file name
#		$2 = app options to reproduce
# 		$3 = title for reporting.  No spaces
#############################################################################
doOnesiteTest () {
	TESTCOUNT=$((TESTCOUNT + 1))
	printf "PHASE %2s  %3s  %s... " ${TESTCOUNT} $1 $3

	if [ "x${2}" != "x" ]; then
		${RRBIN}/importers/onesite/onesiteload ${2} >${1} 2>&1
	fi

	checkPause

	if [ "${FORCEGOOD}" = "1" ]; then
		goldpath
		cp ${1} ${GOLD}/${1}.gold
		echo "DONE"
	elif [ "${SKIPCOMPARE}" = "0" ]; then
		if [ ! -f ${GOLD}/${1}.gold ]; then
			echo "UNSET CONTENT" > ${GOLD}/${1}.gold
			echo "Created a default ${GOLD}/$1.gold for you. Update this file with known-good output."
		fi
		doCompareIgnoreDates "${1}" "${2}" "${3}"
	else
		echo
	fi
}

#############################################################################
# docsvIgnoreDatesTest()
#	just like docsvtest only ignore dates in known good files
#
#	Parameters:
# 		$1 = base file name
#		$2 = app options to reproduce
# 		$3 = title for reporting.  No spaces
#############################################################################
docsvIgnoreDatesTest () {
	TESTCOUNT=$((TESTCOUNT + 1))
	printf "PHASE %2s  %3s  %s... " ${TESTCOUNT} $1 $3

	if [ "x${2}" != "x" ]; then
		${CSVLOAD} $2 >${1} 2>&1
	fi

	checkPause

	if [ "${FORCEGOOD}" = "1" ]; then
		goldpath
		cp ${1} ${GOLD}/${1}.gold
		echo "DONE"
	elif [ "${SKIPCOMPARE}" = "0" ]; then
		if [ ! -f ${GOLD}/${1}.gold ]; then
			echo "UNSET CONTENT" > ${GOLD}/${1}.gold
			echo "Created a default ${GOLD}/$1.gold for you. Update this file with known-good output."
		fi
		doCompareIgnoreDates "${1}" "${2}" "${3}"
	else
		echo
	fi
}


#############################################################################
# doRoomKeyTest()
#	just like docsvtest only for RoomKey
#
#	Parameters:
# 		$1 = base file name
#		$2 = app options to reproduce
# 		$3 = title for reporting.  No spaces
#############################################################################
doRoomKeyTest () {
	TESTCOUNT=$((TESTCOUNT + 1))
	printf "PHASE %2s  %3s  %s... " ${TESTCOUNT} $1 $3

	if [ "x${2}" != "x" ]; then
		${RRBIN}/importers/roomkey/roomkeyload ${2} >${1} 2>&1
	fi

	checkPause

	if [ "${FORCEGOOD}" = "1" ]; then
		goldpath
		cp ${1} ${GOLD}/${1}.gold
		echo "DONE"
	elif [ "${SKIPCOMPARE}" = "0" ]; then
		if [ ! -f ${GOLD}/${1}.gold ]; then
			echo "UNSET CONTENT" > ${GOLD}/${1}.gold
			echo "Created a default ${GOLD}/$1.gold for you. Update this file with known-good output."
		fi
		declare -a out_filters=(
			's/^Time:.*/current time/'
			's/^Date:.*/current time/'
			's/^Import File:.*/import file/'
			's/^Guest Export File:.*/guest file/'
			's/\s+[0-1]?[0-9]\/[0-3]?[0-9]\/[0-9][0-9][^-]*/date/g'
			's/\s+[0-1]?[0-9]\/[0-3]?[0-9]\/20[0-9][0-9][^-]*/date/g'
		)
		cp ${GOLD}/${1}.gold ${GOLD}/${1}.g
		cp ${1} ${1}.g
		for f in "${out_filters[@]}"
		do
			perl -pe "$f" ${GOLD}/${1}.g > ${GOLD}/${1}.t; mv ${GOLD}/${1}.t ${GOLD}/${1}.g
			perl -pe "$f" ${1}.g > ${1}.t; mv ${1}.t ${1}.g
		done
		UDIFFS=$(diff ${1}.g ${GOLD}/${1}.g | wc -l)
		# UDIFFS=$(diff ${1} ${GOLD}/${1}.gold | wc -l)
		if [ ${UDIFFS} -eq 0 ]; then
			if [ ${SHOWCOMMAND} -eq 1 ]; then
				echo "PASSED	cmd: ${RRBIN}/importers/roomkey/roomkeyload ${2}"
			else
				echo "PASSED"
			fi
			rm -f ${1}.g ${GOLD}/${1}.g
		else
			echo "FAILED...   if correct:  mv ${1} ${GOLD}/${1}.gold" >> ${ERRFILE}
			echo "Command to reproduce:  ${RRBIN}/importers/roomkey/roomkeyload ${2}" >> ${ERRFILE}
			echo "Differences in ${1} are as follows:" >> ${ERRFILE}
			# diff ${GOLD}/${1}.gold ${1} >> ${ERRFILE}
			diff ${GOLD}/${1}.g ${1}.g >> ${ERRFILE}
			cat ${ERRFILE}
			failmsg
			if [ "${ASKBEFOREEXIT}" = "1" ]; then
				pause ${1}
			else
				exit 1
			fi
		fi
	else
		echo
	fi
}

########################################
# mysqlverify()
#	Parameters:
# 		$1 = base file name
#		$2 = app options to reproduce
# 		$3 = title
#       $4 = mysql validation query
########################################
mysqlverify () {
# Generate the mysql commands needed to validate...
cat >xxqq <<EOF
use rentroll;
${4}
EOF
	TESTCOUNT=$((TESTCOUNT + 1))
	printf "PHASE %2s  %3s  %s... " ${TESTCOUNT} $1 $3
	#${CSVLOAD} $2 >>${LOGFILE} 2>&1
	CMD1="mysql --no-defaults <xxqq >${1}"
	mysql --no-defaults <xxqq >${1}

	checkPause

	if [ "${FORCEGOOD}" = "1" ]; then
		goldpath
		cp ${1} ${GOLD}/${1}.gold
		echo "DONE"
	elif [ "${SKIPCOMPARE}" = "0" ]; then
		if [ ! -f ${GOLD}/${1}.gold ]; then
			echo "UNSET CONTENT" > ${GOLD}/${1}.gold
			echo "Created a default $1.gold for you. Update this file with known-good output."
		fi

		#--------------------------------------------------------------------
		# The actual data has timestamp information that changes every run.
		# The timestamp can be filtered out for purposes of testing whether
		# or not the web service could be called and can return the expected
		# data.
		#--------------------------------------------------------------------
		declare -a out_filters=(
			's/\s+(20[1-4][0-9]\/[0-1][0-9]\/[0-3][0-9] [0-2][0-9]:[0-5][0-9]:[0-5][0-9]*)(.*)/currentTime/'
			's/\s+(20[1-4][0-9]-[0-1][0-9]-[0-3][0-9] [0-2][0-9]:[0-5][0-9]:[0-5][0-9]*)(.*)/currentTime/'
		)
		cp gold/${1}.gold qqx
		cp ${1} qqy
		for f in "${out_filters[@]}"
		do
			perl -pe "$f" qqx > qqx1; mv qqx1 qqx
			perl -pe "$f" qqy > qqy1; mv qqy1 qqy
		done

		UDIFFS=$(diff qqx qqy | wc -l)
		if [ ${UDIFFS} -eq 0 ]; then
			if [ ${SHOWCOMMAND} -eq 1 ]; then
				echo "PASSED	cmd: mysql --no-defaults <xxqq >${1}"
			else
				echo "PASSED"
			fi
		else
			echo "FAILED...   if correct:  mv ${1} ${GOLD}/${1}.gold" >> ${ERRFILE}
			echo "Command to reproduce:  ${CMD1}" >> ${ERRFILE}
			echo "Differences in ${1} are as follows:" >> ${ERRFILE}
			diff ${GOLD}/${1}.gold ${1} >> ${ERRFILE}
			cat ${ERRFILE}
			if [ "${ASKBEFOREEXIT}" = "1" ]; then
				pause ${1}
			else
				exit 1
			fi
			exit 1
		fi
	else
		echo
	fi
	rm -f qqx qqy
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
		if [ ! -f ${GOLD}/${LOGFILE}.gold -o ! -f ${LOGFILE} ]; then
			echo "Missing file -- Required files for this check: log.gold and log"
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
		cp log llog
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

##########################################################################
# genericlogcheck()
#   Compares the supplied file $1 to gold/$1.gold
#	Parameters:
# 		$1 = base file name
#		$2 = app options to reproduce
# 		$3 = title
##########################################################################
genericlogcheck() {
	TESTCOUNT=$((TESTCOUNT + 1))
	printf "PHASE %2s  %3s  %s... " ${TESTCOUNT} $1 $3

	checkPause

	if [ "${FORCEGOOD}" = "1" ]; then
		goldpath
		cp ${1} ${GOLD}/${1}.gold
		echo "DONE"
	elif [ "${SKIPCOMPARE}" = "0" ]; then
		if [ ! -f ${GOLD}/${1}.gold -o ! -f ${1} ]; then
			echo "Missing file -- Required files for this check: ${1} and ${GOLD}/${1}.gold"
			failmsg
			exit 1
		fi
		UDIFFS=$(diff ${1} gold/${1}.gold | wc -l)
		if [ ${UDIFFS} -eq 0 ]; then
			echo "PASSED"
			passmsg
		else
			echo "FAILED:  differences are as follows:" >> ${ERRFILE}
			diff gold/${1}.gold ${1} >> ${ERRFILE}
			echo >> ${ERRFILE}
			echo "If the new output is correct:  mv ${1} ${GOLD}/${1}.gold" >> ${ERRFILE}
			cat ${ERRFILE}
			failmsg
			if [ "${ASKBEFOREEXIT}" = "1" ]; then
				pause ${1}
			else
				exit 1
			fi
		fi
	else
		echo
	fi
}

#########################################################
# startPhonebookServer()
#	Kills any currently running instances of the server
#   then starts it up again.  The port is set to the
#   default port of 8270.  If you set RRPORT prior
#   to including base.sh to override the port number
#########################################################
startPhonebook() {
	pushd ${PHONEBOOKDIR} >/dev/null
	stopwatchdog
	stopPhonebook
	echo "STARTPHONEBOOKCMD = ${STARTPHONEBOOKCMD}"
	${STARTPHONEBOOKCMD} >pbconsole.log 2>&1 &
	echo "Started phonebook, waiting 2 seconds"
	sleep 2
	echo "continuing"

	ok=$(curl -s http://locahost:8250/v1/version)
	echo "ok = ${ok}"
	popd >/dev/null
}

stopwatchdog() {
	ps -ef | grep activate.sh | awk '{print $2}' | while read -r line ; do
		echo "kill activate.sh pid: ${line}"
		kill -9 ${line}
	done

	ps -ef | grep pbwatchdog | awk '{print $2}' | while read -r line ; do
		echo "kill pbwatchdog pid: ${line}"
		kill -9 ${line}
	done

	killall phonebook
}


#########################################################
# stopPhonebookServer()
#	Kills any currently running instances of the server
#########################################################
stopPhonebook() {
	pushd ${PHONEBOOKDIR} >/dev/null
	echo "stopPhonebook:  calling activate.sh stop"
	results=$(${STOPPHONEBOOKCMD})
	# echo "results = ${results}"
	result=$(echo "${results}" | grep OK | wc -l)
	if [ ${result} -ge 0 ]; then
		echo "phonebook stopped"
		sleep 2
	else
		echo "stopPhonebook:  'activate.sh stop' failed.  result = \"${result}\""
		echo "stopPhonebook:  attempting to shut down forcefully..."
		stopwatchdog
		results=$(${STOPPHONEBOOKCMD})
		result=$(echo "${results}" | grep OK | wc -l)
		if [ ${result} -ge 0 ]; then
	 		exit 1
	 	fi
	 	echo "stopPhonebook:  success"
	fi
	popd >/dev/null
}


########################################
# doPlainPOST()
#   Simulate a POST command to the server and use
#   the supplied file name as the json data
#	Parameters:
# 		$1 = url
#       $2 = json file
# 		$3 = base file name
#		$4 = title
########################################
doPlainPOST () {
	TESTCOUNT=$((TESTCOUNT + 1))
	printf "PHASE %2s  %3s  %s... " ${TESTCOUNT} $3 $4
	CMD="curl -s -X POST ${1} -H \"Content-Type: application/json\" -H \"X-Forwarded-For: 52.90.192.172\" -d @${2}"

	${CMD} >${3} 2>>${LOGFILE}
	echo >>${3}  # this ensures a newline in the file

	checkPause

	if [ "${FORCEGOOD}" = "1" ]; then
		goldpath
		cp ${3} ${GOLD}/${3}.gold
		echo "DONE"
	elif [ "${SKIPCOMPARE}" = "0" ]; then
		if [ ! -f ${GOLD}/${3}.gold ]; then
			echo "UNSET CONTENT" > ${GOLD}/${3}.gold
			echo "Created a default ${GOLD}/$1.gold for you. Update this file with known-good output."
		fi

		#--------------------------------------------------------------------
		# The actual data has timestamp information that changes every run.
		# The timestamp can be filtered out for purposes of testing whether
		# or not the web service could be called and can return the expected
		# data.
		#--------------------------------------------------------------------
		declare -a out_filters=(
			's/(^[ \t]+"LastModTime":).*/$1 TIMESTAMP/'
			's/(^[ \t]+"CreateTS":).*/$1 TIMESTAMP/'
			's/("Expire":)"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}(.*)/$1 TIMESTAMP $2/'
		)
		cp gold/${3}.gold qqx
		cp ${3} qqy
		for f in "${out_filters[@]}"
		do
			perl -pe "$f" qqx > qqx1; mv qqx1 qqx
			perl -pe "$f" qqy > qqy1; mv qqy1 qqy
		done

		UDIFFS=$(diff qqx qqy | wc -l)
		if [ ${UDIFFS} -eq 0 ]; then
			if [ ${SHOWCOMMAND} -eq 1 ]; then
				echo "PASSED	cmd: ${CMD}"
			else
				echo "PASSED"
			fi
		else
			echo "FAILED..." >> ${ERRFILE}
			echo "Differences in ${3} are as follows:" >> ${ERRFILE}
			diff qqx qqy >> ${ERRFILE}
			echo "If correct:  mv ${3} ${GOLD}/${3}.gold" >> ${ERRFILE}
			echo "Command to reproduce:  ${CMD}" >> ${ERRFILE}
			cat ${ERRFILE}
			failmsg
			if [ "${ASKBEFOREEXIT}" = "1" ]; then
				pause ${3}
			else
				if [ ${MANAGESERVER} -eq 1 ]; then
					echo "STOPPING RENTROLL SERVER"
					pkill rentroll
				fi
				exit 1
			fi
		fi
	else
		echo
	fi
	rm -f qqx qqy
}

########################################################################
# dojsonPOST()
#   Simulate a POST command to the server and use
#   the supplied file name as the json data.
#
#	Parameters:
# 		$1 = url
#       $2 = json file
# 		$3 = base file name
#		$4 = title
#       $5 = if present, the name of a property to ignore in the check
########################################################################
dojsonPOST () {
	((TESTCOUNT++))
	printf "PHASE %2s  %3s  %s... " ${TESTCOUNT} $3 $4
	COOK=""
	if [ "${COOKIES}x" != "x" ]; then
		COOK="${COOKIES}"
	fi
	CMD="curl ${COOK} -s -X POST ${1} -H \"Content-Type: application/json\" -d @${2}"
	${CMD} | tee serverreply | ${PYTHON} -m json.tool >${3} 2>>${LOGFILE}

	incStep
	checkPause

	if [ "${FORCEGOOD}" = "1" ]; then
		goldpath
		cp ${3} ${GOLD}/${3}.gold
		echo "DONE"
	elif [ "${SKIPCOMPARE}" = "0" ]; then
		if [ ! -f ${GOLD}/${3}.gold ]; then
			echo "UNSET CONTENT" > ${GOLD}/${3}.gold
			mkdir -p gold
			echo "Created a default ${GOLD}/$1.gold for you. Update this file with known-good output."
		fi

		#--------------------------------------------------------------------
		# The actual data has timestamp and token information that changes every run.
		# The timestamp can be filtered out for purposes of testing whether
		# or not the web service could be called and can return the expected
		# data.
		#--------------------------------------------------------------------
		declare -a out_filters=(
			's/(^[ \t]+"LastModTime":).*/$1 TIMESTAMP/'
			's/(^[ \t]+"CreateTime":).*/$1 TIMESTAMP/'
			's/(^[ \t]+"ApproverDt":).*/$1 TIMESTAMP/'
			's/(^[ \t]+"OwnerDt":).*/$1 TIMESTAMP/'
			's/(^[ \t]+"Token":).*/$1 TOKEN/'
			's/(^[ \t]+"Expire":).*/$1 TIMESTAMP/'
		)
		cp gold/${3}.gold qqx
		cp ${3} qqy
		for f in "${out_filters[@]}"
		do
			perl -pe "$f" qqx > qqx1; mv qqx1 qqx
			perl -pe "$f" qqy > qqy1; mv qqy1 qqy
		done

		if [ "x${5}" != "x" ]; then
			# echo "  [ CHECKING FOR ${5} ]"
			f="s/([ \t]+\"${5}\":).*/\$1 CONFCODE/"
			perl -pe "$f" qqx > qqx1; mv qqx1 qqx
			perl -pe "$f" qqy > qqy1; mv qqy1 qqy
		fi

		UDIFFS=$(diff qqx qqy | wc -l)
		if [ ${UDIFFS} -eq 0 ]; then
			if [ ${SHOWCOMMAND} -eq 1 ]; then
				echo "PASSED	cmd: ${CMD}"
			else
				echo "PASSED"
			fi
		else
			echo "FAILED..." >> ${ERRFILE}
			echo "Differences in ${3} are as follows:" >> ${ERRFILE}
			diff qqx qqy >> ${ERRFILE}
			echo "If correct:  mv ${3} ${GOLD}/${3}.gold" >> ${ERRFILE}
			echo "Command to reproduce:  ${CMD}" >> ${ERRFILE}
			cat ${ERRFILE}
			failmsg
			if [ "${ASKBEFOREEXIT}" = "1" ]; then
				pause ${3}
			else
				if [ ${MANAGESERVER} -eq 1 ]; then
					echo "STOPPING RENTROLL SERVER"
					pkill rentroll
				fi
				exit 1
			fi
		fi
	else
		echo
	fi
	rm -f qqx qqy
}

########################################################################
# dobinPOST()
#   Simulate a POST command to the server and use
#   the supplied file name as the binary data.
#
#   NOTE:   if you want to pass more args to CURL on this, declare
#           an array called CURLARGS .  For example:
#
#           declare -a CURLARGS=('-H' 'Content-Type: multipart/form-data; boundary=------WebKitFormBoundaryNcjHBUJY2uBg9xCJ'
# 			                     '-H' 'Content-Length: 744')
#
#	Parameters:
# 		$1 = url
#       $2 = binary file
# 		$3 = base file name
#		$4 = title
#       $5 = if present, the name of a property to ignore in the check
########################################################################
dobinPOST () {
	((TESTCOUNT++))
	printf "PHASE %2s  %3s  %s... " ${TESTCOUNT} $3 $4
	COOK=""
	if [ "${COOKIES}x" != "x" ]; then
		COOK="${COOKIES}"
	fi
	# CMD="curl ${COOK} ${CURLARGS[@]} -s -X POST ${1}  --data-binary @${2}"
	CMD="curl -v ${COOK} ${CURLARGS[@]} -s -X POST ${1}  --data-binary @${2}"
	if [ ${SHOWCOMMAND} -eq 1 ]; then
		echo ""
		echo "CMD =  ${CMD}"
	fi
	${CMD} | tee serverreply | ${PYTHON} -m json.tool >${3} 2>>${LOGFILE}

	incStep
	checkPause

	if [ "${FORCEGOOD}" = "1" ]; then
		goldpath
		cp ${3} ${GOLD}/${3}.gold
		echo "DONE"
	elif [ "${SKIPCOMPARE}" = "0" ]; then
		if [ ! -f ${GOLD}/${3}.gold ]; then
			echo "UNSET CONTENT" > ${GOLD}/${3}.gold
			mkdir -p gold
			echo "Created a default ${GOLD}/$1.gold for you. Update this file with known-good output."
		fi

		#--------------------------------------------------------------------
		# The actual data has timestamp and token information that changes every run.
		# The timestamp can be filtered out for purposes of testing whether
		# or not the web service could be called and can return the expected
		# data.
		#--------------------------------------------------------------------
		declare -a out_filters=(
			's/(^[ \t]+"LastModTime":).*/$1 TIMESTAMP/'
			's/(^[ \t]+"CreateTime":).*/$1 TIMESTAMP/'
			's/(^[ \t]+"Token":).*/$1 TOKEN/'
			's/(^[ \t]+"Expire":).*/$1 TIMESTAMP/'
		)
		cp gold/${3}.gold qqx
		cp ${3} qqy
		for f in "${out_filters[@]}"
		do
			perl -pe "$f" qqx > qqx1; mv qqx1 qqx
			perl -pe "$f" qqy > qqy1; mv qqy1 qqy
		done

		if [ "x${5}" != "x" ]; then
			# echo "  [ CHECKING FOR ${5} ]"
			f="s/([ \t]+\"${5}\":).*/\$1 CONFCODE/"
			perl -pe "$f" qqx > qqx1; mv qqx1 qqx
			perl -pe "$f" qqy > qqy1; mv qqy1 qqy
		fi

		UDIFFS=$(diff qqx qqy | wc -l)
		if [ ${UDIFFS} -eq 0 ]; then
			if [ ${SHOWCOMMAND} -eq 1 ]; then
				echo "PASSED	cmd: ${CMD}"
			else
				echo "PASSED"
			fi
		else
			echo "FAILED..." >> ${ERRFILE}
			echo "Differences in ${3} are as follows:" >> ${ERRFILE}
			diff qqx qqy >> ${ERRFILE}
			echo "If correct:  mv ${3} ${GOLD}/${3}.gold" >> ${ERRFILE}
			echo "Command to reproduce:  ${CMD}" >> ${ERRFILE}
			cat ${ERRFILE}
			failmsg
			if [ "${ASKBEFOREEXIT}" = "1" ]; then
				pause ${3}
			else
				if [ ${MANAGESERVER} -eq 1 ]; then
					echo "STOPPING RENTROLL SERVER"
					pkill rentroll
				fi
				exit 1
			fi
		fi
	else
		echo
	fi
	rm -f qqx qqy
}


########################################
# dojsonGET()
#   Simulate a GET command to the server and use
#   the supplied file name as the json data
#	Parameters:
# 		$1 = url
# 		$2 = base file name
#		$3 = title
########################################
dojsonGET () {
	TESTCOUNT=$((TESTCOUNT + 1))
	printf "PHASE %2s  %3s  %s... " ${TESTCOUNT} ${2} ${3}
	CMD="curl -s ${1}"
	if [ ${SHOWCOMMAND} = "1" ]; then
		@echo "dojsonGET: Executing command:  ${CMD}"
	fi
	${CMD} | ${PYTHON} -m json.tool >${2} 2>>${LOGFILE}

	checkPause

	if [ "${FORCEGOOD}" = "1" ]; then
		goldpath
		cp ${2} ${GOLD}/${2}.gold
		echo "DONE"
	elif [ "${SKIPCOMPARE}" = "0" ]; then
		if [ ! -f ${GOLD}/${2}.gold ]; then
			echo "UNSET CONTENT" > ${GOLD}/${2}.gold
			echo "Created a default ${GOLD}/$1.gold for you. Update this file with known-good output."
		fi

		#--------------------------------------------------------------------
		# The actual data has timestamp information that changes every run.
		# The timestamp can be filtered out for purposes of testing whether
		# or not the web service could be called and can return the expected
		# data.
		#--------------------------------------------------------------------
		declare -a out_filters=(
			's/(^[ \t]+"LastModTime":).*/$1 TIMESTAMP/'
		)
		cp gold/${2}.gold qqx
		cp ${2} qqy
		for f in "${out_filters[@]}"
		do
			perl -pe "$f" qqx > qqx1; mv qqx1 qqx
			perl -pe "$f" qqy > qqy1; mv qqy1 qqy
		done

		UDIFFS=$(diff qqx qqy | wc -l)
		if [ ${UDIFFS} -eq 0 ]; then
			if [ ${SHOWCOMMAND} -eq 1 ]; then
				echo "PASSED	cmd: ${CMD}"
			else
				echo "PASSED"
			fi
		else
			echo "FAILED..." >> ${ERRFILE}
			echo "Differences in ${2} are as follows:" >> ${ERRFILE}
			diff qqx qqy >> ${ERRFILE}
			echo "If correct:  mv ${2} ${GOLD}/${2}.gold" >> ${ERRFILE}
			echo "Command to reproduce:  ${CMD}" >> ${ERRFILE}
			cat ${ERRFILE}
			failmsg
			if [ "${ASKBEFOREEXIT}" = "1" ]; then
				pause ${2}
			else
				if [ ${MANAGESERVER} -eq 1 ]; then
					echo "STOPPING RENTROLL SERVER"
					pkill rentroll
				fi
				exit 1
			fi
		fi
	else
		echo
	fi
	rm -f qqx qqy
}

########################################
# doValidateFile()
#   Simulate a POST command to the server and use
#   the supplied file name as the json data
#	Parameters:
# 		$1 = base file name
#		$2 = title
########################################
doValidateFile () {
	TESTCOUNT=$((TESTCOUNT + 1))
	printf "PHASE %2s  %3s  %s... " ${TESTCOUNT} $1 $2
	echo >> ${1}	# force a newline at the end of the file, which often doesn't happen with command output

	checkPause

	if [ "${FORCEGOOD}" = "1" ]; then
		cp ${1} ${GOLD}/${1}.gold
		echo "DONE"
	elif [ "${SKIPCOMPARE}" = "0" ]; then
		if [ ! -f ${GOLD}/${1}.gold ]; then
			echo "UNSET CONTENT" > ${GOLD}/${1}.gold
			echo "Created a default ${GOLD}/$1.gold for you. Update this file with known-good output."
		fi

		#--------------------------------------------------------------------
		# The actual data has timestamp information that changes every run.
		# The timestamp can be filtered out for purposes of testing whether
		# or not the web service could be called and can return the expected
		# data.
		#--------------------------------------------------------------------
		declare -a out_filters=(
			's/(^[ \t]+"LastModTime":).*/$1 TIMESTAMP/'
		)
		cp gold/${1}.gold qqx
		cp ${1} qqy
		for f in "${out_filters[@]}"
		do
			perl -pe "$f" qqx > qqx1; mv qqx1 qqx
			perl -pe "$f" qqy > qqy1; mv qqy1 qqy
		done

		UDIFFS=$(diff qqx qqy | wc -l)
		if [ ${UDIFFS} -eq 0 ]; then
			if [ ${SHOWCOMMAND} -eq 1 ]; then
				echo "PASSED	cmd: ${CMD}"
			else
				echo "PASSED"
			fi
		else
			echo "FAILED..." >> ${ERRFILE}
			echo "Differences in ${1} are as follows:" >> ${ERRFILE}
			diff qqx qqy >> ${ERRFILE}
			echo "If correct:  mv ${1} ${GOLD}/${1}.gold" >> ${ERRFILE}
			echo "Command to reproduce:  ${CMD}" >> ${ERRFILE}
			cat ${ERRFILE}
			failmsg
			if [ "${ASKBEFOREEXIT}" = "1" ]; then
				pause ${1}
			else
				if [ ${MANAGESERVER} -eq 1 ]; then
					echo "STOPPING RENTROLL SERVER"
					pkill rentroll
				fi
				exit 1
			fi
		fi
	else
		echo
	fi
	rm -f qqx qqy
}

########################################
# doPlainGET()
#   Simulate a GET command to the server and use
#   the supplied file name as the json data
#	Parameters:
# 		$1 = url
# 		$2 = base file name
#		$3 = title
########################################
doPlainGET () {
	TESTCOUNT=$((TESTCOUNT + 1))
	printf "PHASE %2s  %3s  %s... " ${TESTCOUNT} ${2} ${3}
	CMD="curl -s ${1}"
	${CMD} > ${2} 2>>${LOGFILE}
	echo >> ${2}   # this ensures a newline in the file

	checkPause

	if [ "${FORCEGOOD}" = "1" ]; then
		goldpath
		cp ${2} ${GOLD}/${2}.gold
		echo "DONE"
	elif [ "${SKIPCOMPARE}" = "0" ]; then
		if [ ! -f ${GOLD}/${2}.gold ]; then
			echo "UNSET CONTENT" > ${GOLD}/${2}.gold
			echo "Created a default ${GOLD}/$1.gold for you. Update this file with known-good output."
		fi

		#--------------------------------------------------------------------
		# The actual data has timestamp information that changes every run.
		# The timestamp can be filtered out for purposes of testing whether
		# or not the web service could be called and can return the expected
		# data.
		#--------------------------------------------------------------------
		declare -a out_filters=(
			's/(^[ \t]+"LastModTime":).*/$1 TIMESTAMP/'
		)
		cp gold/${2}.gold qqx
		cp ${2} qqy
		for f in "${out_filters[@]}"
		do
			perl -pe "$f" qqx > qqx1; mv qqx1 qqx
			perl -pe "$f" qqy > qqy1; mv qqy1 qqy
		done

		UDIFFS=$(diff qqx qqy | wc -l)
		if [ ${UDIFFS} -eq 0 ]; then
			if [ ${SHOWCOMMAND} -eq 1 ]; then
				echo "PASSED	cmd: ${CMD}"
			else
				echo "PASSED"
			fi
		else
			echo "FAILED..." >> ${ERRFILE}
			echo "Differences in ${2} are as follows:" >> ${ERRFILE}
			diff qqx qqy >> ${ERRFILE}
			echo "If correct:  mv ${2} ${GOLD}/${2}.gold" >> ${ERRFILE}
			echo "Command to reproduce:  ${CMD}" >> ${ERRFILE}
			cat ${ERRFILE}
			failmsg
			if [ "${ASKBEFOREEXIT}" = "1" ]; then
				pause "${2}"
			else
				if [ ${MANAGESERVER} -eq 1 ]; then
					echo "STOPPING RENTROLL SERVER"
					pkill rentroll
				fi
				exit 1
			fi
		fi
	else
		echo
	fi
	rm -f qqx qqy
}

#############################################################################
# doCasperUITest()
#	UI automation testing in headless mode with casperjs and phantomjs
#
#	Parameters:
# 		$1 = base file name
#		$2 = javascript file from which UI automated testing should be performed with options
# 		$3 = title for reporting.  No spaces
#############################################################################
doCasperUITest () {
	TESTCOUNT=$((TESTCOUNT + 1))
	printf "PHASE %2s  %3s  %s... " ${TESTCOUNT} $1 $3

	if [ "x${2}" != "x" ]; then
		${CASPERTEST} ${2} >${1} 2>&1
	fi

	checkPause

	if [ "${FORCEGOOD}" = "1" ]; then
		goldpath
		cp ${1} ${GOLD}/${1}.gold
		echo "DONE"
	elif [ "${SKIPCOMPARE}" = "0" ]; then
		if [ ! -f ${GOLD}/${1}.gold ]; then
			echo "UNSET CONTENT" > ${GOLD}/${1}.gold
			echo "Created a default ${GOLD}/$1.gold for you. Update this file with known-good output."
		fi
		declare -a out_filters=(
			's/executed in +([0-9]*\.[0-9]*)s,/executed in (EXECUTED)s,/g'
			's/skipped.*/skipped./g'
			's/^Date:.*/current time/'
			's/\s+[0-1]?[0-9]\/[0-3]?[0-9]\/[0-9][0-9][^-]*/date/g'
			's/\s+[0-1]?[0-9]\/[0-3]?[0-9]\/20[0-9][0-9][^-]*/date/g'
		)
		cp ${GOLD}/${1}.gold ${GOLD}/${1}.g
		cp ${1} ${1}.g
		for f in "${out_filters[@]}"
		do
			perl -pe "$f" ${GOLD}/${1}.g > ${GOLD}/${1}.t; mv ${GOLD}/${1}.t ${GOLD}/${1}.g
			perl -pe "$f" ${1}.g > ${1}.t; mv ${1}.t ${1}.g
		done
		UDIFFS=$(diff ${1}.g ${GOLD}/${1}.g | wc -l)
		if [ ${UDIFFS} -eq 0 ]; then
			if [ ${SHOWCOMMAND} -eq 1 ]; then
				echo "PASSED	cmd: ${CASPERTEST} ${2}"
			else
				echo "PASSED"
			fi
			rm -f ${1}.g ${GOLD}/${1}.g
		else
			echo "FAILED...   if correct:  mv ${1} ${GOLD}/${1}.gold" >> ${ERRFILE}
			echo "Command to reproduce:  ${CASPERTEST} ${2}" >> ${ERRFILE}
			echo "Differences in ${1} are as follows:" >> ${ERRFILE}
			diff ${GOLD}/${1}.g ${1}.g >> ${ERRFILE}
			cat ${ERRFILE}
			failmsg
			if [ "${ASKBEFOREEXIT}" = "1" ]; then
				pause ${1}
			else
				exit 1
			fi
		fi
	else
		echo
	fi
}

#------------------------------------------------------------------------------
#  encodeURI encodes data so that it can be passed in a URI.  It
#      does essentially what Javascript's encodeURI does.
#
#  INPUTS
#  $1  The string to encode
#
#  RETURNS
#      the return value is the encoded string.
#
#  USAGE:
#      data=$(encodeURI "4%2F1%2F2019")
#------------------------------------------------------------------------------
encodeURI() {
  local string="${1}"
  local strlen=${#string}
  local encoded=""
  local pos c o

  for (( pos=0 ; pos<strlen ; pos++ )); do
     c=${string:$pos:1}
     case "$c" in
        [-_.~a-zA-Z0-9] ) o="${c}" ;;
        * )               printf -v o '%%%02x' "'$c"
     esac
     encoded+="${o}"
  done
  echo "${encoded}"
}

#------------------------------------------------------------------------------
#  encodeRequest is just like encodeURI except that it saves the output
#      into a file named "request"
#
#  INPUTS
#  $1  The string to encode
#
#  RETURNS
#      nothing, but the encoded string will be in a file named "request"
#------------------------------------------------------------------------------
encodeRequest() {
  local string="${1}"
  local strlen=${#string}
  local encoded=""
  local pos c o

  for (( pos=0 ; pos<strlen ; pos++ )); do
     c=${string:$pos:1}
     case "$c" in
        [-_.~a-zA-Z0-9] ) o="${c}" ;;
        * )               printf -v o '%%%02x' "'$c"
     esac
     encoded+="${o}"
  done
  echo "${encoded}" > request
}


#--------------------------------------------------------------------------
#  Handle command line options...
#--------------------------------------------------------------------------
tdir
while getopts "acfmoprnR:t:" o; do
	echo "o = ${o}"
	case "${o}" in
		a)	ASKBEFOREEXIT=1
			echo "WILL ASK BEFORE EXITING ON ERROR"
			;;
		c | C)
			SHOWCOMMAND=1
			echo "SHOWCOMMAND"
			;;
		r | R)
			doReport
			exit 0
			;;
		p | P)
			PAUSE=1
			echo "PAUSE BETWEEN COMMANDS"
			;;
		f)  SKIPCOMPARE=1
			echo "SKIPPING COMPARES..."
			;;
		m)  MANAGESERVER=0
			echo "SKIPPING SERVER MGMT CMDS..."
			;;
		n)	CREATENEWDB=0
			echo "DATA WILL BE ADDED TO CURRENT DB"
			;;
		o)	FORCEGOOD=1
			echo "OUTPUT OF THIS RUN IS SAVED AS *.GOLD"
			;;
		t) SINGLETEST="${OPTARG}"
			echo "SINGLETEST set to ${SINGLETEST}"
			;;
		*) 	usage
			exit 1
			;;
	esac
done
shift $((OPTIND-1))


rm -f ${ERRFILE}
echo    "Test Name:    ${TESTNAME}" > ${LOGFILE}
echo    "Test Purpose: ${TESTSUMMARY}" >> ${LOGFILE}
echo -n "Date/Time:    " >>${LOGFILE}
date >> ${LOGFILE}
echo >>${LOGFILE}
