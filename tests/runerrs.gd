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

		# variables
		10: %huh

		# exceptions
		21: throw("my double error", 12, 34)
		22: throw("my own error")
		23: throw("my nil error", nil)
		24: throw("my pi error", %pi)
		25: throw(99, %phi)

		# z : nil

		# t : type

		# n : number
		101: 1 + "x"
		102: 1 + main
		103: ?(-3)
		104: 3 @: 7
		105: 3(1,2,e:5)

		121: 10 to []
		122: nil to 10
		123: "x" to 10
		124: 1 to nil
		125: 1 to []
		126: 1 to 10 by 0

		140: gmean(0,1,2)
		141: gmean(3,1,-1)
		142: hmean(-1,3)
		143: hmean(2,0)

		# s : string
		201: "x" || main
		210: ?"abcd" := 5
		221: @"efgh"
		222: "efgh" @: "ijkl"
		223: (^s := "lmno") @: [666]

		240: char(-1)
		241: char(123456x)
		242: ord("ab")
		243: left("ab",5,"")
		244: center("x",7,"")
		245: right("xyz",26,"")
		246: map("abc", "123", "3210")

		# f : file
		301: file("x", "q")
		302: file("/no/such/file")
		303: file("/bin", "w")
		304: file("/bin/ls", "w")
		305: file("/dev/null").write()
		306: !file("/dev/null", "w")
		311: file("/dev/null").close().read()
		312: file("/dev/null").close().close()
		315: remove("/no/such/file")

		# c : channel
		341: select { 666 @: 0 : 0}

		# m : methodvalue

		# p : procedure
		361: echolist(a:12)
		362: echo(1, a:2)
		363: echo(b:3, b:4)
		364: echo(d:5)
		365: printf(s:"hello")

		# L : list
		401: 57 ||| []
		402: [] ||| 58
		403: channel()[3]
		404: [].sort(-1)

		# S : set
		441: set([1,2]) ++ 441
		442: set([1,2]) -- 442
		443: set([1,2]) ** 443
		447: 447 ++ set([1,2])
		448: 448 -- set([1,2])
		449: 449 ** set([1,2])
		451: set().put(external([1,2,3]))

		# T : table

		# R : record
		501: r.d
		502: r().huh
		503: [].huh
		511: r().echo(e:5)
		512: r().echo(self:6)
		521: r(x:7)
		522: r(1,2,3,4)
		541: constructor()
		542: constructor(%pi)
		543: constructor("r", "a", "2", "c")
		551: tuple(1,2,3)
		552: tuple(a:1,b:2,a:3)

		# X : external
		901: regex("(")

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
