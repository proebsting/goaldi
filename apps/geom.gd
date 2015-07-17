#  geom.gd -- draw points, discs, lines, rects in multiple sizes

procedure main() {
	^w := canvas()
	^x := -120
	^y := -120
	every ^d := 0.25 to 7 by 0.25 do {
		w.Size := 1
		w.Line(x, y, x+ 5, y)
		w.Size := d
		w.Point(x + 15, y)
		every ^dx := 0 to 50 by d + 5 do {
			w.Disc(x + 75 - dx, y, d)
		}
		w.Line(x + 85, y, x + 110, y)
		w.Rect(x + 120, y - d/2, 25, d)
		y +:= d + 5
	}
	while @w.Events ~=== "stop"
}
