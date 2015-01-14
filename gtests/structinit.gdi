#SRC: Goaldi original
#
#  test structure initialization  e0 { e1:v1 ...}

record rectangle(x,y,w,h,)

procedure main() {

	^a := list(26,"-"){ 1:"a", 5:"e", 9:"i", 15:"o", 21:"u"}
	write(a, " ", image(a))
	a{ "25":"y" }	# 25 gets converted to number
	write(a, " ", image(a))

	^c := table(){ "California":"Berkeley", "Arizona":"Tucson"}
	write(c, " ", image(c))
	c{ "Massachusetts":"Cambridge" }
	write(c, " ", image(c))

	^r := rectangle(){ "y":3, "x":5, "h":1, "w":2 }
	write(r, " ", image(r))
	r{ "x":8, "y":9 }
	write(r, " ", image(r))

	^s := "word"
	write(image(s))
	s{ 2:"i", 1:"b" }
	write(image(s))
}
