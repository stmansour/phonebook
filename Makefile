all: clean phonebook

clean:
	rm -f phonebook

phonebook: *.go
	go vet
	golint
	go build
