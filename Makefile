.PHONY: all run

all: DOOM1.gwa

DOOM1.gwa: DOOM1.wad
	glbsp -v5 DOOM1.wad

run:
	cd cmd/doom && go run main.go -wad ../../DOOM1