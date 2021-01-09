package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

func main() {
	// open image
	if len(os.Args) < 2 {
		err := errors.New("no filename given")
		chk(err)
	}
	filename := os.Args[1]
	mf, err := os.Open(filename)
	chk(err)

	// decode image
	m, mt, err := image.Decode(mf)
	chk(err)

	driver.Main(func(s screen.Screen) {
		// create window sized to image or max
		width, height := 1280, 720
		mw := m.Bounds().Dx()
		mh := m.Bounds().Dy()
		title := fmt.Sprintf("%s (%s/%d*%d)", filepath.Base(filename), mt, mw, mh)
		if mw < width {
			width = mw
		}
		if mh < height {
			height = mh
		}
		w, err := s.NewWindow(&screen.NewWindowOptions{
			Width:  width,
			Height: height,
			Title:  title,
		})
		chk(err)
		defer w.Release()

		// set default points
		wr := m.Bounds()
		x, y := 0, 0

		// create buffer
		b, err := s.NewBuffer(m.Bounds().Max)
		chk(err)
		defer b.Release()

		for {
			// wait for next event and handle
			switch e := w.NextEvent().(type) {

			// window close
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					return
				}

			// window resize
			case size.Event:
				wr = e.Bounds()

			// other paint
			case paint.Event:
				// fill window as black
				w.Fill(wr, color.Black, draw.Src)

				// resize and draw image to buffer in centre
				rm := resize.Thumbnail(uint(wr.Dx()), uint(wr.Dy()), m, resize.Lanczos3)
				if rm.Bounds().Dx() < wr.Dx() {
					x = (wr.Dx() - rm.Bounds().Dx()) / 2
				}
				if rm.Bounds().Dy() < wr.Dy() {
					y = (wr.Dy() - rm.Bounds().Dy()) / 2
				}
				draw.Draw(b.RGBA(), b.Bounds(), rm, image.Point{}, draw.Src)

				// upload buffer to window and publish
				w.Upload(image.Point{x, y}, b, b.Bounds())
				w.Publish()
			}
		}
	})
}

func init() {
	err := loadPlugins()
	chk(err)
}

func chk(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "eiv: %s\n", err)
		os.Exit(1)
	}
}
