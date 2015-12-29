pbadduser: *.go
	go build

clean:
	rm -f pbadduser

install: pbadduser
	cp pbadduser /usr/local/accord/bin

package: pbadduser
	cp pbadduser ../../tmp/phonebook/
	cp *.1 ../../tmp/phonebook/man/man1/
	@echo "*** Packaging completed in pbadduser ***"

test:
	@echo "*** Testing completed in pbadduser ***"

manpage:
	nroff -man pbadduser.1
	groff -man -Tascii pbadduser.1
