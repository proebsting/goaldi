//  fpainter.go -- canvas painting functions and methods

package runtime

import (
	"math"
)

//  Declare methods
var PainterMethods = MethodTable([]*VProcedure{
	DefMeth((*VPainter).Color, "color", "k", "set drawing color"),
	DefMeth((*VPainter).Turn, "turn", "d", "alter orientation by d degrees"),
})

//  MakeCanvas(w,h,d) creates a new canvas and returns a painter value.
func MakeCanvas(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("canvas", args)
	if len(args) == 0 {
		return Return(NewPainter(-1, -1, -1))
	}
	w := IntVal(ProcArg(args, 0, NewNumber(10*72)))
	h := IntVal(ProcArg(args, 1, NewNumber(float64(w))))
	d := FloatVal(ProcArg(args, 2, ONE))
	if w < 1 {
		panic(NewExn("Invalid width", w))
	}
	if h < 1 {
		panic(NewExn("Invalid height", h))
	}
	if d <= 0 {
		panic(NewExn("Invalid density", d))
	}
	return Return(NewPainter(w, h, d))
}

//  C.color(r,g,b,a) sets the drawing color for canvas c.
//  With no arguments, the color remains unchanged.
//  The current or updated color value is returned.
func (v *VPainter) Color(args ...Value) (Value, *Closure) {
	defer Traceback("color", args)
	k := ProcArg(args, 0, NilValue)
	if k != NilValue {
		if _, ok := k.(VColor); !ok {
			k, _ = Color(nil, args...)
		}
		v.VColor = k.(VColor)
	}
	return Return(v.VColor)
}

//  C.turn(d) adjusts the current orientation by d degrees.
//  If d is nil, the orientation remains unchanged.
//  The current or updated orientation is returned.
func (v *VPainter) Turn(args ...Value) (Value, *Closure) {
	defer Traceback("turn", args)
	d := ProcArg(args, 0, NilValue)
	if d != NilValue {
		v.Aim = math.Mod(v.Aim+FloatVal(d), 360)
	}
	return Return(NewNumber(v.Aim))
}
