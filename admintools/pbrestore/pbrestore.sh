#!/bin/bash

FULL=0
GET="/usr/local/accord/bin/getfile.sh"
RESTORE="/usr/local/accord/testtools/restoreMySQLdb.sh"
DATABASE="accord"

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
	${GET} ${DATABASE}/db/${DATABASE}db.tar
	tar xvf ${DATABASE}db.tar
	echo "Done."

	echo "Extracting pictures"
	gunzip pictures.tar.gz
	tar xvf pictures.tar
	echo "Done."
	
	echo "Extracting data"
	gunzip ${DATABASE}db.sql.gz
	${RESTORE} ${DATABASE} ${DATABASE}db.sql
	echo "Done."

	echo "Cleaning up..."
	rm -f pictures.tar ${DATABASE}db.sql ${DATABASE}db.tar
	echo "Done."
}

##############################################################
#   RESTORE - DATA ONLY
##############################################################
restoreData() {
	echo "Retrieving backup data from Artifactory"
	${GET} ${DATABASE}/db/${DATABASE}db.sql.gz
	echo "Done."

	echo "Extracting data"
	gunzip ${DATABASE}db.sql.gz
	${RESTORE} ${DATABASE} ${DATABASE}db.sql
	echo "Done."

	echo "Cleaning up..."
	rm -f ${DATABASE}db.sql
	echo "Done."
}


##############################################################
#   MAIN ROUTINE
##############################################################
while getopts ":fhN:" o; do
    case "${o}" in
        f)
            FULL=1
            ;;
        h)
			usage
            ;;
        N)
			DATABASE=${OPTARG}
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
