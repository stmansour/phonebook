#!/bin/bash
DBNAME=accordtest
HOST=localhost
PORT=8250
PEOPLEOPT=""

usage() {
    cat << ZZEOF
Phonebook Database Role / Fieldperm Initialization Script
Usage:   apply.sh [OPTIONS]

OPTIONS:
-p port           (default is 8250)
-h hostname       (default is localhost)
-N database name  (default is accordtest)

Examples:
Command to create roles in accordtest:
	bash$  apply.sh 

Command to create roles in a database named 'accord':
	bash$  apply.sh -n accord

ZZEOF
	exit 0
}

while getopts ":p:h:N:r" o; do
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
        r)
            PEOPLEOPT="-r"
	    	echo "PEOPLEOPT set to: ${PEOPLEOPT}"
            ;;           
        *)
            usage
            ;;
    esac
done
shift $((OPTIND-1))

./roleinit -N ${DBNAME} ${PEOPLEOPT}
