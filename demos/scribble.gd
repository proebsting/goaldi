#  draw on the screen

procedure main() {
	^w := canvas()	# open a window
	w.Size := 3		# use a fat pen
	^x := 0
	^y := 0
	while ^e := @w.Events do {
		if e.Action = 1 then {		# if a drag event
			w.Line(x, y, e.X, e.Y)	# draw a line
		}
		x := e.X
		y := e.Y
	}
}

