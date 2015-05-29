#!/bin/sh
#
#  Usage:  release.sh version-label  (e.g. release.sh v47)
#
#  Makes a release package containing the currently built executable.
#  Before running this, run "make full" and "make accept".

VERSION=${1?"Version number required"}

U=`uname`
case $U in
	Darwin) UNAME=Mac;;
	*)		UNAME=$U;;
esac
VNAME="Goaldi-$UNAME-$VERSION"

set -e
rm -rf $VNAME $VNAME.tgz
mkdir $VNAME
cp README.adoc $VNAME
cp LICENSE.adoc $VNAME
cp INSTALL.adoc $VNAME
cp goaldi $VNAME/goaldi
(
	file $VNAME/goaldi
	echo `date` /`whoami`
	uname -n -s -m
) >$VNAME/MANIFEST
chmod 755 $VNAME $VNAME/[a-z]*
chmod 644 $VNAME/[A-Z]*
echo
echo MANIFEST:
cat $VNAME/MANIFEST
tar cfz $VNAME.tgz $VNAME
echo
tar tvfz $VNAME.tgz
rm -rf $VNAME
echo
chmod 644 $VNAME.tgz
ls -l $VNAME.tgz
