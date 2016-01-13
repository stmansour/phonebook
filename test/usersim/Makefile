usersim: *.go
	go vet
	golint
	go build

test: usersim
	./functest1.sh

clean:
	rm -f usersim

packageqa:
	mkdir -p ../../tmp/phonebookqa/test/usersim
	cp usersim *.txt *.sh ../../tmp/phonebookqa/test/usersim/
	cp *.1 ../../tmp/phonebookqa/man/man1/
	@echo "*** Packaging completed in usersim ***"
