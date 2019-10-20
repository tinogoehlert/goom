package main

import (
	"flag"
	"fmt"
	"image/color"
	"math"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/tinogoehlert/goom"
	"github.com/tinogoehlert/goom/geometry"
	"github.com/ttacon/chalk"
	"golang.org/x/image/colornames"
)

var gameData *goom.GameData
var ssectName = "SSECTORS"

func dist(x1, y1, x2, y2 float64) float64 {
	first := math.Pow(float64(x1-x2), 2)
	second := math.Pow(float64(y1-y2), 2)
	return math.Sqrt(first + second)
}

func Perp(v, b pixel.Vec, length float64) pixel.Vec {
	perp := pixel.V(v.X-b.X, v.Y-b.Y)

	angle := math.Atan2(v.Y, -v.X)

	perp.X += math.Cos(angle) * length
	perp.Y -= math.Sin(angle) * length

	return perp

}

func inter(a1, a2, b1, b2 pixel.Vec) (pixel.Vec, bool) {
	tmp := (b2.X-b1.X)*(a2.Y-a1.Y) - (b2.Y-b1.Y)*(a2.X-a1.X)

	if tmp == 0 {
		return pixel.V(0, 0), false
	}

	mu := ((a1.X-b1.X)*(a2.Y-a1.Y) - (a1.Y-b1.Y)*(a2.X-a1.X)) / tmp

	return pixel.V(
		b1.X+(b2.X-b1.X)*mu,
		b1.Y+(b2.Y-b1.Y)*mu,
	), true
}

func cross(a, b pixel.Vec) pixel.Vec {
	c := a.Sub(b)
	return pixel.V(-c.Y, c.X)
}

func cross2(a, b pixel.Vec) pixel.Vec {
	c := a.Sub(b)
	return pixel.V(c.Y, -c.X)
}

func norm(v pixel.Vec) pixel.Vec {
	var length = v.X * v.X
	length += v.Y * v.Y
	length = math.Sqrt(length)

	if length == 0.0 {
		return v
	}

	d := 1.0 / length
	v.X *= d
	v.Y *= d
	return v
}

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
	e1m1 := gameData.Level("E1M1")
	camPos := pixel.ZV
	camZoom := float64(1)

	down := float32(6)

	camPos.X = float64(float32(e1m1.Things[0].X) / down)
	camPos.Y = float64(float32(e1m1.Things[0].Y) / down)

	player := e1m1.Things[0]
	ppos := pixel.V(float64(player.X/down), float64(player.Y/down))
	nppos := ppos
	y, x := math.Sincos(float64(player.Angle) * math.Pi / 180)
	pdir := pixel.V(x, y)

	for !win.Closed() {
		win.Clear(color.Black)

		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)

		for _, line := range e1m1.LinesDefs[:] {
			imd.SetColorMask(colornames.White)

			var (
				ls  = e1m1.Vert(uint32(line.Start))
				le  = e1m1.Vert(uint32(line.End))
				vls = pixel.V(float64(ls.X()/down), float64(ls.Y()/down))
				vle = pixel.V(float64(le.X()/down), float64(le.Y()/down))
			)

			imd.Push(vls, vle)
			imd.Line(1)

		}

		col := false
		for _, line := range e1m1.LinesDefs[:] {
			imd.SetColorMask(colornames.Aqua)

			var (
				ls  = e1m1.Vert(uint32(line.Start)).DivScalar(down)
				le  = e1m1.Vert(uint32(line.End)).DivScalar(down)
				mid = geometry.V2((ls.X()+le.X())/2, (ls.Y()+le.Y())/2)
			)

			tangent := le.Sub(ls).Normalize()
			xt := tangent.X()
			yt := tangent.Y()

			ct := tangent.CrossVec2()

			vmid := pixel.V(float64(mid.X()), float64(mid.Y()))
			//vct := pixel.V(float64(ct.X()), float64(ct.Y()))
			vtan := pixel.V(float64(tangent.X()), float64(tangent.Y()))

			imd.Push(
				vmid,
				vmid.Add(vtan.Scaled(10)),
			)
			imd.Line(1)

			lineLen := ls.DistanceTo(le)
			imd.SetColorMask(colornames.Yellow)

			d := ls.Dot(ct)
			sd := ls.Dot(tangent)

			pd := float32(nppos.X)*ct.X() + float32(nppos.Y)*ct.Y() - d
			mul := float32(1.0)
			xNudge := float32(0.0)
			yNudge := float32(0.0)
			radius := float32(1)
			if pd >= -radius && pd <= radius {
				if pd < 0 {
					pd = -pd
					mul = -1.0
				}

				psd := float32(nppos.X)*xt + float32(nppos.Y)*yt - sd
				if psd >= 0.0 && psd <= lineLen {
					fmt.Println("the center of the seg")
					toPushOut := radius - pd + 0.001
					xNudge = ct.X() * toPushOut * mul
					yNudge = ct.Y() * toPushOut * mul
					col = true
				} else {
					var (
						tmpxd float32
						tmpyd float32
					)
					if psd <= 0.0 {
						//fmt.Println("sd smaller")
						tmpxd = float32(nppos.X) - ls.X()
						tmpyd = float32(nppos.Y) - ls.Y()
					} else {
						//fmt.Println("sd greater")
						tmpxd = float32(nppos.X) - le.X()
						tmpyd = float32(nppos.Y) - le.Y()
					}

					distSqr := tmpxd*tmpxd + tmpyd*tmpyd
					if distSqr < radius*radius {
						fmt.Println("Hit either corner of the linedef")
						dist := float32(math.Sqrt(float64(distSqr)))
						toPushOut := radius - dist + 0.001
						xNudge = tmpxd / dist * toPushOut
						yNudge = tmpyd / dist * toPushOut
						col = true
					}
				}

				if col {
					ppos.X += float64(xNudge)
					ppos.Y += float64(yNudge)
					break
				}
			}
		}

		if !col {
			ppos = nppos
		}
		imd.SetColorMask(colornames.Red)
		imd.Push(ppos)
		imd.Circle(3, 0)

		dir := ppos.Add(pdir.Scaled(5))

		imd.Push(ppos, dir)

		imd.Line(1)
		imd.Draw(win)

		win.Update()

		if win.Pressed(pixelgl.KeyRight) {
			camPos.X += 0.2 * (camZoom + 0.2)
		}
		if win.Pressed(pixelgl.KeyLeft) {
			camPos.X -= 0.2 * (camZoom + 0.2)
		}

		if win.Pressed(pixelgl.KeyI) {
			nppos = ppos.Add(pdir.Scaled(0.05))
		}
		if win.Pressed(pixelgl.KeyK) {
			nppos = ppos.Add(pdir.Scaled(-0.05))
		}

		if win.Pressed(pixelgl.KeyJ) {
			player.Angle += 0.3
			y, x := math.Sincos(float64(player.Angle) * math.Pi / 180)
			pdir = pixel.V(x, y)
		}
		if win.Pressed(pixelgl.KeyL) {
			player.Angle -= 0.3
			y, x := math.Sincos(float64(player.Angle) * math.Pi / 180)
			pdir = pixel.V(x, y)
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

		time.Sleep(25)
		imd.Clear()

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
