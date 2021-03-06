#SRC: goaldi original
#
#  select test
#
#  results are deterministic if the random number generator is predictable

procedure main() {
	# local i, n, c1, c2, c3, c9
	local i
	local n
	local c1
	local c2
	local c3
	local c9

	every c1 | c2 | c9 := channel(1)
	c3 := channel(5)
	every i := !40 do {
		writes(i, ". ")
		select {
			n := @c1 : { write("c1 got ", n) }
			n := @c2 : { write("c2 got ", n); c1.put(n) }
			n := @c3 : { write("c3 got ", n); c2 @: n }
			c9 @: i   : { write("c9 sent ", i) }
			default   : {
				if ?4 === 0 then {
					write("c9 got ", @c9)
				} else {
					write("sending ", i)
					?[c1, c2, c3] @: i
				}
			}
		}
	}
	every c3 @: 77 | 88 | 99
	drain("c1", c1)
	drain("c2", c2)
	drain("c3", c3)
	drain("c9", c9)
	select {
		n := @c1 : write("oops: closed c1 returned ", n)
		n := @c9 : write("oops: closed c9 returned ", n)
		default  : write("ok: got default when files closed")
	}

	select {
		n := @c1 : write("oops: closed c1 returned ", n)
		n := @c9 : write("oops: closed c9 returned ", n)
	} | write("ok: no-default select failed as expected")

	write(select{} | "ok: empty select failed as expected")
}

procedure drain(name, ch) {
	ch.close()
	every writes(" ", "   drain" | name | ":" | !ch | "\n")
}
