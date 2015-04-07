#SRC: goaldi original
#  demo of 32-bit hash functions

procedure main() {
	# local files, lines, s
	local files
	local lines
	local s

	files := [adler32(), crc32(), fnv32(), fnv32a()]
	lines := ["", "tyger", "tyger", "burning", "bright", ""]
	report("[init]", files)
	every s := !lines do {
		every (!files).writes(s)
		report(s, files)
	}
}

procedure report(s, files) {
	printf("%-8s", s)
	every printf("  %10.0f", hashvalue(!files))
	printf("\n")
}
