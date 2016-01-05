#!/bin/bash
# A watchdog process for phonebook.  Checks the status command
# and validates that the process is running.  It also does a 
# database backup every night at 12:15am PST.  It does a full
# backup every Saturday night at 12:15am PST.

PHONEBOOKSTOPPED=0                 # switcher for mail sending
CHECKINGPERIOD=10                  # in sec
MYEMAIL="sman@stevemansour.com"    # where to send info
LOGFILE="pbwatchdog.log"
BACKUPCOMPLETED=0

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

    #---------------------------------------------------------------------------
    # Check for database backup needed
    # Do the backups at 00:15:00 PST every day
    # 00:15:00 PST (Silicon Valley CA) == 08:15:00 GMT == 02:15:00 CST (Austin)
    #---------------------------------------------------------------------------
    HR=$(date +%H)
    MN=$(date +%M)
    if [ ${HR} == "08" and ${MN} == "15" and ${BACKUPCOMPLETED} -eq 0 ]; then
        ./pbbkup >dailyDBbackup.log &
        DOW=$(date +%u)
        if [ ${dow} == "6" ]; then
            ./pbbkup -f >weeklyDBfullbackup.log &
        fi
        BACKUPCOMPLETED=1
    else
        BACKUPCOMPLETED=0
    fi

    #---------------------------------------------------------------------------
    # Wait for a bit, then do it all again...
    #---------------------------------------------------------------------------
    sleep ${CHECKINGPERIOD}
done
