departments: *.go
	go vet
	@if [[ -f "${GOPATH}/bin/golint" ]]; then golint; else staticcheck; fi
	go build

clean:
	rm -f departments

package:
	cp departments ../../tmp/phonebook/

packageqa:
	mkdir -p ../../tmp/phonebookqa/dbtools/departments
	cp departments depts.csv ../../tmp/phonebookqa/dbtools/departments/
secure:
	@rm -f config.json confdev.json confprod.json
