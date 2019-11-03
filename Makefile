.PHONY: all clean tidy test test-run run

export GO111MODULE=auto

FILES = DOOM1.gwa go.mod
TARGETS = $(FILES)

all: $(TARGETS) tidy test

clean:
	rm -f $(FILES)
	rm -rf .test

tidy:
	go mod tidy

test: $(TARGETS)
	go test -v ./...

go.mod:
	go mod init github.com/tinogoehlert/goom
	go get -u ./...

DOOM1.gwa: DOOM1.WAD
	glbsp -v5 DOOM1.WAD

test-run: TEST=-test
test-run: run

run: $(TARGETS)
	cd cmd/doom && go run main.go -iwad $(CURDIR)/DOOM1 $(TEST)
