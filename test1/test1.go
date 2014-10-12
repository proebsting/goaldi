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

	//#%#% this code doesn't check for thrown exceptions or even failures

	f.Println("testing construction, imaging, a few operators")

	f.Print("22+33: ")
	f.Println(V(22).(IAdd).Add(V(33)))
	f.Print("7*11: ")
	f.Println(NewNumber(7).Mul(NewNumber(11)))
	f.Print("strings: ")
	f.Println(NewString("19").Mul(NewString("3").ToNumber()))

	return Fail()
}
