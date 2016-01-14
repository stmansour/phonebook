#!/bin/bash

DBTOOLSDIR="../../dbtools"
USERSIMDIR="."
USERSIM="${USERSIMDIR}/usersim"

DBNAME="accordtest"
DBUSER="ec2-user"

PBHOST="localhost"
PBPORT="8250"

PHONEBOOKDIR="../.."
STARTPHONEBOOKCMD="./activate.sh -N ${DBNAME} -T start"
STOPPHONEBOOKCMD="./activate.sh stop"

initDB() {
	pushd ${DBTOOLSDIR} >/dev/null
	./apply.sh > /dev/null
	popd >/dev/null
	${USERSIM} -f -u 300 -c 70 -C 70
}

stopPhonebook() {
	L=$(ps -ef | grep phonebook | grep -v grep | wc -l)
	if [ ${L} -gt 0 ]; then
		pushd ${PHONEBOOKDIR} >/dev/null
		results=$(${STOPPHONEBOOKCMD})
		# echo "results = ${results}"
		result=$(echo "${results}" | grep OK | wc -l)
		if [ ${result} -ge 0 ]; then
			echo "phonebook stopped"
			sleep 5
		else 
			echo "phonebook did not stop properly.  result = \"${result}\""
			exit 1
		fi
		popd >/dev/null
	fi
}

startPhonebook() {
	pushd ${PHONEBOOKDIR} >/dev/null
	if [ ! -d "./images" ]; then
		/usr/local/accord/bin/getfile.sh jenkins-snapshot/phonebook/latest/pbimages.tar.gz >phonebook.log 2>&1
		gunzip -f pbimages.tar.gz >phonebook.log 2>&1
		tar xvf pbimages.tar >phonebook.log 2>&1
	fi
	
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

L=$(ps -ef | grep phonebook | grep -v grep | wc -l)
if [ ${L} -gt 0 ]; then
	echo "Could not stop running instance of phonebook..."
	ps -ef | grep phonebook | grep -v grep 
	exit 1
fi

echo "Creating new db: ${DBNAME}"
initDB

echo "Starting phonebook service"
startPhonebook

L=$(ps -ef | grep phonebook | grep -v grep | wc -l)
if [ ${L} -ne 1 ]; then
	echo "Could not find one and only one running instance of phonebook..."
	ps -ef | grep phonebook | grep -v grep 
	exit 1
fi

usrsimout=$(${USERSIM})

echo "results: ${usrsimout}"
rs=$(echo "${usrsimout}" | grep "fail: 0" | wc -l)
if [ ${rs} -eq 1 ]; then
	exit 0
else
	exit 1
fi