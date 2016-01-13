db:
	echo "./schema: build done"

clean:
	rm -f schema.sql
	echo "./schema: clean done"

package:
	echo "./schema: package done"

packageqa:
	mkdir -p ../../tmp/phonebookqa/dbtools/schema
	cp apply.sh tables.sql ../../tmp/phonebookqa/dbtools/schema
	echo "./schema: packageqa done"