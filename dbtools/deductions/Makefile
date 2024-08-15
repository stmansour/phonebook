deductions: *.go
	go vet
	@if [[ -f "${GOPATH}/bin/golint" ]]; then golint; else staticcheck; fi
	go build

clean:
	rm -f deductions

package:
	cp deductions ../../tmp/phonebook/

packageqa:
	mkdir -p ../../tmp/phonebookqa/dbtools/deductions
	cp deductions ../../tmp/phonebookqa/dbtools/deductions/
secure:
	@rm -f config.json confdev.json confprod.json
