#SRC: icon/case.icn
#
#  test control structures

record rec(a)

procedure main() {

	if 1 < 92 then write("okay1")
	write(if 2 < 17 then "okay2" else "oops3")

	local i := 3
	while i <= 5 do
		write("i=", i+:= 1, " [while/do]")
	repeat {
		write("i=", i+:= 1, " [repeat/until]")
	} until i > 7

	while i <= 5 do
		write("i=", i+:= 1, " [while/do OOPS]")
	repeat {
		write("i=", i+:= 1, " [repeat/until #2]")
	} until i > 7

	repeat {
		write("i=", i+:= 1, " [repeat #3a]")
		if i > 18 then break
		write("i=", i+:= 1, " [repeat #3b]")
		if i < 14 then continue
		write("i=", i+:= 1, " [repeat #3c]")
	}

	every writes(!"abcde") do writes(" ")
	every writes(!"fghij\n")

	local cx := create !"pqrst\n"
	while writes(@cx)			# while without do

	local r := rec(45)
	local c := create 1 | 2
	local t := table()
	t["a"] := "aaa"
	t["x"] := "xyz"
	local L := [nil, 0, 1.0, 2, %e, %pi, "", "0", "1", "2",
		rec, %stdin, main, write, rec(), r, c, cx, t, []]
	L.put(L.pop)	# append the "pop" method
	L.put(L)		# and L itself

	write()
	every local x := !L do {
		local s := case x of {
			"":		"\"\""
			0.0:	"0.0"
			1.0:	"1.0"
			2:		"2"
			%pi:	"%pi"
			"1":	"\"1\""
			nil:	"nil"
			main:	"main"
			write:	"write"
			%stdin:	"%stdin"
			rec:	"rec"
			rec():	"rec()"	# shouldn't ever match
			r:		"r"
			c:		"c"
			t:		"t"
			cx:		"cx"
			L:		"L"
			L.pop:	"L.pop"	# won't match, distinct methodvalue
			default:  "default"
		}
		printf("%-10s : %-10s : %s\n", s, string(x), image(x))
	}
}
