#  demo sprites

global win		# main display window
global boxy		# a simple sprite figure

procedure main() {

	# make a sprite
	boxy := canvas(50, 50, 3)
	rtgl(boxy, 0, 0, +25, +25, "red")
	rtgl(boxy, 0, 0, +25, -25, "teal")
	rtgl(boxy, 0, 0, -25, +25, "green")
	rtgl(boxy, 0, 0, -25, -25, "purple")
	rtgl(boxy, -10, -10, 20, 20, "yellow")
	rtgl(boxy, -3, -3, 6, 6, "navy")

	# make the main canvas, splat a fixed sprite at the origin, draw axes
	win := canvas()
	win.Overlay(0, 0, boxy)							# copy sprite to origin
	rtgl(win, -25, -25, 50, 50, "#FFFD")			# whitewash incompletely
	every !4 do win.turn(90).copy().Forward(150)	# draw 4 axes

	# start the animations as asynchronous processes
	randomize()
	create stones()
	create drifter(0, 0, 0.2)
	create drifter(0, 0, 0.3)
	create drifter(0, 0, 0.4)
	create drifter(0, 0, 0.5)
	create drifter(0, 0, 0.6)
	create drifter(0, 0, 0.7)

	# wait until window is closed
	while (@win.Events).Action ~=== "stop"
}

#  draw a ring of standing stones
procedure stones() {
	^w := win.copy().color("orange")
	every ^i := -8 to 23 do {
		^a := i * %pi / 16
		^x := 120 * cos(a)
		^y := 120 * sin(a)
		w.Rect(x-5, y-5, 10, 10)
		sleep(0.01)
	}
}

#  animate a drifting sprite
#  (basically a random walk with some inertia and a little gravity)
procedure drifter(x, y, z) {
	# start some sprites drifting
	^e := win.AddSprite(boxy.Canvas, x, y, z)
	sleep(2)
	^dx := ^dy := ^ddx := ^ddy := 0
	repeat {
		sleep(0.05 * z)
		ddx := 0.3 * (?3 - 1)
		ddy := 0.3 * (?3 - 1)
		dx +:= ddx - 0.002 * x
		dy +:= ddy - 0.002 * y
		x +:= dx
		y +:= dy
		e.MoveTo(x, y, z)
	}
}

#  draw a rectangle of color k at (x,y,w,h) on canvas c
procedure rtgl(c, x, y, w, h, k) {
    c.copy().color(k).Rect(x, y, w, h)
}
