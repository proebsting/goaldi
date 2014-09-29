#!/bin/sh
#
#  goaldi [options] file [arg...] -- compile and execute Goaldi program
#
#  Options -l -t -v -A -J -T are passed along to the interpreter.
#
#  Assumes that gtran and gexec are in the search path.

XFLAGS=ltvAJT
USAGE="usage: $0 [-$XFLAGS] file [arg...]"

#  process options
XOPTS=
while getopts $XFLAGS C; do
    case $C in
	[ltvAJT]) XOPTS="$XOPTS -$C";;
    ?)
	echo 1>&2 $USAGE; exit 1;;
    esac
done
shift $(($OPTIND - 1))
test $# -lt 1 && echo 1>&2 $USAGE && exit 1

I=$1
shift

export COEXPSIZE=300000
exec gtran preproc $I : yylex : parse : ast2ir : optim -O : json_File : stdout |
  gexec $XOPTS "$@"
