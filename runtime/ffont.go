//  ffont.go -- font functions and methods

package runtime

import (
	"fmt"
)

var _ = fmt.Printf // enable debugging

var dftSize = NewNumber(DefaultFontSize) // default point size

//  Declare methods
var FontMethods = MethodTable([]*VProcedure{})

//	Font(name,ptsize) loads a font at a particular point size.
//	The only name that currently works is "mono".
func Font(env *Env, args ...Value) (Value, *Closure) {
	name := ToString(ProcArg(args, 0, EMPTY)).ToUTF8()
	ptsize := FloatVal(ProcArg(args, 1, dftSize))
	return Return(NewFont(name, ptsize))
}
