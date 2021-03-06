#SRC: goaldi original
#
#   test of L.sort() and M.sort()

procedure main() {
	# local a, b, c, n, l, m
	local a
	local b
	local c
	local n
	local l
	local m

	a := [3,1,4,1,5,9,2,6,5,3,5]
	write("a1: ", image(a))
	write("a2: ", image(a.sort()))

	b := [:!"cowabunga,dude!":]
	write("b1: ", image(b))
	write("b2: ", image(b.sort()))

	m := table()
	every n := !8 do
		m[n * 97 % 61] := n * 71 % 43
	write("m0: ", image(m))
	show("l!:", l := [:!m:].sort())
	show("l0:", l.sort())
	show("l1:", l.sort(1))
	show("l2:", l.sort(2))
	show("m1:", m.sort(1))
	show("m2:", m.sort(2))
}

procedure show(label, kvlist)  {
	local kv
	writes(label)
	every kv := !kvlist do
		writes("  ", kv.key, ":", kv.value)
	return write()
}
