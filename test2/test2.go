package main

import (
	f "fmt"
	. "goaldi"
)

func main() {
	Run(NewProcedure("main", gmain), nil)
}

//  procedure gmain()
func gmain(env *Env, args ...Value) (Value, *Closure) {

	f.Println("testing string<->number and assignment operators")
	var a Value = NewNumber(3)
	ta := Trapped(&a)
	var b Value = NewString("5")
	tb := Trapped(&b)
	var c *Closure
	f.Println("sums:  ", a, b, a.(IAdd).Add(b), b.(IAdd).Add(a))
	f.Printf("begin:  a=%v b=%v\n", a, b)
	Swap(ta, tb)
	f.Printf("swap:   a=%v b=%v\n", a, b)
	_, c = RevSwap(ta, tb)
	f.Printf("rswap:  a=%v b=%v\n", a, b)
	_, c = c.Resume()
	f.Printf("resume: a=%v b=%v\n", a, b)
	_, c = c.Resume()
	f.Printf("resume: a=%v b=%v\n", a, b)

	return Fail()
}
