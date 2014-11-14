//  fchannel.go -- channel functions and methods

package goaldi

//  Declare methods
var ChannelMethods = map[string]interface{}{
	"type":   VChannel.Type,
	"copy":   VChannel.Copy,
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
	v := <-c
	if v == nil { // if closed
		return Fail()
	} else {
		return Return(v)
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
func (c VChannel) Buffer(args ...Value) (Value, *Closure) {
	defer Traceback("C.buffer", args)
	i := int(ProcArg(args, 0, ONE).(Numerable).ToNumber().Val())
	r := NewChannel(i)
	go func() {
		for {
			r <- <-c
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

//  VChannel.Take() implements the unary '@' operator
func (c VChannel) Take() Value {
	return <-c
}
