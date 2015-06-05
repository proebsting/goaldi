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
	r := ProcArg(args, 0, NilValue)
	if s, ok := r.(*VString); ok {
		if k, ok := ColorMeaning[s.ToUTF8()]; ok {
			return Return(k)
		}
	}
	//#%#% TODO: handle numeric arguments
	panic(NewExn("Unrecognized color name", r))
}
