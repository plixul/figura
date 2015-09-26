package main

import (
	"image"
	"image/color"
	"image/draw"
	"log"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"
)

var (
	images *glutil.Images
)

func onStart(glctx gl.Context) {
	images = glutil.NewImages(glctx)
}

func onStop(glctx gl.Context) {
	images.Release()
}

func onPaint(glctx gl.Context, sz size.Event) {

	glctx.ClearColor(128, 0, 128, 1)
	glctx.Clear(gl.COLOR_BUFFER_BIT)

	m := images.NewImage(100, 100)

	draw.Draw(m.RGBA, m.RGBA.Bounds(), &image.Uniform{color.RGBA{0, 0, 255, 255}}, image.Point{}, draw.Src)

	m.Upload()

	m.Draw(
		sz,
		geom.Point{0, 0},
		geom.Point{100, 0},
		geom.Point{0, 100},
		m.RGBA.Bounds(),
	)

}

func main() {

	log.Println("Starting the app")

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
