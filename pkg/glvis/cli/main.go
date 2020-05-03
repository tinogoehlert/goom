package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/tinogoehlert/goom/pkg/glvis"

	"github.com/ttacon/chalk"
)

func fatalf(format string, a ...interface{}) {
	fmt.Println(fmt.Sprintf(chalk.Red.Color(format), a...))
	os.Exit(1)
}

func main() {
	wadFile := flag.String("wad", "", "iwad or pwad file")
	gwaFile := flag.String("gwa", "", "glbsp gwa")
	flag.Parse()

	if *wadFile == "" {
		fatalf("no wad file given")
	}

	wadBuff, err := ioutil.ReadFile(*wadFile)
	if err != nil {
		fatalf("could not load WAD file: %s", err.Error())
	}

	if *gwaFile == "" {
		*gwaFile = strings.Replace(*wadFile, path.Ext(*wadFile), ".gwa", 1)
	}

	gwaBuff, err := ioutil.ReadFile(*gwaFile)
	if err != nil {
		fatalf("could not load GWA file %s: %s", *gwaFile, err.Error())
	}
	fmt.Println(*gwaFile)
	gv := glvis.NewGLVis()
	defer gv.Free()
	out := gv.BuildVis(wadBuff, gwaBuff)

	info, _ := os.Stat(*gwaFile)

	ioutil.WriteFile(*gwaFile, out, info.Mode())
}
