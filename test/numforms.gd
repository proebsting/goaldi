#SRC: Goaldi original
#
# test numeric conversion in both translator and runtime system

procedure main() {
	write("            want    transl       string    runconv")
	trynum(      0,        0, "0")
	trynum(      1,       1b, "1b")
	trynum(      7,       7o, "7o")
	trynum(      8,      9r8, "9r8")
	trynum(     15,      0Fx, "0Fx")
	trynum(     42,       42, "42")
	trynum(     42,      042, "042")
	trynum(     42,      042, " \t 042\t ")
	trynum(     42, 2r101010, "2r101010")
	trynum(     42,  101010b, "101010b")
	trynum(     42,     8r52, "8r52")
	trynum(     42,      52o, "52o")
	trynum(     42,    16r2A, "16r2A")
	trynum(     42,      2Ax, "2Ax")
	trynum(     42,    23r1J, "23r1J")
	trynum(   1295,    36rZz, "36rZz")
	trynum(  27183,    27183, "27183")
	trynum(    210,11010010b, "11010010b")
	trynum(  10039,   23467o, "23467o")
	trynum( 524095,   7Ff3Fx, "7Ff3Fx")
	trynum(6.02e23,   602e21, "602e21")
	trynum(6303265,  602e21x, "602e21x")
	trynum(  .0123, 0.123e-1, "0.123e-1")
	trynum(   1123, 1.123e+3, "1.123e+3")
	trynum(  .0123,  .123e-1, ".123e-1")
	trynum(   1230,  .123e+4, ".123e+4")
	trynum(123.456,  123.456, "123.456")
	trynum(    789,     789., "789.")
}

procedure trynum(want, n, s) {
	local sn := number(s) | nil
	writes(if n === want & sn === want then "Okay: " else "ERROR:")
	write(right(want,10), right(n,10), right(image(s), 14), right(sn, 10))
	return
}
