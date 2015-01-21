//  vchannel.go -- VChannel, the Goaldi type "channel"

package goaldi

import (
	"fmt"
	"strings"
)

//  VChannel implements a Goaldi channel, which just wraps a Go channel.
type VChannel chan Value

//  NewChannel -- construct a new Goaldi channel
func NewChannel(i int) VChannel {
	return VChannel(make(chan Value, i))
}

//  ChannelType is the channel instance of type type.
var ChannelType = NewType("c", rChannel, Channel, ChannelMethods,
	"channel", "size", "create channel")

//  VChannel.String -- default conversion to Go string returns "c:size"
func (c VChannel) String() string {
	return fmt.Sprintf("c:%d", cap(c))
}

//  VChannel.GoString -- convert to Go string for image() and printf("%#v")
func (c VChannel) GoString() string {
	return fmt.Sprintf("channel(%d)", cap(c))
}

//  VChannel.Type -- return the channel type
func (c VChannel) Type() IRank {
	return ChannelType
}

//  VChannel.Copy returns itself
func (c VChannel) Copy() Value {
	return c
}

//  VChannel.Import returns itself
func (v VChannel) Import() Value {
	return v
}

//  VChannel.Export returns itself.
func (v VChannel) Export() interface{} {
	return v
}

//  CoSend(chan, value) sends a co-expression result to a channel.
//  Returns chan if successful, nil if channel had been closed.
//  Panics on any other error.
func CoSend(ch VChannel, v Value) VChannel {
	result := ch
	defer func() {
		r := recover()
		if r != nil {
			result = nil
			if !strings.HasSuffix(fmt.Sprint(r), "send on closed channel") {
				panic(r) // not what we expected
			}
		}
	}()
	ch <- v
	return result
}
