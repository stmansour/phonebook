TOP=..
THISDIR=db

lib: *.go
	go vet
	golint
	go build
	go install

clean:
	go clean

package:
	@echo "package completed in ${THISDIR}"
secure:
	@rm -f config.json confdev.json confprod.json
