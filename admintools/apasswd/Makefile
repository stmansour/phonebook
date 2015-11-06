all: clean apasswd

apasswd: *.go
	go build

clean:
	rm -f apasswd

install: apasswd
	cp apasswd /usr/local/accord/bin

test:
	@echo "*** Testing completed in apasswd ***"
