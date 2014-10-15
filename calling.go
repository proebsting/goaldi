//  calling.go -- procedure call / suspension / resumption interface
//
//  In the Go implementation, a Goaldi procedure or operation result
//  is the pair (Value, *Closure) with this meaning:
//
//  Value  *Closure   Interpretation
//  -----  -------   ---------------------------------
//   nil     nil     procedure failed
//  value    nil     procedure returned a value
//  value   resume   procedure suspended and can be resumed

package goaldi

//  execution environment
type Env struct {
	//#%#% to be determined
	//#%#% dynamic variables?
	//#%#% thread id for monitoring?
}

type Closure struct {
	Go Resumer // start or resumption function, depending on context
}

//  Procedure function prototype
type Procedure func(env *Env, args ...Value) (Value, *Closure)

//  Procedure resumption prototype
type Resumer func() (Value, *Closure)

//  Resume() executes the entry point in a Closure to produce the next result.
//  If the pointer is nil, failure is produced.
//  But watch out for nils that aren't.
func (c *Closure) Resume() (Value, *Closure) {
	if c == nil {
		return Fail()
	}
	return c.Go()
}

//  Fail returns a failure indicator
func Fail() (Value, *Closure) {
	return nil, nil
}

//  Return returns a simple value as a duo
func Return(v Value) (Value, *Closure) {
	return v, nil
}
