package main

import (
	f "fmt"
	. "goaldi"
)

func main() {
	Run(gmain)
}

//  procedure gmain()
func gmain(args []Value) (Value, *Closure) {

	//#%#% this code doesn't check for thrown exceptions or even failures

	var a Value = NewNumber(3); ta := Trapped(&a)
	var b Value = NewString("5"); tb := Trapped(&b)
	var c Value = NewNil(); tc := Trapped(&c)
	f.Println(a, ta, b, tb, c, tc)
	d, _ := ta.Add(tb)
	f.Println(d)
	tc.Assign(NewNumber(7.3))
	f.Println(c)

	return Fail()
}
