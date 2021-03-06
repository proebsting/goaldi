#SRC: Goaldi original
#
#   test named arguments in procedure and method calls

record r()

procedure main() {

	try(show)        # Goaldi procedure
	try(r().show)    # Goaldi method
	write()
	showlist(4,5,6)
	showlist(args:[1,2,3])

	# library methods
	write()
	write(image(channel(3).buffer(size:5)))
	write(image([].sort(i:1)))
	printf(x:[%phi, %e, %pi], fmt:"%.4f  %.4f  %.4f\n")
	%stdout.writeb(s:"stdout writeb\n")
}

procedure try(p) {
	write()
	write(image(p))
	write()
	p()
	p(10,20,30,40)
	p(a:11, b:21, c:31, d:41)
	p(d:42, c:32, b:22, a:12)
	p(c:33, a:13, d:43, b:23)
	write()
	p(14, c:34)
	p(b:25)
	p(d:46, a:16)
	p(17, d:47, b:27)
}

procedure show(a,b,c,d) {
	write("a=", a, "  b=", b, "  c=", c, "  d=", d)
}

procedure r.show(a,b,c,d) {
	write("a=", a, "  b=", b, "  c=", c, "  d=", d)
}

procedure showlist(args[]) {
	writes("args:")
	every writes(" ", !args)
	write()
}
