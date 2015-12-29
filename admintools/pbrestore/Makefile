pbrestore: *.sh

clean:
	rm -f pbrestore
	@echo "*** Clean completed in pbrestore ***"

install: pbrestore
	@echo "*** Packaging completed in pbrestore ***"

package: pbrestore
	cp pbrestore ../../tmp/phonebook/
	cp *.1 ../../tmp/phonebook/man/man1/
	@echo "*** Packaging completed in pbrestore ***"

test:
	@echo "*** Testing completed in pbrestore ***"

manpage:
	nroff -man pbrestore.1
	#groff -man -Tascii pbrestore.1
