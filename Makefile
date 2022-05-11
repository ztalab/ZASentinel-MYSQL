PROG=bin/za-mysql

SRCS=cmd/proxy/main.go

CFLAGS = -ldflags "-s -w "

INSTALL_PREFIX=/usr/local

build:
	if [ ! -d "./bin/" ]; then \
		mkdir bin; \
	fi
	go build $(CFLAGS) -o $(PROG) $(SRCS)

install:
	cp $(PROG) $(INSTALL_PREFIX)/bin

race:
	if [ ! -d "./bin/" ]; then \
    	mkdir bin; \
    fi
	go build $(CFLAGS) -race -o $(PROG) $(SRCS)

clean:
	rm -rf ./bin

run:
	go run --race cmd/main.go -c config/config.yaml