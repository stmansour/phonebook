pbsetusername: *.go
	go build

clean:
	rm -f pbsetusername

install: pbsetusername
	cp pbsetusername /usr/local/accord/bin

package: pbsetusername
	cp pbsetusername ../../tmp/phonebook/
	cp *.1 ../../tmp/phonebook/man/man1/
	@echo "*** Packaging completed in pbsetusername ***"

test:
	@echo "*** Testing completed in pbsetusername ***"

manpage:
	#nroff -man pbsetusername.1
	groff -man -Tascii pbsetusername.1
