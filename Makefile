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
	cd tmp;deployfile.sh phonebook.tar.gz jenkins-snapshot/phonebook/latest
