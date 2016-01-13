pbbkup: *.sh

clean:
	rm -f pbbkup
	@echo "*** Clean completed in pbbkup ***"

install: pbbkup
	@echo "*** Packaging completed in pbbkup ***"

package: pbbkup
	cp pbbkup ../../tmp/phonebook/
	cp *.1 ../../tmp/phonebook/man/man1/
	@echo "*** Packaging completed in pbbkup ***"

packageqa: pbbkup
	cp pbbkup ../../tmp/phonebookqa/
	cp *.1 ../../tmp/phonebookqa/man/man1/
	@echo "*** Packaging completed in pbbkup ***"

test:
	@echo "*** Testing completed in pbbkup ***"

manpage:
	nroff -man pbbkup.1
	#groff -man -Tascii pbbkup.1
