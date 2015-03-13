#SRC: goaldi original
#
#  Table operations

procedure main() {
	# local t, u, c, i, k, kv, l
	local t
	local u
	local c
	local i
	local k
	local kv
	local l

	t := table()
	println("\t\t\t\t", type(t), t, *t, image(t))
	ck(t)
	every i := !4 do {
		t[i] := "abcd"[i]
		ck(t)
	}
	every c := !"efgh" do {
		t[toupper(c)] := c
		ck(t)
	}
	every i := 5-!4 do {
		t[i] := "wxyz"[i]
		ck(t)
	}
	println("\t\t\t\t", type(t), t, *t, image(t))
	u := t.copy()
	write("\t\t\t\t t vs t: ", if t === t then "identical" else "distinct")
	write("\t\t\t\t t vs u: ", if t === u then "identical" else "distinct")
	u.delete(2).delete("G")
	ck(u)
	ck(t) # should be unchanged
	println("\t\t\t\t", type(t), t, *t, image(t))
	every k := 3 | "G" | 1 | "H" | 2 | 4 | "F" | "E" do {
		t.delete(k)
		ck(t)
	}
	println("\t\t\t\t", type(t), t, *t, image(t))
	ck(u)
	l := []
	every l.put(!u)
	ck(u)
	every kv := !l.sort() do
		writes(" ", kv.key, ":", kv.value)
	write()

	write()
	t := table("#")
	t { 2 to 3 : "l", 4 | "F": "a", "E": "m"}
	ck(t)

	#%#% random portion disabled
	#  every t[!4 | !"EFGH"] := ?"abcdefghijklmnopqrstuvwxyz"
	#  ck(t)
	#  writes("\t\t\t")
	#  every !12 do
	#     kv := ?t & writes(" ", kv.key, ":", kv.value)
	#  write()
}

procedure ck(t) {		#: show table indexed by 1..4 and "E".."H"
	# local k, kv
	local k
	local kv

	writes(*t, " ")
	every k := !"1234" | !4 | !"EFGH" do {
		writes(t.member(k) | "-")
	}
	every writes(" " | t[!4 | !"EFGH"])
	writes(" ")
	every kv := !t.sort() do
		writes(" ", kv.key, ":", kv.value)
	write()
}
