#SRC: icon/wordcnt.icn (extensively rewritten)
#
#	Word Counter
#
#	A word is a string of one or more Unicode "letters".

procedure main(filename) {

	local f := file(\filename) | %stdin		# input file
	local words := table(0)					# table for talling counts
	local rx := regex(`\pL+`)				# expr to match words

	while local line := f.read() do {		# read line
		local matches := rx.FindAllString(line, -1)		# find words
		every local w := !\matches do {					# for each (if any)
			words[w] +:= 1								# bump the tally
		}
	}
	every local kv := !words.sort() do				# for each key/value pair
		printf("%6.0f  %s\n", kv.value, kv.key)		# print count and word
}
