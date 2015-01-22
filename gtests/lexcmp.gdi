#SRC: icon/lexcmp.icn

# lexical comparison test

procedure main() {
	write("    s1    s2    <<   <<=    ==   ~==   >>=    >>")
	every (local s := "" | "a" | "b" | "c" | "x" | 2 | "") &
		(local t := "" | "a" | "c" | "x" | "2") do {
		wr(s)
		wr(t)
		wr(s << t  | nil)
		wr(s <<= t | nil)
		wr(s == t  | nil)
		wr(s ~== t | nil)
		wr(s >>= t | nil)
		wr(s >> t  | nil)
		write()
		}
	}

procedure wr(s) {
	printf("%6v", \s | "---")
	return
}
