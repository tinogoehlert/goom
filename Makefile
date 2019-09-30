.PHONY: all clean run

export GO111MODULE=auto

TARGETS = DOOM1.gwa go.mod

all: $(TARGETS)

clean:
	rm $(TARGETS)

go.mod:
	go mod init github.com/tinogoehlert/goom
	go get -u ./...

DOOM1.gwa: DOOM1.wad
	glbsp -v5 DOOM1.wad

run: all
	cd cmd/doom && go run main.go -iwad ../../DOOM1
