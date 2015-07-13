# colors demo

global x := -150
global y := -140
global dy := 16

procedure main() {
	^c := [ "aqua", "black", "blue", "brown", "fuchsia", "gold", "gray",
		"green", "lime", "maroon", "navy", "olive", "orange", "purple", "red",
		"silver", "slate", "teal", "yellow",]
	^w := canvas()
	w.VFont := font("mono", dy)
	every w.color(!c) do
		showcolor(w)
	x := 0
	y := -140
	every ncolor(w, 0 to 1 by 0.125)
	ncolor(w, .5, 1, 0)	# yellow-green
	ncolor(w, 0, 1, .75)	# better aqua (i.e. not cyan)
	ncolor(w, .2, .8, 1)	# sky blue
	ncolor(w, .75, .5, 1)	# violet
	ncolor(w, .5, 0, 1)		# purple
	ncolor(w, .75, 0, 0)	# dark red
	ncolor(w, 1, .5, 0)		# better orange
	ncolor(w, .75, .5, 0)	# tan

	while (@w.Events).Action ~=== "stop"
}

procedure ncolor(w, a[]) {
	w.color(color ! a)
	showcolor(w)
}

procedure showcolor(w) {
	w.Rect(x, y, 20, -(dy - 2))
	w.Text(x + 25, y, string(w.VColor))
	y +:= dy
}
