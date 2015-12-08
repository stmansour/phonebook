#!/bin/bash

FULL=0
GET="/usr/local/accord/bin/getfile.sh"
RESTORE="/usr/local/accord/testtools/restoreMySQLdb.sh"

##############################################################
#   USAGE
##############################################################
usage() {
    cat << ZZEOF
ACCORD Phonebook database restore utility
Usage:   pbrestore.sh [OPTIONS]

OPTIONS:
-f 	full restore -- includes pictures. Default is data only
-h	print this help message

Examples:
Command to restore data only
	bash$  ./pbrestore.sh 

Command to restore data and pictures:
	bash$  ./pbrestore.sh -f

Command to get help
    bash$  ./pbrestore.sh -h

ZZEOF
	exit 0
}

##############################################################
#   RESTORE - FULL
##############################################################
restoreFull() {
	echo "Retrieving backup data from Artifactory"
	${GET} accord/db/accorddb.tar
	tar xvf accorddb.tar
	echo "Done."

	echo "Extracting pictures"
	gunzip pictures.tar.gz
	tar xvf pictures.tar
	echo "Done."
	
	echo "Extracting data"
	gunzip accorddb.sql.gz
	${RESTORE} accord accorddb.sql
	echo "Done."

	echo "Cleaning up..."
	rm -f pictures.tar accorddb.sql accorddb.tar
	echo "Done."
}

##############################################################
#   RESTORE - DATA ONLY
##############################################################
restoreData() {
	echo "Retrieving backup data from Artifactory"
	${GET} accord/db/accorddb.sql.gz
	echo "Done."

	echo "Extracting data"
	gunzip accorddb.sql.gz
	${RESTORE} accord accorddb.sql
	echo "Done."

	echo "Cleaning up..."
	rm -f accorddb.sql
	echo "Done."
}


##############################################################
#   MAIN ROUTINE
##############################################################
while getopts ":fh" o; do
    case "${o}" in
        f)
            FULL=1
            ;;
        h)
			usage
            ;;
        *)
			echo "UNRECOGNIZED OPTION:  ${o}"
            usage
            ;;
    esac
done
shift $((OPTIND-1))

if [ ${FULL} -eq 0 ]; then
	echo "Restore data..."
	restoreData
elif [ ${FULL} -eq 1 ]; then
	echo "Restore data and pictures..."
	restoreFull
fi
