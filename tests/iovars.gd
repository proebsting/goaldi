#SRC: goaldi original
#  test effectiveness of setting dynamic variables %stdin and %stdout

procedure main() {

	genout("0")

	with %stdout := file("iovars1.tmp", "w") do {
		write("writing iovars1.tmp")
		genout("1")
		%stdout.close()
	}
	with %stdout := file("iovars2.tmp", "w") do {
		write("writing iovars2.tmp")
		genout("2")
		%stdout.close()
	}

	genout("3")

	with %stdin := file("iovars1.tmp") do {
		readall()
		%stdin.close()
	}

	with %stdin := file("iovars2.tmp") do {
		readall()
		%stdin.close()
	}

	remove("iovars1.tmp")
	remove("iovars2.tmp")
}

procedure genout(label) {
	^g := "g" || label || ": "
	write(g)
	write(g, "genout:")
	write(g, "%stdout = ", image(%stdout))
	writes(g, "part of line")
	write(" ... and the rest")
	print(g, "print", "a", "b", "c", "\n")
	println(g, "println", "a", "b", "c")
	%stdout.write(g, "explicit %stdout write")
	write(g)
	%stdout.flush()
}

procedure readall() {
	write()
	write("%stdin = ", image(%stdin))
	while write("> ", read())
	write()
}
