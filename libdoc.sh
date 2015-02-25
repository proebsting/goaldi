#!/bin/sh
#
#  libdoc.sh -- extract Goaldi library documentation
#
#  This script uses "goaldi -E" to list the standard library contents,
#  runs "godoc" on each procedure, and mechanically edits the output.
#
#  Only procedures are listed.  Methods, at least for now, are discarded.

goaldi -E /dev/null 2>/dev/null |
sed '
	1,/-------/d
	/^$/,$d
	s/ -- /#/
	s/    */#/
	s/\([a-z0-9]\)\.\([A-Za-z]\)/\1#\2/
' | while IFS='#' read NAME DESCR PKG FUNC; do
	echo
	echo "$NAME -- $DESCR    ($PKG.$FUNC)"
	echo
	(echo; godoc $PKG $FUNC) | sed "
		1,/func $FUNC/d
		/^func (/,/^$/d
		/^type /,/^$/d
		# annoying special case:
		/^func NewSource(seed int64)/,/^$/d
	"
	echo "--------------------------------------------------------------------"
done
