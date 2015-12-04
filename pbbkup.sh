#!/bin/bash

FULL=0
DEPLOY="/usr/local/accord/bin/deployfile.sh"

##############################################################
#   USAGE
##############################################################
usage() {
    cat << ZZEOF
ACCORD Phonebook database backup utility
Usage:   pbbkup.sh [OPTIONS]

OPTIONS:
-f 	full backup -- includes pictures. Default is data only
-h	print this help message

Examples:
Command to backup data only
	bash$  ./pbbkup.sh 

Command to backup data and pictures:
	bash$  ./pbbkup.sh -f

Command to get help
    bash$  ./pbbkup.sh -h

ZZEOF
	exit 0
}

##############################################################
#   BACKUP - FULL
##############################################################
bkupFull() {
	tar cvf pictures.tar pictures
	gzip pictures.tar
	mysqldump accord > accorddb.sql
	gzip accorddb.sql
	tar cvf accorddb.tar pictures.tar.gz accorddb.sql.gz

	${DEPLOY} accorddb.tar accord/db

	rm -f pictures.tar.gz accorddb.sql.gz accorddb.tar
}

##############################################################
#   BACKUP - DATA ONLY
##############################################################
bkupData() {
	mysqldump accord > accorddb.sql
	gzip accorddb.sql
	${DEPLOY} accorddb.sql.gz accord/db
	rm -f accorddb.sql.gz
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
	echo "Backing up data..."
	bkupData
elif [ ${FULL} -eq 1 ]; then
	echo "Backing up data and pictures..."
	bkupFull
fi

echo
echo "Done."