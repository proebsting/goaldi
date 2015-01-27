#SRC: icon/record.icn

record simple(f)
record rec(f1, f2)

procedure main() {

	local a := rec()
	a.f1 := 1
	a.f2 := 2
	write("a1 ", a.f1, " ", a.f2)
	a := rec(3)
	a.f2 := 4
	write("a2 ", a.f1, " ", a.f2)
	a := rec(5,6)
	write("a3 ", a.f1, " ", a.f2)
	a.f1 := 7
	a.f2 := 8
	write("a4 ", a.f1, " ", a.f2)
	a := rec(9,10)
	write("a5 ", a.f1, " ", a.f2)
	a := rec(11, 12)
	every write("!a ", !a)
	every !a := 13
	write("a6 ", a.f2)

	local b := simple(14)
	write("*b ", *b)
	write("?b ", ?b)
	?b := 15
	write("!b ",!b)

	b := rec(3, 7)
	every write("b[n] ", b[1 to 3])
	every write("b[s] ", b["f" || (1 to 3)])

	a := rec(1, 2)
	b := rec(3, 4)
	a.f1 +:= 10
	a.f2 +:= 20
	every !b +:= 70
	every writes(" ", !a | !b | "\n")

	local c := b.copy()
	b.f2 +:= 3
	write(image(a))
	write(image(b))
	write(image(c))

	write("simple.f: ", simple.f)
	write("rec.f1: ", rec.f1)
	write("rec.f2: ", rec.f2)
}
