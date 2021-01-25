#!/bin/bash

#==========================================================================
#  This script performs SQL schema changes on the test databases that are
#  saved as SQL files in the test directory. It loads them, performs the
#  ALTER commands, then saves the sql file.
#
#  If the test file uses its own database saved as a .sql file, make sure
#  it is listed in the dbs array
#==========================================================================

MODFILE="dbqqqmods.sql"
MYSQL="mysql --no-defaults"
MYSQLDUMP="mysqldump --no-defaults"
DBNAME="accord"

#=====================================================
#  Retain prior changes as comments below
#  Examples:
# ALTER TABLE Property MODIFY LeaseCommencementDt RentCommencementDt DATETIME NOT NULL DEFAULT '1970-01-01 00:00:00';
# ALTER TABLE Property DROP COLUMN HQAddress, DROP COLUMN HQAddress2, DROP COLUMN HQPostalCode, DROP COLUMN HQCountry;
# ALTER TABLE Property ADD TermRemainingOnLeaseUnits SMALLINT NOT NULL DEFAULT 0 AFTER TermRemainingOnLease;
# ALTER TABLE Property CHANGE CreateTS CreateTime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;
#=====================================================

#=====================================================
#  Put modifications to schema in the lines below
#=====================================================

# ALTER TABLE people ADD CreateTime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP AFTER LastModBy;
# ALTER TABLE people ADD CreateBy BIGINT NOT NULL DEFAULT 0 AFTER CreateTime;

#=====================================================
#  END of modification history
#=====================================================

cat > "${MODFILE}" << LEOF
LEOF

#=====================================================
#  Put dir/sqlfilename in the list below
#=====================================================
declare -a dbs=(
	"../directory.sql"
	"../testdb.sql"
	"../accord.sql"
	"ws/dbws.sql"
)

for f in "${dbs[@]}"
do
	echo "DROP DATABASE IF EXISTS ${DBNAME}; CREATE DATABASE ${DBNAME}; USE ${DBNAME}; GRANT ALL PRIVILEGES ON ${DBNAME}.* TO 'ec2-user'@'localhost';" | ${MYSQL}
	echo -n "${f}: loading... "
	${MYSQL} ${DBNAME} < ${f}
	echo -n "updating... "
	${MYSQL} ${DBNAME} < ${MODFILE}
	echo -n "saving... "
	${MYSQLDUMP} ${DBNAME} > ${f}
	echo "done"
done
