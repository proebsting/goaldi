#SRC: goaldi original
#  regex demo

procedure main() {
	# local c, v, p
	local c
	local v
	local p

	rex("a(x*)b(y|z)c", "-axxxbyc-", "-abzc-")
	rex("(a|bcdef|g|ab|c|d|e|efg|fg)*", "abcdefg", )
	rex(`\d+(\.\d*)?(e\d+)?`, "5", "2.71", "3e9", "x59", "16r99", "eleven")
	v := "([aeiou]*)"
	c := "([bcdfghj-np-tv-z]*)"
	p := "p" ||  v || c || "ch"
	rex(p, "punch", "patch", "peach", "pitch", "porch", "pooch", "prunch", )
}

procedure rex(expr, s[]) {
	local e
	if e := regex(expr) then every try(e, !s) else write("FAILED: ", expr)
}

procedure try(re, s) {
	local a
	writes(re, " : ", image(s), " :")
	if a := \re.FindStringSubmatch(s) then {
		every writes(" ", image(!a) | "\n")
	} else {
		write(" [no match]")
	}
	return
}
