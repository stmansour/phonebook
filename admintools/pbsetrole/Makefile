pbsetrole: *.go
	go build

clean:
	rm -f pbsetrole

install: pbsetrole
	cp pbsetrole /usr/local/accord/bin

package: pbsetrole
	cp pbsetrole ../../tmp/phonebook/
	@echo "*** Packaging completed in pbsetrole ***"

test:
	@echo "*** Testing completed in pbsetrole ***"
