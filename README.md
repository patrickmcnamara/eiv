# eiv (extensible image viewer)

`eiv` is a GUI image viewer written in Go that is extensible with plugins.
By default it supports gif, jpeg, png, bmp, ccitt, tiff, vp8, vp8l, and webp.

![The Starry Night displayed in eiv](./misc/thestarrynight.png)

# Installation

`go get -u github.com/patrickmcnamara/eiv`

# Usage

`eiv $IMAGEFILENAME`

# Plugins

To create a plugin for `eiv` for an image format `abc`, create (or download) a library that can decode that format into an [image.Image](https://pkg.go.dev/image#Image) in Go and [register the format](https://pkg.go.dev/image#RegisterFormat) in the init function of your library.
Next, create a `package main` program in Go that imports your library and does not have a `main` function.
Build this using the plugin build mode for Go with an output name with the extension `.eivp` (e.g. `go build -buildmode=plugin -o abc.eivp`).
Put this plugin into your config directory under `$CONFIGDIR/eiv/plugin`.
Voilà.
`eiv` can now display your images.
