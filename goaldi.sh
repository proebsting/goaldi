#!/bin/sh
#
#  goaldi [options] file [arg...] -- compile and execute Goaldi program
#
#  To see options, run with no arguments.
#  Assumes that gtran and gexec are in the search path.

FLAGS=acdNltvADEFJPT
TMP=/tmp/gd.$$.gir

#  define usage abort
usage() {
	exec >&2
	cat <<==EOF==
Usage: $0 [-$FLAGS] file [arg...]
  -N  no optimization
  -c  compile only, producing IR on file.gir
  -a  compile only, producing IR on file.gir and assembly listing on file.gia
  -d  compile only, producing Dot directives on file.dot
==EOF==
	gexec -? 2>&1 | sed -n 's/=false: /  /p'
	exit 1
}

#  process options
XOPTS=
AFLAG=
CFLAG=
DFLAG=
OPT=": optim -O"
while getopts $FLAGS C; do
    case $C in
	a)			AFLAG=$C;;
	c)			CFLAG=$C;;
	d)			DFLAG=$C;;
	N)			OPT="";;
	[ltvADEFJPT])	XOPTS="$XOPTS -$C";;
	?)			usage;;
    esac
done
shift $(($OPTIND - 1))
test $# -lt 1 && usage

I=$1
B=${I%.*}
DOT="gtran cat $I : yylex : parse : ast2ir $OPT : dot_File : stdout"
TRAN="gtran cat $I : yylex : parse : ast2ir $OPT : json_File : stdout"
shift

export COEXPSIZE=300000

if [ -n "$AFLAG" ]; then	# -a: produce file.gir and file.gia, then quit
    $TRAN >$B.gir && gexec $XOPTS -l -A $B.gir >$B.gia
	exit
fi

if [ -n "$CFLAG" ]; then	# -c: produce file.gir, and quit
    exec $TRAN >$B.gir
fi

if [ -n "$DFLAG" ]; then	# -d: produce file.dot, and quit
    exec $DOT >$B.dot
fi

#  translate and execute
trap 'X=$?; rm -f $TMP; exit $X' 0 1 2 15
if $TRAN >$TMP; then
	exec gexec $XOPTS $TMP "$@"
else
	exit 125		# exit 125 designates compilation error
fi
