# palm trees gpx demo: penwidth, colors, text, turtles, touching, stopping

procedure main() {
	^n := 0
	^c := [ "aqua", "black", "blue", "brown", "fuchsia", "gold", "gray",
		"green", "lime", "maroon", "navy", "olive", "orange", "purple", "red",
		"silver", "slate", "teal", "yellow",]
	^w := canvas()
	w.color("white")
	w.Forward(-100)
	w.turn(-90)
	w.turn(180 / *c)
	w.Size := 20
	every w.color(!c) do {
		w.Forward(40)
		w.turn(360 / *c)
	}
	w.Size := 3
	while ^e := @w.Events do {
		if e.Action ~== "release" then {
			w.color(?c)
		}
		if e.Action == "stop" then {
			stop("GOODBYE")
		}
		w.Point(e.X, e.Y)
		if e.Action == "release" then {
			n +:= 1
			w.Text(e.X - 10, e.Y + 10, n)
			w.Goto(e.X, e.Y, -90)
			w.Forward(50)
			every ^i := !10 do {
				w.turn(24 + ?24)
				w.Forward(20)
				w.Forward(-20)
			}
		}
	}
}

