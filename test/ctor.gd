#SRC: Goaldi original
#
#  constructor test

global L
global V

procedure main() {
	L := []
	V := []
	local c
	c := constructor("empty")
	show(c, c())
	c := constructor("point", "xpos", "ypos")
	show(c, c(3, 5))
	show(c, c(2, 3))
	show(c, c(ypos:8, xpos:4))
	c := constructor("rect", "x", "y", "w", "h")
	show(c, c(4,3,2,1))
	show(c, c(w:6, h:4, x:1, y:3))
	every write("L: ", image(!L.sort()))
	every write("V: ", image(!V.sort()))
}

procedure show(c, v) {
	L.put(c)
	V.put(v)
	write(image(c), " : ", *c, " : ", image(v))
	every ^i := 1 to *c do {
		^s := c[i] | "[missing]"
		write("   c[", i, "] == ", image(s),
			"    c[", image(s), "] = ", c[s] | "[failed]",
			"    v[", i, "] = ", image(v[i]) | "[failed]",
			"    v[", image(s), "] = ", image(v[s]) | "[failed]")
	}
}
