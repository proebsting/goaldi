package main

import (
	f "fmt"
	. "goaldi"
)

func main() {
	Run(gmain)
}

//  procedure gmain()
func gmain(env *Env, args ...Value) (Value, *Closure) {

	//#%#% this code doesn't check for thrown exceptions or even failures

	f.Println("testing construction, imaging, a few operators")

	i := NewNumber(23)
	s := NewString("45.0")
	n := NewNil()
	p := NewProcedure("gmain", gmain)
	f.Printf("%v %v %v\n", i, Type(i), Image(i))
	f.Printf("%v %v %v\n", s, Type(s), Image(s))
	f.Printf("%v %v %v\n", n, Type(n), Image(n))
	f.Printf("%v %v %v\n", p, Type(p), Image(p))
	f.Println(i.ToString())
	f.Println(s.ToNumber())

	f.Print("22+33: ")
	f.Println(V(22).(IAdd).Add(V(33)))
	f.Print("7*11: ")
	f.Println(NewNumber(7).Mul(NewNumber(11)))
	f.Print("strings: ")
	f.Println(NewString("19").Mul(NewString("3").ToNumber()))

	return Fail()
}
