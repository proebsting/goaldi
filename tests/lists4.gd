#SRC: Goaldi original
#
#   Tests some problems seen in accessing reversed lists
#   when used in rvalue contexts.

procedure main() {
	^L := [1,2,3]	# initially a "normal" list
	show(L)			# show it
	show(copy(L))	# show a copy
	write()
	L.push(0)		# now it's "reversed" internally
	show(L)			# and you SHOULDN'T be able to tell that
	show(copy(L))	# nor on a copy
}

procedure show(L) {
	local i
	write(image(L))
	every writes(" ", !L | "\n")			# takes rvalue path
	every writes(" ", (!L + 0) | "\n")		# takes lvalue path
	every i := 1 to *L do writes(" ", L[i])
	write()
	every i := 1 to *L do writes(" ", L[i] + 0)
	write()
}
