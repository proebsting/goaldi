//  sprite.go -- code dealing with sprites

//  A sprite contains parameters for overlaying an image on a parent.
//  The displayed screen is rendered by prefix traversal of a tree of
//  one or more sprites.

//  #%#%  A sprite can only be placed on an *app* canvas.

package graphics

import (
	"fmt"
	"golang.org/x/mobile/exp/f32"
)

type Sprite struct {
	Parent   *VPainter  // destination on which sprite is overlain
	Source   *Canvas    // image to be drawn over parent
	X, Y     float32    // location on parent
	Scale    float32    // scaling
	Xform    f32.Affine // transformation for placement on parent
	Children []*Sprite  // subnodes in tree
}

//  Sprite.String() produces a printable representation of a Sprite.
func (p *Sprite) String() string {
	return fmt.Sprintf("Sprite(%v,%v,%d)", p.Source, p.Xform, len(p.Children))
}

//  NewSprite(dst, src, x, y, scale) creates and initializes a new sprite.
//  The src sprite is displayed with its origin over (x,y) of dst.
func NewSprite(dst *VPainter, src *Canvas, x, y, scale float32) *Sprite {
	e := &Sprite{Parent: dst, Source: src}
	e.MoveTo(x, y, scale)
	return e
}

//  VPainter.AddSprite(src, x, y, scale) creates a sprite on a canvas.
func (p *VPainter) AddSprite(src *Canvas, x, y, scale float32) *Sprite {
	src.MakeDisplayable()
	e := NewSprite(p, src, x, y, scale)
	p.Canvas.Sprite.Children = append(p.Canvas.Sprite.Children, e)
	return e
}

//  Sprite.MoveTo(x,y,scale) sets the location of a sprite on its parent.
//  The center (#%#%??) of the sprite is aligned with (x,y).
//
//  Note that this does not expose the full generality of possible transforms:
//  There is no provision for rotation or skew.
func (e *Sprite) MoveTo(x, y, scale float32) {
	e.X = x
	e.Y = y
	e.Scale = scale
	if p := e.Parent; p != nil {
		x += float32(p.ToPx(p.Dx))
		y += float32(p.ToPx(p.Dy))
		scale *= float32(e.Parent.PixPerPt / e.Source.PixPerPt)
	}
	v := e.Source
	m := &e.Xform
	m.Identity()
	m.Translate(m, x, y)
	m.Scale(m, scale, scale)
	m.Translate(m, float32(-v.Width)/2, float32(-v.Height)/2)
}
