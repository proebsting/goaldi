//  canvas.go -- image canvas code.

package runtime

import (
	"fmt"
	"golang.org/x/mobile/exp/gl/glutil"
	"image"
	"image/draw"
)

//  A Canvas is a grid of pixels forming an image.
type Canvas struct {
	*App                  // associated app if app canvas, else nil
	Width         int     // width in pixels
	Height        int     // height in pixels
	PixPerPt      float64 // density in pixels/point
	*glutil.Image         // underlying image
}

//  Canvas.String() produces a printable representation of a Canvas struct.
func (s *Canvas) String() string {
	a := "-"
	if s.App != nil {
		a = "A"
	}
	return fmt.Sprintf("Canvas(%s,%dx%dx%.2f)",
		a, s.Width, s.Height, s.PixPerPt)
}

//  NewCanvas creates a new Canvas.
//  If w or h is negative, the canvas is connected to the app window.
func NewCanvas(w int, h int, ppp float64) *Canvas {
	var c = &Canvas{}
	if w >= 0 && h >= 0 {
		// simple offscreen scratch canvas
		c.setup(w, h, ppp)
	} else {
		// application canvas
		c.setup(AppSize())
		AppCanvas(c)
	}
	return c
}

//  C.setup initializes a canvas struct.
func (c *Canvas) setup(w int, h int, ppp float64) {
	im := glutil.NewImage(w, h)
	draw.Draw(im, im.Bounds(), image.White, image.Point{}, draw.Src) // erase
	c.Width = w
	c.Height = h
	c.PixPerPt = ppp
	c.Image = im
}
