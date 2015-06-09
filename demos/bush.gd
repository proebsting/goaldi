#  an early demo of Goaldi turtle graphics
#  draws a multi-colored bush

procedure main() {
	^w := canvas()
	w.color("white")
	w.Forward(-120)
	w.turn(-90)
	w.color("silver")
	every !72 do {
		w.turn(2.5)
		w.Forward(10)
		w.turn(2.5)
		sleep(0.01)
	}
	w.turn(90)
	randomize()
	bush(w, 3, 8)
	sleep()
}

procedure bush(w, z, n) {
	/static clist := ["black", "brown", "red", "orange", "gold", "gray"]
	w.color(?clist)
	w.Size := z
	w.turn(?90 - 45)
	w.Forward(8 + ?17)
	sleep(0.002)
	if n > 0 then {
		bush(w.copy(), .93 * z, n-1)
		bush(w.copy(), .87 * z, n-1)
	}
}
