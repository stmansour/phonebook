TOP=..
BINDIR=${TOP}/tmp/phonebook
THISDIR="ws"

ws: *.go
	@touch fail
	go vet
	if [[ -f "${GOPATH}/bin/golint" ]]; then golint; else staticcheck; fi
	go build
	go test
	go install
	@rm -f fail
	@echo "*** ${THISDIR}: completed build ***"

webdoc:
	rm -rf doc
	${WEBDOC} 
	@cp ${WEBDOCDIR}/doc.css doc/
	@cp ${WEBDOCDIR}/tester.html doc/
	rm -rf ${BINDIR}/doc
	cp -r doc ${BINDIR}/

test:
	@touch fail
	go test
	@rm -f fail
	@echo "*** ${THISDIR}: completed testing - ALL TESTS PASS ***"

clean:
	rm -f ver.go
	go clean
	@rm -f fail
	@echo "*** ${THISDIR}: completed clean ***"

package:
	#@cp -r doc ${BINDIR}/
	@echo "*** ${THISDIR}: completed package ***"
secure:
	@rm -f config.json confdev.json confprod.json
