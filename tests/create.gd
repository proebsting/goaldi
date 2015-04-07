#SRC: goaldi original
#   simple test of create (and channels)

procedure main() {
	# local x, y
	local x
	local y
	x := create !10
	y := create !"abcdefg"
	while write(@x, ". ", @y) # n.b. consumes 8 before failing
	while write("+ ", @x | @y)
	create unused()	# test create result not used
	sleep(0.001)	# don't exit before cx runs
}

procedure unused() {
	write("unused here")
}
