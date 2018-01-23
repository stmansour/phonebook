lib: ver.go *.go
	go vet
	golint
	go build
	go install

ver.go:
	./mkver.sh

clean:
	go clean
