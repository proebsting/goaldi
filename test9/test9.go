package main

import (
	f "fmt"
	g "goaldi"
)

// text implementation of a Value struct outside the Goaldi package

type MyValue struct {
	value int
}

func (v *MyValue) String() string {
	return f.Sprintf("my(%d)", v.value)
}

func main() {
	g.Run(gmain)
}

//  procedure gmain()
func gmain(args []g.Value) (g.Value, *g.Closure) {

	//#%#% this code doesn't check for thrown exceptions or even failures
	var v g.Value = &MyValue{}
	f.Println(v)
	v.(*MyValue).value = 19
	f.Println(v)
	f.Println("Expect PANIC:")
	v.(g.IAdd).Add(g.NewNumber(23)) // Add not impl by MyValue
	return g.Fail()
}
