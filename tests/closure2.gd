#SRC: goaldi original
#
#	closures with lots of variables of various kinds

global g
global r

procedure main() {
	static s

	local a := 100
	local b := 200
	local c := 300
	g := 1000
	s := 2000

	local p := procedure(x) {
		static t
		local b
		println("P0.", g, s, a, b, c, x, t, b)
		/t := 400
		g +:= 100
		s +:= 200
		a +:= 10
		b := 7389
		t +:= 1
		println("P1.", g, s, a, b, c, x, t, b)
	}

	println("M0.", g, s, a, b, c)
	p(1000)
	println("M1.", g, s, a, b, c)

	local q := procedure(y) {
		static u
		local d
		println("Q1.", g, s, a, b, c, y, u, d)
		/u := 400
		g +:= 10
		s +:= 20
		d := 700
		b +:= 10
		r := procedure(z) {
			static v
			println("R1.", g, s, a, b, c, y, u, d, z, v)
			/v := 0
			g +:= 1
			s +:= 2
			c +:= 10
			d +:= 1
			v +:= 1
			println("R1.", g, s, a, b, c, y, u, d, z, v)
		}
		u +:= 12
		println("Q1.", g, s, a, b, c, y, u, d)
	}

	println("M2.", g, s, a, b, c)
	p(3000)
	println("M3.", g, s, a, b, c)
	q(4000)
	println("M4.", g, s, a, b, c)
	r(5000)
	println("M5.", g, s, a, b, c)

	println("M6.", g, s, a, b, c)
	p(7000)
	println("M7.", g, s, a, b, c)
	q(9000)
	println("M8.", g, s, a, b, c)
	r(9000)
	println("M9.", g, s, a, b, c)
}
