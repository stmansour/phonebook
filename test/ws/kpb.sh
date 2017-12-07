#!/bin/bash


ps -ef | grep activate.sh | awk '{print $2}' | while read -r line ; do
	echo "processing: ${line}"
	kill -9 ${line}
done

ps -ef | grep pbwatchdog | awk '{print $2}' | while read -r line ; do
	echo "processing: ${line}"
	kill -9 ${line}
done

killall phonebook
