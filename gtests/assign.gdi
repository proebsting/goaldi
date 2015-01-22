#SRC: goaldi original
# tests swapping and reversible assignment
procedure main() {

	write()
	local a := " algebra "
	local b := " botany "
	local c := " civics "
	write(1, a, b, c)
	a :=: b
	write(2, a, b, c)
	a :=: b :=: c
	write(3, a, b, c)
	a <- b <- c & write(4, a, b, c) & (1 < 0)
	write(5, a, b, c)
	a <-> b <-> c & write(6, a, b, c) & (1 < 0)
	write(7, a, b, c)
}
