//  draw.go -- canvas drawing operations

//#%#% a crude first hack.
//#%#% will need a good rewrite with error checking, clipping, etc.

//#%#% don't call line/rect/point from user code -- wrong coord system

package runtime

import (
	"image"
	"image/draw"
	"math"
)

//  VCanvas.DwForward(d) draws a line by moving the pen forward d units.
func (v *VCanvas) DwForward(d float64) {
	d = d * float64(v.PixPerPt) //#%#% no coord sys yet, scale sensibly
	s, c := math.Sincos(v.Aim * (math.Pi / 180))
	x := v.Xloc + d*c
	y := v.Yloc + d*s
	v.DwLine(v.Xloc, v.Yloc, x, y)
	v.Xloc = x
	v.Yloc = y
}

//  VCanvas.DwLine(x1, y1, x2, y2) draws a line.
//  #%#% in a really dumb way. should stroke, not draw a zillion points.
func (v *VCanvas) DwLine(x1, y1, x2, y2 float64) {
	dx := x2 - x1
	dy := y2 - y1
	n := int(math.Ceil(math.Max(math.Abs(dx), math.Abs(dy))))
	dx /= float64(n)
	dy /= float64(n)
	for i := 0; i <= n; i++ {
		v.DwPoint(x1, y1)
		x1 += dx
		y1 += dy
	}
}

//  VCanvas.DwPoint(x, y) draws a point.
//  #%#% in a really dumb way. should cache the pen. and it should be round.
func (v *VCanvas) DwPoint(x, y float64) {
	v.DwRect(x-v.Size/2, y-v.Size/2, v.Size, v.Size)
}

//  VCanvas.DwRect(x1, y1, x2, y2) draws a rectangle.
func (v *VCanvas) DwRect(x, y, w, h float64) {
	if w < 0 {
		x, w = x+w, -w
	}
	if h < 0 {
		y, h = y+h, -h
	}
	r := image.Rect(int(x+0.5), int(y+0.5), int(x+w+0.5), int(y+h+0.5))
	draw.Draw(v.Surface.Image, r,
		image.NewUniform(v.VColor), image.Point{}, draw.Src)
}
