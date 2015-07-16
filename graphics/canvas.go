//  canvas.go -- image canvas code.

package graphics

import (
	"fmt"
	"golang.org/x/mobile/exp/gl/glutil"
	"image"
	"image/draw"
)

//  A Canvas is a grid of pixels forming an image.
type Canvas struct {
	*App               // associated app if app canvas, else nil
	*Sprite            // placement on screen, and overlain sprites
	Width      int     // width in pixels
	Height     int     // height in pixels
	PixPerPt   float64 // density in pixels/point
	draw.Image         // underlying image
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

//  NewCanvas creates a new Canvas of size w x h points at density ppp.
//  If w or h is negative, size and density are set by the app window.
func NewCanvas(w, h, ppp float64) *Canvas {
	var c = &Canvas{}
	if w >= 0 && h >= 0 {
		// simple offscreen scratch canvas
		// (not a GL image, because that doesn't work if not an app)
		ww := int(w*ppp + 0.5)
		hh := int(h*ppp + 0.5)
		c.Image = image.NewRGBA(image.Rect(0, 0, ww, hh))
		c.setup(ww, hh, ppp)
	} else {
		// application canvas
		w, h, ppp := AppSize()
		c.Image = glutil.NewImage(w, h)
		c.setup(w, h, ppp)
		AppCanvas(c)
	}
	return c
}

//  Canvas.setup initializes a canvas struct.
func (c *Canvas) setup(w int, h int, ppp float64) {
	im := c.Image
	draw.Draw(im, im.Bounds(), image.White, image.Point{}, draw.Src) // erase
	c.Width = w
	c.Height = h
	c.PixPerPt = ppp
	c.Image = im
	c.Sprite = NewSprite(nil, c, 0, 0, 1)
}

//  Canvas.MakeDisplayable() makes a canvas useable in an app context.
//  This means changing its image to a GL image if it is not one already.
func (c *Canvas) MakeDisplayable() {
	if _, ok := c.Image.(*glutil.Image); !ok { // if not alread a GL image
		im := glutil.NewImage(c.Width, c.Height)
		draw.Draw(im, im.Bounds(), c.Image, image.Point{}, draw.Src)
		c.Image = im
	}
}
