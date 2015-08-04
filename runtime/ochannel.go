//  ochannel.go -- operations on channels

package runtime

import (
	"fmt"
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

//  GetChan(c) receives and imports a value from a Goaldi or external channel,
//  or fails if no value is available.
func GetChan(c interface{} /*anychan*/) (Value, *Closure) {
	cv := reflect.ValueOf(c)
	cases := []reflect.SelectCase{
		reflect.SelectCase{Dir: reflect.SelectDefault},
		reflect.SelectCase{Dir: reflect.SelectRecv, Chan: cv},
	}
	i, v, ok := reflect.Select(cases)
	if i > 0 && ok {
		return Return(Import(v.Interface())) // got a value
	} else {
		return Fail() // no value available
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

//  A Selector struct implements a select statement.
type Selector struct {
	cases   []reflect.SelectCase // cases for reflect.Select
	nCases  int                  // number of cases expected
	defCase int                  // index of default case, if any
}

//  NewSelector(n) creates a selector to hold n select cases.
func NewSelector(n int) *Selector {
	return &Selector{make([]reflect.SelectCase, 0, n+1), n, -1}
}

//  Selector.SendCase(ch, x) adds a "send" case.
func (s *Selector) SendCase(ch Value, x Value) {
	if _, ok := ch.(VChannel); !ok {
		// not a Goaldi channel; convert data value to best Go type
		x = Export(x)
	}
	s.cases = append(s.cases, reflect.SelectCase{
		Dir:  reflect.SelectSend,
		Chan: channelValue(ch),
		Send: reflect.ValueOf(x)})
}

//  Selector.RecvCase(ch) adds a "receive" case.
func (s *Selector) RecvCase(ch Value) {
	s.cases = append(s.cases, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: channelValue(ch)})
}

//  Selector.DfltCase() adds a "default" case.
func (s *Selector) DfltCase() {
	s.defCase = len(s.cases)
	s.cases = append(s.cases, reflect.SelectCase{
		Dir: reflect.SelectDefault})
}

//  Selector.Execute() runs the select loop.
//  It returns the selected index and also, for receive, the associated value.
//  It returns (-1,nil) to fail.
func (s *Selector) Execute() (int, Value) {
	if len(s.cases) != s.nCases {
		panic(Malfunction(fmt.Sprintf(
			"Expected %d cases but got %d", s.nCases, len(s.cases))))
	}
	if s.defCase < 0 { // if we need to add a default case
		s.cases = append(s.cases,
			reflect.SelectCase{Dir: reflect.SelectDefault})
	}
	// repeat until we get anything other than a read on a closed channel
	for {
		// call select through the reflection interface
		i, v, recvOK := reflect.Select(s.cases)
		// select has returned, having chosen case i
		if i == s.nCases {
			// this is the default case we added, because there was none
			return -1, nil // so the select expression fails
		}
		chosen := s.cases[i]
		if chosen.Dir != reflect.SelectRecv {
			// send or default chosen
			return i, nil
		} else if recvOK {
			// return choice number and received value
			return i, Import(v.Interface())
		} else {
			// a closed channel was selected
			s.cases[i].Chan = hungChannel // disable this case and retry
		}
	}
}

//  used for disabling one branch of a select
var hungChannel = reflect.ValueOf(make(chan interface{}))

//  get and validate a channel value, returning a reflect.Value
func channelValue(ch Value) reflect.Value {
	cv := reflect.ValueOf(ch)
	if cv.Kind() != reflect.Chan {
		panic(NewExn("Not a channel", ch))
	}
	return cv
}
