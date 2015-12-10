#!/bin/bash
# A watchdog process for phonebook.  Checks the status command
# and validates that the process is running

PHONEBOOKSTOPPED=0                 # switcher for mail sending
CHECKINGPERIOD=20                  # in sec
MYEMAIL="sman@stevemansour.com"    # where to send info
LOGFILE="pbwatchdog.log"

while [ 1=1 ];
do

    if [ ! "$(pidof phonebook)" ]
    then
        if [ "${PHONEBOOKSTOPPED}" = "0" ]; then
            # echo "WARNING! Phonebook crashed!" | mail -s "Phonebook crash." "$MYEMAIL"
            echo "ALERT: no process id for phonebook" >> ${LOGFILE}
            PHONEBOOKSTOPPED=1
        fi
        ./activate.sh start
    else
        if [ "${PHONEBOOKSTOPPED}" = "1" ]; then
            # echo "Phonebook was successfully restarted." | mail -s "Phonebook restarted." "$MYEMAIL"
            echo "ALERT: Phonebook was successfully restarted." >>${LOGFILE}
            PHONEBOOKSTOPPED=0
        fi
    fi
    sleep ${CHECKINGPERIOD}

done
