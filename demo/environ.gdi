#  demo of external array:
#  get environment from Go,
#  stomp a few entries randomly,
#  and print it out

procedure main() {
	local e := environ()
	write("Unix environment (", *e, " entries):")
	every local i := 11 to 100 by 13 do
		e[i] := "===============[REDACTED]==============="
	every !5 do
		?e := "===============[STOMPED]==============="
	every !5 do
	# every write(">> ", !e)
	i := 0
	while (i +:= 1) & write(i, ".  ", e[i])
}
