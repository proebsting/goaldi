package main

import (
	"fmt"
	. "goaldi"
	"math"
	"math/rand"
	"strings"
)

func main() {
	Run(gmain)
}

//  procedure gmain()
func gmain(env *Env, args ...Value) (Value, *Closure) {

	fmt.Println("testing calls of Go library functions")

	// make Goaldi procedures corresponding to Go library functions
	var Sqrt Value = GoProcedure("Sqrt", math.Sqrt)
	var Max Value = GoProcedure("Max", math.Max)
	var IsNaN Value = GoProcedure("IsNaN", math.IsNaN)
	var NaN Value = GoProcedure("NaN", math.NaN)
	var Intn Value = GoProcedure("Intn", rand.Intn)
	var Seed Value = GoProcedure("Seed", rand.Seed)
	var EqualFold Value = GoProcedure("EqualFold", strings.EqualFold)
	var Replace Value = GoProcedure("Replace", strings.Replace)
	var ToUpper Value = GoProcedure("ToUpper", strings.ToUpper)
	var Trim Value = GoProcedure("Trim", strings.Trim)
	var Print Value = GoProcedure("Print", fmt.Print)
	var Println Value = GoProcedure("Println", fmt.Println)
	var Printf Value = GoProcedure("Printf", fmt.Printf)

	// call them and print results (value,closure)
	fmt.Println(Sqrt, r(Sqrt.(ICall).Call(env, NewNumber(2))))
	fmt.Println(Max, r(Max.(ICall).Call(env, NewNumber(7), NewNumber(5))))
	fmt.Println(NaN, r(NaN.(ICall).Call(env)))
	fmt.Println(IsNaN, r(IsNaN.(ICall).Call(env, NewNumber(33))))
	fmt.Println(IsNaN, r(IsNaN.(ICall).Call(env, ZERO.Div(ZERO))))
	fmt.Println(Intn, r(Intn.(ICall).Call(env, NewNumber(10000))))
	fmt.Println(Intn, r(Intn.(ICall).Call(env, NewNumber(10000))))
	fmt.Println(Seed, r(Seed.(ICall).Call(env, NewNumber(1))))
	fmt.Println(Intn, r(Intn.(ICall).Call(env, NewNumber(10000))))
	fmt.Println(ToUpper, r(ToUpper.(ICall).Call(env, NewString("WasMixed"))))
	fmt.Println(Replace, r(Replace.(ICall).Call(env, // example from GoDoc
		V("oink oink oink"), V("k"), V("ky"), V(2))))
	fmt.Println(Trim, r(Trim.(ICall).Call(env,
		NewString("  a b c  "), NewString(" "))))
	fmt.Println(EqualFold, r(EqualFold.(ICall).Call(env,
		NewString("mixedCase"), NewString("Mixedcase"))))
	Print.(ICall).Call(env, V(11), V(22), V(33))
	Println.(ICall).Call(env, V(11), V(22), V(33))
	Println.(ICall).Call(env, NewString("car"), NewNumber(54))
	Printf.(ICall).Call(env, V("%.3s %6.4f\n"), V("cowboy"), V(3.1415926535))
	return Fail()
}

//  r(Value, Closure) -- return value, ignore closure
func r(v Value, c *Closure) Value {
	return v
}
