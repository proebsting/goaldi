#!/bin/sh
#
#	goaldi [options] file.gd... [--] [arg...] -- compile and run Goaldi program
#
#	To see options, run with no arguments.
#	This script assumes that gexec is in the search path.

FLAGS=acNltvADEPT

#  define the usage abort function
usage() {
	exec >&2
	cat <<==EOF==
Usage: $0 [-$FLAGS] file.gd... [--] [arg...]
  -N  no optimization
  -c  compile only, producing IR on file.gir
  -a  compile only, producing IR on file.gir and assembly listing on file.gia
==EOF==
	# add option descriptions from back end (gexec)
	gexec -.!! -? 2>&1 | sed -n -e '/-.!!/d' -e 's/=false: /  /p'
	exit 1
}

#  process options
XOPTS=
WHAT=x
NFLAG=
while getopts $FLAGS C; do
	case $C in
	a)			WHAT=$C;;
	c)			WHAT=$C;;
	N)			NFLAG=-$C;;
	[ltvADEPT])	XOPTS="$XOPTS -$C";;
	?)			usage;;
	esac
done

#  collect source file names:
#  the first argument, always, plus any following that end in ".gd"
shift $(($OPTIND - 1))		# remove flag arguments
test $# -lt 1 && usage		# require at least one file argument
SRCS=$1						# save that argument
shift						# and remove from execution parameters

while [ "$1" != "${1%.gd}" ]; do	# while name ends in .gd
	SRCS="$SRCS $1"				# add to list
	shift						# and remove from execution parameters
done

#  remove a "--" separator argument if present
if [ "$1" = "--" ]; then
	shift
fi

#  make scratch directory for temporary files, and arrange its deletion
SCR=/tmp/goaldi.$$
trap 'X=$?; rm -rf $SCR; exit $X' 0 1 2 15
mkdir $SCR

#  compile the source files
OBJS=
QUIT=:
for F in $SRCS; do
	B=${F%.*}
	case $WHAT in
		a)	# -a: produce file.gir and file.gia
			gexec $NFLAG $F >$B.gir && gexec -.!! $XOPTS -l -A $B.gir >$B.gia
			QUIT="exit $?"
			;;
		c)	# -c or nothing: produce file.gir
			gexec $NFLAG $F >$B.gir
			QUIT="exit $?"
			;;
		x)	# no flag: produce temporary file.gir for later execution
			O=$SCR/${B##*/}.gir
			OBJS="$OBJS $O"
			gexec $NFLAG $F >$O || QUIT="exit 1"
			;;
	esac
done

$QUIT	# exit if nothing more to do, or if errors in compilation

# execute compiled files
gexec -.!! $XOPTS $SCR/*.gir -- "$@"
