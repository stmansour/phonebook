usersim: *.go
	go vet
	golint
	go build

test: usersim
	./functest1.sh

clean:
	rm -f usersim
