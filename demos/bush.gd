#  an early demo of Goaldi turtle graphics
#  draws a multi-colored bush; redraws on a click or touch

global clist := ["black", "brown", "green", "blue", "gray", "purple"]
global dlist := ["black", "brown", "red", "orange", "gold", "gray"]

procedure main() {
	randomize()
	^w := canvas()
	repeat {
		moon(w)
		branch(w, 3, 8)
		while (@w.Events).Action ~= 2
		w.Reset()
		clist :=: dlist
	}
}

procedure moon(w) {
	w.Reset()
	w.color("white")
	w.Forward(-120)
	w.turn(-90)
	w.color("silver")
	every !72 do {
		w.turn(2.5)
		w.Forward(10)
		w.turn(2.5)
		sleep(0.001)
	}
	w.turn(90)
}

procedure branch(w, z, n) {
	w.color(?clist)
	w.Size := z
	w.Forward(8 + ?17)
	sleep(0.002)
	if n > 0 then {
		w.turn(?60 - 30)
		branch(w.copy(), .93 * z, n-1)
		w.turn(?60 - 30)
		branch(w.copy(), .87 * z, n-1)
	}
}
