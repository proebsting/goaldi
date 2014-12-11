//  fchannel.go -- channel functions and methods

package goaldi

import (
	"reflect"
)

//  Declare methods
var ChannelMethods = map[string]interface{}{
	"type":   VChannel.Type,
	"copy":   VChannel.Copy,
	"string": VChannel.String,
	"image":  VChannel.GoString,
	"get":    VChannel.Get,
	"put":    VChannel.Put,
	"close":  VChannel.Close,
	"buffer": VChannel.Buffer,
}

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

//  Channel(i) returns a new channel with buffer size i
func Channel(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("channel", args)
	i := int(ProcArg(args, 0, ZERO).(Numerable).ToNumber().Val())
	return Return(NewChannel(i))
}

//  VChannel.Get() reads the next value from a channel
func (c VChannel) Get(args ...Value) (Value, *Closure) {
	defer Traceback("C.get", args)
	v, ok := <-c
	if ok {
		return Return(v) // got a value
	} else {
		return Fail() // channel was closed
	}
}

//  VChannel.Put(e...) writes values to a channel
func (c VChannel) Put(args ...Value) (Value, *Closure) {
	defer Traceback("C.put", args)
	for _, v := range args {
		c <- v
	}
	return Return(c)
}

//  VChannel.Close() closes the channel
func (c VChannel) Close(args ...Value) (Value, *Closure) {
	defer Traceback("C.close", args)
	close(c)
	return Return(c)
}

//  VChannel.Buffer(i) interposes a buffer of size i in front of a channel.
//  NOTE: a Go nil is indistinguishable from EOF and is treated as such.
//  (This only applies to Go values.  There is no problem with Goaldi nil.)
func (c VChannel) Buffer(args ...Value) (Value, *Closure) {
	defer Traceback("C.buffer", args)
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

//  Buffer(i, C) is the static version of C.Buffer(i).
//  This is useful in the Goaldi form buffer(i, create e).
func Buffer(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("buffer", args)
	i := ProcArg(args, 0, ONE)
	c := ProcArg(args, 1, NilValue)
	return c.(VChannel).Buffer(i)
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

//  TakeChan(c) receives and imports a value from an external channel
func TakeChan(c interface{} /*anychan*/) Value {
	v, ok := reflect.ValueOf(c).Recv()
	if ok {
		return Import(v) // got a value
	} else {
		return nil // fail: channel was closed
	}
}
