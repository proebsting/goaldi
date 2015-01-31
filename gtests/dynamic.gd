#SRC: Goaldi original
#
#  test dynamic variables

global i
global j

procedure main() {
	show("a"); dynamic %x := 100
	show("b")
	every i := !2 do {
		show("   c");  %x +:= 1
		show("   d");  dynamic %y := 200
		show("   e");  %y +:= 2
		show("   f")
		showxyz("   g", %x, %y, "--")
		every j := !3 do {
			show("      h");  dynamic %x := 300
			show("      i");  %x +:= 3
			show("      j");  %y +:= 4
			show("      k");  dynamic %z := 400
			show("      l")
			showxyz("      m", %x, %y, %z)
		}
		dynamic %z := 500
		show("   u")
	}	
	show("v")
}

procedure show(label) {
	showxyz(label, xval(), yval(), zval())
}

procedure showxyz(label, x, y, z) {
	write(label, ":  i=", i, "  j=", j, "  %x=", x, "  %y=", y, "  %z=", z)
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
	return "(" || string(e)[-3:0]	
}
