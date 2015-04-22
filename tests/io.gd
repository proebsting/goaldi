#SRC: goaldi original
#  i/o test

procedure main() {
	# local f, s
	local f
	local s

	# simple reading
	write("a. ", read())
	write("b. ", @%stdin)
	write("c. ", !%stdin)
	every write("d. ", !%stdin \ 3)

	# open and read, including binary
	f := file("io.dat")
	write("e. ", @f)
	write("f. ", !f)
	write("g. ", f.read())
	write("h. ", f.get())
	write("i. ", image(f.readb(3)))
	write("j. ", image(f.readb(4)))
	write("k. ", image(f.readb(5)))
	write("l. ", image(f.get()))
	%stdout.put(f.get(), f.get(), f.get(), f.get())
	write("m.")
	every 1 to 3 do
		%stdout @: @f
	f.close()

	# open and write, including binary writes to make CRLF and raw CR
	write()
	f := file("io1.tmp", "w")
	f.write("first line normal")
	f.write("raw\r    CR embedded in this line")
	f.writes("line ending in CRLF\r\n")
	f.write("another normal line")
	f.print(12, 34, 5)	# spaces, no newline
	f.println(6, 78, 90)	# adjoins previous, spaces, newline
	f.flush()
	# extended character sets
	f.write("Latin1: naïve Häagen-Dazs Frusen Glädjé")
	f.write("Latin1: na\xEFve H\xE4agen-Dazs Frusen Gl\xE4dj\xE9")
	f.writeb("Latin1: na\xC3\xAFve H\xC3\xA4agen-Dazs Frusen Gl\xC3\xA4dj\xC3\xA9\n")
	f.write("Unicode: ✔§⌘±∮π€♻★☯♖☂☮♫¶")
	f.write("Unicode:  ♠ A K Q  ♥ A K Q  ♦ A K Q J  ♣ K J 9")
	f.write("Unicode:  \u2660 A K Q  \u2665 A K Q  \U2666 A K Q J  \u2663 K J 9")
	f.writeb("Unicode:  \xE2\x99\xA0 A K Q  \xE2\x99\xA5 A K Q  \xE2\x99\xA6 A K Q J  \xE2\x99\xA3 K J 9\n")
	f.write("another normal line")
	f.writes("unterminated line")
	f.close()

	# read back that file as normal text
	f := file("io1.tmp")
	while show(@f)
	f.close()

	# read back that file in binary
	# (non-ASCII chars look strange because UTF-8 is not decoded)
	write()
	f := file("io1.tmp")
	show(f.readb(1000))
	f.close()

	# test failure to open
	file("/no/such/file/exists", "f") | write("[open failed as expected]")

	# test bidirectional appending I/O
	write()
	file("io2.tmp", "w").write("abcde\nfghij").close()
	f := file("io2.tmp", "rwa")
	write("skip: ", @f)
	f.write("klmno")
	write("skip: ", @f)
	f.write("pqrst")
	f.close()
	f := file("io2.tmp")
	every write("reread: ", !f)
	f.close()

}

procedure show(s) {
	write(*s, ": ", image(s))
	return
}
