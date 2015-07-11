//  sprite.go -- code dealing with sprites

//  A sprite contains parameters for overlaying an image on a parent.
//  The displayed screen is rendered by prefix traversal of a tree of
//  one or more sprites.

package runtime

import (
	"fmt"
	"golang.org/x/mobile/exp/f32"
)

type Sprite struct {
	Parent   *Canvas    // destination on which sprite is overlain
	Source   *Canvas    // image to be drawn over parent
	Xform    f32.Affine // transformation for placement on parent
	Children []*Sprite  // subnodes in tree
}

//  Sprite.String() produces a printable representation of a Sprite.
func (p *Sprite) String() string {
	return fmt.Sprintf("Sprite(%v,%v,%d)", p.Source, p.Xform, len(p.Children))
}

//  NewSprite(dst, src, x, y, scale) creates and initializes a new sprite.
func NewSprite(dst, src *Canvas, x, y, scale float32) *Sprite {
	e := &Sprite{Parent: dst, Source: src}
	e.MoveTo(x, y, scale)
	return e
}

//  Canvas.AddSprite(src, x, y, scale) creates a sprite on a canvas.
func (c *Canvas) AddSprite(src *Canvas, x, y, scale float32) *Sprite {
	//#%#% recompute ppp in case of ppp disagreements?
	e := NewSprite(c, src, x, y, scale)
	c.Sprite.Children = append(c.Sprite.Children, e)
	return e
}

//  Sprite.MoveTo(x,y,scale) sets the location of a sprite on its parent.
//
//  Note that this does not expose the full generality of possible transforms:
//  There is no provision for rotation or skew.
func (e *Sprite) MoveTo(x, y, scale float32) {
	m := &e.Xform
	m.Identity()
	m.Scale(m, scale, scale)
	m.Translate(m, x, y)
}
