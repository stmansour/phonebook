TOP=../..
DIST=tmp
db:
	./apply.sh
	echo "./schema: build done"

clean:
	rm -f schema.sql
	echo "./schema: clean done"

package:
	echo "./schema: package done"

packageqa:
	mkdir -p "${TOP}"/"${DIST}"/phonebookqa/dbtools/schema
	cp apply.sh tables.sql "${TOP}"/"${DIST}"/phonebookqa/dbtools/schema
	echo "./schema: packageqa done"secure:
	@rm -f config.json confdev.json confprod.json

secure:
	@rm -f config.json confdev.json confprod.json
