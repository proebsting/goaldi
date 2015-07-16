//  vpainter.go -- VPainter, the Goaldi type "painter"

package graphics

import (
	"fmt"
	g "goaldi/runtime"
)

//  VPainter implements a Goaldi painter for drawing on a canvas.
//  Units are measured in points, as in the Go app package.
//  Floating values are float64, for use with the math package.
type VPainter struct {
	*Canvas            // underlying canvas
	VColor             // drawing color
	*VFont             // text font
	Dx, Dy     float64 // offset to coordinate origin
	Xloc, Yloc float64 // drawing location
	Aim        float64 // orientation in degrees
	Size       float64 // drawing width
}

//  NewPainter -- construct a new Goaldi canvas and return a painter.
//  If w or h is negative, an app canvas is created and installed.
func NewPainter(w, h, d float64) *VPainter {
	p := &VPainter{}
	p.Canvas = NewCanvas(w, h, d)
	p.VFont = NewFont("mono", DefaultFontSize)
	return p.Reset()
}

const rPainter = 32         // declare sort ranking
var _ g.ICore = &VPainter{} // validate implementation

//  PainterType is the painter instance of type type.
var PainterType = g.NewType("painter", "P", rPainter, MakeCanvas, PainterMethods,
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
func (c *VPainter) Type() g.IRank {
	return PainterType
}

//  VPainter.Copy returns a new painter sharing the same underlying canvas
func (c *VPainter) Copy() g.Value {
	new := *c
	return &new
}

//  VPainter.Before compares two painters for sorting
func (a *VPainter) Before(b g.Value, i int) bool {
	return false // no ordering defined
}

//  VPainter.Import returns itself
func (p *VPainter) Import() g.Value {
	return p
}

//  VPainter.Export returns itself.
func (p *VPainter) Export() interface{} {
	return p
}

//  VPainter.ToPx scales a point value to a pixel value
func (p *VPainter) ToPx(n float64) int {
	return int(p.PixPerPt*n + 0.5)
}
