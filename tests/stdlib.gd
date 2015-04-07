#SRC: goaldi original
#  test miscellaneous library functions

procedure main()  {
	testprint()
	write("\nstrings:")
	every teststring("aBc" | "d33" | 47 | 3)
	write("\nconversion:")
	every testcnv(nil | "" | "abc" | "12" | "23.4" | 0 | 1 |
		%phi | %e | %pi | %stdin | %stdout | %stderr)
	testcommand()
	exit()
}

procedure testprint() {
	writes("ab", 34, "ef", %phi)
	write("gh", 90, "kl", %pi)
	print("mn", 37, "qr", %phi)
	println("st", 25, "uv", 0)
	write("543210")
	printf("%10.3f %g %.0f %s\n", %phi, %pi, 12345, "abcde")
	fprintf(%stdout, "%10.3f %g %.0f %s\n", %phi, %pi, 12345, "abcde")
	write(image(sprintf("%.4f", %e)))
}

procedure teststring(v) {
	writes(v, ":")
	apply(equalfold, v, "3")
	apply(repl, v, "3")
	apply(toupper, v)
	apply(tolower, v)
	apply(trim, v, "3")
	write()
	return
}

procedure testcnv(v) {
	writes(v, ":")
	every apply(type | image | number | string, v)
	write()
	return
}

procedure apply(p, x, y) {
	local v := (if \y then p(x,y) else p(x)) | "--"
	writes(" ", string(p)[3:0], "()", v)
	return
}

procedure testcommand() {
	write("\ncommand():")
	^c := command("echo", "hello", "world")
	c.Stdout := %stdout
	c.Stderr := %stderr
	write("command: ", c.Path, " ", c.Args)
	^r := c.Run()
	write("result:  ", image(r))
	write("state:   ", c.ProcessState)
}
