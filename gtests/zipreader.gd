#SRC: goaldi original
#  zip file reader demo

procedure main(fname) {
	# local zr, zf, f
	local zr
	local zf
	local f
	/fname := "zipreader.dat"
	zr := zipreader(fname)
	if /zr then stop("cannot open ", fname)
	write(fname)				# show archive name
	write("" ~== zr.Comment)	# show comment if present
	every zf := !zr.File do
		showfile(zf)
	zr.Close()
	write("[end]")
}

procedure showfile(zf) {
	# local f, h
	local f
	local h
	write(repl("-", 60))
	h := zf.FileHeader
	write(h.Name, ":  ", h.UncompressedSize64, " bytes")
	f := zf.Open()	# u.c. "Open": Go method on zip file reader object
	contents(f)		# show contents
	f.close()		# l.c. "close": Goaldi method on Goaldi file value
}

procedure contents(f) {
	local i
	every i := !5 do
		write(@f) | return fail
	@f & write("   ...   ")
}