#!/bin/sh  
#
#  runtest [name...] -- test Goaldi translator and intepreter

#  check for necessary binaries
GOBIN=${GOPATH%%:*}/bin
GTRAN=$GOBIN/goaldi
GTRAN=$GOBIN/gtran
GEXEC=$GOBIN/gexec
ls -l $GOALDI $GTRAN $GEXEC || exit

#  if no test files specified, run them all
if [ $# = 0 ]; then
   set - *.std
fi

#  loop through the chosen tests
FAILURES=
for F in $*; do
    F=`basename $F .std`
    F=`basename $F .gdi`
    rm -f $F.gir $F.out $F.err
    printf "%-12s" $F:
	if test -r $F.dat; then
		exec <$F.dat
	else
		exec </dev/null
	fi
    if goaldi $F.gdi >$F.out 2>$F.err; then
	if cmp -s $F.std $F.out; then
	    echo "ok"
	    rm $F.out
	    test -s $F.err || rm $F.err
		rm -f $F*.tmp
	else
	    echo "output differs"
	    FAILURES="$FAILURES $F"
	fi
    else
	echo "error exit"
	FAILURES="$FAILURES $F"
    fi
done

echo ""
if [ "x$FAILURES" = "x" ]; then
   echo "gtests: all tests passed"
   echo ""
   exit 0
else
   echo "gtests failed: $FAILURES"
   echo ""
   exit 1
fi
