all: clean auser

auser: *.go
	go build

clean:
	rm -f auser

install: auser
	cp auser /usr/local/accord/bin

test:
	@echo "*** Testing completed in auser ***"
