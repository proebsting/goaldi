//  fpainter.go -- canvas painting functions and methods

package graphics

import (
	"fmt"
	g "goaldi/runtime"
	"math"
)

var _ = fmt.Printf // enable debugging

//  Declare methods
var PainterMethods = g.MethodTable([]*g.VProcedure{
	g.DefMeth((*VPainter).Color, "color", "k", "set drawing color"),
	g.DefMeth((*VPainter).Turn, "turn", "d", "alter orientation by d degrees"),
})

//  MakeCanvas(w,h,d) creates a new canvas and returns a painter value.
func MakeCanvas(env *g.Env, args ...g.Value) (g.Value, *g.Closure) {
	defer g.Traceback("canvas", args)
	if len(args) == 0 {
		return g.Return(NewPainter(-1, -1, -1))
	}
	w := g.FloatVal(g.ProcArg(args, 0, g.NewNumber(10*72)))
	h := g.FloatVal(g.ProcArg(args, 1, g.NewNumber(w)))
	d := g.FloatVal(g.ProcArg(args, 2, g.ONE))
	if w < 1 {
		panic(g.NewExn("Invalid width", w))
	}
	if h < 1 {
		panic(g.NewExn("Invalid height", h))
	}
	if d <= 0 {
		panic(g.NewExn("Invalid density", d))
	}
	return g.Return(NewPainter(w, h, d))
}

//  C.color(r,g,b,a) sets the drawing color for canvas c.
//  With no arguments, the color remains unchanged.
//  The current or updated color value is returned.
func (v *VPainter) Color(args ...g.Value) (g.Value, *g.Closure) {
	defer g.Traceback("color", args)
	k := g.ProcArg(args, 0, g.NilValue)
	if k != g.NilValue {
		if _, ok := k.(VColor); !ok {
			k, _ = Color(nil, args...)
		}
		v.VColor = k.(VColor)
	}
	return g.Return(v)
}

//  C.turn(d) adjusts the current orientation by d degrees.
//  If d is nil, the orientation remains unchanged.
//  The current or updated orientation is returned.
func (v *VPainter) Turn(args ...g.Value) (g.Value, *g.Closure) {
	defer g.Traceback("turn", args)
	d := g.ProcArg(args, 0, g.NilValue)
	if d != g.NilValue {
		v.Aim = math.Mod(v.Aim+g.FloatVal(d), 360)
	}
	return g.Return(v)
}
