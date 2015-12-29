#!/bin/bash
cd /usr/local/share/man/man1
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
	rm -f ${i}.1
	ln -s /usr/local/accord/phonebook/man/man1/${i}.1 
done