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
	w(^p := point(0,0))
	w(^c := circle(2,1,3))
	w(^s := square(1,2,3))
	w(^r := rect(5,3,4,3))
	every (p|c|s|r).exhibit()
}

procedure circle.exhibit() {	# overrides point.exhibit
	write("CIRCLE(", self.x, ",", self.y, ",", self.r, ") ",
		type(self), "  ", image(self))
}

procedure point.exhibit() {
	write("at ", self.x, ",", self.y, ":  ", self.type(), "  ",self.image())
	case self.type() of {
		point:	nil
		circle:	write("  r=", self.r)
		square:	write("  w=", self.w)
		rect:   write("  size=(", self.w, ",", self.h, ")")
	}
}

procedure w(x) {
	write(image(x))
}
