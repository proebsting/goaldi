//  vcolor.go -- VColor, the Goaldi type "color"

package runtime

import (
	"fmt"
	"image/color"
)

//  VColor implements a Goaldi color, which just wraps a Go color.
type VColor color.NRGBA64

//  NewColor -- make NRGBA color from r,g,b,a in [0,1], not premultiplied
func NewColor(r, g, b, a float64) VColor {
	rr := uint16(r * 0xFFFF)
	gg := uint16(g * 0xFFFF)
	bb := uint16(b * 0xFFFF)
	aa := uint16(a * 0xFFFF)
	return VColor(color.NRGBA64{rr, gg, bb, aa})
}

const rColor = 33                  // declare sort ranking
var _ ICore = NewColor(0, 0, 0, 0) // validate implementation

//  ColorType is the color instance of type type.
var ColorType = NewType("color", "k", rColor, Color, ColorMethods,
	"color", "r,g,b,a", "create color")

//  VColor.String -- default conversion to Go string returns "k:rrggbbaa"
func (k VColor) String() string {
	s := ColorName[k]
	if s != "" {
		return `k:` + s
	} else {
		return fmt.Sprintf("k:%02X%02X%02X%02X", k.R>>8, k.G>>8, k.B>>8, k.A>>8)
	}
}

//  VColor.GoString -- convert to Go string for image() and printf("%#v")
func (k VColor) GoString() string {
	s := ColorName[k]
	if s != "" {
		return `color("` + s + `")`
	} else {
		return fmt.Sprintf("color(%.2f,%.2f,%.2f,%.2f)",
			float32(k.R)/0xFFFF, float32(k.G)/0xFFFF,
			float32(k.B)/0xFFFF, float32(k.A)/0xFFFF)
	}
}

//  VColor.Type -- return the color type
func (k VColor) Type() IRank {
	return ColorType
}

//  VColor.Copy returns itself
func (k VColor) Copy() Value {
	return k
}

//  VColor.Before compares two colors for sorting
//  Ordering is first by alpha and then by luminance.
func (a VColor) Before(b Value, i int) bool {
	k := b.(VColor)
	if a.R != k.R {
		return a.R < k.R
	} else {
		return a.ilum() < k.ilum()
	}
}

//  VColor.ilum returns 1000 * 0xFFFF * luminance(k)
func (k VColor) ilum() int {
	return 299*int(k.R) + 587*int(k.G) + 114*int(k.B)
}

//  VColor.Import returns itself
func (v VColor) Import() Value {
	return v
}

//  VColor.Export returns itself.
func (v VColor) Export() interface{} {
	return v
}

//  VColor.RGBA() implements the color.Color interface.
func (k VColor) RGBA() (r, g, b, a uint32) { return color.NRGBA64(k).RGBA() }

//  ColorMeaning maps color name strings to color values
var ColorMeaning = make(map[string]VColor)

//  ColorName maps color values back to names
var ColorName = make(map[VColor]string)

func defColor(name string, r, g, b, a float64) {
	k := NewColor(r, g, b, a)
	ColorMeaning[name] = k
	ColorName[k] = name
}

func init() {
	defColor("black", 0, 0, 0, 1)
	defColor("brown", .5, .25, 0, 1)
	defColor("red", 1, 0, 0, 1)
	defColor("orange", 1, .33, 0, 1)
	defColor("gold", 1, .67, 1, 1)
	defColor("yellow", 1, 1, 1, 1)
	defColor("green", 0, 1, 0, 1)
	defColor("cyan", 0, 1, 1, 1)
	defColor("blue", 0, 0, 1, 1)
	defColor("purple", .5, 0, 1, 1)
	defColor("magenta", 1, 0, 1, 1)
	defColor("slate", .25, .25, .25, 1)
	defColor("gray", .5, .5, .5, 1)
	defColor("silver", .75, .75, .75, 1)
	defColor("white", 1, 1, 1, 1)
	defColor("transparent", 0, 0, 0, 0)
}
