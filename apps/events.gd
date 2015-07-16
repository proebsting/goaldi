#  events.gd -- event tracer
#
#  echo events to stdout until window is closed

procedure main() {
	local e
	local w := canvas()
	w.Text(-100, -50, "(echoing events to stdout)")
	write(w.Canvas)
	repeat {
		write(e := @w.Events)
	} until e.Action === "stop"
}
