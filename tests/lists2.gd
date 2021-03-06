#SRC: icon/lists.icn
#
#  List test from Icon

procedure main() {
	# local i, x, y, z
	local i
	local x
	local y
	local z

	limage("a", list())
	limage("b", list(2))
	limage("c", list(,3))
	limage("d", list(4,5))
	limage("d", list(6,7))
	limage("e", [])
	limage("f", [nil])
	limage("g", [1])
	limage("h", [2,3,4,5])
	limage("i", [1,2,3] ||| [4,5,6,7,8])

	x := [1,2,3]
	x.push();				limage("-", x)
	x.put();					limage("-", x)
	x.push(nil);			limage("A", x)
	x.put(nil);				limage("B", x)
	write("\t", image(x.pop()));	limage("C", x)
	write("\t", image(x.get()));	limage("D", x)
	write("\t", image(x.pull()));	limage("E", x)
	x.push(4);			limage("F", x)
	x.push(5,6,7);		limage("G", x)
	x.push(8,9).push(10,11);	limage("H", x)
	x.put(12);			limage("I", x)
	x.put(13,14,15);		limage("J", x)
	x.put(16,17).put(18,19);	limage("K", x)
	x.push(20,21).put(22,23);	limage("L", x)
	every !x := 7;		limage("M", x)

	x := [1,2,3,4,5]

	every i := 0 to *x+3 do
		x[i] := i;
	limage("N", x)

	every i := -*x-3 to 0 do
		x[i] := i;
	limage("O", x)

	x := [1]
	write("\t", ?x)
	?x := 2
	limage("P", x)
	write(x[0] | "ok failure 0")
	write(x[2] | "ok failure 2")
	write(x[-2] | "ok failure -2")
	x.get()
	write(x.get() | "ok failure on get")
	write(x.pop() | "ok failure on pop")
	write(x.pull() | "ok failure on pull")

	x := [1,2,3,4,5,6,7,8,9]
	limage("p", x)
	limage("q", x[1:0])
	limage("r", x[2:5])
	limage("s", x[-3:5])
	limage("t", x[-5:-1])
	limage("u", x[-3+:6]) | write("u. ok wraparound failed")
	limage("v", x[3-:6]) | write("v. ok wraparound failed")

	write()
	y := copy(x)		# ensure that copies are distinct
	every !x +:= 10
	every !y +:= 20
	limage("x", x)
	limage("y", y)

	z := x ||| y
	limage("z", z)
	every !x +:= 10
	every !y +:= 20
	every !z +:= 50
	limage("x", x)
	limage("y", y)
	limage("z", z)

}

procedure limage(label, lst) {
	writes(label, ". [", *lst, "]")
	every writes(" ", image(!lst))
	write()
	return
}
