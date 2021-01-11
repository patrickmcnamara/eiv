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

	"github.com/nfnt/resize"
)

const (
	maxw = 3840
	maxh = 2160
	bw   = 1280
	bh   = 720
)

func main() {
	// open image file
	if len(os.Args) < 2 {
		err := errors.New("no filename given")
		chk(err)
	}
	filename := os.Args[1]
	mf, err := os.Open(filename)
	chk(err)

	// decode image and resize to max
	m, mt, err := image.Decode(mf)
	chk(err)
	omw, omh := m.Bounds().Dx(), m.Bounds().Dy()
	m = resize.Thumbnail(maxw, maxh, m, resize.Lanczos3)

	// create initial resized image
	rm := resize.Thumbnail(bw, bh, m, resize.NearestNeighbor)
	rmw, rmh := rm.Bounds().Dx(), rm.Bounds().Dy()

	driver.Main(func(s screen.Screen) {
		// create window sized to resized image
		title := fmt.Sprintf("%s (%s/%d*%d)", filepath.Base(filename), mt, omw, omh)
		wnd, err := s.NewWindow(&screen.NewWindowOptions{
			Width:  rmw,
			Height: rmh,
			Title:  title,
		})
		chk(err)
		wr := image.Rect(0, 0, rmw, rmh)

		// create initial buffer and draw image to it
		buf, err := s.NewBuffer(rm.Bounds().Size())
		chk(err)
		draw.Draw(buf.RGBA(), buf.Bounds(), rm, image.Point{}, draw.Src)

		for {
			// wait for next event and handle
			switch e := wnd.NextEvent().(type) {

			// window close
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					return
				}

			// window resize (or close on macOS)
			case size.Event:
				// update window rectangle
				wr = e.Bounds()

				// check to close window on macOS
				if wr.Empty() {
					return
				}

			// other paint
			case paint.Event:
				// if image size has changed since last paint
				if !m.Bounds().In(wr) {
					// release old buffer
					buf.Release()

					// create new resized image and buffer
					rm = resize.Thumbnail(uint(wr.Dx()), uint(wr.Dy()), m, resize.NearestNeighbor)
					buf, err = s.NewBuffer(rm.Bounds().Size())
					chk(err)

					// draw image to buffer
					draw.Draw(buf.RGBA(), buf.Bounds(), rm, image.Point{}, draw.Src)
				}

				// calculate new starting point
				sp := image.Pt((wr.Dx()-rm.Bounds().Dx())/2, (wr.Dy()-rm.Bounds().Dy())/2)

				// fill window as black
				wnd.Fill(wr, color.Black, draw.Src)

				// upload buffer to window at starting point and publish
				wnd.Upload(sp, buf, buf.Bounds())
				wnd.Publish()
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
