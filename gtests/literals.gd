#SRC: Goaldi original
#
#  test literals

procedure main() {

	# should interpret excapes in quoted string
	^s := "\b\d\e\f\l\n\r\t\v\'\"\\\067\130\x58\u0058\^H"
	write(image(s))
	every ^c := !s do
		write(ord(c), " ", image(c))

	# these escapes should have no effect
	write(image("\a\c\g\h\i\j\k\m\o\p\q\s\u\w\x\y\z"))

	# should not interpret excapes in raw-quoted string
	s := `\b\d\e\f\l\n\r\t\v\'\"\\\067\130\x58\u0058\^H`
	write(*s, " ", image(s))
	write(*s, " ", s)

	# try multi-line raw-quoted string
	write(`line 1
		line 2
line 3
		line 4`)
}
