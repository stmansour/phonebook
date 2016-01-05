#!/bin/bash
DOM=$(date +%d)
FULL=0
DEPLOY="/usr/local/accord/bin/deployfile.sh"
DATABASE="accord"

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
	mysqldump ${DATABASE} > ${DATABASE}db.sql
	gzip ${DATABASE}db.sql
	tar cvf ${DATABASE}db.tar pictures.tar.gz ${DATABASE}db.sql.gz

	${DEPLOY} ${DATABASE}db.tar ${DATABASE}/db

	rm -f pictures.tar.gz ${DATABASE}db.sql.gz ${DATABASE}db.tar
}

##############################################################
#   BACKUP - DATA ONLY
##############################################################
bkupData() {
	mysqldump ${DATABASE} > ${DATABASE}db.sql
	gzip ${DATABASE}db.sql
	${DEPLOY} ${DATABASE}db.sql.gz ${DATABASE}/db/${DOM}
	rm -f ${DATABASE}db.sql.gz
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
	echo "Backing up data on database ${DATABASE}..."
	bkupData
elif [ ${FULL} -eq 1 ]; then
	echo "Backing up data and pictures for database ${DATABASE}..."
	bkupFull
fi

echo
echo "Done."