DIRS=jobtitles departments roleinit deductions prod

db:
	for dir in $(DIRS); do make -C $$dir;done
	echo "./dbtools  done."

clean:
	for dir in $(DIRS); do make -C $$dir clean;done
	echo "./dbtools: clean done."

package:
	for dir in $(DIRS); do make -C $$dir package;done
	echo "./dbtools: package done."

packageqa:
	for dir in $(DIRS); do make -C $$dir packageqa;done
	cp apply.sh ../tmp/phonebookqa/dbtools/
	echo "./dbtools: packageqa done."

test:
	for dir in $(DIRS); do make -C $$dir package;done
	echo "./dbtools: test done."
secure:
	for dir in $(DIRS); do make -C $${dir} secure;done
	@rm -f config.json confdev.json confprod.json
