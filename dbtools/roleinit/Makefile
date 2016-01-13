roleinit: *.go
	go build

clean:
	rm -f roleinit

package:
	cp roleinit ../../tmp/phonebook/

packageqa:
	mkdir -p ../../tmp/phonebookqa/dbtools/roleinit
	cp roleinit apply.sh ../../tmp/phonebookqa/dbtools/roleinit
