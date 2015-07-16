#  scribble.gd -- a very simple interactive demo
#
#  use mouse/finger to draw on the screen

procedure main() {
	^w := canvas()	# open a window
	w.Size := 3		# use a fat pen
	^x := 0
	^y := 0
	while ^e := @w.Events do {
		if e.Action == "drag" then {
			w.Line(x, y, e.X, e.Y)	# draw a line
		}
		x := e.X
		y := e.Y
	}
}

