//  fcanvas.go -- canvas functions and methods

package runtime

import ()

//  Declare methods
var CanvasMethods = MethodTable([]*VProcedure{
	DefMeth((*VCanvas).Color, "color", "k", "set drawing color"),
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
	return Return(NewCanvas(w, h, float32(d)))
}

//  C.color(k) sets the drawing color for canvas c.
//  If k is nil, the color remains unchanged.
//  The current or updated value is returned.
func (v *VCanvas) Color(args ...Value) (Value, *Closure) {
	defer Traceback("color", args)
	k := ProcArg(args, 0, NilValue)
	if k != NilValue {
		v.k = k.(VColor)
	}
	return Return(v.k)
}
