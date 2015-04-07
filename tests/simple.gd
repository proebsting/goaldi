#SRC: goaldi original
procedure main() {
local i
local s
local n
local p
local x
	# local i, s, n, p, x

	println("hello", 47, 3.14159)

	i := sqrt(2)
	s := "abc"
	n := nil
	p := main
	every x := i | s | n | p do
		println(x, type(x), image(x))

	every write (1 | 2 | 3)
	every write (4 to 6)
	write(p1())
	every write(p1())
	every write(p2())
	p3(11,12,13)
	every writes (!12 | "\n")
	every writes ((!12 \ 5) | "\n")
	every 1 to 5 do write(?0)
	every writes (" ", (!27) | "\n")
	every writes (" ", ?(!27) | "\n")
	every writes (" ", (!27 & ?100) | "\n")
	every i := -5 to +5 do
		write (i, i("a","b","c","d") | "--")

	writes("seq")
	every writes(" ", ":" | seq() \ 3)
	every writes(" ", ":" | seq(5) \ 3)
	every writes(" ", ":" | seq(10, 2) \ 3)
	every writes(" ", ":" | seq(, 17) \ 3)
	every writes(" ", ":" | seq(2.5, .375) \ 5)
	write()

	every writes(" a", !3 | !3)
	every writes(" b", !3 ~| !3)
	every writes(" c", !1 ~| !3)
	every writes(" d", !0 ~| !3)
	every writes(" e", !3 ~| !0)
	every writes(" f", 1 | 2 ~| 3 | 4 ~| 5 | 6)
	every writes(" g", no() | 2 ~| 3 | 4 ~| 5 | 6)
	every writes(" h", no() | no() ~| 3 | 4 ~| 5 | 6)
	every writes(" i", no() | no() ~| no() | 4 ~| 5 | 6)
	every writes(" j", no() | no() ~| no() | no() ~| 5 | 6)
	every writes(" k", no() | no() ~| no() | no() ~| no() | 6)
	every writes(" l", no() | no() ~| no() | no() ~| no() | no())
	write()

	p4()
	p4(1)
	p4(2,3)
	p4(4,5,6)
	p4(7,8,9,10)
	p4(11,22,31,41,59,26,535)
}

procedure p1() {
	return 7
}

procedure p2() {
	suspend 7 | 8 | 9
}

procedure p3(a,b,c) {
	write(a,b,c)
}

procedure p4(a,b,c[])
{
	write("p4: ", image(a), " ", image(b), " ", image(c))
}

procedure no(){}	# always just fails
