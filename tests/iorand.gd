procedure main() {
	^L := list()
	^f := file("iorand.tmp", "crw")
	while ^line := read() do {
		L.push(^n := f.where())
		show(n, line)
		f.write(line)
	}
	^eof := f.where()
	show(eof, "[EOF]")
	write()
	every ^n := !L do {
		f.seek(n)
		show(n, f.read())
	}
	write()
	f.seek()
	show(1, !f)
	write()
	every ^i := eof to 1 by -20 do {
		f.seek(i)
		show(i, f.read())
	}
	write()
	every i := 0 to -eof by -20 do {
		f.seek(i)
		show(i, f.read())
	}
}

procedure show(n, s) {
	return write(right(n,5), ".   ", s)
}
