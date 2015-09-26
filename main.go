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
	"golang.org/x/mobile/event/touch"
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

	glctx.ClearColor(250, 250, 250, 1)
	glctx.Clear(gl.COLOR_BUFFER_BIT)

	height := sz.HeightPt / 6
	width := sz.WidthPt / 6

	m := images.NewImage(int(height.Px(sz.PixelsPerPt)), int(width.Px(sz.PixelsPerPt)))

	draw.Draw(m.RGBA, m.RGBA.Bounds(), &image.Uniform{color.RGBA{248, 90, 96, 1}}, image.Point{}, draw.Src)

	m.Upload()

	m.Draw(
		sz,
		geom.Point{0, 0},
		geom.Point{width, 0},
		geom.Point{0, height},
		m.RGBA.Bounds(),
	)

}

func main() {

	log.Println("starting the app")

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
