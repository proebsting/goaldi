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
	GoLib(adler32.New, "adler32", "", "create Adler-32 checksum engine")
	GoLib(crc32.NewIEEE, "crc32", "", "create IEEE CRC-32 checksum engine")
	GoLib(fnv.New32, "fnv32", "", "create 32-bit FNV-1 checksum engine")
	GoLib(fnv.New32a, "fnv32a", "", "create 32-bit FNV-1a checksum engine")
	GoLib(hashvalue, "hashvalue", "", "return current value of checksum engine")
}

//  hashvalue(f) returns the current value of the hash engine f.
func hashvalue(f *VFile) uint32 {
	return f.Writer.(hash.Hash32).Sum32()
}
