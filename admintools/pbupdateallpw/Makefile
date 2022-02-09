pbupdateallpw: *.go config.json
	go vet
	if [[ -f "${GOPATH}/bin/golint" ]]; then golint; else staticcheck; fi
	go build

clean:
	rm -f pbupdateallpw conf*.json

config.json:
	@/usr/local/accord/bin/getfile.sh accord/db/confdev.json
	@cp confdev.json config.json

install: pbupdateallpw
	cp pbupdateallpw /usr/local/accord/bin

installman:
	cp pbupdateallpw.1 /usr/local/share/man/man1/

package: pbupdateallpw
	cp pbupdateallpw ../../tmp/phonebook/
	cp *.1 ../../tmp/phonebook/man/man1/
	@echo "*** Packaging completed in pbupdateallpw ***"

packageqa: pbupdateallpw
	cp pbupdateallpw ../../tmp/phonebookqa/
	cp *.1 ../../tmp/phonebookqa/man/man1/
	@echo "*** Packaging completed in pbupdateallpw ***"

test:
	@echo "*** Testing completed in pbupdateallpw ***"
secure:
	@rm -f config.json confdev.json confprod.json
