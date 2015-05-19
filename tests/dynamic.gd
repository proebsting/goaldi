#SRC: Goaldi original
#
#  test dynamic variables

procedure main() {
	show("a")
	with %x := 12 do {
		show("b")
		with %y := 23 do {
			show("c")
			showxyz("C", %x, %y, "--")
			with %x := 15 do {
				show("d")
				with %z := 35 do {
					show("e")
					showxyz("E", %x, %y, %z)
				}
				with %z := 37 do {
					show("f")
				}
			}
			with %x := 17, %y := 25 do {
				show("g")
			}
			with %x := 19, %z := 39 do {
				show("h")
			}
		}
		show("v")
	}
	show("w")
	with %x := 555, %y := 666, %z := 777 do {
		show ("x")
		showxyz("X", %x, %y, %z)
	}
	show("z")

	write(with %foo := 1 do { 2 })
	# what should the following do?
	# every write(with %foo := 3 | 4  do { 5 | 6 })
}

procedure show(label) {
	showxyz(label, xval(), yval(), zval())
}

procedure showxyz(label, x, y, z) {
	write(label, ":  %x=", x, "  %y=", y, "  %z=", z)
}

procedure xval() {
	catch nope
	return %x
}

procedure yval() {
	catch nope
	return %y
}

procedure zval() {
	catch nope
	return %z
}

procedure nope(e) {
	# extract variable name from exception message
	return "(" || string(e)[-4:-2] || ")"
}
