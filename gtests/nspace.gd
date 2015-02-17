#SRC: Goaldi original
#INCL: nspack.gd
#
#  nspace.gd -- namespace test

global g1 := note("g1", 12)
global g2 := pack::note("g2", 24)

procedure main() {
	write("g1 = ", g1)
	write("g2 = ", g2)
	write("pack::gval = ", pack::gval)
	write("pack::glen = ", pack::glen)
	write("pack::n3(3) = ", pack::n3(3))
	note("a",20)
	pack::note("b", 21)
	pack::run()
	^x := pack::r(98,76)
	x.show()
}

initial { write("main initial") }

procedure ilen(x) {
	^s := image(x)
	write("ilen(", s, ") = ", *s)
	return *s
}

procedure note(label, value) {
	write("----- note(", label, ",", value, ")")
	return value
}
