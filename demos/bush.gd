#  an early demo of Goaldi turtle graphics
#  draws a multi-colored bush

procedure main() {
	^w := canvas()
	w.color("white")
	w.forward(-95)
	w.turn(-90)
	w.color("silver")
	every !36 do {
		w.forward(20)
		w.turn(10)
		sleep(0.01)
	}
	w.turn(90)
	randomize()
	bush(w, 8)
	sleep(10)
}

procedure bush(w, n) {
	/static clist := 
		["black", "brown", "red", "orange", "green", "blue", "purple", "gray"]
	w.color(?clist)
	w.turn(?90 - 45)
	w.forward(12 + ?10)
	sleep(0.002)
	if n > 0 then {
		bush(w.copy(), n-1)
		bush(w.copy(), n-1)
	}
}
