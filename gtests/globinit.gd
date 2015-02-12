#SRC: Goaldi original
#
#  global initialization test -- checks proper sequencing
#
#  execution should be ordered so that a through m form a Fibonacci sequence

initial { printall("init1") }

initial { %stdin := %stdout }

global d := b + c
global j := h + i
global t := a + d + g + j + l + m
global l := j + k
global c := a + b
global x := 0
global y := v + t + u
global h := f + g
global i := g + h

initial { t +:= 1; u -:= 5; printall("init2") }

global e := c + d
global z := a + b + c + d + e + f + g + h + i + j + k + l + m
global g := e + f
global y
global a := 1
global u := b + d + f + h + k + m
global k := i + j
global w := 77
global f := d + e
global b := 1
global v := c + i + e
global m := k + l

initial { v -:= 1; w +:= 3; printall("init3") }

procedure main() {
	printall("main")
	write("aa=", aa, " bb=", bb, " cc=", cc, " dd=", dd)
	%stdin.println("done")	# appears on stdout due to initial{} reassignment
}

initial { x := 407; y := reverse(y); printall("init4") }

procedure printall(label) {
	println(label)
	println("    a-m:", a, b, c, d, e, f, g, h, i, j, k, l, m)
	println("    t-z:", t, u, v, w, x, y, z)
}

initial { z := 6789; printall("init5") }

# test dependency involving a procedure call
# (from the Go language reference page)

global aa := show("aa", cc + bb)
global bb := show("bb", ff())
global cc := show("cc", ff())
global dd := show("dd", 3)

procedure ff() {
	dd +:= 1
	write("ff returning ", dd)
	return dd
}

procedure show(label, value) {
	write(label, " := ", value)
	return value
}

# test dependency on mutually recursive procedures

global rr := show("rr", r1(5))

procedure r1(n) {
	if n > 100 then return n
	return r2(2 * n)
}

procedure r2(n) {
	return r1(3 * n)
}
