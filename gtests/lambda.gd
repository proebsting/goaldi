#SRC: goaldi original
#
#	lambda test
#	also procedure type test

procedure main() {

	write(" main: ", image(main), " : ", image(type(main)),
		" : ", image(type(main)(main) | "[procedure constructor failed]"))

	local a := 7
	write(" a = ", a)

	local by3 := lambda(i, j) { a := i to j by 3 }
	write(" by3 = ", image(by3), " : ", type(main) === type(by3) | "[FAILED]")
	every writes(" ", by3(1, 20) | "\n")
	write(" a = ", a)

	local by7 := lambda(i, j)  local a := i to j by 7
	write(" by7 = ", image(by7), " : ", type(main) === type(by7) | "[FAILED]")
	every writes(" ", by7(21, 50) | "\n")
	write(" a = ", a)

	write(type(main)("OOPS") | "done")
}
