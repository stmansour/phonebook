usersim: *.go config.json
	go vet
	if [[ -f "${GOPATH}/bin/golint" ]]; then golint; else staticcheck; fi
	go build

test: usersim
	./functest1.sh

clean:
	rm -f usersim conf*.json usersim.log usersim.out errors.txt log

config.json:
	@/usr/local/accord/bin/getfile.sh accord/db/confdev.json
	@cp confdev.json config.json

packageqa:
	mkdir -p ../../tmp/phonebookqa/test/usersim
	cp usersim *.txt *.sh ../../tmp/phonebookqa/test/usersim/
	cp *.1 ../../tmp/phonebookqa/man/man1/
	@echo "*** Packaging completed in usersim ***"
secure:
	@rm -f config.json confdev.json confprod.json
