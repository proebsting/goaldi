//  fcolor.go -- color functions and methods

package runtime

import (
	"fmt"
)

var _ = fmt.Printf // enable debugging

//  Declare methods
var ColorMethods = MethodTable([]*VProcedure{})

//	Color(r,g,b,a) creates and returns a new color.
//
//	With one argument:  r is a color name, or a grayscale value in (0, 1).
//
//	With two arguments: r is a grayscale value; g is an alpha value in (0, 1).
//
//	With three arguments: r,g,b are color components in (0, 1).
//
//	With four arguments:  r,g,b,a are color components in (0, 1).
func Color(env *Env, args ...Value) (Value, *Closure) {
	x := ProcArg(args, 0, NilValue)
	if s, ok := x.(*VString); ok {
		u := s.ToUTF8()
		if k, ok := ColorMeaning[u]; ok {
			return Return(k)
		}
		if _, err := ParseNumber(u); err != nil {
			panic(NewExn("Unrecognized color name", x))
		}
	}
	var r, g, b, a float64
	switch len(args) {
	case 1:
		r = FloatVal(args[0])
		g = r
		b = r
		a = 1.0
	case 2:
		r = FloatVal(args[0])
		g = r
		b = r
		a = FloatVal(args[1])
	case 3:
		r = FloatVal(args[0])
		g = FloatVal(args[1])
		b = FloatVal(args[2])
		a = 1.0
	case 4:
		r = FloatVal(args[0])
		g = FloatVal(args[1])
		b = FloatVal(args[2])
		a = FloatVal(args[3])
	}
	if r < 0 || r > 1 {
		panic(NewExn("Color value out of range", args[0]))
	}
	if g < 0 || g > 1 {
		panic(NewExn("Color value out of range", args[1]))
	}
	if b < 0 || b > 1 {
		panic(NewExn("Color value out of range", args[2]))
	}
	if a < 0 || a > 1 {
		panic(NewExn("Alpha value out of range", a))
	}
	return Return(NewColor(r, g, b, a))
}
