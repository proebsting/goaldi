#SRC: goaldi original
#
#	label test

procedure main() {

	every:outer local i := 10 to 90 by 10 do {
		every:inner local j := i + 1 to i + 9 do {
			if j = 14 then continue:inner
			if j = 24 then continue
			if j = 34 then continue:outer
			writes(" ", j)
			if j = 16 then break:inner
			if j = 26 then break
			if j = 46 then break:outer
		}
	}
	write()

	local ii := 10
	repeat:outer {
		local j := ii + 1
		repeat:inner {
			if j = 14 then continue:inner
			if j = 24 then continue
			if j = 34 then continue:outer
			writes(" ", j)
			if j = 16 then break:inner
			if j = 26 then break
			if j = 46 then break:outer
		} until (j +:= 1) > ii + 9
	} until (ii +:= 10) > 90
	write()

	ii := 0
	while:outer (ii +:= 10) < 90 do {
		local j := ii
		while:inner (j +:= 1) < ii + 10 do {
			if j = 14 then continue:inner
			if j = 24 then continue
			if j = 34 then continue:outer
			writes(" ", j)
			if j = 16 then break:inner
			if j = 26 then break
			if j = 46 then break:outer
		}
	}
	write()
}
