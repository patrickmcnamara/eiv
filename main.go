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

const (
	bw = 1280
	bh = 720
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

	// decode image
	m, mt, err := image.Decode(mf)
	chk(err)

	driver.Main(func(s screen.Screen) {
		// create initial resized image
		rm := resize.Thumbnail(bw, bh, m, resize.Lanczos3)

		// get image sizes
		mw, mh := m.Bounds().Dx(), m.Bounds().Dy()
		rmw, rmh := rm.Bounds().Dx(), rm.Bounds().Dy()

		// create window sized to resized image
		title := fmt.Sprintf("%s (%s/%d*%d)", filepath.Base(filename), mt, mw, mh)
		wnd, err := s.NewWindow(&screen.NewWindowOptions{
			Width:  rmw,
			Height: rmh,
			Title:  title,
		})
		chk(err)
		wr := image.Rect(0, 0, rmw, rmh)

		// create initial buffer
		buf, err := s.NewBuffer(rm.Bounds().Size())
		chk(err)

		// calculate initial buffer starter point
		sp := image.Pt((wr.Dx()-rm.Bounds().Dx())/2, (wr.Dy()-rm.Bounds().Dy())/2)

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
				// release old buffer
				buf.Release()

				// create new resized image and buffer
				rm = resize.Thumbnail(uint(wr.Dx()), uint(wr.Dy()), m, resize.NearestNeighbor)
				buf, err = s.NewBuffer(rm.Bounds().Size())
				chk(err)

				// resize and draw image to buffer in centre
				draw.Draw(buf.RGBA(), buf.Bounds(), rm, image.Point{}, draw.Src)

				// calculate new starting point
				sp = image.Pt((wr.Dx()-rm.Bounds().Dx())/2, (wr.Dy()-rm.Bounds().Dy())/2)

				// fill window as black
				wnd.Fill(wr, color.Black, draw.Src)

				// upload buffer to window and publish
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
