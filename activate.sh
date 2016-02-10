#!/bin/bash
# chkconfig: 345 99 01
# description: activation script to start/stop Accord Phonebook
#
# processname: phonebook
# pidfile: /var/run/phonebook/phonebook.pid


HOST=localhost
PORT=8250
STARTPBONLY=0
WATCHDOGOPTS=""
QA=0

DBNAME="accord"
DBUSER="ec2-user"
IAM=$(whoami)


usage() {
    cat <<ZZEOF
Phonebook activation script.
Usage:   activate.sh [OPTIONS] CMD

This is the Accord Phonebook activation script. It is designed to work in two environments.
First, it works with Plum - Accord's test environment automation infrastructure
Second, it can work as a service script in /etc/init.d

OPTIONS:
-b 			 (enable db backups, default is off)
-p port      (default is 8250)
-h hostname  (default is localhost)
-N dbname    (default is accord)
-q           (start as qa version)
-T           (use this option to indicate testing rather than production)

CMD is one of: start | stop | status | restart | ready | reload | condrestart | images

Examples:
Command to start phonebook:
	bash$  activate.sh start 

Command to start phonebook for testing purposes:
	bash$  activate.sh -T start 

	If you do testing in the context of the source code tree and you don't 
	use -T, you may see messages like this:

	  ./activate.sh: line 101: ./pbwatchdog: No such file or directory

Command to stop phonebook:
	bash$  activate.sh Stop

Command to see if phonebook is ready for commands... the response
will be "OK" if it is ready, or something else if there are problems:

    bash$  activate.sh ready
    OK
ZZEOF
}

updateImages() {
	/usr/local/accord/bin/getfile.sh jenkins-snapshot/phonebook/latest/pbimages.tar.gz
	gunzip pbimages.tar.gz
	tar xvf pbimages.tar
}

stopwatchdog() {
	# make sure we can find it
    pidline=$(ps -ef | grep pbwatchdog | grep -v grep)
    if [ "${pidline}" != "" ]; then
        lines=$(echo "${pidline}" | wc -l)
        if [ $lines -gt "0" ]; then
            pid=$(echo "${pidline}" | awk '{print $2}')
            $(kill $pid)
        fi          
    fi      
}

start() {
	if [ 0 -eq ${QA} ]; then
		if [ ${IAM} == "root" ]; then
			chown -R ec2-user *
			chmod u+s phonebook pbwatchdog
			if [ $(uname) == "Linux" -a ! -f "/etc/init.d/phonebook" ]; then
				cp ./activate.sh /etc/init.d/phonebook
				chkconfig --add phonebook
			fi
		fi

		if [ "${STARTPBONLY}" -ne "1" ]; then
			if [ ! -d "./images" ]; then
				/usr/local/accord/bin/getfile.sh jenkins-snapshot/phonebook/latest/pbimages.tar.gz >phonebook.log 2>&1
				gunzip -f pbimages.tar.gz >phonebook.log 2>&1
				tar xvf pbimages.tar >phonebook.log 2>&1
			fi
			if [ ! -f "/usr/local/share/man/man1/pbbkup.1" ]; then
				./installman.sh >phonebook.log 2>&1
			fi
		fi
		./phonebook -N ${DBNAME} >pbconsole.out 2>&1 &
		if [ ${IAM} == "root" ]; then
			if [ ! -d /var/run/phonebook ]; then
				mkdir /var/run/phonebook
			fi
			echo $! >/var/run/phonebook/phonebook.pid
			touch /var/lock/phonebook
		fi

		# give phonebook a few seconds to start up before initiating the watchdog
		sleep 5
		if [ "${STARTPBONLY}" -ne "1" ]; then
			stopwatchdog
			./pbwatchdog ${WATCHDOGOPTS} >pbwatchdogstartup.out 2>&1 &
		fi
	elif [[ ${QA} -eq 1 ]]; then
		# echo "STARTING FOR QA"
		pushd ./test/usersim >/dev/null 2>&1
		./fntest.sh >fntest.log 2>&1 &
		popd >/dev/null 2>&1
	fi
}

stop() {
	stopwatchdog
	curl -s http://${HOST}:${PORT}/extAdminShutdown/
	if [ ${IAM} == "root" ]; then
		sleep 6
		rm -f /var/run/phonebook/phonebook.pid /var/lock/phonebook
	fi
}

status() {
	ST=$(curl -s http://${HOST}:${PORT}/status/ | wc -c)
	case "${ST}" in
	"2")
		exit 0
		;;
	"0")
		# phonebook is not responsive or not running.  Exit status as described in 
		# http://refspecs.linuxbase.org/LSB_3.1.0/LSB-Core-generic/LSB-Core-generic/iniscrptact.html
		if [ ${IAM} == "root" -a -f /var/run/phonebook/phonebook.pid ]; then
			exit 1
		fi
		if [ ${IAM} == "root" -a -f /var/lock/phonebook ]; then
			exit 2
		fi
		exit 3
		;;
	esac
}

reload() {
	ST=$(curl -s http://${HOST}:${PORT}/status/)
}

while getopts ":p:qih:N:Tb" o; do
    case "${o}" in
       b)
            WATCHDOGOPTS="-b"
	    	# echo "WATCHDOGOPTS set to: ${WATCHDOGOPTS}"
            ;;
       h)
            HOST=${OPTARG}
            echo "HOST set to: ${HOST}"
            ;;
        N)
            DBNAME=${OPTARG}
            # echo "DBNAME set to: ${DBNAME}"
            ;;
        p)
            PORT=${OPTARG}
	    	# echo "PORT set to: ${PORT}"
            ;;
        q)
            QA=1
	    	# echo "PORT set to: ${PORT}"
            ;;
        T)
            STARTPBONLY=1
	    	# echo "STARTPBONLY set to: ${STARTPBONLY}"
            ;;
        *)
            usage
            ;;
    esac
done
shift $((OPTIND-1))

if [ ${IAM} == "root" ]; then
	cd ~ec2-user/apps/phonebook
fi

for arg do
	# echo '--> '"\`$arg'"
	cmd=$(echo ${arg}|tr "[:upper:]" "[:lower:]")
    case "$cmd" in
    "images")
		updateImages
		echo "Images updated"
		;;
	"start")
		start
		echo "OK"
		exit 0
		;;
	"stop")
		stop
		echo "OK"
		exit 0
		;;
	"ready")
		ST=$(curl -s http://${HOST}:${PORT}/status/)
		echo "${ST}"
		exit 0
		;;
	"status")
		status
		;;
	"restart")
		stop; sleep 10; start
		echo "OK"
		exit 0
		;;
	"reload")
		reload
		exit 0
		;;
	"condrestart")
		if [ -f /var/lock/phonebook ] ; then
			stop
			# avoid race
			sleep 3
			start
		fi
		;;
	*)
		echo "Unrecognized command: $arg"
		usage
		exit 1
		;;
    esac
done
