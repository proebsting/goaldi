#SRC: goaldi original
#  channel test

procedure main() {
	local ch
	write("[new1]")
	ch := channel(5)
	try(ch) # should fail, nothing sent yet
	ch.put("algebub")
	ch @: "biolozy"
	ch.put("chemixtry")
	try(ch)
	try(ch)
	ch.close()
	write("[closed]")
	try(ch)	# should get pending value
	try(ch) # should fail
	try(ch) # should fail

	write("[new2]")
	ch := channel()
	ch := ch.buffer(3)
	ch @: "one"
	ch.put("two").put("three")
	# ch.put("four") would deadlock
	try(ch)
	try(ch)
	try(ch)

	write("[new3]")
	ch := buffer(4, create(!6))
	drain(ch)

	write("===")
	write("ch1 === ch1:  ", if ch === ch then "identical" else "distinct")
	write("ch1 === ch2:  ", if ch === create(1) then "identical" else "distinct")
}

#  try reading one value from channel, showing size
procedure try(ch) {
	write(image(ch), " size=", *ch, "   =>   ", image(ch.get()) | "[failed]")
}

#  drain channel and print, without showing size  (more deterministic)
procedure drain(ch) {
	while write(image(ch), "  =>  ", image(@ch))
	write(image(ch), "  =>  [failed]")
}
