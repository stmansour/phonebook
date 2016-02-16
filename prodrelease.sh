#!/bin/bash
# Accord Phonebook release to production script
# This script updates existing production nodes with the latest code

declare arr=(
	"pb1"
	"pb2"
	)

for f in "${arr[@]}"
do
	echo "Updating host: $f"
    ssh $f "cd apps/phonebook; sudo ./updatePhonebook.sh"
done
