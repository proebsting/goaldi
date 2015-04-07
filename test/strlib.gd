#SRC: goaldi original
procedure main() {
	# local i, c1, c2, n1, n2, el, f
	local i
	local c1
	local c2
	local n1
	local n2
	local el
	local f

	el := "argon, boron, carbon, freon, krypton, silicon, teflon"
	write("repl:    ", repl("la, ", 11), "hey Jude")
	every write("reverse: ", reverse("abcde"[1+:(0 to 6)]) | "--")
	every write("tolower: ", tolower("AbCdE"))
	every write("toupper: ", toupper("AbCdE"))
	writes("fields: "); every writes(" ", image(!fields(el)) | "\n")
	writes("split:  "); every writes(" ", image(!split(el, ", ")) | "\n")

	write()
	every i := 0 | 1 | 33 | 100 | 200 | 300 | 1000 | 10000 | 100000 do {
		c1 := char(i)
		n1 := ord(c1)
		c2 := char(n1)
		n2 := ord(c2)
		println("char/ord:", i, image(c1), n1, image(c2), n2)
	}

	local s
	local pad
	local proc
	local w
	write()
	every s := "" | "*" | "xy" | "abc" do {
		every proc := left | center | right do {
			writes("pad:")
			every pad := "-" | "123" do {
				writes(" ")
				every w := 0 to 7 do {
					writes(" ", proc(s,w,pad))
				}
			}
			write()
		}
	}

	write()
	tryquote(`"abc"`)
	tryquote("`abc`")
	tryquote(`"abc\tdef"`)
	tryquote("`abc\tdef`")
	tryquote("`t0±Δt`")
	tryquote(`"t0±Δt"`)
	tryquote("abc")
	tryquote(`"ab`)
	tryquote("`ab")
	tryquote(`"ab\fyz"`)
	tryquote(`"ab\kyz"`)
	tryquote(`"ab\"`)

	write()
	write("map: ", map("aBcDeF"))
	write("map: ", map("AbCdEf"))
	write("map: ", map("aBcDeF", "abcdefghijklmnopqrstuvwxyz"))
	write("map: ", map("AbCdEf", "abcdefghijklmnopqrstuvwxyz"))
	write("map: ", map("aBcDeF", , "12345678901234567890123456"))
	write("map: ", map("AbCdEf", , "12345678901234567890123456"))
	write("map: ", map("aBcDeF", "abcdef", "!@#$%^"))
	write("map: ", map("AbCdEf", "abcdef", "!@#$%^"))
	write("map: ", map("", "abcdef", "!@#$%^"))
	write("map: ", map("abcdef", "aa", "bc"))
	write("map: ", map("Capitals Make A Title Or Slogan More Important"))
	write("map: ", map("but not too many!!!!", "abmnotuy", "ABMNOTUY"))
	write("map: ", map("If you can read this you can get a good job", "aeiou", ""))
	write("map: ", map("♠♥♦♣"))
	write("map: ", map("SDHC♠♥♦♣","♠♥♦♣","SHDC"))
	write("map: ", map("SDHC♠♥♦♣", "SHDC", "♠♥♦♣"))
	write("map: ", map("123456", "654321", "abcdef"))
	write("map: ", map("124578", "12345678", "03:56:42"))
	write("map: ", map("Hh:Mm:Ss", "HhMmSs", "035642"))
	write("map: ", map("123321", "123", "abc"))
}

procedure tryquote(a) {
	^b := unquote(a) | "[FAILED]"
	^c := quote(b)
	write("quoting: ", a, " => ", image(b), " => " ,c)
}
