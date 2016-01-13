jobtitles: *.go
	go vet
	golint
	go build

clean:
	rm -f jobtitles

package:
	cp jobtitles ../../tmp/phonebook/

packageqa:
	mkdir -p ../../tmp/phonebookqa/dbtools/jobtitles
	cp jobtitles jobtitles.csv ../../tmp/phonebookqa/dbtools/jobtitles
