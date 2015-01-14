//  fchannel.go -- channel functions and methods

package goaldi

import (
	"reflect"
)

//  Declare methods
var ChannelMethods = MethodTable([]*VProcedure{
	DefMeth(VChannel.Type, "type", "", "return channel type"),
	DefMeth(VChannel.Copy, "copy", "", "return channel value"),
	DefMeth(VChannel.String, "string", "", "return short string"),
	DefMeth(VChannel.GoString, "image", "", "return string image"),
	DefMeth(VChannel.Get, "get", "", "read from channel"),
	DefMeth(VChannel.Put, "put", "x", "send to channel"),
	DefMeth(VChannel.Close, "close", "", "close channel"),
	DefMeth(VChannel.Buffer, "buffer", "n", "create buffer"),
})

//  VChannel.Field implements method calls
func (m VChannel) Field(f string) Value {
	return GetMethod(ChannelMethods, m, f)
}

//  Declare static function
func init() {
	LibProcedure("buffer", Buffer)
}

//  init() declares the constructor function
func init() {
	// Goaldi procedures
	LibProcedure("channel", Channel)
}

//  Declare methods on Go channels
var GoChanMethods = MethodTable([]*VProcedure{
	DefMeth(GoChanGet, "get", "", "read from channel"),
	DefMeth(GoChanPut, "put", "x", "send to channel"),
	DefMeth(GoChanClose, "close", "", "close channel"),
})

//  Channel(i) returns a new channel with buffer size i
func Channel(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("channel", args)
	i := int(ProcArg(args, 0, ZERO).(Numerable).ToNumber().Val())
	return Return(NewChannel(i))
}

//  VChannel.Get() reads the next value from a channel
func (c VChannel) Get(args ...Value) (Value, *Closure) {
	defer Traceback("c.get", args)
	return c.Take(), nil
}

//  GoChanGet(c) returns the next value from a Go channel
func GoChanGet(c Value, args ...Value) (Value, *Closure) {
	defer Traceback("c.get", args)
	return TakeChan(c), nil
}

//  VChannel.Take() implements the unary '@' operator for a Goaldi channel
func (c VChannel) Take() Value {
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

//  VChannel.Put(e...) writes values to a Goaldi channel
func (c VChannel) Put(args ...Value) (Value, *Closure) {
	defer Traceback("c.put", args)
	for _, v := range args {
		c <- v
	}
	return Return(c)
}

//  GoChanPut(c, e...) writes values to a Go channel
func GoChanPut(c Value, args ...Value) (Value, *Closure) {
	defer Traceback("c.put", args)
	for _, v := range args {
		Send(c, v)
	}
	return Return(c)
}

//  Send(x,v) sends value v to the Goaldi or external channel x.
//  It panics on an inappropriate argument type.
func Send(x Value, v Value) Value {
	if c, ok := x.(VChannel); ok { // if a Goaldi channel
		c <- v // no conversion or reflection needed
		return v
	}
	cv := reflect.ValueOf(x)
	if cv.Kind() != reflect.Chan {
		panic(NewExn("Not a channel", x))
	}
	cv.Send(reflect.ValueOf(Export(v)))
	return v
}

//  VChannel.Close() closes the channel
func (c VChannel) Close(args ...Value) (Value, *Closure) {
	defer Traceback("c.close", args)
	close(c)
	return Return(c)
}

//  GoChanClose(c) closes a Go channel
func GoChanClose(c Value, args ...Value) (Value, *Closure) {
	defer Traceback("c.close", args)
	reflect.ValueOf(c).Close()
	return Return(c)
}

//  VChannel.Buffer(i) interposes a buffer of size i in front of a channel.
func (c VChannel) Buffer(args ...Value) (Value, *Closure) {
	defer Traceback("c.buffer", args)
	i := int(ProcArg(args, 0, ONE).(Numerable).ToNumber().Val())
	r := NewChannel(i)
	go func() {
		for {
			v, ok := <-c // get value from input side
			if !ok {     // if input channel was closed
				close(r) // then close output channel
				return   // and return (killing this thread)
			}
			if CoSend(r, v) == nil { // send output; if closed,
				return // return (killing this thread)
			}
		}
	}()
	return Return(r)
}

//  Buffer(i, C) is the static version of c.Buffer(i).
//  This is useful in the Goaldi form buffer(i, create e).
func Buffer(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("buffer", args)
	i := ProcArg(args, 0, ONE)
	c := ProcArg(args, 1, NilValue)
	return c.(VChannel).Buffer(i)
}
