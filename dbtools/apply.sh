#!/bin/bash
DBNAME=accord
HOST=localhost
PORT=8250
PEOPLEOPT=""
SCHEMAFILE="schema.sql"

usage() {
    cat <<ZZEOF
Create New Phonebook Database
Usage:   apply.sh [OPTIONS]

OPTIONS:
-p port           (default is 8250)
-h hostname       (default is localhost)
-N database name  (default is accordtest)

Examples:
Command to create roles in accordtest:
	bash$  apply.sh

Command to create roles in a database named 'accord':
	bash$  apply.sh -N accord

ZZEOF
	exit 0
}


while getopts ":p:h:N:" o; do
    case "${o}" in
        h)
            HOST=${OPTARG}
            echo "HOST set to: ${HOST}"
            ;;
        p)
            PORT=${OPTARG}
	    	echo "PORT set to: ${PORT}"
            ;;
        N)
            DBNAME=${OPTARG}
	    	echo "DBNAME set to: ${DBNAME}"
            ;;
        *)
            usage
            ;;
    esac
done
shift $((OPTIND-1))

echo "schema";      pushd ../db/schema; ./apply.sh -N ${DBNAME}    ; popd
echo "roleinit";    pushd roleinit;     ./apply.sh -N ${DBNAME}    ; popd
echo "jobtitles";   pushd jobtitles;    ./jobtitles -N ${DBNAME}   ; popd
echo "deductions";  pushd deductions;   ./deductions -N ${DBNAME}  ; popd
echo "departments"; pushd departments;  ./departments -N ${DBNAME} ; popd
