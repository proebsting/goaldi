#SRC: goaldi original
#
#  test generation of values by a loop expression

procedure main() {
	local n := 74
	every writes(" ",
		"GO:" |
		(every local i := 10 to 50 by 10 do {
			yield i
			if i == (20 | 40) then yield i + 5
			if i == 30 then yield i+4 to i+6
		}) |
		(repeat {
			yield 61 to 64
			break
		}) |
		(while n < 78 do yield n +:= 1) |
		"DONE\n")
}
