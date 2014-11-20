//  xhttp.go -- HTTP interface *** SAMPLE EXTENSION ***

package goaldi

import (
	"io"
	"net/http"
	"net/url"
)

//  declare new procedures for use from Goaldi
func init() {
	LibGoFunc("htopen", htopen)
	LibGoFunc("htget", htget)
	LibGoFunc("htpost", htpost)
}

//  htopen(url) returns a file for reading the body of a web file
func htopen(u string) io.Reader {
	resp, err := http.Get(u)
	if err != nil {
		return nil
	}
	return resp.Body
}

//  htget(url) returns an HTTP response R for reading a web file.
//  Given this object R:
//  	The operation *R returns a file for reading the body.
func htget(u string) *htresp {
	resp, err := http.Get(u)
	if err != nil {
		return nil
	}
	return &htresp{resp}
}

//  htpost(url, name, k1, v1, k2, v2, ...) posts a form and returns response R.
func htpost(u string, kv ...string) *htresp {
	data := url.Values{}
	for i := 0; i < len(kv); i += 2 {
		data.Add(kv[i], kv[i+1])
	}
	resp, err := http.PostForm(u, data)
	if err != nil {
		return nil
	}
	return &htresp{resp}
}

//  htresp is our version of an http.Response with added Goaldi methods
type htresp struct {
	Resp *http.Response
}

//  htresp.Size() returns the underlying file.  (*H: unary operator abuse.)
func (h *htresp) Size() Value {
	return Import(h.Resp.Body) // Import() converts Go value to Goaldi value
}

//  htresp.Dispense() generates the headers as name:value.  (!H unary operator.)
func (h *htresp) Dispense(unused IVariable) (Value, *Closure) {
	//  range over the headers and feed into a channel
	ch := make(chan *VString)
	go func() {
		for k, l := range h.Resp.Header {
			for _, v := range l {
				ch <- NewString(k + ":" + v)
			}
		}
		close(ch)
	}()
	//  read from the channel in a suspendable closure
	var f *Closure
	f = &Closure{func() (Value, *Closure) {
		v := <-ch
		if v == nil {
			return nil, nil
		} else {
			return v, f
		}
	}}
	return f.Resume()
}
