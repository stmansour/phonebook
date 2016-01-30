pbsetpw: *.go
	go vet
	golint
	go build

clean:
	rm -f pbsetpw

install: pbsetpw
	cp pbsetpw /usr/local/accord/bin

package: pbsetpw
	cp pbsetpw ../../tmp/phonebook/
	cp *.1 ../../tmp/phonebook/man/man1/
	@echo "*** Packaging completed in pbsetpw ***"

packageqa: pbsetpw
	cp pbsetpw ../../tmp/phonebookqa/
	cp *.1 ../../tmp/phonebookqa/man/man1/
	@echo "*** Packaging completed in pbsetpw ***"

test:
	@echo "*** Testing completed in pbsetpw ***"

manpage:
	#nroff -man pbsetpw.1
	groff -man -Tascii pbsetpw.1
