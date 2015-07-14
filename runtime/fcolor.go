//  fcolor.go -- color functions and methods

package runtime

import (
	"fmt"
	"strconv"
)

var _ = fmt.Printf // enable debugging

//  Declare methods
var ColorMethods = MethodTable([]*VProcedure{})

//	Color(r,g,b,a) creates and returns a new color.
//
//	With one argument:  r is a color name, a grayscale value in (0, 1), or
//	a hexadecimal specification (#k, #kk, #rgb, #rgba, #rrggbb, #rrggbbaa).
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
		if len(u) > 0 && u[0] == '#' {
			return hexcolor(u)
		}
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

//  hexcolor(s) returns a color specified by hex digits beginning with '#'.
func hexcolor(s string) (Value, *Closure) {
	var w int     // width of each field (1 or 2 chars)
	var d float64 // corresponding denominator for scaling
	switch len(s) {
	case 2, 4, 5: // #k, #rgb, #rgba
		w = 1
		d = 15
	case 3, 7, 9: // #kk, #rrggbb, #rrggbbaa
		w = 2
		d = 255
	default:
		panic(NewExn("Invalid length for color specification", s))
	}
	r := hexv(s, w, 0) // red component or grayscale value
	g := r             // assume grayscale
	b := r
	a := d            // assume alpha absent
	if len(s) > 2*w { // if green and blue specified
		g = hexv(s, w, 1) // green
		b = hexv(s, w, 2) // blue
	}
	if len(s) > 4*w { // if alpha specified
		a = hexv(s, w, 3) // alpha
	}
	return Return(NewColor(r/d, g/d, b/d, a/d))
}

//  hexv(s, w, i) interprets s[1+iw:1+(i+1)w] as a hex value.
func hexv(s string, w, i int) float64 {
	v, err := strconv.ParseInt(s[1+i*w:1+(i+1)*w], 16, 32)
	if err != nil {
		panic(NewExn("Invalid hexadecimal value", s))
	}
	return float64(v)
}
