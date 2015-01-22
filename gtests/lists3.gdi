#SRC: goaldi original
#
#  Tests random list selection and assignment.
#  Sensitive to randomness implementation.

procedure main() {
	# local l, c
	local l
	local c

	l := []
	show(l)
	every l.put(!"abcdefghijklmn") do
		show(l)
	every c := !"pqrstuvwxyz" do
		?l := c & show(l)
}

#  show random samples and then whole list
procedure show(l) {
	every 1 to 12 do
		writes(?l | "-")
	every writes("  " | !l | "\n")
}
