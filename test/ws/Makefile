all:
	@echo "*** Completed in svc ***"

clean:
	rm -f rentroll.log log llog *.g ./gold/*.g err.txt [a-z][0-9] [a-z] [a-z][a-z1-9] qq? fail request w.[a-z] restore.sql
	@echo "*** CLEAN completed in svc ***"

test:
	@touch fail
	./functest.sh
	@echo "*** TEST completed in svc ***"
	@rm -f fail

package:
	@echo "*** PACKAGE completed in svc ***"

get:
	@touch fail
	../../tools/bashtools/getdb.sh
	@rm -f fail
