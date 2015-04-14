#!/bin/sh
#
#  libdoc.sh -- extract Goaldi library documentation
#
#  This script uses "goaldi -l -E" to list the standard library contents,
#  runs "godoc" on each referenced package, then runs a Goaldi program
#  to produce the final output.
#
#  Note that libdoc.gd has an "exclusion list" to suppress certain procedures
#  such as sample extensions.  This list may need manual updating.

TMP1=/tmp/libdoc.$$a
TMP2=/tmp/libdoc.$$b
trap 'rm -f $TMP1 $TMP2; exit' 0 1 2 15

#  get the Goaldi procedure listing
goaldi -l -E /dev/null >$TMP1

#  extract a list of referenced packages
PKGS=`goaldi -l -E /dev/null 2>/dev/null |
	sed -n '/ -- /s/.*  \([a-zA-Z0-9/]*\)\.[^.]*$/\1/p' |
	sort |
	uniq`

#  get the documetation for those packages
for P in $PKGS; do
	godoc $P >>$TMP2
done

#  now process everything
./libdoc.gd $TMP1 $TMP2
