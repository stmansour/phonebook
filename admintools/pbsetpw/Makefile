pbsetpw: *.go
	go build

clean:
	rm -f pbsetpw

install: pbsetpw
	cp pbsetpw /usr/local/accord/bin

package: pbsetpw
	cp pbsetpw ../../tmp/phonebook/
	@echo "*** Packaging completed in pbsetpw ***"

test:
	@echo "*** Testing completed in pbsetpw ***"
