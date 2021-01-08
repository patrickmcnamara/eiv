package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
	"path/filepath"

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
		// create window
		title := fmt.Sprintf("%s (%s)", filepath.Base(filename), mt)
		w, err := s.NewWindow(&screen.NewWindowOptions{
			Width:  m.Bounds().Dx(),
			Height: m.Bounds().Dy(),
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
				switch e.To {
				case lifecycle.StageDead:
					return
				default:
					continue
				}
			// window resize
			case size.Event:
				wr = e.Bounds()
			// other paint
			case paint.Event:
			default:
				continue
			}

			// fill window as black
			w.Fill(wr, color.Black, draw.Src)

			// draw image to buffer in centre
			if m.Bounds().Dx() < wr.Dx() {
				x = (wr.Dx() - m.Bounds().Dx()) / 2
			}
			if m.Bounds().Dy() < wr.Dy() {
				y = (wr.Dy() - m.Bounds().Dy()) / 2
			}
			draw.Draw(b.RGBA(), b.Bounds(), m, image.Point{}, draw.Src)

			// upload buffer to window and publish
			w.Upload(image.Point{x, y}, b, b.Bounds())
			w.Publish()
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
