//  ochannel.go -- operations on channels

package runtime

import (
	"reflect"
)

//  VChannel.Take(lval) implements the unary '@' operator for a Goaldi channel.
func (c VChannel) Take(lval Value) Value {
	v, ok := <-c
	if ok {
		return v // got a value
	} else {
		return nil // fail: channel was closed
	}
}

//  TakeChan(c) receives and imports a value from a Goaldi or external channel
func TakeChan(c interface{} /*anychan*/) Value {
	v, ok := reflect.ValueOf(c).Recv()
	if ok {
		return Import(v.Interface()) // got a value
	} else {
		return nil // fail: channel was closed
	}
}

//  DispenseChan(c) implements @c for a Goaldi or external channel
func DispenseChan(c interface{} /*anychan*/) (Value, *Closure) {
	var f *Closure
	f = &Closure{func() (Value, *Closure) {
		v := TakeChan(c)
		if v != nil {
			return v, f
		} else {
			return Fail()
		}
	}}
	return f.Resume()
}

//  VChannel.Send(v) implements the '@:' operator for a Goaldi channel.
func (c VChannel) Send(lval Value, v Value) Value {
	c <- v
	return v
}
