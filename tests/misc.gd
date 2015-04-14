#SRC: icon/misc.icn

record message(
	who,	# something
	gap,	# something else
	what,	# something more
)

procedure main() {
	# local i, x
	local i
	local x

	x := 1
	x +:= |1		# tickled optimizer bug.
	write(x)

	x := table()
	write(x[])

	x := "o"
	write("a" & "b")
	write("c" | "d")
	write(\"e")
	write(!"f")
	write(\nil | "g")
	write(/nil & "h")
	write("i" || "jk")
	write(23 || "skidoo")
	write(x, .x, x := "b")

	every write( (1|2)("hello", "mom"), "!")
	every write ! [ (1|2) ! ["hello", "mom"], "!"]
	#write ! message("hello")
	#write ! message("hello", " ", "pop")
	every i := -4 to 4 do
		write("i=", i, ": ", i("a","b","c") | "failed")

	every write(seq() \ 3)
	every write(seq(4) \ 3)
	every write(seq(,4) \ 3)
	every write(seq(10,20) \ 3)

	write("repl: ", repl("",5), repl("x",3), repl("foo",0), repl("xyz",4))
	write("reverse: ", reverse(""), reverse("x"), reverse("ab"), reverse(12345));
	every i := 0 to 255 do
		if (ord(char(i)) ~= i) then write("char/ord oops ", i)
	writes("char: ")
	every writes(char((64 to 126) | 10))

	evaluation("1234567890", "abcdefghi")

	every write(image(nullsuspend()))

	every write(tstreturn())

	write("done")
	exit()
	write("oops!")
	dummy()
}

procedure tstreturn() {
	return fn()
}

procedure fn() {
	suspend "OK to get here"
	write("Should not get here when called from a 'return'")
}

# These got different results under
# Icon's (odd) two-pass argument evaluation process.
procedure evaluation(a,b) {
	# local x,y
	local x
	local y

	write("argument evaluation test")
	write(x, x:=1)
	write(x:=2, x:=3)
	write(a, a := 3)
	write(b[2], b[2] := "q")
	write(b[2:3], b[1:4] := "qwerty")
	y := [1,2,3,4]
	write(y[1], y[1] := 3)
	x := 7
	write(x[2], y[2] := 3)
	y := table()
	write(y[3], y[3] := 7)
	x := y
	write(x[5], y[5] := 8)
}

procedure dummy() {
	image(every 1) | 2	# this triggered a problem once upon a time.
}

procedure args(x[]) {	# later replaced by proc("args",0)
	local s
	s := ""
	every s ||:= image(!x) do
		s ||:= " "
	return s[1:-1] | ""
}

procedure nullsuspend() {
	suspend
	suspend
}
