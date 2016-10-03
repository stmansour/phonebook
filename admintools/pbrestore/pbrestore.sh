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
	echo "GRANT ALL PRIVILEGES ON ${DATABASE} TO 'ec2-user'@'localhost' WITH GRANT OPTION;" >> restore.sql
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

#####################################################################################
#  UPDATE SCHEMA
#  After the initial release of Phonebook, the schema was updated.  This routine
#  checks the schema of the database just restored.  If the schema needs to be
#  updated, it will make the updates.
#####################################################################################
updateSchema() {
# Generate the mysql commands needed to validate...
cat >xxqq <<EOF
use accord;
describe classes;
EOF

	mysql ${MYSQLOPTS} <xxqq >xxqqout
	rm -f xxqq
	HASCOCODE=$(grep CoCode xxqqout | wc -l)
	rm -f xxqqout
	if [ ${HASCOCODE} -ne 1 ]; then
cat >xxqq <<EOF
use accord;
ALTER TABLE classes ADD CoCode MEDIUMINT NOT NULL DEFAULT 0 AFTER ClassCode;
EOF
		mysql ${MYSQLOPTS} <xxqq >xxqqout
		./roleinit >roleinit.log
	fi
}

#####################################################################################
#  UPDATE COCODE
#  For the first month, the links will not be in place between classes and companies
#  In order for the sample data to work correctly, this is being introduced as a
#  TEMPORARY hack to establish some important links. This function should be 
#  removed on or about November 1, 2016
#####################################################################################
updateCoCode() {
cat >xxqq1 <<EOF
use accord;
update classes set CoCode=24 where ClassCode=10;
update classes set CoCode=24 where ClassCode=12;
EOF
		mysql ${MYSQLOPTS} <xxqq1
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

#------------------------------------------------------
# this will update the schema only if necessary...
#------------------------------------------------------
updateSchema

# Remove this line after November 1, 2016
updateCoCode
