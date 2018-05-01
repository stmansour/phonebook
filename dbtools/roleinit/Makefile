roleinit: *.go config.json
	go build

clean:
	rm -f roleinit conf*.json

config.json:
	@/usr/local/accord/bin/getfile.sh accord/db/confdev.json
	@cp confdev.json config.json

package:
	cp roleinit ../../tmp/phonebook/

packageqa:
	mkdir -p ../../tmp/phonebookqa/dbtools/roleinit
	cp roleinit apply.sh ../../tmp/phonebookqa/dbtools/roleinit
secure:
	@rm -f config.json confdev.json confprod.json
