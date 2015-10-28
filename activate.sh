#!/bin/bash
# activation script for phonebook

HOST=localhost
PORT=8250

usage() {
    cat << ZZEOF
Phonebook activation script.
Usage:   activate.sh [OPTIONS] CMD

OPTIONS:
-p port      (default is 8250)
-h hostname  (default is localhost)

CMD is one of: start | stop | ready | test | teststatus

cmd is case insensitive

Examples:
Command to start phonebook:
	bash$  activate.sh START 

Command to stop phonebook:
	bash$  activate.sh Stop

Command to see if phonebook is ready for commands... the response
will be "OK" if it is ready, or something else if there are problems:

    bash$  activate.sh ready
    OK
ZZEOF
	exit 0
}


while getopts ":p:ih:" o; do
    case "${o}" in
        h)
            HOST=${OPTARG}
            echo "HOST set to: ${HOST}"
            ;;
        p)
            PORT=${OPTARG}
	    	echo "PORT set to: ${PORT}"
            ;;
        *)
            usage
            ;;
    esac
done
shift $((OPTIND-1))

for arg do
	# echo '--> '"\`$arg'"
	cmd=$(echo ${arg}|tr "[:upper:]" "[:lower:]")
    case "$cmd" in
	"start")
		#===============================================
		# START
		# Add the command to start your application...
		#===============================================
		./phonebook >phonebook.log 2>&1 &
		echo "OK"
		exit 0
		;;
	"stop")
		#===============================================
		# STOP
		# Add the command to terminate your application...
		#===============================================
		curl -s http://${HOST}:${PORT}/shutdown/
		echo "OK"
		exit 0
		;;
	"ready")
		#===============================================
		# READY
		# Is your application ready to receive commands?
		#===============================================
		ST=$(curl -s http://${HOST}:${PORT}/status/)
			echo "${ST}"
			exit 0
		;;
	*)
		echo "Unrecognized command: $arg"
		usage
		;;
    esac
done
