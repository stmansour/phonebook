all: clean phonebook

clean:
	cd dbtools;make clean
	cd admintools;make clean
	rm -rf phonebook tmp Phonebook.log pbimages.tar* x.sh*
	go clean

phonebook: *.go
	cd admintools;make
	go vet
	golint
	go build

dbmake:
	#cd ../dir/obfuscate;./obfuscate
	mysqldump accord > testdb.sql

dbrestore:
	restoreMySQLdb.sh accord testdb.sql

package: phonebook
	#cd admintools;make
	rm -rf tmp
	mkdir -p tmp/phonebook
	cp phonebook activate.sh testdb.sql *.css *.html pbbkup.sh pbrestore.sh tmp/phonebook/
	cd admintools;make package
	cd dbtools;make package
	cd tmp;tar cvf phonebook.tar phonebook; gzip phonebook.tar

publish: package
	cd tmp;/usr/local/accord/bin/deployfile.sh phonebook.tar.gz jenkins-snapshot/phonebook/latest

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
	
