#SRC: Goaldi original
#
#  Show examples of all types, presented multiple ways,
#  and test sorting and various type and value methods

record Example(	# one example for display
	value,		# example value
	type,		# value.type() (for sorting by type)
	gtype,		# corresponding global type if any
)

record Point(x,y)								# a simple illustrative record
procedure Point.dist() { return hypot(self.x, self.y) }	# and a method for it


procedure main() {

	# make a list of examples with associated global type values
	^E := []
	add(E, nil)
	add(E, type(), type)
	add(E, 17, number)
	add(E, %pi, number)
	add(E, 6.02214129e23, number)
	add(E, "abcd", string)
	add(E, %stdin, file)
	add(E, channel(3), channel)
	add(E, Point, constructor)
	add(E, ^P := Point(7,5), Point)
	add(E, P.dist, methodvalue)
	add(E, main)
	add(E, ^L := [2,3,5,7,11], list)
	add(E, ^T := table(){"Fe":"Iron","Au":"Gold"}, table)
	add(E, !T.sort())	# table element
	add(E, tuple(w:6,h:4))
	add(E, duration(3600+120+3), external)

	# show values various ways, checking universal methods in the process
	write()
	write("Examples sorted by value, showing presentation options:")
	E := E.sort(Example.value)
	write()
	^format := "%-4s %-15s %-30s %s\n"
	printf(format, "ch", "x.string()", "x.image()", "printf(\"%v\")")
	printf(format, "--", "----------", "---------", "------------")
	every ^x := !E do {
		^v := x.value
		^s := check(string, v, v.string())
		^i := check(image, v, v.image())
		^t := check(type, v, v.type())
		^f := sprintf("%v", v)
		if f[1+:2] == "0x" then		# if hex address
			f := "0xXXXXXX"			# hide actual value for reproducibility
		printf(format, t.char(), s, i, f)
	}

	write()
	write("Examples sorted by type, showing type information:")
	# n.b. stable sort keeps ordering reproducible within type
	E := E.sort(Example.type)
	write()
	format := "%-4s %-15s %-15s %-15s %s\n"
	printf(format, "ch", "x.string()", "x.type()", "t.name()", "global")
	printf(format, "--", "----------", "--------", "--------", "------")
	every x := !E do {
		^v := x.value
		^t := x.type
		printf(format, t.char(), v.string(), t, t.name(), t===x.gtype | "")
	}
	write()
}

procedure add(E, v, g) {			#: add global type and sample value
	return E.put(Example(v, type(v), g))
}

procedure check(p, x, s) {			#: validate p(x) === s
	^t := p(x)
	if t ~=== s then
		write("MISMATCH: ", p, "(", x, ")===", t, " ~=== ", s)
	return s
}
