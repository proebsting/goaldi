record inner extends outer (c)
record outer (a,b)
record innie extends outer (e, f, g)

record point(x,y)
record circle extends point(r)
record square extends point(w)
record rect extends square(h)

procedure main() {
	every w(outer | inner | innie)
	every w(point | circle | square | rect)
	w(outer(1,2))
	w(inner(3,4,5))
	w(innie(6,7,8,9,0))
	w(point(0,0))
	w(circle(2,1,3))
	w(square(1,2,3))
	w(rect(5,3,4,3))
}

procedure w(x) {
	write(image(x))
}
