package main

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"math/rand"
	"time"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"
)

type tile struct {
	color color.Color
}

type tiles [48]tile

var (
	images *glutil.Images
	t      tiles
)

func onStart(glctx gl.Context) {
	images = glutil.NewImages(glctx)
}

func onStop(glctx gl.Context) {
	images.Release()
}

func onPaint(glctx gl.Context, sz size.Event) {

	glctx.ClearColor(236, 240, 241, 1)
	glctx.Clear(gl.COLOR_BUFFER_BIT)

	pt := sz.WidthPt / 8

	m := images.NewImage(int(pt.Px(sz.PixelsPerPt)), int(pt.Px(sz.PixelsPerPt)))

	ps := (sz.HeightPt / 2) - (pt * 5)

	for c := 1; c <= 6; c++ {

		for r := 1; r <= 8; r++ {

			draw.Draw(m.RGBA, m.RGBA.Bounds(), &image.Uniform{t[(c*r)-1].color}, image.Point{}, draw.Src)

			m.Upload()

			m.Draw(
				sz,
				geom.Point{(geom.Pt(c) * pt), (geom.Pt(r) * pt) + ps},
				geom.Point{(geom.Pt(c) * pt) + pt, (geom.Pt(r) * pt) + ps},
				geom.Point{(geom.Pt(c) * pt), (geom.Pt(r) * pt) + (pt + ps)},
				m.RGBA.Bounds(),
			)

		}

	}

}

func main() {

	log.Println("starting the app")

	colors := []color.RGBA{
		{52, 152, 219, 1},
		{231, 76, 60, 1},
		{52, 73, 94, 1},
		{46, 204, 113, 1},
	}

	for c := 1; c <= 6; c++ {

		for r := 1; r <= 8; r++ {

			t[(c*r - 1)].color = colors[random(0, 4)]

		}

	}

	app.Main(func(a app.App) {

		var glctx gl.Context

		visible, sz := false, size.Event{}

		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					visible = true
					glctx, _ = e.DrawContext.(gl.Context)
					onStart(glctx)
				case lifecycle.CrossOff:
					visible = false
					onStop(glctx)
				}
			case size.Event:
				sz = e
			case touch.Event:
				log.Println("touch event")
			case paint.Event:
				onPaint(glctx, sz)
				a.Publish()
				if visible {
					a.Send(paint.Event{})
				}
			}

		}

	})

}

func random(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min) + min
}
