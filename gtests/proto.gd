#SRC: icon/proto.icn
#  This program contains samples of all the basic syntactic forms in Icon.
#  (Now modified somewhat for Goaldi but not necessarily complete.)

record three(x,y,z)
record zero()
record one(z)

global line
global count

procedure main() {
	write()
}

procedure expr1(a, b) {
	local x
	local y
	local i
	local j

	static e1
	/e1 := 0
	()
	{}
	();()
	[]
	[,]
	x.y
	x[i]
	x[i:j]
	x[i+:j]
	x[i-:j]
	(,,,)
	x(,,,)
	x!y
	not x
	|x
	!x
	*x
	+x
	-x
	.x
	/x
	=x
	?x
	\x
	@x
}

procedure expr2(a, b[]) {
	local x
	local y
	local i
	local j
	local k
	local c1
	local c2
	local s1
	local s2
	x \ i
	x @ y
	i ^ j
	i * j
	i / j
	i % j
	c1 ** c2
	i + j
	i - j
	c1 ++ c2
	c1 -- c2
	s1 || s2
	x ||| y
	i < j
	i <= j
	i = j
	i >= j
	i > j
	i ~= j
	s1 << s2
	s1 == s2
	s1 >>= s2
	s1 >> s2
	s1 ~== s2
	x === y
	x ~=== y
	x | y
	x ~| y
	i to j
	i to j by k
	x := y
	x <- y
	x :=: y
	x <-> y
	i +:= j
	i -:= j
	i *:= j
	i /:= j
	i %:= j
	i ^:= j
	i <:= j
	i <=:= j
	i =:= j
	i >=:= j
	i ~=:= j
	c1 ++:= c2
	c1 --:= c2
	c1 **:= c2
	s1 ||:= s2
	s1 <<:= s2
	s1 <<=:= s2
	s1 ==:= s2
	s1 >>=:= s2
	s1 >>:= s2
	s1 ~==:= s2
	x |||:= y
	x ===:= y
	x ~===:= y
	x &:= y
	x @:= y
	x & y
	create x
	return
	return x
	suspend x
	suspend x do y
}

procedure expr4() {
	^i; ^j; ^s; ^x
	local e; local e1; local e2; local e3

	while e1 do break
	# while e1 do break e2
	while e1 do continue
	case e of {
		x:   return fail
		(i > j) | 1    :  return
		}
	case *s of {
		1:   1
		default:  return fail
		}
	if e1 then e2
	if e1 then e2 else e3
	repeat e
	repeat e1 until e2
	while e1
	while e1 do e2
	every e1
	every e1 do e2
}

procedure expr9() {
	^x
	x
	local X_
	nil
	"abc"
	"\n"
	"^a"
	"\001"
	"\x01"
	1
	999999
	36ra1
	3.5
	2.5e4
	4e-10
	.127
}
