#SRC: Goaldi original
#
#  test literals

procedure main() {

	# should interpret excapes in quoted string
	^s := "\b\d\e\f\l\n\r\t\v\'\"\\"
	write(image(s))
	every ^c := !s do
		write(ord(c), " ", image(c))

	# this gets error in both Goaldi and Jcon
	# write(image("\a\c\g\h\i\j\k\m\o\p\q\s\u\w\x\u\z"))

	# should not interpret excapes in raw-quoted string
	s := `\b\d\e\f\l\n\r\t\v\'\"\\`
	write(*s, " ", image(s))
	write(*s, " ", s)

	# try multi-line raw-quoted string
	write(`line 1
		line 2
line 3
		line 4`)
}
