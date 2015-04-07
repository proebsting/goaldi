procedure main() {
	local t
	write(image(tuple()))
	write(image(tuple(a:1,b:3)))
	write(image(tuple(a:2,b:1)))
	write(image(tuple(a:1,b:3,c:5)))
	t := tuple(key:3, value:5)
	write(image(t))
	write(t.key, " : ", t.value)
	t := tuple(x:3, y:5, w:2, h:1)
	write(image(t))
	write(t.x, " ", t.y, " ", t.w, " ", t.h)
	every write(image(!t))
}
