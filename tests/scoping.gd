#SRC: goaldi original
#
# test nested scoping

procedure main() {
	local x := 1
	static y := 2
	write("00. ", x, " ", y)
	every local i := !3 do {
		write(i, "a. ", x, " ", y)
		local x := 10
		/static y := 20
		write(i, "b. ", x, " ", y)
		x +:= 2
		y +:= 3
		write(i, "c. ", x, " ", y)
	}
	write("99. ", x, " ", y)
	local L := []
	every local j := 1 to 4 do {
		local x := j
		local f := procedure () { return x }
		L.put(f)
	}
	every write("f: ", (!L)())
}
