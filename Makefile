BUILD = winbox-export-parser
VERSION	?= 0.0

all: clean format build

clean:
	rm -f $(BUILD)

format:
	go fmt

build:
	CGO_ENABLED=0 go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD) *.go
