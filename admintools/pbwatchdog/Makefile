pbwatchdog: *.sh

clean:
	rm -f pbwatchdog
	@echo "*** Clean completed in pbwatchdog ***"

install: pbwatchdog
	@echo "*** Packaging completed in pbwatchdog ***"

package: pbwatchdog
	cp pbwatchdog ../../tmp/phonebook/
	cp *.1 ../../tmp/phonebook/man/man1/
	@echo "*** Packaging completed in pbwatchdog ***"

packageqa: pbwatchdog
	cp pbwatchdog ../../tmp/phonebookqa/
	cp *.1 ../../tmp/phonebookqa/man/man1/
	@echo "*** Packaging completed in pbwatchdog ***"

test:
	@echo "*** Testing completed in pbwatchdog ***"

manpage:
	nroff -man pbwatchdog.1
	#groff -man -Tascii pbwatchdog.1
