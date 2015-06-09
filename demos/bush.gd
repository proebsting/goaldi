#  an early demo of Goaldi turtle graphics
#  draws a multi-colored bush; redraws on a click or touch

procedure main() {
	randomize()
	^w := canvas()
	repeat {
		bush(w)
		while (@w.Events).Action ~= 2
		w.Reset()
	}
}

procedure bush(w) {
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
	branch(w, 3, 8)
}

procedure branch(w, z, n) {
	/static clist := ["black", "brown", "red", "orange", "gold", "gray"]
	w.color(?clist)
	w.Size := z
	w.turn(?90 - 45)
	w.Forward(8 + ?17)
	sleep(0.002)
	if n > 0 then {
		branch(w.copy(), .93 * z, n-1)
		branch(w.copy(), .87 * z, n-1)
	}
}
