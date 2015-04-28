#!/bin/sh
#
#  Usage:  release.sh version-label
#
#  Makes a release package from the currently built executable.
#  Before running this, run "make full" and "make accept".

VERSION=${1?"Version number required"}
GOBIN=${GOPATH%%:*}/bin

U=`uname`
case $U in
	Darwin) UNAME=Mac;;
	*)		UNAME=$U;;
esac
VNAME="Goaldi-$UNAME-$VERSION"

STABLE=`head -3 tran/stable-gtran | awk 'NR==2 {print $4}'`
GHEAD=`git rev-parse HEAD`
if [ "$STABLE" != "$GHEAD" ]; then
	echo 1>&2
	echo 1>&2 WARNING: possible mismatch
	echo 1>&2 WARNING: stable-gtran: $STABLE
	echo 1>&2 WARNING: current-head: $GHEAD
fi

rm -rf $VNAME $VNAME.tgz
mkdir $VNAME
cp README.adoc $VNAME
cp LICENSE.adoc $VNAME
cp INSTALL.adoc $VNAME
cp $GOBIN/goaldi $VNAME
(
	file $VNAME/goaldi
	echo $GHEAD
	echo `date` /`whoami`
) >$VNAME/MANIFEST
echo
echo MANIFEST:
cat $VNAME/MANIFEST
tar cfz $VNAME.tgz $VNAME
echo
tar tvfz $VNAME.tgz
rm -rf $VNAME
echo
ls -l $VNAME.tgz
