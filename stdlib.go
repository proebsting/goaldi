//  stdlib.go -- definition of standard library

//  #%#% this initial set is for testing and illustration; it is NOT final!

package goaldi

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
)

var StdLib = []*VProcedure{

	GoProcedure("abs", math.Abs),
	GoProcedure("min", math.Min),
	GoProcedure("max", math.Max),
	GoProcedure("log", math.Log),
	GoProcedure("sqrt", math.Sqrt),

	GoProcedure("intn", rand.Intn),
	GoProcedure("seed", rand.Seed),

	GoProcedure("equalfold", strings.EqualFold),
	GoProcedure("replace", strings.Replace),
	GoProcedure("toupper", strings.ToUpper),
	GoProcedure("tolower", strings.ToLower),
	GoProcedure("trim", strings.Trim),

	GoProcedure("print", fmt.Print),
	GoProcedure("println", fmt.Println),
	GoProcedure("printf", fmt.Printf),
	GoProcedure("fprint", fmt.Fprint),
	GoProcedure("fprintln", fmt.Fprintln),
	GoProcedure("fprintf", fmt.Fprintf),

	GoProcedure("exit", os.Exit),
	GoProcedure("remove", os.Remove),
}
