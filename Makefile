.PHONY: build install clean

BINARY_NAME=kodia

# Build the binary
build:
	go build -o bin/$(BINARY_NAME) ./kodia

# Install the binary to GOPATH/bin
install:
	go install ./kodia

# Remove build artifacts
clean:
	rm -rf bin/
