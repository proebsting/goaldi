#  colors.gd -- colors demo
#
#  shows all the named colors (except white)
#  and illustrates examples of other specification forms.

global x
global y
global dy := 13

procedure main() {
	^w := canvas()
	w.VFont := font("mono", dy)

	x := -125	# left column
	y := -125	# reset to top
	every showcolor(w, !names)

	x := 0		# right column
	y := -125	# reset to top
	showcolor(w, .333)			# darkish gray
	showcolor(w, .667, 1)		# lightish gray
	showcolor(w, .5, 1, 0)		# yellow-green
	showcolor(w, 0, 1, .75)		# better aqua (i.e. not cyan)
	showcolor(w, .2, .8, 1)		# sky blue
	showcolor(w, .75, .5, 1)	# violet
	showcolor(w, .5, 0, 1)		# better purple
	showcolor(w, .75, 0, 0)		# dark red
	showcolor(w, 1, .5, 0)		# better orange
	showcolor(w, .75, .5, 0)	# tan
	showcolor(w, "white")		# (for spacing)
	showcolor(w, "#D")			# #k
	showcolor(w, "#A3")			# #kk
	showcolor(w, "#FDB")		# #rgb
	showcolor(w, "#FDB8")		# #rgba
	showcolor(w, "#6F6BFE")		# #rrggbb
	showcolor(w, "#6F6BFE80")	# #rrggbbaa

	while (@w.Events).Action ~=== "stop"	# wait until closed
}

procedure showcolor(w, a[]) {
	y +:= dy
	w.color(color ! a)
	w.Rect(x, y, 40, -(dy - 2))
	w.Text(x + 45, y, string(w.VColor)[3:0])
}

global names := [
	"fuchsia",	# magenta
	"red",
	"orange",
	"gold",
	"yellow",
	"lime",
	"aqua",	# cyan (turquoise)
	"blue",
	"navy",
	"teal",
	"green",
	"olive",
	"brown",
	"maroon",
	"purple",
	"black",
	"slate",
	"gray",
	"silver",
]
