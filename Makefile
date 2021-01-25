DIRS=lib db dbtools admintools test
DIST=tmp

phonebook: *.go config.json
	for dir in $(DIRS); do make -C $$dir;done
	go vet
	golint
	go build

.PHONY:  test

clean:
	for dir in $(DIRS); do make -C $$dir clean;done
	rm -rf phonebook pbbkup pbrestore pbwatchdog tmp Phonebook.log pbimages.tar* *.out *.log x.sh* ver.go conf*.json
	go clean

config.json:
	@/usr/local/accord/bin/getfile.sh accord/db/confdev.json
	@cp confdev.json config.json

test: package
	cd test;make test
	@echo "*** TESTING COMPLETE, ALL TESTS PASSED ***"

try:	clean phonebook package
	mysql --no-defaults accord < accord.sql


all: clean phonebook test

build: clean phonebook

dbmake:
	#cd ../dir/obfuscate;./obfuscate
	mysqldump accord > testdb.sql

dbrestore:
	restoreMySQLdb.sh accord testdb.sql

package: phonebook
	rm -rf tmp
	mkdir -p tmp/phonebook
	mkdir -p tmp/phonebook/man/man1/
	cp *.1 tmp/phonebook/man/man1/
	cp config.json tmp/phonebook/
	cp phonebook sbsetup.sh activate.sh update.sh updatePhonebook.sh testdb.sql *.css *.html  tmp/phonebook/
	for dir in $(DIRS); do make -C $$dir package;done

accorddb:
	mysql --no-defaults accord < accord.sql

packageqa: phonebook
	rm -rf tmp
	mkdir -p tmp/phonebookqa/man/man1/
	cp *.1 tmp/phonebookqa/man/man1/
	cp phonebook sbsetup.sh activate.sh updatePhonebook.sh testdb.sql *.css *.html  tmp/phonebookqa/
	cd admintools;make packageqa
	cd dbtools;make packageqa
	cd test;make packageqa
	cd tmp;tar cvf phonebookqa.tar phonebookqa; gzip phonebookqa.tar

install: package
	if [ ! -d /usr/local/accord ]; then mkdir -p /usr/local/accord; fi
	tar -C /usr/local/accord -xzvf ./tmp/phonebook.tar.gz


publish: package
	if [ -f tmp/phonebook/config.json ]; then mv tmp/phonebook/config.json tmp/config.json; fi
	cd tmp;tar cvf phonebook.tar phonebook; gzip phonebook.tar
	cd tmp;/usr/local/accord/bin/deployfile.sh phonebook.tar.gz jenkins-snapshot/phonebook/latest
	if [ -f tmp/config.json ]; then mv tmp/config.json tmp/phonebook/config.json; fi

publishqa: packageqa
	cd tmp;/usr/local/accord/bin/deployfile.sh phonebookqa.tar.gz jenkins-snapshot/phonebook/latest

publishprod: package
	cd tmp;/usr/local/accord/bin/deployfile.sh phonebook.tar.gz accord/phonebook/1.0

# Handling the images must be done from the development workstation.
# make pkimages pubimages
# This will update artifactory with the images needed by phonebook
# activate.sh will take care of downloading the images
pkgimages:
	tar cvf pbimages.tar images; gzip pbimages.tar

pubimages: pkgimages
	/usr/local/accord/bin/deployfile.sh pbimages.tar.gz jenkins-snapshot/phonebook/latest

cert:
	openssl genrsa -des3 -out server.key 1024
	openssl req -new -key server.key -out server.csr
	cp server.key server.key.org
	openssl rsa -in server.key.org -out server.key
	openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt

secure:
	for dir in $(DIRS); do make -C $${dir} secure;done
	@rm -f config.json confdev.json confprod.json

tarzip:
	cd ${DIST};if [ -f ./phonebook/config.json ]; then mv ./phonebook/config.json .; fi
	cd ${DIST};rm -f phonebook.tar*;tar czf phonebook.tar.gz phonebook
	cd ${DIST};if [ -f ./config.json ]; then mv ./config.json ./phonebook/config.json; fi

snapshot: package tarzip
	cd ${DIST}; /usr/local/accord/bin/snapshot.sh phonebook.tar.gz

release:
	/usr/local/accord/bin/release.sh phonebook
