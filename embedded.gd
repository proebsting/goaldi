#  sample embedded app for testing

procedure main(args[]) {
	writes("running embedded app with args:")
	every writes(" ", image(!args) | "\n")
}
