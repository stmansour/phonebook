pbadduser: *.go config.json
	go vet
	golint
	go build

config.json:
	@/usr/local/accord/bin/getfile.sh accord/db/confdev.json
	@cp confdev.json config.json

clean:
	rm -f pbadduser conf*.json

install: pbadduser
	cp pbadduser /usr/local/accord/bin

package: pbadduser
	cp pbadduser ../../tmp/phonebook/
	cp defaultProfileImage.png ../../tmp/phonebook/
	cp *.1 ../../tmp/phonebook/man/man1/
	@echo "*** Packaging completed in pbadduser ***"

packageqa: pbadduser
	cp pbadduser ../../tmp/phonebookqa/
	cp *.1 ../../tmp/phonebookqa/man/man1/
	@echo "*** Packaging completed in pbadduser ***"

test:
	@echo "*** Testing completed in pbadduser ***"

manpage:
	nroff -man pbadduser.1
	groff -man -Tascii pbadduser.1
secure:
	@rm -f config.json confdev.json confprod.json
