//  prb.go -- parallel buffer implementation for concurrent input

package ir

import (
	"bytes"
	"io"
	"time"
)

type prBuffer struct { // ParallelReader buffer
	rdr io.Reader     // underlying Reader
	buf *bytes.Buffer // memory buffer
	avb chan int      // available bytes & buffer mutex
}

const BLINK = 1 * time.Millisecond // sleep time when blocked (empty / EOF)

//  A ParallelReader spawns a separate goroutine to fill a hidden buffer,
//  allowing the underlying reader (e.g. a decompressor) to run concurrently.
//  Any error from the underlying reader causes a panic.
func ParallelReader(rdr io.Reader, bufsize int) io.ReadCloser {
	prb := &prBuffer{rdr, bytes.NewBuffer(make([]byte, 0, bufsize)),
		make(chan int, 1)}
	go prb.filler(bufsize)
	return prb
}

//  prBuffer.filler runs in the background to fill the hidden buffer.
func (prb *prBuffer) filler(bufsize int) {
	rbuffer := make([]byte, bufsize)
	first := true // first-time flag
	for {
		n, err := prb.rdr.Read(rbuffer) // read some data
		if err == io.EOF {
			break // break on EOF
		}
		if err != nil {
			panic(err) // throw error
		}
		if n == 0 {
			continue // retry empty read
		}
		// we have read n > 0 bytes; copy to shared buffer
		if first {
			first = false // clear first-time flag
		} else {
			<-prb.avb // get interlock
		}
		prb.buf.Write(rbuffer[:n]) // write to shared buffer
		prb.avb <- prb.buf.Len()   // release with count
	}
	// EOF was read
	for {
		<-prb.avb          // get interlock
		n := prb.buf.Len() // is buffer empty?
		if n == 0 {
			prb.avb <- -1 // signal EOF
			return        // and exit
		}
		prb.avb <- prb.buf.Len() // show data available
		time.Sleep(BLINK)
	}
}

//  prBuffer.Read gets data from the hidden buffer.
func (prb *prBuffer) Read(p []byte) (int, error) {
	n := <-prb.avb // get interlock
	for n == 0 {   // while buffer is empty
		prb.avb <- n // release interlock
		time.Sleep(BLINK)
		n = <-prb.avb // reclaim interlock
	}
	if n < 0 { // if EOF
		prb.avb <- n     // release interlock
		return 0, io.EOF // return EOF indication
	}
	if n < len(p) { // if less available than supplied buffer
		p = p[:n] // shorten buffer
	}
	n, _ = prb.buf.Read(p)   // get available data (will not block)
	prb.avb <- prb.buf.Len() // release interlock with updated count
	return n, nil            // return results
}

//  prBuffer.Close closes the underlying Reader, if possible, and frees memory.
func (prb *prBuffer) Close() error {
	r := prb.rdr                    // save underlying Reader
	*prb = prBuffer{}               // zero struct to clear references
	if c, ok := r.(io.Closer); ok { // if the Reader has a Closer
		return c.Close()
	} else {
		return nil
	}
}
