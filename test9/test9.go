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
	return f.Sprintf("MyValue(%d)", v.value)
}

func main() {
	g.Run(g.NewProcedure("main", gmain), nil)
}

//  procedure gmain()
func gmain(env *g.Env, args ...g.Value) (g.Value, *g.Closure) {

	// boilerplate prologue
	var ev interface{}
	var ln = "100"
	defer func() {
		if p := recover(); p != nil {
			panic(g.Catch(p, ev, "test.gdi", ln, "gmain", args))
		}
	}()

	f.Println("testing panic, offending value, traceback")
	var p g.Value = g.NewProcedure("gsubr", gsubr)
	return p.(g.ICall).Call(env, g.NewNumber(23), g.NewString("skidoo"))
}

//  procedure gsubr()
func gsubr(env *g.Env, args ...g.Value) (g.Value, *g.Closure) {

	// boilerplate prologue
	var ev interface{}
	var ln = "200"
	defer func() {
		if p := recover(); p != nil {
			panic(g.Catch(p, ev, "test.gdi", ln, "gsubr", args))
		}
	}()

	var v g.Value = &MyValue{}
	f.Println(v)
	v.(*MyValue).value = 19
	f.Println(v)
	f.Println("Expect PANIC:")
	ln = "222"
	ev = v
	v.(g.IAdd).Add(g.NewNumber(23)) // Add not impl by MyValue
	return g.Fail()
}
