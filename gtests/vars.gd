#SRC: goaldi original
global a
global g

procedure main() {
	a := "Ahoy"
	g := "Gladiola"
	p()
	q()
	p()
	q()
	println("m:", a, g)
}

procedure p() {
	local a
	static s
	println("p:", \a | "--", g, \s | "--")
	a := "Ain't gonna see this"
	g := "Gorgonzola"
	s := "Sarasota"
}

procedure q() {
	local a
	static t
	a := "Asparagus"
	println("q:", a, g, \t | "--")
	a := "Ain't gonna see this either"
	g := "Gouda"
	t := "Turnip"
}
