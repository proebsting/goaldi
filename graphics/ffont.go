//  ffont.go -- font functions and methods

package graphics

import (
	"fmt"
	g "goaldi/runtime"
)

var _ = fmt.Printf // enable debugging

var dftSize = g.NewNumber(DefaultFontSize) // default point size

//  Declare methods
var FontMethods = g.MethodTable([]*g.VProcedure{})

//	Font(name,ptsize) loads a font at a particular point size.
//	The only name that currently works is "mono".
func Font(env *g.Env, args ...g.Value) (g.Value, *g.Closure) {
	name := g.ToString(g.ProcArg(args, 0, g.EMPTY)).ToUTF8()
	ptsize := g.FloatVal(g.ProcArg(args, 1, dftSize))
	return g.Return(NewFont(name, ptsize))
}
