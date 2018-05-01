DIRS = pbadduser pbbkup pbrestore pbsetpw pbsetrole pbsetusername pbupdateallpw pbwatchdog

tools:
	for dir in $(DIRS); do make -C $$dir;done

clean:
	for dir in $(DIRS); do make -C $$dir clean;done


package:
	mkdir -p ../tmp/phonebook/man/man1
	cp installman.sh ../tmp/phonebook/
	for dir in $(DIRS); do make -C $$dir package;done

packageqa:
	mkdir -p ../tmp/phonebookqa/man/man1
	cp installman.sh ../tmp/phonebookqa/
	for dir in $(DIRS); do make -C $$dir packageqa;done

test:
	for dir in $(DIRS); do make -C $$dir test;done
secure:
	for dir in $(DIRS); do make -C $${dir} secure;done
	@rm -f config.json confdev.json confprod.json
