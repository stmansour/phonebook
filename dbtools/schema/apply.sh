#!/bin/bash
DBNAME=accordtest
HOST=localhost
PORT=8250
PEOPLEOPT=""
SCHEMAFILE="schema.sql"

usage() {
    cat <<ZZEOF
Phonebook Database SCHEMA Initialization Script
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

genGrants() {
    cat >${SCHEMAFILE} <<ZZEOF1
-- ACCORD PHONEBOOK DATABASE
-- mysql> show grants for 'ec2-user'@'localhost';
-- +-----------------------------------------------------------------------------+
-- | Grants for ec2-user@localhost                                               |
-- +-----------------------------------------------------------------------------+
-- | GRANT USAGE ON *.* TO 'ec2-user'@'localhost'                                |
-- | GRANT ALL PRIVILEGES ON accordtest.* TO 'ec2-user'@'localhost'              |
-- | GRANT ALL PRIVILEGES ON accordtest.accordtest TO 'ec2-user'@'localhost'     |
-- +-----------------------------------------------------------------------------+

DROP DATABASE IF EXISTS ${DBNAME};
CREATE DATABASE ${DBNAME};
USE ${DBNAME};
GRANT ALL PRIVILEGES ON ${DBNAME} TO 'ec2-user'@'localhost';
GRANT ALL PRIVILEGES ON ${DBNAME}.* TO 'ec2-user'@'localhost';
ZZEOF1

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

genGrants
cat tables.sql >>${SCHEMAFILE}

mysql < ${SCHEMAFILE}
