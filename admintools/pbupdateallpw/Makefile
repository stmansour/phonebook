pbupdateallpw: *.go
	go build

clean:
	rm -f pbupdateallpw

install: pbupdateallpw
	cp pbupdateallpw /usr/local/accord/bin

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
