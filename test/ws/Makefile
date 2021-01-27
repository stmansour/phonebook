all:
	@echo "*** Completed in svc ***"

clean:
	rm -f rentroll.log log llog *.g ./gold/*.g err.txt [a-z][0-9] [a-z] [a-z][a-z1-9] qq? fail request w.[a-z] pbconsole.* *.json loginrequest serverreply
	@echo "*** CLEAN completed in svc ***"

test: config.json
	@touch fail
	./functest.sh
	@echo "*** TEST completed in svc ***"
	@rm -f fail

config.json:
	/usr/local/accord/bin/getfile.sh accord/db/confdev.json
	cp confdev.json config.json

package:
	@echo "*** PACKAGE completed in svc ***"

get:
	@touch fail
	../../tools/bashtools/getdb.sh
	@rm -f fail
secure:
	@rm -f config.json confdev.json confprod.json
