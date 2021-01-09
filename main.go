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
		w, h := bw, bh
		mw, mh := m.Bounds().Dx(), m.Bounds().Dy()
		if mw > bw {
			h = mh * bw / mw
			w = bw
		}
		if h > bh {
			w = w * bh / h
			h = bh
		}
		title := fmt.Sprintf("%s (%s/%d*%d)", filepath.Base(filename), mt, mw, mh)
		wnd, err := s.NewWindow(&screen.NewWindowOptions{
			Width:  w,
			Height: h,
			Title:  title,
		})
		chk(err)
		defer wnd.Release()
		wr := image.Rect(0, 0, w, h)

		// create buffer
		buf, err := s.NewBuffer(m.Bounds().Max)
		chk(err)
		defer buf.Release()

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
				wr = e.Bounds()
				if wr.Empty() {
					return
				}

			// other paint
			case paint.Event:
				// fill window as black
				wnd.Fill(wr, color.Black, draw.Src)

				// resize and draw image to buffer in centre
				rm := resize.Thumbnail(uint(wr.Dx()), uint(wr.Dy()), m, resize.Lanczos3)
				sp := image.Pt((wr.Dx()-rm.Bounds().Dx())/2, (wr.Dy()-rm.Bounds().Dy())/2)
				draw.Draw(buf.RGBA(), buf.Bounds(), rm, image.Point{}, draw.Src)

				// upload buffer to window and publish
				wnd.Upload(sp, buf, rm.Bounds())
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
