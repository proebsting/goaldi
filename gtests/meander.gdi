#SRC: icon/meander.icn
#
#		M E A N D E R I N G   S T R I N G S
#

#  This main procedure accepts specifications for meandering strings
#  from standard input with the alphabet separated from the length by
#  a colon.

procedure main() {
	while local line := read() do {
		local f := split(line, ":")
		if (*f = 2) & (local n := integer(f[2])) then {
			local alpha := f[1]
			write("meander(", alpha, ",", n, "): ")
			write(meander(alpha,n))
		} else {
			stop("erroneous input: ", line)
		}
	}
}

procedure meander(alpha,n) {
	local i := local k := *alpha
	local t := n-1
	local result := repl(alpha[1],t)
	while local c := alpha[i] do {
		if contains(result, result[-t:0] || c) ~= 0 then i -:= 1 else {result ||:= c; i := k}
	}
	return result
}
