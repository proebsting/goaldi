//  draw.go -- canvas drawing operations

//#%#% a crude first hack.
//#%#% will need a good rewrite with error checking, clipping, etc.

package runtime

import (
	"fmt"
	"image"
	"image/draw"
	"math"
)

var _ = fmt.Printf // enable debugging

//  VCanvas.Reset() establishes or reestablises initial conditions:
//		origin = center of surface
//		current location = origin
//		orientation = towards top
//		drawing size = 1 pt
//		color = black
func (v *VCanvas) Reset() *VCanvas {
	v.Dx = float64(v.Width) / (2 * v.PixPerPt) // offset to origin
	v.Dy = float64(v.Height) / (2 * v.PixPerPt)
	v.Xloc = 0 // current location
	v.Yloc = 0
	v.Aim = -90                          // orientation = towards top
	v.Size = 1                           // drawing width = 1 pt
	v.VColor = NewColor(1, 1, 1, 1)      // color = white
	v.Rect(-v.Dx, -v.Dy, 2*v.Dx, 2*v.Dy) // clear the canvas
	v.VColor = NewColor(0, 0, 0, 1)      // color = black
	return v
}

//  VCanvas.Goto(x,y,o) sets the current location to (x, y).
//  If o is a number (not nil) it sets the current orientation to that.
func (v *VCanvas) Goto(x, y float64, o interface{}) *VCanvas {
	v.Xloc = x
	v.Yloc = y
	if o != nil {
		v.Aim = o.(float64)
	}
	return v
}

//  VCanvas.Forward(d) draws a line by moving the pen forward d units.
func (v *VCanvas) Forward(d float64) *VCanvas {
	s, c := math.Sincos(v.Aim * (math.Pi / 180))
	x := v.Xloc + d*c
	y := v.Yloc + d*s
	v.Line(v.Xloc, v.Yloc, x, y)
	v.Xloc = x
	v.Yloc = y
	return v
}

//  VCanvas.Line(x1, y1, x2, y2) draws a line.
//  #%#% in a really dumb way. should stroke, not draw a zillion points.
func (v *VCanvas) Line(x1, y1, x2, y2 float64) *VCanvas {
	dx := x2 - x1
	dy := y2 - y1
	dmax := math.Max(math.Abs(dx), math.Abs(dy))
	n := int(math.Ceil(float64(v.PixPerPt) * dmax))
	dx /= float64(n)
	dy /= float64(n)
	for i := 0; i <= n; i++ {
		v.Point(x1, y1)
		x1 += dx
		y1 += dy
	}
	return v
}

//  VCanvas.Point(x, y) draws a point.
//  #%#% in a really dumb way. should cache the pen. and it should be round.
func (v *VCanvas) Point(x, y float64) *VCanvas {
	v.Rect(x-v.Size/2, y-v.Size/2, v.Size, v.Size)
	return v
}

//  VCanvas.Rect(x, y, w, h) draws a rectangle.
func (v *VCanvas) Rect(x, y, w, h float64) *VCanvas {
	if w < 0 {
		x, w = x+w, -w
	}
	if h < 0 {
		y, h = y+h, -h
	}
	x = x + v.Dx
	y = y + v.Dy
	r := image.Rect(v.ToPx(x), v.ToPx(y), v.ToPx(x+w), v.ToPx(y+h))
	draw.Draw(v.Surface.Image, r,
		image.NewUniform(v.VColor), image.Point{}, draw.Src)
	return v
}

//  VCanvas.Text(x, y, s) draws a string of text characters.
func (v *VCanvas) Text(x, y float64, s string) *VCanvas {
	v.VFont.Typeset(v, v.ToPx(x+v.Dx), v.ToPx(y+v.Dy), s)
	return v
}