pbadduser: *.go
	go build

clean:
	rm -f pbadduser

install: pbadduser
	cp pbadduser /usr/local/accord/bin

package: pbadduser
	cp pbadduser ../../tmp/phonebook/
	@echo "*** Packaging completed in pbadduser ***"

test:
	@echo "*** Testing completed in pbadduser ***"
