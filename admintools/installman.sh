#!/bin/bash
TARGET=/usr/local/share/man/man1
DIR=$(pwd)

# for i in "${pages1[@]}"
for i in "${DIR}/man/man1/*.1"
do
	rm -f ${TARGET}/${i}.1
	ln -s ${i} ${TARGET}/${i}.1 
done