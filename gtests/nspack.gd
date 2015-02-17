#  nspack.gd -- a supplemental package for namespace testing

package pack

record r(a,b)

procedure r.show() { write(image(self)) }

global gval := note("gvinit", 7)

global glen := ilen(n3)

initial { write("pack initial") }

procedure run() {
	write("pack::run here")
	write("n3(4) = ", n3(4))
	note("run", 47)
	r(12,34).show()
}

procedure n3(n) {
	return n * n * n
}

procedure note(label, value) {
	write("pack::note(", label, ",", value, ")")
	return value
}
