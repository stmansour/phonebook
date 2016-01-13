pbsetrole: *.go
	go build

clean:
	rm -f pbsetrole

install: pbsetrole
	cp pbsetrole /usr/local/accord/bin

package: pbsetrole
	cp pbsetrole ../../tmp/phonebook/
	cp *.1 ../../tmp/phonebook/man/man1/
	@echo "*** Packaging completed in pbsetrole ***"

packageqa: pbsetrole
	cp pbsetrole ../../tmp/phonebookqa/
	cp *.1 ../../tmp/phonebookqa/man/man1/
	@echo "*** Packaging completed in pbsetrole ***"

test:
	@echo "*** Testing completed in pbsetrole ***"

manpage:
	#nroff -man pbsetrole.1
	groff -man -Tascii pbsetrole.1
