.PHONY: all clean test test-run run

export GO111MODULE=auto

TARGETS = DOOM1.gwa go.mod

all: $(TARGETS) test

clean:
	rm $(TARGETS)

test: $(TARGETS)
	go test -v github.com/tinogoehlert/goom

go.mod:
	go mod init github.com/tinogoehlert/goom
	go get -u ./...

DOOM1.gwa: DOOM1.wad
	glbsp -v5 DOOM1.wad

test-run: TEST=-test
test-run: run

run: $(TARGETS)
	cd cmd/doom && go run main.go -iwad ../../DOOM1 $(TEST)
