DIRS = pbadduser pbbkup pbrestore pbsetpw pbsetrole pbsetusername pbupdateallpw pbwatchdog

tools:
	for dir in $(DIRS); do make -C $$dir;done

clean:
	for dir in $(DIRS); do make -C $$dir clean;done


package:
	mkdir -p ../tmp/phonebook/man/man1
	for dir in $(DIRS); do make -C $$dir package;done

test:
	for dir in $(DIRS); do make -C $$dir test;done
