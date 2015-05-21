#SRC: goaldi original
#   simple test of create (and channels)

procedure main() {
	local x := create !10
	local y := create !"abcdefg"
	local e := create evens()
	local o := create odds()
	local c := create %current | image(%current) | (lambda() %current)()
	while write(@x, ". ", @y) # n.b. consumes 8 before failing
	while write("+ ", @x | @y)
	while write("e ", @e)
	while write("o ", @o)
	while write("c ", @c)	#%#% SHOULD write three channels (no nils)
	create unused()	# test create result not used
	sleep(0.001)	# don't exit before cx runs
}

procedure evens() {
	suspend seq(0,2) \ 10
}

procedure odds() {
	every %current @: seq(1,2) \ 10
}

procedure unused() {
	write("unused here")
}
