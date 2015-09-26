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

const (
	gridColumnsCount = 6
	gridrowsCount    = 8
	tileSize         = 60
)

type tile struct {
	columnOffset geom.Pt
	rowOffset    geom.Pt
	x            geom.Pt
	y            geom.Pt
	image        *glutil.Image
	color        color.Color
}

func (t *tile) Draw(sz size.Event) {

	draw.Draw(t.image.RGBA, t.image.RGBA.Bounds(), &image.Uniform{t.color}, image.Point{}, draw.Src)

	t.image.Upload()

	t.image.Draw(
		sz,
		geom.Point{t.columnOffset + t.x, t.rowOffset + t.y},
		geom.Point{t.columnOffset + (t.x + tileSize), t.rowOffset + t.y},
		geom.Point{t.columnOffset + t.x, t.rowOffset + (t.y + tileSize)},
		t.image.RGBA.Bounds(),
	)

}

type tiles [gridColumnsCount * gridrowsCount]tile

var (
	images *glutil.Images
	grid   tiles
)

func onStart(glctx gl.Context, sz size.Event) {
	images = glutil.NewImages(glctx)

	colors := []color.RGBA{
		{52, 152, 219, 1},
		{231, 76, 60, 1},
		{52, 73, 94, 1},
		{46, 204, 113, 1},
	}

	for c := 0; c < gridColumnsCount; c++ {
		for r := 0; r < gridrowsCount; r++ {
			t := &grid[((c+1)*(r+1))-1]
			t.image = images.NewImage(int(sz.PixelsPerPt*tileSize), int(sz.PixelsPerPt*tileSize))
			t.color = colors[random(0, 4)]
		}
	}
}

func onStop(glctx gl.Context) {
	images.Release()
}

func onPaint(glctx gl.Context, sz size.Event) {

	glctx.ClearColor(236, 240, 241, 1)
	glctx.Clear(gl.COLOR_BUFFER_BIT)

	columnOffset := (sz.WidthPt / 2) - ((geom.Pt(gridColumnsCount) * tileSize) / 2)
	rowOffset := (sz.HeightPt / 2) - ((geom.Pt(gridrowsCount) * tileSize) / 2)

	for c := 0; c < gridColumnsCount; c++ {

		for r := 0; r < gridrowsCount; r++ {

			t := &grid[((c+1)*(r+1))-1]

			t.columnOffset = columnOffset
			t.rowOffset = rowOffset
			t.x = geom.Pt(c) * tileSize
			t.y = geom.Pt(r) * tileSize
			t.Draw(sz)

		}

	}

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
					onStart(glctx, sz)
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
