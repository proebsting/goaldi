//  fchannel.go -- channel functions and methods

package goaldi

//  Declare methods
var ChannelMethods = map[string]interface{}{
	"type":  VChannel.Type,
	"copy":  VChannel.Copy,
	"image": VChannel.GoString,
	"get":   VChannel.Get,
	"put":   VChannel.Put,
	"close": VChannel.Close,
}

//  VChannel.Field implements method calls
func (m VChannel) Field(f string) Value {
	return GetMethod(ChannelMethods, m, f)
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

//  VChannel.Take() implements the unary '@' operator
func (c VChannel) Take() Value {
	return <-c
}
