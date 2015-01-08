#SRC: goaldi original
#
#	test panic recovery ("catch p")

procedure main () {
	try("failure", noresult)
	try("nil", nilresult)
	try("panic value", errresult)
	try("raspberry", myrasp)
	try("catch message", mycatch)
	try("rethrow", rethrow)
	try("custom panic", altpanic)
	try("type conversion error", 5)
	try("17", suspender)
}

#	try(label, proc) -- with tracing, force error and report result of catch
procedure try(label, rproc) {
	catch tryfailed
	write("expect ", label, ":")
	local v := boom(rproc) | "[FAILED]"
	write("   got ", v)
}

#   report try failure (panic not caught, or rethrown)
procedure tryfailed(e) {
	write("   UNCAUGHT PANIC: ", e)
}

#	register rproc, force error
procedure boom(rproc) {
	catch errresult	# superseded unless rproc is invalid
	catch rproc
	2 to 1 by 0
}

#	return raspberry
procedure myrasp(e) {
	return "pbpbpbpbpttttt"
}

#	return catch message showing exception
procedure mycatch(e) {
	return "caught: " || string(e)
}

#	re-throw panic
procedure rethrow(e) {
	write("   caught panic; now reissuing")
	throw(e)
}

#	throw a different exception instead
procedure altpanic(e) {
	write("   caught panic; throwing another")
	throw("CUSTOM PANIC")
}

#	try suspending (shoudn't resume)
procedure suspender(e) {
	suspend 17 to 23 do
		write("RESUMED?!")
}
