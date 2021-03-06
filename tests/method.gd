#SRC: goaldi original

# simple method test

record point()
record square(w)
record circle(r, color)

procedure main() {

	# define values
	local c1 := circle(2)
	local c2 := circle(7, "red")
	local s1 := square(4)
	local s2 := square(6)

	#  try methods
	every (point() | c1 | s1 | c2 | s2) . draw (1 | 5, 5 | 1)

	#  try methodvalue
	local m := c1.draw
	write("value: ", m, " : ", image(m), " : ", methodvalue(m) | "[failed]")
	write("type:  ", type(m), " ", m.type(), " ", type(m) === methodvalue)
	c1.color := "purple"
	m(3,4)

	#  check methodvalue comparisons
	compare(c1.draw,s1.draw)
	compare(c1.draw,c2.draw)
	compare(c1.draw,c1.draw)

	write(methodvalue("FAIL") | "done")
}

procedure point.draw(x, y) {
	show(self, x, y, "P")
}

procedure circle.draw(x, y) {
	show(self, x, y, "C")
	self.r +:= 1
}

procedure square.draw(x, y) {
	show(self, x, y, "Q")
}

procedure show(o, x, y, c) {
	printf("at %.0f,%.0f: %s: %#v\n", x, y, c, o)
}

procedure compare(m1, m2) {
	writes(if m1 === m2 then "SAME:" else "DIFF:")
	write("  ", image(m1), " : ", image(m2))
}
