deductions: *.go
	go vet
	golint
	go build

clean:
	rm -f deductions

package:
	cp deductions ../../tmp/phonebook/

packageqa:
	mkdir -p ../../tmp/phonebookqa/dbtools/deductions
	cp deductions ../../tmp/phonebookqa/dbtools/deductions/
