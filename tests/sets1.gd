#SRC: Goaldi original
#
#   Test set operations

procedure main() {
	testset("empty",    set())
	testset("pidigits", ^S := set([3,1,4,1,5]))
	testset("nodelete", S.delete())
	testset("delete",   S.delete(3,7,5))
	testset("noput",    S.put())
	testset("put",      S.put(4, 7,9))
	testset("delete",   S.delete(9, 1))
	testset("put",      S.put(2,4,6,8))
	testset("delete",   S.delete(2,5,7,9))
	testset("S@:x",     S @: 3 & S @: 1 & S)
	testset("strings",	set(["three","one","four","one","five"]))
	testset("mixed",	set([,1,"two",channel(3),%stdin,type,main]))
	every ^S2 := set([] | [0,2,4,6,8]) do {
		every ^S3 := set([] | [0,3,6,9]) do {
			write("S2 = ", image(S2), "  S3 = ", image(S3))
			testset("S2 ++ S3", S2 ++ S3)
			testset("S2 -- S3", S2 -- S3)
			testset("S2 ** S3", S2 ** S3)
		}
	}
}

#   print set contents and run some tests
procedure testset(label, S) {
	#  show label, short string, size, and image
	writes(left((label || ":"), 10), S, " (", *S, ") ", image(S), " :")
	#  look for and print small numbers (two different ways)
	every writes(" ", S.member(0 to 9))
	writes(" :")
	every writes(" ", S[0 to 9])
	write()

	#  run some tests and print error if results don't match
	cksame("not self", S, S)
	cksame("copy(S)", S, copy(S))
	cksame("S.copy()", S, S.copy())
	cksame("S.sort()", S, set(S.sort()))
	^L := []
	every L.put(!S)
	cksame("!S", S, set(L))
	L := []
	^S2 := S.copy()
	while L.put(@S2)
	cksame("@S", S, set(L))
	L := []
	S2 := set()
	while *S2 < *S do
		S2.put(?S)
	cksame("?S", S, S2)
}

#	check that two sets are the same by comparing their images
procedure cksame(label, S1, S2) {
	^im1 := image(S1)
	^im2 := image(S2)
	if im1 ~== im2 then {
		write("   ERROR: ", label, ": ", im1, " ~=== ", im2)
	}
}
