#SRC: goaldi original
#
#  provoke and print, as examples, several run-time errors
#
#  the numbering of cases is just help to correlate output with program code

record r(a,b,c)

procedure main() {
	every try(0 to 1000)
}

procedure try(i) {
	catch lambda(e) {
		write(i, ". ", e)
	}
	provoke(i)
}

procedure provoke(i) {
	case i of {

		10: %huh

		20: "x" || main
		21: 1 + "x"
		22: 1 + main
		23: ?(-3)
		24: select { 666 @: 0 : 0}

		32: ?"abcd" := 5
		33: @"efgh"

		40: constructor()
		41: constructor(%pi)
		42: constructor("r", "a", "2", "c")
		43: tuple(1,2,3)
		44: tuple(a:1,b:2,a:3)

		50: nil to 10
		51: "x" to 10
		52: 1 to nil
		53: 1 to []
		54: 1 to 10 by 0

		60: channel()[3]
		61: 3 @: 7
		62: [].sort(-1)
		63: r().huh
		64: [].huh


		69: echolist(a:12)
		70: echo(1, a:2)
		71: echo(b:3, b:4)
		72: echo(d:5)
		73: r().echo(e:5)
		74: r().echo(self:6)
		75: r(x:7)
		76: r(1,2,3,4)
		77: r.d
		78: printf(s:"hello")
		79: 3(1,2,e:5)

		80: file("x", "q")
		81: file("/no/such/file")
		82: file("/bin", "w")
		83: file("/bin/ls", "w")
		84: file("/dev/null").write()
		85: !file("/dev/null", "w")
		86: file("/dev/null").close().read()
		87: file("/dev/null").close().close()

		120: gmean(0,1,2)
		121: gmean(3,1,-1)
		122: hmean(-1,3)
		123: hmean(2,0)

		130: char(-1)
		131: char(123456x)
		132: ord("ab")
		133: left("ab",5,"")
		134: center("x",7,"")
		135: right("xyz",26,"")
		136: map("abc", "123", "3210")

		140: remove("/no/such/file")

		150: regex("(")

		195: throw("my double error", 12, 34)
		196: throw("my own error")
		197: throw("my nil error", nil)
		198: throw("my pi error", %pi)
		199: throw(99, %phi)

		default: return fail
	}
	write(i, ". unexpected survival")
}

procedure echo(a,b,c) {
	write("echo: a=", a, " b=", b, " c=", c)
}

procedure r.echo(a,b,c) {
	write("r.echo: a=", a, " b=", b, " c=", c)
}

procedure echolist(a[]) {
	writes("echolist:")
	every writes(" ", !a)
	write()
}
