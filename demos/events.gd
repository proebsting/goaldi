#  echo events to stdout until window is closed

procedure main() {
	local e
	local w := canvas()
	write(w.Canvas)
	repeat {
		write(e := @w.Events)
	} until e.Action === "stop"
}
