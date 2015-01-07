#SRC: goaldi original
#
#  simple test of closures

procedure main() {
	local p
	local q
	local r

	local mkproc := procedure(name) {
		/static pcount := 0
		local pnum := (pcount +:= 1)
		local n := 0
		return procedure(arg)  {
			static t
			/t := 0
			n +:= 1
			t +:= 1
			write("p#", pnum, ":  ", name, "(", arg, ")",
				"  call #", n, ", total=", t)
		}
	}

	p := mkproc("p")
	q := mkproc("q")
	q("00")
	p(11)
	p(22)
	q(33)
	q(44)
	r := mkproc("r")
	q(55)
	p(66)
	r(77)
	r(88)
}
