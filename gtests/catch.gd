#SRC: goaldi original
#
#	test panic recovery ("catch p")

procedure main () {
	catch lambda(e) write("main caught ", e, "; exiting")

	try("failure", noresult)
	try("nil", nilresult)
	try("panic value", errresult)
	try("raspberry", myrasp)
	try("catch message", mycatch)
	try("rethrow", rethrow)
	try("custom panic", altpanic)
	try("type conversion error", 5)
	try("17", suspender)

	write()
	write("dp1. ", image(doubleplay(errresult)))	# print exception
	write("dp2. ", image(doubleplay(noresult)))		# fail
	write("dp3. ", image(doubleplay(nilresult)))	# print nil
	write("dp4. ", image(doubleplay(nil)))			# abort (caught by main)
	write("dp5. (oops)")							# not reached
}

#	set exception handler twice, then raise exception
procedure doubleplay(handler) {
	catch nilresult		# will be overridden
	catch handler		# this one counts -- from our argument
	3 to 5 by 0			# provoke exception
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
	write("   catch ", image(catch rproc))
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
