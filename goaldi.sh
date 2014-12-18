#!/bin/sh
#
#  goaldi [options] file [arg...] -- compile and execute Goaldi program
#
#  -c	compile only, producing IR on file.gir (interpreter options ignored)
#  -d	compile only, producing Dot directives on file.dot
#  -N	no optimization
#
#  Options -l -t -v -A -D -F -J -P -T are passed along to the interpreter.
#
#  Assumes that gtran and gexec are in the search path.

FLAGS=cdNltvADFJPT
USAGE="usage: $0 [-$FLAGS] file [arg...]"
TMP=/tmp/gdi.$$.gir

#  process options
XOPTS=
CFLAG=
DFLAG=
OPT=": optim -O"
while getopts $FLAGS C; do
    case $C in
	c)			CFLAG=$C;;
	d)			DFLAG=$C;;
	N)			OPT="";;
	[ltvADFJPT])	XOPTS="$XOPTS -$C";;
	?)			echo 1>&2 $USAGE; exit 1;;
    esac
done
shift $(($OPTIND - 1))
test $# -lt 1 && echo 1>&2 $USAGE && exit 1

I=$1
DOT="gtran preproc $I : yylex : parse : ast2ir $OPT : dot_File : stdout"
TRAN="gtran preproc $I : yylex : parse : ast2ir $OPT : json_File : stdout"
shift

export COEXPSIZE=300000

if [ -n "$DFLAG" ]; then	# -d: produce file.dot, and quit
    exec $DOT >${I%.*}.dot
fi

if [ -n "$CFLAG" ]; then	# -c: produce file.gir, and quit
    exec $TRAN >${I%.*}.gir
fi

#  translate and execute
trap 'rm -f $TMP; exit' 0 1 2 15
$TRAN >$TMP && exec gexec $XOPTS $TMP "$@"
