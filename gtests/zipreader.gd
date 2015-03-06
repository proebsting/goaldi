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
	write(repl("-", 60))
	local h := zf.FileHeader
	write(h.Name, ":  ", h.UncompressedSize64, " bytes")
	local retv := zf.Open()		# u.c. "Open": Go method on zip file
	throw(\retv[2]) 			# handle error from Open
	local f := retv[1]			# extract file result
	contents(f)					# show contents
	f.close()					# l.c. "close": Goaldi method on Goaldi file
}

procedure contents(f) {
	local i
	every i := !5 do
		write(@f) | return fail
	@f & write("   ...   ")
}
