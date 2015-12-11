pbsetusername: *.go
	go build

clean:
	rm -f pbsetusername

install: pbsetusername
	cp pbsetusername /usr/local/accord/bin

package: pbsetusername
	cp pbsetusername ../../tmp/phonebook/
	@echo "*** Packaging completed in pbsetusername ***"

test:
	@echo "*** Testing completed in pbsetusername ***"
