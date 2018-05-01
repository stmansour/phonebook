pbsetpw: *.go config.json
	go vet
	golint
	go build

clean:
	rm -f pbsetpw conf*.json

config.json:
	@/usr/local/accord/bin/getfile.sh accord/db/confdev.json
	@cp confdev.json config.json

install: pbsetpw
	cp pbsetpw /usr/local/accord/bin

installman:
	cp pbsetpw.1 /usr/local/share/man/man1/

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
secure:
	@rm -f config.json confdev.json confprod.json
