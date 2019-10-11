SRCS=*.go

EXECUTABLE=bin/jacksquat

LINUX=$(EXECUTABLE)-linux
LINUX_AMD64=$(LINUX)-amd64

all: linux-amd64

linux-amd64: $(LINUX_AMD64)

test:
	go test -race -coverprofile=coverage.txt -covermode=atomic jacksquat*.go
	go tool cover -html=coverage.txt -o coverage.html

$(LINUX_AMD64): $(SRCS)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $@ $(SRCS)

clean:
	rm -rf bin coverage.txt coverage.html
