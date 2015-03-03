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
	every i := 0 | 1 | 33 | 100 | 200 | 300 | 1000 | 10000 | 100000 do {
		c1 := char(i)
		n1 := ord(c1)
		c2 := char(n1)
		n2 := ord(c2)
		println("char/ord:", i, image(c1), n1, image(c2), n2)
	}

	local proc
	local pad
	local wid
	local s
	every proc := left | right do {
		write()
		write(image(proc), ":")
		every pad := nil | "=" | "123" do
			every s := "" | "X" | "Tucson" do
				every wid := 1 | 5 | 10 do
					write("p(", image(s), ",", wid, ",", image(pad), ") => ",
						image(proc(s,wid,pad)))
	}
}
