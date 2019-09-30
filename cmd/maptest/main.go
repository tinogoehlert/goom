package main

import (
	"flag"
	"fmt"
	"image/color"
	"math/rand"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/tinogoehlert/goom"
	"github.com/ttacon/chalk"
	"golang.org/x/image/colornames"
)

var gameData *goom.GameData
var ssectName = "SSECTORS"

func run() {
	cfg := pixelgl.WindowConfig{
		Title:     "GOOM - maptest " + ssectName,
		Bounds:    pixel.R(0, 0, 800, 600),
		Resizable: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		fmt.Println(chalk.Red.Color(err.Error()))
		os.Exit(2)
	}
	imd := imdraw.New(nil)
	imd.Color = colornames.White
	e1m1 := gameData.Level("E1M1")
	camPos := pixel.ZV
	camZoom := float64(1)

	down := float32(6)
	rand.Seed(time.Now().UnixNano())

	for _, ssect := range e1m1.SubSectors(ssectName)[:] {
		var (
			r = rand.Float64()
			g = rand.Float64()
			b = rand.Float64()
		)

		fseg := ssect.Segments()[0]
		for _, seg := range ssect.Segments() {
			imd.Color = pixel.RGB(r, g, b)

			var (
				f = e1m1.Vert(fseg.StartVert())
				s = e1m1.Vert(seg.StartVert())
				e = e1m1.Vert(seg.EndVert())
			)

			imd.Push(
				pixel.V(float64(f.X()/down), float64(f.Y()/down)),
				pixel.V(float64(s.X()/down), float64(s.Y()/down)),
				pixel.V(float64(e.X()/down), float64(e.Y()/down)),
			)
			imd.Line(1)

			if seg.LineDef() != -1 {
				line := e1m1.LinesDefs[seg.LineDef()]
				var (
					ls = e1m1.Vert(uint32(line.Start))
					le = e1m1.Vert(uint32(line.End))
				)
				imd.Color = pixel.RGB(1, 1, 1).Mul(pixel.Alpha(0.7))
				imd.Push(
					pixel.V(float64(ls.X()/down), float64(ls.Y()/down)),
					pixel.V(float64(le.X()/down), float64(le.Y()/down)),
				)
				imd.Line(1)
			}
		}
	}

	camPos.X = float64(float32(e1m1.Things[0].X) / down)
	camPos.Y = float64(float32(e1m1.Things[0].Y) / down)

	for !win.Closed() {
		win.Clear(color.Black)

		if win.Pressed(pixelgl.KeyRight) {
			camPos.X += 0.2 * (camZoom + 0.2)
		}
		if win.Pressed(pixelgl.KeyLeft) {
			camPos.X -= 0.2 * (camZoom + 0.2)
		}

		if win.Pressed(pixelgl.KeyDown) {
			camPos.Y += 0.2
		}
		if win.Pressed(pixelgl.KeyUp) {
			camPos.Y -= 0.2
		}
		if win.Pressed(pixelgl.KeyW) {
			camZoom += 0.002
		}
		if win.Pressed(pixelgl.KeyS) {
			camZoom -= 0.002
		}
		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)
		imd.Draw(win)
		win.Update()
		time.Sleep(20)
	}
}

func main() {
	path := flag.String("wad", "DOOM1", "WAD file to load")
	ssects := flag.String("ssect", "SSECTORS", "use this secs")
	flag.Parse()

	ssectName = *ssects
	var err error
	gameData, err = goom.LoadGameData(*path + ".wad")
	if err != nil {
		fmt.Println(chalk.Red.Color(err.Error()))
		os.Exit(2)
	}
	pixelgl.Run(run)
}
