//  xhash.go -- hashing interface *** SAMPLE EXTENSION ***
//
//  adler32(), crc32(), fnv32(), and fnv32a() each create a new hashing engine.
//  The return value is a file, and data written to the file updates the
//  running checksum.
//
//  Given one of these special files,
//  	hashvalue(f)  returns a numeric hash value of the data written so far.

package goaldi

import (
	"hash"
	"hash/adler32"
	"hash/crc32"
	"hash/fnv"
)

//  declare new procedures for use from Goaldi
func init() {
	LibGoFunc("adler32", adler32.New)
	LibGoFunc("crc32", crc32.NewIEEE)
	LibGoFunc("fnv32", fnv.New32)
	LibGoFunc("fnv32a", fnv.New32a)
	LibGoFunc("hashvalue", hashvalue)
}

//  hashvalue(f) returns the current value of the hash engine f.
func hashvalue(f *VFile) uint32 {
	return f.Writer.(hash.Hash32).Sum32()
}
