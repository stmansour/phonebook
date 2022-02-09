lib: ver.go *.go
	go vet
	if [[ -f "${GOPATH}/bin/golint" ]]; then golint; else staticcheck; fi
	go build
	go install

ver.go:
	./mkver.sh

clean:
	rm -f ver.go
	go clean
secure:
	@rm -f config.json confdev.json confprod.json
