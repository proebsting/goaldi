procedure main() {
	local t

	# tuple construction
	write(image(tuple()))
	write(image(tuple(a:1,b:3)))
	write(image(tuple(a:2,b:1)))
	write(image(tuple(a:1,b:3,c:5)))

	# tuple operations
	write()
	t := tuple(key:3, value:5)
	write(image(t))
	write(t.key, " : ", t.value)
	t := tuple(x:3, y:5, w:2, h:1)
	write(image(t))
	write(t.x, " ", t.y, " ", t.w, " ", t.h)
	every writes(" ", image(!t) | "\n")
	every writes(" ", image(t[!4]) | "\n")
	every writes(" ", image(t[!*t]) | "\n")

	# inspecting the tuple type
	write("t.type():")
	^y := t.type()
	every writes(" ", image(!y) | "\n")
	every writes(" ", image(y[!4]) | "\n")
	every writes(" ", image(y[!*y]) | "\n")
	every writes(" ", image(y[!"xywh"]) | "\n")

	# use tuple to protect an unhashable value
	write()
	^L := [1,2,3]
	^t1 := tuple(k:external(L))
	^t2 := tuple(k:external(L))			# a distinct value
	^t3 := tuple(k:external([1,2,3]))	# also distinct
	^t4 := tuple(k:external([4,5,6,7]))
	^S := set([t1,t2,t3,t4])
	write("S:  ", image(S))
	L := [: image(!S) :].sort()	# for reproducibility
	every write("!S: ", !L)
	# all should be of the same type (tuple(k))
	write(image(type(t1)))
	write(if type(t1) === type(t2) then "t1 === t2" else "t1 ~=== t2")
	write(if type(t2) === type(t3) then "t2 === t3" else "t2 ~=== t3")
	write(if type(t3) === type(t4) then "t3 === t4" else "t3 ~=== t4")
}
