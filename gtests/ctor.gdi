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
	c := constructor("point", "x", "y")
	show(c, c(3, 5))
	show(c, c(2, 3))
	show(c, c(y:8, x:4))
	c := constructor("rect", "x", "y", "w", "h")
	show(c, c(4,3,2,1))
	show(c, c(w:6, h:4, x:1, y:3))
	every write("L: ", image(!L.sort()))
	every write("V: ", image(!V.sort()))
}

procedure show(c, v) {
	L.put(c)
	V.put(v)
	write(image(c), " : ", image(v))
}
