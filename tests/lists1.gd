#SRC: goaldi original
#
#  Tests most list operations and methods except random selection.

procedure main() {
	# local a, b, c, i, j, l, l2
	local a
	local b
	local c
	local i
	local j
	local l
	local l2

	write("A:")
	show(list())
	show(list(3))
	show(list(5, 9))
	show(list(, "X"))

	write("B:")
	l := ["a","b","3"]
	write("\ttype:", type(l), "  size:", *l, "  print:", l, "  image:", image(l))
	show(l)

	write("C:")
	l := list()
	show(l)
	every l.put(!"def") do show(l)
	every l.push(!"cba") do show(l)
	every writes((l.get | l.pull | l.pop)(), " : ") do show(l)
	every l.put(!"ghi") do show(l)
	every writes((l.get | l.pull | l.pop)(), " : ") do show(l)
	show(l)
	l.push(3, 2, 1) & show(l)
	l.put(7, 8, 9) & show(l)

	write("D:")
	every (l := []).put(!"1yly5pmno")
	show(l)

	write("E:")
	while writes(@l & l.pull() & l.pop(), " : ") do
		l.push(*l) & l.put(*l) & show(l)
	show(l)

	write("F:")
	show([])
	show([7])
	show(l := [3, 1, 4, 1, 5, 9])
	show(l[2+:4])
	show(l[2:6].put(1).put(6))
	show(l)	# should be unchanged

	write("G:")
	every (l := []).push(!"fedcba")
	show(l)
	every i := -2 to 2 do
		every j := 3 to 5 do
			writes(i,":",j) & show(l[i:j])

	write("H:")
	l := [2,3,4]
	l2 := [7,8,9]
	show(l ||| l2)
	show(l.push(1) ||| l2)
	show(l ||| l2.push(5))
	show(l.put(5) ||| l2)

	write("I:")
	show(l2 := copy(l))
	show(l2.get() & l2.pull() & @l2 & l2)
	show(l)	# should be unchanged

	write("J:")
	show([:!"wxyz":])
	show([: 3 * !7 % 10 :])

	write("K:")
	a := [3,1,4,1,5,9,2,6,5,]
	write(image(a))
	write(image(a.sort()))
	every (b := []).put(!"cowabunga!")
	write(image(b))
	write(image(b.sort()))
	c := [3, "x", nil, 5.5, a, "q", 7, %stdin, "t", main, 9, ]
	write(image(c))
	write(image(c.sort()))

	write("L:")
	a := [2,7,1,8]
	b := list(3,a)	# should be 3 distinct lists in Goaldi
	b[1].put(3)
	b[2].put(2,8)
	b[3].put(0,9)
	every write(image(b | !b))

	write("M:")
	write([0][1] := 7)  # test asgmt to rvalue L[1] derived from lvalue [9]
	write(?[2,3,4] := 8)
	write(![5,6,7] := 9)

	write("N:")
	a := [5,3,0,9]
	a @: "E"
	a @: "A"
	a @: "9"
	write(image(a))

	write ! ["all ", "done", "!"]
}

procedure show(l) {
	local i
	writes("\t", *l, ":  ")
	every i := -9 to 9 do
		writes(l[i] | "-", " ")
	writes(" : ")
	every writes(" ", !l | "\n")
}
