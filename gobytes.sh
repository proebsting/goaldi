#!/bin/sh
#
#  gobytes -- turn binary data into a byte array for embedding in Go code.
#
#  usage:  gobytes pkgname varname <binaryfile >file.go

echo package ${1-main}
echo var ${2-data} ' = []byte{'
od -v -An -tu1 | sed 's/  *\([0-9][0-9]*\)/\1,/g'
echo '}'
