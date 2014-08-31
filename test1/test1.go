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

	i := V(23)
	s := V("45.0")
	n := V(nil)
	f.Println(i, s, n)
	f.Println(i.AsString(), i.AsNumber())
	f.Println(s.AsString(), s.AsNumber())

	f.Print("22+33: ")
	f.Println(V(22).Add(V(33)))
	f.Print("7*11: ")
	f.Println(V(7).Mult(V(11)))
	f.Print("strings: ")
	f.Println(V("19").Mult(V("3")))

	f.Print("boom: ")
	f.Println(V(31).Add(V(nil)))

	return Fail()
}
