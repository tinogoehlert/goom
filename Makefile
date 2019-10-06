.PHONY: all clean test test-run run

export GO111MODULE=auto
export DOOM_TEST=$(CURDIR)/.test

FILES = DOOM1.gwa go.mod
TARGETS = $(FILES) $(DOOM_TEST)

all: $(TARGETS) test

clean:
	rm -f $(FILES)
	rm -f $(DOOM_TEST)/*.mid
	rm -f $(DOOM_TEST)/*.mus

test: $(TARGETS)
	go test -v . ./audio

go.mod:
	go mod init github.com/tinogoehlert/goom
	go get -u ./...

DOOM1.gwa: DOOM1.wad
	glbsp -v5 DOOM1.wad

$(DOOM_TEST):
	mkdir -p $@

test-run: TEST=-test
test-run: run

run: $(TARGETS)
	cd cmd/doom && go run main.go -iwad ../../DOOM1 $(TEST)
