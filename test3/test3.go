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
func gmain(args ...Value) (Value, *Closure) {

	fmt.Println("testing calls of Go library functions")

	// make Goaldi procedures corresponding to Go library functions
	var Sqrt Value = GoProcedure("Sqrt", math.Sqrt)
	var Intn Value = GoProcedure("Intn", rand.Intn)
	var ToUpper Value = GoProcedure("ToUpper", strings.ToUpper)

	// call them and print results (value,closure)
	fmt.Print(Sqrt, ": ")
	fmt.Println(Sqrt.(ICall).Call(NewNumber(2)))
	fmt.Print(Intn, ": ")
	fmt.Println(Intn.(ICall).Call(NewNumber(10000)))
	fmt.Print(ToUpper, ": ")
	fmt.Println(ToUpper.(ICall).Call(NewString("WasMixed")))

	return Fail()
}
