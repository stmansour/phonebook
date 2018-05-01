departments: *.go
	go vet
	golint
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
