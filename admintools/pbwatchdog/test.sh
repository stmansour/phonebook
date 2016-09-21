CHECKINGPERIOD=10                  # in sec
PICSYNCINTERVAL=60                 # in sec
PICSYNCCOUNT=0

while [ 1 == 1 ]; do
    ((PICSYNCCOUNT += CHECKINGPERIOD))
    if ((PICSYNCCOUNT >= PICSYNCINTERVAL)); then
	echo "Hit the interval"
	PICSYNCCOUNT=0
	exit 1
    fi
done
