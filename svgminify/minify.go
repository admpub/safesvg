package svgminify

import (
	"io"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/svg"
)

const MIME_SVG = `image/svg+xml`

var defaultMinifier = NewMinifier()

func NewMinifier() *minify.M {
	m := minify.New()
	m.AddFunc(`text/css`, css.Minify)
	m.AddFunc(MIME_SVG, svg.Minify)
	return m
}

func MinifyReader(w io.Writer, r io.Reader) error {
	return defaultMinifier.MinifyMimetype([]byte(MIME_SVG), w, r, nil)
}

func Minify(b []byte) ([]byte, error) {
	return defaultMinifier.Bytes(MIME_SVG, b)
}
