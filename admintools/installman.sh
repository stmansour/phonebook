#!/bin/bash
TARGET=/usr/local/share/man/man1
DIR=$(pwd)
pages1=("pbadduser"
	"pbbkup"
	"pbrestore"
	"pbsetpw"
	"pbsetrole"
	"pbsetusername"
	"pbsetallpw"
	"pbwatchdog"
	)


for i in "${pages1[@]}"
do
	rm -f ${TARGET}/${i}.1
	ln -s ${DIR}/man/man1/${i}.1 ${TARGET}/${i}.1 
done