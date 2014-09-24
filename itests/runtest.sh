#!/bin/sh  
#
#  runtest [name...] -- test Goaldi translator and intepreter
#
#  Initially just tests that a program can be translated and loaded;
#  does not check for correct output.

#  check for necessary binaries
JTRAN=../tran/jtran
TERP=$GOPATH/bin/terp
ls -l $JTRAN $TERP || exit

#  define jtran arguments
export COEXPSIZE=300000		# need 250000 for v9/ipl/farb.icn!
jtran() {
    $JTRAN preproc $1 : yylex : parse : ast2ir : optim -O : json_File : stdout
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
	if $TERP $F.gir >$F.out 2>>$F.err; then
	    echo "ok"
	    test -z $F.err && rm $F.err
	else
	    echo "terp failed"
	    TFAIL="$TFAIL $F"
	fi
    else
    	echo "jtran failed"
	JFAIL="$JFAIL $F"
    fi
done

echo ""
if [ "x$TFAIL$JFAIL" = "x" ]; then
   echo "All tests passed."
   echo ""
   exit 0
else
   echo "jtran failed: $JFAIL"
   echo "terp failed:  $TFAIL"
   echo ""
   exit 1
fi
