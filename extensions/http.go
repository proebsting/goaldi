//  http.go -- HTTP interface extension to Goaldi
//
//  This extension adds a few random operations for dealing with HTTP files.
//  It is illustrative but incomplete.
//
//  htfile() opens a URL for reading as a file.
//  htget() is similar but returns a pointer to an http.Response struct.
//  This can be queried to read the headers as well as the body.
//  htpost() posts a form and returns the response struct for inspection.

package extensions

import (
	g "github.com/proebsting/goaldi/runtime"
	"io"
	"net/http"
	"net/url"
)

//  declare new procedures for use from Goaldi
func init() {
	g.GoLib(htfile, "htfile", "url", "open URL and return file")
	g.GoLib(htget, "htget", "url", "get URL and return response")
	g.GoLib(htpost, "htpost", "url,name,kv[]", "post form and return response")
}

//  htfile(url) returns a file for reading the body of a web file.
//  It returns nil if the URL cannot be opened.
func htfile(u string) io.Reader {
	resp, err := http.Get(u)
	if err != nil {
		return nil
	}
	return resp.Body
}

//  htget(url) returns an HTTP response object R for reading a web file.
//  It returns nil if the url cannot be opened.
//
//  Given an object R returned by htget():
//  	The operation *R returns a file for reading the body.
//  	The operation !R generates the header files.
//
//  (These use Goaldi operators instead of method calls just to show how.)
func htget(u string) *htresp {
	resp, err := http.Get(u)
	if err != nil {
		return nil
	}
	return &htresp{resp}
}

//  htpost(url, k1, v1, k2, v2, ...) posts a form and returns response R.
//  The (ki, vi) arguments are key-value pairs for supplying parameters.
//  htpost() returns nil if the url cannot be opened.
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

//  htresp.Size() returns the underlying file for reading.
//  This is a perversion of the *H unary operator for illustrative purposes.
func (h *htresp) Size() g.Value {
	return g.Import(h.Resp.Body) // Import() converts Go value to Goaldi value
}

//  htresp.Dispense() generates the headers as name:value strings.
//  This implements the !H unary operator.
func (h *htresp) Dispense(unused g.Value) (g.Value, *g.Closure) {
	//  range over the headers and feed into a channel
	ch := make(chan *g.VString)
	go func() {
		for k, l := range h.Resp.Header {
			for _, v := range l {
				ch <- g.NewString(k + ":" + v)
			}
		}
		close(ch)
	}()
	//  read from the channel in a suspendable closure
	var f *g.Closure
	f = &g.Closure{func() (g.Value, *g.Closure) {
		v := <-ch
		if v == nil {
			return nil, nil
		} else {
			return v, f
		}
	}}
	return f.Resume()
}
