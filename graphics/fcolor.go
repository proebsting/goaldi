//  fcolor.go -- color functions and methods

package graphics

import (
	"fmt"
	g "goaldi/runtime"
	"strconv"
)

var _ = fmt.Printf // enable debugging

//  Declare methods
var ColorMethods = g.MethodTable([]*g.VProcedure{})

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
func Color(env *g.Env, args ...g.Value) (g.Value, *g.Closure) {
	x := g.ProcArg(args, 0, g.NilValue)
	if s, ok := x.(*g.VString); ok {
		u := s.ToUTF8()
		if len(u) > 0 && u[0] == '#' {
			return hexcolor(u)
		}
		if k, ok := ColorMeaning[u]; ok {
			return g.Return(k)
		}
		if _, err := g.ParseNumber(u); err != nil {
			panic(g.NewExn("Unrecognized color name", x))
		}
	}
	var rr, gg, bb, aa float64
	switch len(args) {
	case 1:
		rr = g.FloatVal(args[0])
		gg = rr
		bb = rr
		aa = 1.0
	case 2:
		rr = g.FloatVal(args[0])
		gg = rr
		bb = rr
		aa = g.FloatVal(args[1])
	case 3:
		rr = g.FloatVal(args[0])
		gg = g.FloatVal(args[1])
		bb = g.FloatVal(args[2])
		aa = 1.0
	case 4:
		rr = g.FloatVal(args[0])
		gg = g.FloatVal(args[1])
		bb = g.FloatVal(args[2])
		aa = g.FloatVal(args[3])
	}
	if rr < 0 || rr > 1 {
		panic(g.NewExn("Color value out of range", args[0]))
	}
	if gg < 0 || gg > 1 {
		panic(g.NewExn("Color value out of range", args[1]))
	}
	if bb < 0 || bb > 1 {
		panic(g.NewExn("Color value out of range", args[2]))
	}
	if aa < 0 || aa > 1 {
		panic(g.NewExn("Alpha value out of range", aa))
	}
	return g.Return(NewColor(rr, gg, bb, aa))
}

//  hexcolor(s) returns a color specified by hex digits beginning with '#'.
func hexcolor(s string) (g.Value, *g.Closure) {
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
		panic(g.NewExn("Invalid length for color specification", s))
	}
	rr := hexv(s, w, 0) // red component or grayscale value
	gg := rr            // assume grayscale
	bb := rr
	aa := d           // assume alpha absent
	if len(s) > 2*w { // if green and blue specified
		gg = hexv(s, w, 1) // green
		bb = hexv(s, w, 2) // blue
	}
	if len(s) > 4*w { // if alpha specified
		aa = hexv(s, w, 3) // alpha
	}
	return g.Return(NewColor(rr/d, gg/d, bb/d, aa/d))
}

//  hexv(s, w, i) interprets s[1+iw:1+(i+1)w] as a hex value.
func hexv(s string, w, i int) float64 {
	v, err := strconv.ParseInt(s[1+i*w:1+(i+1)*w], 16, 32)
	if err != nil {
		panic(g.NewExn("Invalid hexadecimal value", s))
	}
	return float64(v)
}
