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

packageqa: pbrestore
	cp pbrestore ../../tmp/phonebookqa/
	cp *.1 ../../tmp/phonebookqa/man/man1/
	@echo "*** Packaging completed in pbrestore ***"

test:
	@echo "*** Testing completed in pbrestore ***"

manpage:
	nroff -man pbrestore.1
	#groff -man -Tascii pbrestore.1
secure:
	@rm -f config.json confdev.json confprod.json
