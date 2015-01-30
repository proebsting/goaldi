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
	%stdin.println("done")	# appears on stdout due to initial{} reassignment
}

initial { x := 407; y := reverse(y); printall("init4") }

procedure printall(label) {
	println(label)
	println("    a-m:", a, b, c, d, e, f, g, h, i, j, k, l, m)
	println("    t-z:", t, u, v, w, x, y, z)
}

initial { z := 6789; printall("init5") }
