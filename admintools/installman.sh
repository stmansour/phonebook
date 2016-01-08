#!/bin/bash
TARGET=/usr/local/share/man/man1
DIR=$(pwd)

echo "DIR = ${DIR}"
pushd man/man1
for i in *.1
do
	rm -f ${TARGET}/${i}.1
	ln -s ${DIR}/man/man1/${i} ${TARGET}/${i}.1 
done
popd