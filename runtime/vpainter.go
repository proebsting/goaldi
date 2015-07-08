//  vpainter.go -- VPainter, the Goaldi type "painter"

package runtime

import (
	"fmt"
)

//  VPainter implements a Goaldi painter for drawing on a canvas.
//  Units are measured in points, as in the Go app package.
//  Floating values are float64, for use with the math package.
type VPainter struct {
	*Surface           // underlying surface
	VColor             // drawing color
	*VFont             // text font
	Dx, Dy     float64 // offset to coordinate origin
	Xloc, Yloc float64 // drawing location
	Aim        float64 // orientation in degrees
	Size       float64 // drawing width
}

//  NewPainter -- construct a new Goaldi canvas and return a painter
func NewPainter(w int, h int, d float64) *VPainter {
	v := &VPainter{}
	if w < 0 || h < 0 {
		v.Surface = AppSurface()
		d = v.PixPerPt
	} else {
		v.Surface = MemSurface(w, h, d)
	}
	v.VFont = NewFont("mono", DefaultFontSize)
	return v.Reset()
}

const rPainter = 32       // declare sort ranking
var _ ICore = &VPainter{} // validate implementation

//  PainterType is the painter instance of type type.
var PainterType = NewType("painter", "P", rPainter, Canvas, PainterMethods,
	"canvas", "width,height,density", "create canvas and return painter")

//  VPainter.String -- default conversion to Go string returns "C:nnxnn"
func (c *VPainter) String() string {
	return fmt.Sprintf("C:%dx%d", c.Width, c.Height)
}

//  VPainter.GoString -- convert to Go string for image() and printf("%#v")
func (c *VPainter) GoString() string {
	return fmt.Sprintf("painter(%d,%d,%.2f)", c.Width, c.Height, c.PixPerPt)
}

//  VPainter.Type -- return the painter type
func (c *VPainter) Type() IRank {
	return PainterType
}

//  VPainter.Copy returns a new painter sharing the same underlying surface
func (c *VPainter) Copy() Value {
	new := *c
	return &new
}

//  VPainter.Before compares two painters for sorting
func (a *VPainter) Before(b Value, i int) bool {
	return false // no ordering defined
}

//  VPainter.Import returns itself
func (v *VPainter) Import() Value {
	return v
}

//  VPainter.Export returns itself.
func (v *VPainter) Export() interface{} {
	return v
}

//  VPainter.ToPx scales a point value to a pixel value
func (v *VPainter) ToPx(n float64) int {
	return int(v.PixPerPt*n + 0.5)
}
