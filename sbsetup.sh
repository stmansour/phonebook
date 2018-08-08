#!/bin/bash
#------------------------------------------------------------------------------
#  sbsetup.sh
#  Set up a sandbox node before phonebook is launched for the first time.
#  The tasks in this file need only be done once, before phonebook is launched.
#
#  Tasks:
#	1. Pull down the development config file for sandboxes
#
#  Notes:
#   1. Env values:
#      0 = development
#      1 = production
#      2 = QA
#------------------------------------------------------------------------------

HOST=localhost
PROGNAME="phonebook"
PORT=8270
WATCHDOGOPTS=""
GETFILE="/usr/local/accord/bin/getfile.sh"
DATABASENAME="accord"
DBUSER="ec2-user"
IAM=$(whoami)
OS=$(uname)
MYSQL="mysql --no-defaults"
#--------------------------------------------------------------
#  For QA, Sandbox, and Production nodes, go through the
#  laundry list of details...
#  1. Set up permissions for the database on QA and Sandbox nodes
#  2. Install a database with some data for testing
#  3. For PDF printing, install wkhtmltopdf
#--------------------------------------------------------------
setupAppNode() {
	#---------------------
	# database
	#---------------------
    rm -rf ${DATABASENAME}db*  >log.out 2>&1
    ${GETFILE} accord/db/${DATABASENAME}db.sql.gz  >log.out 2>&1
    gunzip ${DATABASENAME}db.sql  >log.out 2>&1
	echo "DROP USER IF EXISTS 'ec2-user'@'localhost';CREATE USER 'ec2-user'@'localhost';GRANT ALL PRIVILEGES ON *.* TO 'ec2-user'@'localhost' WITH GRANT OPTION;" | mysql > log.out 2>&1
    echo "DROP DATABASE IF EXISTS ${DATABASENAME}; CREATE DATABASE ${DATABASENAME}; USE ${DATABASENAME};" > restore.sql
    echo "source ${DATABASENAME}db.sql" >> restore.sql
    ${MYSQL} < restore.sql  >log.out 2>&1
}

commonSetup() {
    if [ ${IAM} == "root" ]; then
        chown -R ec2-user *
        chmod u+s phonebook pbwatchdog picsync.sh
        if [ $(uname) == "Linux" -a ! -f "/etc/init.d/phonebook" ]; then
            cp ./activate.sh /etc/init.d/phonebook
            chkconfig --add phonebook
        fi
    fi
    #-----------------------------------------------------------------
	#  If no config.json exists, pull the development environment
	#  version and use it.  The Env values mean the following:
	#    0 = development environment
	#    1 = production environment
	#    2 = QA environment
	#-----------------------------------------------------------------
	if [ ! -f ./config.json ]; then
		${GETFILE} accord/db/confdev.json  >log.out 2>&1
		mv confdev.json config.json
	fi
    if [ ! -d "./images" ]; then
        /usr/local/accord/bin/getfile.sh jenkins-snapshot/phonebook/latest/pbimages.tar.gz >phonebook.log 2>&1
        gunzip -f pbimages.tar.gz >phonebook.log 2>&1
        tar xvf pbimages.tar >phonebook.log 2>&1
    fi
    if [ ! -f "/usr/local/share/man/man1/pbbkup.1" ]; then
        ./installman.sh >phonebook.log 2>&1
    fi
}

commonSetup
setupAppNode
