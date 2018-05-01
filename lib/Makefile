lib: ver.go *.go
	go vet
	golint
	go build
	go install

ver.go:
	./mkver.sh

clean:
	go clean
secure:
	@rm -f config.json confdev.json confprod.json
