#!/bin/bash
DOM=$(date +%d)
FULL=0
GET="/usr/local/accord/bin/getfile.sh"
RESTORE="/usr/local/accord/testtools/restoreMySQLdb.sh"
DATABASE="accord"
MYSQLOPTS=""

##############################################################
#   USAGE
##############################################################
usage() {
    cat <<ZZEOF
ACCORD Phonebook database restore utility
Usage:   pbrestore.sh [OPTIONS]

OPTIONS:
-d  day-of-month. Default is today's date. Note that if the daily
                  backup has not been performed yet, this would restore
                  last month's data.
-f 	full restore -- includes pictures. Default is data only
-h	print this help message
-n  force --no-defaults onto all mysql commands

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

restoreMySQLdb() {
DB=$1
DBfile=$2

echo "DROP DATABASE IF EXISTS ${DB}; CREATE DATABASE ${DB}; USE ${DB};" > restore.sql
echo "source ${DBfile}" >> restore.sql
echo "GRANT ALL PRIVILEGES ON accord TO 'ec2-user'@'localhost' WITH GRANT OPTION;" >> restore.sql
mysql ${MYSQLOPTS} < restore.sql

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
	
	echo "DROP DATABASE IF EXISTS ${DATABASE}; CREATE DATABASE ${DATABASE}; USE ${DATABASE};" > restore.sql
	echo "source ${DATABASE}db.sql" >> restore.sql
	echo "GRANT ALL PRIVILEGES ON accord TO 'ec2-user'@'localhost' WITH GRANT OPTION;" >> restore.sql
	mysql ${MYSQLOPTS} < restore.sql
	echo "Done."

	echo "Cleaning up..."
	rm -f pictures.tar ${DATABASE}db.sql ${DATABASE}db.tar
	echo "Done."
}

##############################################################
#   RESTORE - DATA ONLY
##############################################################
restoreData() {
	echo "Retrieving backup data from Artifactory: ${DATABASE}/db/${DOM}/accorddb.sql.gz"
	${GET} ${DATABASE}/db/${DOM}/${DATABASE}db.sql.gz
	echo "Done."

	if [ ! -f ${DATABASE}db.sql.gz ]; then
		echo "failed to download ${DATABASE}/db/${DOM}/accorddb.sql.gz"
		exit 1
	fi

	echo "Extracting data"
	gunzip ${DATABASE}db.sql.gz

	echo "DROP DATABASE IF EXISTS ${DATABASE}; CREATE DATABASE ${DATABASE}; USE ${DATABASE};" > restore.sql
	echo "source ${DATABASE}db.sql" >> restore.sql
	echo "GRANT ALL PRIVILEGES ON accord TO 'ec2-user'@'localhost' WITH GRANT OPTION;" >> restore.sql
	mysql ${MYSQLOPTS} < restore.sql
	echo "Done."

	echo "Cleaning up..."
	rm -f ${DATABASE}db.sql
	echo "Done."
}


##############################################################
#   MAIN ROUTINE
##############################################################
while getopts ":d:fhnN:" o; do
    case "${o}" in
        d)
            DOM=${OPTARG}
            if [ ${DOM} -gt 31 ]; then
            	echo "Largest value for DOM is 31."
            	exit 1
            fi
            if [ ${DOM} -lt 1 ]; then
            	echo "Small value for DOM is 1."
            	exit 1
            fi
            ;;
        f)
            FULL=1
            ;;
        h)
			usage
            ;;
        N)
			DATABASE=${OPTARG}
            ;;
        n)
			MYSQLOPTS="--no-defaults"
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
