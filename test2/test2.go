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

	var a Value = NewNumber(3)
	ta := Trapped(&a)
	var b Value = NewString("5")
	tb := Trapped(&b)
	var c Value = NewNil()
	tc := Trapped(&c)
	f.Println(a, ta, b, tb, c, tc)
	av := ta.Deref()
	bv := tb.Deref()
	d := av.(IAdd).Add(bv)
	f.Println(d)
	e := bv.(IAdd).Add(av)
	f.Println(e)
	tc.Assign(NewNumber(7.3))
	f.Println(c)

	return Fail()
}
