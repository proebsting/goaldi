//  hash.go -- hashing interface extension to Goaldi
//
//  This straightforward extension adds 32-bit hashing functions to Goaldi.
//
//  adler32(), crc32(), fnv32(), and fnv32a() each create a new hashing engine.
//  The returned value is a file, and data written to the file updates the
//  running checksum.
//
//  Given one of these special files,
//  	hashvalue(f)  returns a numeric hash value of the data written so far.

package extensions

import (
	g "github.com/proebsting/goaldi/runtime"
	"hash"
	"hash/adler32"
	"hash/crc32"
	"hash/fnv"
)

// declare new procedures for use from Goaldi
func init() {
	g.GoLib(adler32.New, "adler32", "", "create Adler-32 checksum engine")
	g.GoLib(crc32.NewIEEE, "crc32", "", "create IEEE CRC-32 checksum engine")
	g.GoLib(fnv.New32, "fnv32", "", "create 32-bit FNV-1 checksum engine")
	g.GoLib(fnv.New32a, "fnv32a", "", "create 32-bit FNV-1a checksum engine")
	g.GoLib(hashvalue, "hashvalue", "", "return accumulated checksum value")
}

// hashvalue(f) returns the current value of the hash engine f.
func hashvalue(f hash.Hash32) uint32 {
	return f.Sum32()
}
