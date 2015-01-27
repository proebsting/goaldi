#SRC: icon/wordcnt.icn
#     (extensively rewritten)
#
#	W O R D   C O U N T I N G
#
#	This program tabulates the words in standard input and writes the results.
#	The definition of a "word" is naive.

procedure main() {
	# local line, words, rx, w, kv
	local line
	local words
	local rx
	local w
	local kv

	words := table()
	rx := regex("[\\W]+")
	while line := read() do {
		line := rx.ReplaceAllString(line, " ")
		every w := !fields(line) do
			(\words[w] +:= 1) | (words[w] := 1)
	}
	every kv := !words.sort() do
	printf("%6.0f  %s\n", kv.value, kv.key)
}
