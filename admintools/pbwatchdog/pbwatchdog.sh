#!/bin/bash
# A watchdog process for phonebook.  Checks the status command
# and validates that the process is running

PHONEBOOKSTOPPED=0                 # switcher for mail sending
CHECKINGPERIOD=10                  # in sec
MYEMAIL="sman@stevemansour.com"    # where to send info
LOGFILE="pbwatchdog.log"

logtimestamp() {
    date >> ${LOGFILE}
}

while [ 1=1 ];
do
    if [ ! "$(pidof phonebook)" ]
    then
        if [ "${PHONEBOOKSTOPPED}" = "0" ]; then
            # echo "WARNING! Phonebook crashed!" | mail -s "Phonebook crash." "$MYEMAIL"
            echo -n "ALERT: no process id for phonebook at " >> ${LOGFILE}
            logtimestamp
            PHONEBOOKSTOPPED=1
        fi
            echo -n "calling ./activate.sh start " >> ${LOGFILE}
            logtimestamp
            ./activate.sh start >> ${LOGFILE}
    else
        if [ "${PHONEBOOKSTOPPED}" = "1" ]; then
            # echo "Phonebook was successfully restarted." | mail -s "Phonebook restarted." "$MYEMAIL"
            echo -n "ALERT: Phonebook was successfully restarted at " >>${LOGFILE}
            logtimestamp
            PHONEBOOKSTOPPED=0
        fi
    fi
    sleep ${CHECKINGPERIOD}
done
