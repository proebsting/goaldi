//  vfont.go -- VFont, the Goaldi type "font"

package graphics

import (
	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/truetype"
	"fmt"
	g "goaldi/runtime"
	"golang.org/x/mobile/exp/font"
	"image"
)

//  VFont implements a Goaldi font.
type VFont struct {
	Name           string  // requested name
	Ptsize         float64 // requested point size
	*truetype.Font         // underlying font
}

const DefaultFontSize = 12

//  NewFont -- make a new font
func NewFont(name string, ptsize float64) *VFont {
	var f []byte
	switch name {
	case "", "default":
		{
			f = font.Default()
		}
	case "mono", "fixed":
		{
			f = font.Monospace()
		}
	default:
		{
			panic(g.NewExn("Unrecognized typeface", name))
		}
	}
	font, err := freetype.ParseFont(f)
	if err != nil {
		panic(err)
	}
	return &VFont{name, ptsize, font}
}

const rFont = 34         // declare sort ranking
var _ g.ICore = &VFont{} // validate implementation

//  FontType is the font instance of type type.
var FontType = g.NewType("font", "y", rFont, Font, FontMethods,
	"font", "name,ptsize", "load font")

//  VFont.String -- convert to string
func (y *VFont) String() string {
	return fmt.Sprintf("y:%s", y.Name)
}

//  VFont.GoString -- convert to Go string for image() and printf("%#v")
func (y *VFont) GoString() string {
	return fmt.Sprintf("font(%s,%.1f)", y.Name, y.Ptsize)
}

//  VFont.Type -- return the font type
func (y *VFont) Type() g.IRank {
	return FontType
}

//  VFont.Copy returns itself
func (y *VFont) Copy() g.Value {
	return y
}

//  VFont.Before compares two fonts for sorting
//  Ordering is first by name and then by size.
func (a *VFont) Before(b g.Value, i int) bool {
	y := b.(*VFont)
	if a.Name != y.Name {
		return a.Name < y.Name
	} else {
		return a.Ptsize < y.Ptsize
	}
}

//  VFont.Import returns itself
func (y *VFont) Import() g.Value {
	return y
}

//  VFont.Export returns itself.
func (y *VFont) Export() interface{} {
	return y
}

//  VFont.Typeset(painter, x, y, s) draws a string in an image.
func (f *VFont) Typeset(v *VPainter, x, y int, s string) {
	cx := freetype.NewContext()
	cx.SetFont(f.Font)
	cx.SetSrc(image.NewUniform(v.VColor))
	cx.SetHinting(freetype.FullHinting)
	cx.SetDst(v.Image)
	cx.SetClip(v.Image.Bounds())
	cx.SetDPI(72 * v.PixPerPt)
	cx.SetFontSize(f.Ptsize)
	base := freetype.Pt(x, y)
	if _, err := cx.DrawString(s, base); err != nil {
		panic(err)
	}
}
