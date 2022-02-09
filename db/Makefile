DIRS=schema
TOP=..
THISDIR=db

lib: *.go
	for dir in $(DIRS); do make -C $$dir;done
	go vet
	# if [[ -f "${GOPATH}/bin/golint" ]]; then golint; else staticcheck; fi
	go build
	go install

clean:
	for dir in $(DIRS); do make -C $$dir clean;done
	go clean

package:
	for dir in $(DIRS); do make -C $$dir package;done
	@echo "package completed in ${THISDIR}"

secure:
	for dir in $(DIRS); do make -C $$dir secure;done
	@rm -f config.json confdev.json confprod.json
