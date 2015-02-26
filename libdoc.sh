#!/bin/sh
#
#  libdoc.sh -- extract Goaldi library documentation
#
#  This script uses "goaldi -E" to list the standard library contents,
#  runs "godoc" on each procedure, and mechanically edits the output.
#
#  Only procedures are listed.  Methods, at least for now, are discarded.

UL12='____________'
HBAR="$UL12$UL12$UL12$UL12$UL12$UL12"

goaldi -E /dev/null 2>/dev/null |
sed '
	# delete everything except procedure listing
	1,/-------/d
	/^$/,$d
	# reformat each line as: name(args)#descr#package#func
	s/ -- /#/
	s/    */#/
	s/\([a-z0-9]\)\.\([A-Za-z]\)/\1#\2/
' |
while IFS='#' read NAME DESCR PKG FUNC; do
	# process one library procedure
	echo $HBAR
	echo
	echo "$NAME -- $DESCR    ($PKG.$FUNC)"
	echo
	(echo; godoc $PKG $FUNC) | sed '
		# delete everything before desired function
		1,/func '$FUNC'/d
		# delete other functions, methods, and types
		/^func /,$d
		/^type /,$d
	' | 
	uniq |		# combine blank lines
	sed '$d'	# delete final blank line
done
echo $HBAR
