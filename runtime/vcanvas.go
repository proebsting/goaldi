//  vcanvas.go -- VCanvas, the Goaldi type "canvas"

package runtime

import (
	"fmt"
)

//  VCanvas implements a Goaldi canvas, a Surface pointer plus local attributes
type VCanvas struct {
	*Surface         // underlying surface
	VColor           // drawing color
	Xloc     float64 // drawing location
	Yloc     float64
	Aim      float64 // orientation in degrees
	Size     float64 // drawing width in pixels
}

//  NewCanvas -- construct a new Goaldi canvas
func NewCanvas(w int, h int, d float32) *VCanvas {
	v := &VCanvas{}
	if w < 0 || h < 0 {
		v.Surface = AppSurface()
	} else {
		v.Surface = MemSurface(w, h, d)
	}
	v.VColor = NewColor(0, 0, 0, 1) // color = black
	v.Xloc = float64(v.Width) / 2   // origin at center
	v.Yloc = float64(v.Height) / 2
	v.Aim = -90                  // orientation = towards top
	v.Size = float64(v.PixPerPt) // drawing width = 1 pt
	return v
}

const rCanvas = 32       // declare sort ranking
var _ ICore = &VCanvas{} // validate implementation

//  CanvasType is the canvas instance of type type.
var CanvasType = NewType("canvas", "C", rCanvas, Canvas, CanvasMethods,
	"canvas", "width,height,density", "create canvas")

//  VCanvas.String -- default conversion to Go string returns "C:nnxnn"
func (c *VCanvas) String() string {
	return fmt.Sprintf("C:%dx%d", c.Width, c.Height)
}

//  VCanvas.GoString -- convert to Go string for image() and printf("%#v")
func (c *VCanvas) GoString() string {
	return fmt.Sprintf("canvas(%d,%d,%.2f)", c.Width, c.Height, c.PixPerPt)
}

//  VCanvas.Type -- return the canvas type
func (c *VCanvas) Type() IRank {
	return CanvasType
}

//  VCanvas.Copy returns a new canvas sharing the same underlying surface
func (c *VCanvas) Copy() Value {
	new := *c
	return &new
}

//  VCanvas.Before compares two canvases for sorting
func (a *VCanvas) Before(b Value, i int) bool {
	return false // no ordering defined
}

//  VCanvas.Import returns itself
func (v *VCanvas) Import() Value {
	return v
}

//  VCanvas.Export returns itself.
func (v *VCanvas) Export() interface{} {
	return v
}
