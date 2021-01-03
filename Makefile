.PHONY: all test clean

all: tor-tools

tor-tools:
	go build -o tor-tools ./cmd

test:
	@true

clean:
	rm -f ./tor-tools
