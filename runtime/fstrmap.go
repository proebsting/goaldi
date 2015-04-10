//  fstrmap.go -- string mapping function

//#%#% initially naive.  no caching.  many opportunities for optimization here.
//#%#% might want to have a separate proc for creating a mapping table.

package runtime

import ()

func init() {
	DefLib(Map, "map", "s,from,into", "map characters")
}

const MAPSIZE = 128 // initial mapping table size
const MMARGIN = 128 // extra margin to allow when growing the mapping table

//  map(s,from,into) produces a new string that result from mapping the
//  individual characters of a source string.
//  Each character of s that appears in the "from" string is replaced by
//  the corresponding character of the "into" string.  If there is no
//  corresponding character, because "into" is shorter, then the character
//  from s is discarded.
func Map(env *Env, args ...Value) (Value, *Closure) {
	defer Traceback("map", args)

	// get arguments as rune arrays
	s := ProcArg(args, 0, NilValue).(Stringable).ToString().ToRunes()
	from := ProcArg(args, 1, UCASE).(Stringable).ToString().ToRunes()
	into := ProcArg(args, 2, LCASE).(Stringable).ToString().ToRunes()
	if len(into) > len(from) {
		panic(NewExn("Map: *into > *from", RuneString(into)))
	}

	// build a mapping table ctable
	//	 an entry value of -1 means delete
	//	 the default entry value of 0 means no mapping
	//	 store result+1 in entries that are to be mapped
	// start with size 128 and grow as needed
	ctable := make([]rune, MAPSIZE)
	for i := 0; i < len(from); i++ {
		f := from[i]
		if int(f) >= len(ctable) {
			cnew := make([]rune, f+MMARGIN)
			copy(cnew, ctable)
			ctable = cnew
		}
		if i < len(into) {
			ctable[f] = into[i] + 1
		} else {
			ctable[f] = -1
		}
	}

	// compute the result
	j := 0 // j is the output index
	for i := 0; i < len(s); i++ {
		c := s[i]
		if int(c) < len(ctable) { // if entry is in table
			t := ctable[c] // get entry value
			if t < 0 {
				continue // discard input character
			} else if t > 0 {
				c = t - 1 // map to new character
			} // else leave alone
		}
		s[j] = c // save result character
		j++      // bump store index
	}
	return Return(RuneString(s[:j]))
}
