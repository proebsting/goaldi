//  fchannel.go -- channel functions and methods

package runtime

import (
	"reflect"
)

//  Declare methods
var ChannelMethods = MethodTable([]*VProcedure{
	DefMeth(VChannel.Get, "get", "", "read from channel"),
	DefMeth(VChannel.Put, "put", "x", "send to channel"),
	DefMeth(VChannel.Close, "close", "", "close channel"),
	DefMeth(VChannel.Buffer, "buffer", "size", "interpose buffer"),
})

//  Declare static functions
func init() {
	DefLib(Buffer, "buffer", "size,c", "interpose buffer before channel")
}

//  Declare methods on Go channels
var GoChanMethods = MethodTable([]*VProcedure{
	DefMeth(GoChanGet, "get", "", "read from channel"),
	DefMeth(GoChanPut, "put", "x", "send to channel"),
	DefMeth(GoChanClose, "close", "", "close channel"),
})

//  channel(size) creates and returns a new channel with the given buffer size.
func Channel(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("channel", args)
	i := int(ProcArg(args, 0, ZERO).(Numerable).ToNumber().Val())
	return Return(NewChannel(i))
}

//  c.get() reads the next value from channel c.
func (c VChannel) Get(args ...Value) (Value, *Closure) {
	defer Traceback("c.get", args)
	return c.Take(nil), nil
}

//  GoChanGet(c) returns the next value from a Go channel
func GoChanGet(c Value, args ...Value) (Value, *Closure) {
	defer Traceback("c.get", args)
	return TakeChan(c), nil
}

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

//  c.put(e...) writes its argument values, in order, to channel c.
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
		GoChanSend(c, v)
	}
	return Return(c)
}

//  VChannel.Send(v) implements the '@:' operator for a Goaldi channel.
func (c VChannel) Send(lval Value, v Value) Value {
	c <- v
	return v
}

//  GoChanSend(x,v) sends value v to the Goaldi channel x.
//  It panics on an inappropriate argument type.
func GoChanSend(x Value, v Value) Value {
	cv := reflect.ValueOf(x)
	if cv.Kind() != reflect.Chan {
		panic(NewExn("Not a channel", x))
	}
	cv.Send(reflect.ValueOf(Export(v)))
	return v
}

//  c.close() closes the channel c.
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

//  c.buffer(size) returns a channel that interposes a buffer of the given size
//  before the channel c.
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

//  buffer(size, c) returns a channel that interposes a buffer of the given size
//  before the channel c.
//  This is useful in the Goaldi form buffer(size, create e)
//  to provide buffering of the results produced by an asynchronous thread.
func Buffer(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("buffer", args)
	i := ProcArg(args, 0, ONE)
	c := ProcArg(args, 1, NilValue)
	return c.(VChannel).Buffer(i)
}
