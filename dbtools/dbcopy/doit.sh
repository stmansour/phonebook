#!/bin/bash
# This script drops the current

echo "Dropping existing 'accordtest' database and creating"
echo "new 'accordtest' based on new table schemas"
pushd ../schema
./apply.sh
popd

echo "Copying data from database 'accord' to database 'accordtest'"
mysql < migrate.sql