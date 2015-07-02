//  fcanvas.go -- canvas functions and methods

package runtime

import (
	"math"
)

//  Declare methods
var CanvasMethods = MethodTable([]*VProcedure{
	DefMeth((*VCanvas).Color, "color", "k", "set drawing color"),
	DefMeth((*VCanvas).Turn, "turn", "d", "alter orientation by d degrees"),
})

//  canvas(w,h,d) creates and returns a new canvas.
func Canvas(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("canvas", args)
	if len(args) == 0 {
		return Return(NewCanvas(-1, -1, -1))
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
	return Return(NewCanvas(w, h, d))
}

//  C.color(r,g,b,a) sets the drawing color for canvas c.
//  With no arguments, the color remains unchanged.
//  The current or updated color value is returned.
func (v *VCanvas) Color(args ...Value) (Value, *Closure) {
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
func (v *VCanvas) Turn(args ...Value) (Value, *Closure) {
	defer Traceback("turn", args)
	d := ProcArg(args, 0, NilValue)
	if d != NilValue {
		v.Aim = math.Mod(v.Aim+FloatVal(d), 360)
	}
	return Return(NewNumber(v.Aim))
}
