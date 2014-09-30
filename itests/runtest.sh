#!/bin/sh  
#
#  runtest [name...] -- test Goaldi translator and intepreter
#
#  Initially just tests that a program can be translated and loaded;
#  does not check for correct output.

#  check for necessary binaries
GOBIN=${GOPATH%%:*}/bin
GTRAN=$GOBIN/gtran
GEXEC=$GOBIN/gexec
ls -l $GTRAN $GEXEC || exit

#  define gexec arguments
TARGS="-l -v -A"

#  define jtran arguments
export COEXPSIZE=300000		# need 250000 for v9/ipl/farb.icn!
jtran() {
    $GTRAN preproc $1 : yylex : parse : ast2ir : optim -O : json_File : stdout
};

#  if no test files specified, run them all
if [ $# = 0 ]; then
   set - *.std
fi

#  loop through the chosen tests
JFAIL=
TFAIL=
for F in $*; do
    F=`basename $F .std`
    F=`basename $F .icn`
    rm -f $F.gir $F.out $F.err
    printf "%-12s" $F:
    if jtran $F.icn >$F.gir 2>$F.err; then
	if $GEXEC $TARGS $F.gir >$F.out 2>>$F.err; then
	    echo "ok"
	    test -z $F.err && rm $F.err
	else
	    echo "gexec failed"
	    TFAIL="$TFAIL $F"
	fi
    else
    	echo "jtran failed"
	JFAIL="$JFAIL $F"
    fi
done

echo ""
if [ "x$TFAIL$JFAIL" = "x" ]; then
   echo "itests: all tests passed"
   echo ""
   exit 0
else
   echo "gtran failed: $JFAIL"
   echo "gexec failed:  $TFAIL"
   echo ""
   exit 1
fi
