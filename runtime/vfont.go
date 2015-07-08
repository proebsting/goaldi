//  vfont.go -- VFont, the Goaldi type "font"

package runtime

import (
	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/truetype"
	"fmt"
	"golang.org/x/mobile/exp/font"
	"image"
)

//  VFont implements a Goaldi font.
type VFont struct {
	name           string  // requested name
	ptsize         float64 // requested point size
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
			panic(NewExn("Unrecognized typeface", name))
		}
	}
	font, err := freetype.ParseFont(f)
	if err != nil {
		panic(err)
	}
	return &VFont{name, ptsize, font}
}

const rFont = 34       // declare sort ranking
var _ ICore = &VFont{} // validate implementation

//  FontType is the font instance of type type.
var FontType = NewType("font", "y", rFont, Font, FontMethods,
	"font", "name,ptsize", "load font")

//  VFont.String -- convert to string
func (y *VFont) String() string {
	return fmt.Sprintf("y:%s", y.name)
}

//  VFont.GoString -- convert to Go string for image() and printf("%#v")
func (y *VFont) GoString() string {
	return fmt.Sprintf("font(%s,%.1f)", y.name, y.ptsize)
}

//  VFont.Type -- return the font type
func (y *VFont) Type() IRank {
	return FontType
}

//  VFont.Copy returns itself
func (y *VFont) Copy() Value {
	return y
}

//  VFont.Before compares two fonts for sorting
//  Ordering is first by alpha and then by luminance.
func (a *VFont) Before(b Value, i int) bool {
	y := b.(*VFont)
	if a.name != y.name {
		return a.name < y.name
	} else {
		return a.ptsize < y.ptsize
	}
}

//  VFont.Import returns itself
func (y *VFont) Import() Value {
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
	cx.SetFontSize(f.ptsize)
	base := freetype.Pt(x, y)
	if _, err := cx.DrawString(s, base); err != nil {
		panic(err)
	}
}
