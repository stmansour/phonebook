#!/bin/bash
# A watchdog process for phonebook.  Checks the status command
# and validates that the process is running.  It also does a 
# database backup every night at 12:15am PST.  It does a full
# backup every Saturday night at 12:15am PST.

PHONEBOOKSTOPPED=0                 # switcher for mail sending
CHECKINGPERIOD=10                  # in sec
MYEMAIL="sman@stevemansour.com"    # where to send info
LOGFILE="pbwatchdog.log"           # where to log messages
DBBACKUP=0                         # backups are OFF by default
BACKUPCOMPLETED=0

TRIGGERHR="08"
TRIGGERMN="15"
RESETMN=$((TRIGGERMN+1))

logtimestamp() {
    date >> ${LOGFILE}
}

usage() {
    cat <<ZZEOF
Phonebook pbwatchdog
Usage:   pbwatchdog [OPTIONS]

OPTIONS:
-b      turn on database backups

Examples:
Command to start the watchdog without database backups (for testing):
    bash$  ./pbwatchdog

Command to start the watchdog with database backups turned on:
    bash$  ./pbwatchdog -b

ZZEOF
    exit 0
}

while getopts ":b" o; do
    case "${o}" in
        b)
            DBBACKUP=1
            echo "DBBACKUP set to: ${DBBACKUP}"
            ;;
        *)
            usage
            ;;
    esac
done
shift $((OPTIND-1))

TSTART=$(date)
   cat >${LOGFILE} <<ZEOF
Phonebook PBWATCHDOG Initiated at ${TSTART}
Logfile = ${LOGFILE}
Database backups = ${DBBACKUP}
    Trigger at ${TRIGGERHR}:${TRIGGERMN} UTC
ZEOF


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
            echo -n "calling ./activate.sh -T start " >> ${LOGFILE}
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
    if [ ${DBBACKUP} -gt 0 ]; then
        HR=$(date +%H)
        MN=$(date +%M)
        # echo "HR = ${HR}, MN = ${MN}"
        if [ ${HR} = ${TRIGGERHR} -a ${MN} = ${TRIGGERMN} -a ${BACKUPCOMPLETED} -eq 0 ]; then
        # echo "TIME TO DO A BACKUP"
            ./pbbkup >dailyDBbackup.log &
            DOW=$(date +%u)
            if [ ${DOW} = "2" ]; then
                ./pbbkup -f >weeklyDBfullbackup.log &
            fi
            BACKUPCOMPLETED=1
        fi

        if [ ${HR} = ${TRIGGERHR} -a ${MN} = ${RESETMN} -a ${BACKUPCOMPLETED} -eq 1 ]; then
            BACKUPCOMPLETED=0
        fi
    fi

    #---------------------------------------------------------------------------
    # Wait for a bit, then do it all again...
    #---------------------------------------------------------------------------
    touch ${LOGFILE}
    sleep ${CHECKINGPERIOD}
done
