pbsetusername: *.go config.json
	go vet
	golint
	go build

clean:
	rm -f pbsetusername conf*.json

config.json:
	@/usr/local/accord/bin/getfile.sh accord/db/confdev.json
	@cp confdev.json config.json

install: pbsetusername
	cp pbsetusername /usr/local/accord/bin

package: pbsetusername
	cp pbsetusername ../../tmp/phonebook/
	cp *.1 ../../tmp/phonebook/man/man1/
	@echo "*** Packaging completed in pbsetusername ***"

packageqa: pbsetusername
	cp pbsetusername ../../tmp/phonebookqa/
	cp *.1 ../../tmp/phonebookqa/man/man1/
	@echo "*** Packaging completed in pbsetusername ***"

test:
	@echo "*** Testing completed in pbsetusername ***"

manpage:
	#nroff -man pbsetusername.1
	groff -man -Tascii pbsetusername.1
secure:
	@rm -f config.json confdev.json confprod.json
