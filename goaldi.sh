#!/bin/sh
#
#  goaldi [options] file [arg...] -- compile and execute Goaldi program
#
#  -c	compile only, producing IR on file.gir (interpreter options ignored)
#
#  Options -l -t -v -A -J -P -T are passed along to the interpreter.
#
#  Assumes that gtran and gexec are in the search path.

FLAGS=cltvAJPT
USAGE="usage: $0 [-$FLAGS] file [arg...]"
TMP=/tmp/gdi.$$.gir

#  process options
XOPTS=
CFLAG=
while getopts $FLAGS C; do
    case $C in
	c)		CFLAG=$C;;
	[ltvAJPT])	XOPTS="$XOPTS -$C";;
	?)		echo 1>&2 $USAGE; exit 1;;
    esac
done
shift $(($OPTIND - 1))
test $# -lt 1 && echo 1>&2 $USAGE && exit 1

I=$1
TRAN="gtran preproc $I : yylex : parse : ast2ir : optim -O : json_File : stdout"
shift

export COEXPSIZE=300000
if [ -n "$CFLAG" ]; then
    exec $TRAN >${I%.*}.gir
fi

trap 'rm -f $TMP; exit' 0 1 2 15
$TRAN >$TMP && exec gexec $XOPTS $TMP "$@"
