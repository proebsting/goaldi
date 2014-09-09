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

	i := NewNumber(23)
	s := NewString("45.0")
	n := NewNil()
	f.Printf("%v %#v\n", i, i)
	f.Printf("%v %#v\n", s, s)
	f.Printf("%v %#v\n", n, n)
	f.Println(i.ToString())
	f.Println(s.ToNumber())

	f.Print("22+33: ")
	f.Println(V(22).(IMath).Add(V(33).(IMath)))
	f.Print("7*11: ")
	f.Println(NewNumber(7).Mult(NewNumber(11)))
	f.Print("strings: ")
	f.Println(NewString("19").Mult(NewString("3").ToNumber()))

	return Fail()
}
