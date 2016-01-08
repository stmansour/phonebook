usersim: *.go
	go vet
	golint
	go build

test: usersim
	./functest1.sh

clean:
	rm -f usersim

package:
	mkdir ../../tmp/phonebook/usersim
	cp usersim *.txt ../../tmp/phonebook/usersim/
	cp *.1 ../../tmp/phonebook/man/man1/
	@echo "*** Packaging completed in usersim ***"
