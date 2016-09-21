#!/bin/bash

# Sept 21, 2016
# I'm reviewing this directory about 9 months after the production database was created.
# I think this script was used to create the first production database. 

OUTFILE="newproddb.sql"
DATABASE="accord"

echo "DROP DATABASE IF EXISTS ${DATABASE}; CREATE DATABASE ${DATABASE}; USE ${DATABASE};" > ${OUTFILE}
echo "source accorddb.sql" >> ${OUTFILE}
echo "GRANT ALL PRIVILEGES ON accord TO 'ec2-user'@'localhost' WITH GRANT OPTION;" >> ${OUTFILE}
echo "GRANT ALL PRIVILEGES ON accord TO 'adbuser'@'localhost' WITH GRANT OPTION;" >> ${OUTFILE}

mysql -h phbk.cjkdwqbdvxyu.us-east-1.rds.amazonaws.com -P 3306 -u adbuser -p < ${OUTFILE}
