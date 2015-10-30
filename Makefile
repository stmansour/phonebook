all: clean phonebook

clean:
	rm -rf phonebook tmp

phonebook: *.go
	go vet
	golint
	go build

package: phonebook
	rm -rf tmp
	mkdir -p tmp/phonebook
	cp phonebook activate.sh testdb.sql *.css *.html tmp/phonebook/
	#cp -r images tmp/phonebook/
	cd tmp;tar cvf phonebook.tar phonebook; gzip phonebook.tar

publish: package
	cd tmp;/usr/local/accord/bin/deployfile.sh phonebook.tar.gz jenkins-snapshot/phonebook/latest

# Handling the images must be done from the development workstation.
# make pkimages pubimages
# This will update artifactory with the images needed by phonebook
# activate.sh will take care of downloading the images
pkgimages:
	tar cvf pbimages.tar images; gzip pbimages.tar

pubimages:
	/usr/local/accord/bin/deployfile.sh pbimages.tar.gz jenkins-snapshot/phonebook/latest