#SRC: icon/sorting.icn

# test sorting and copying


global letters
global digits

procedure main(args) {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits := "0123456789"
	listtest()
	# rectest()  # disabled: records don't sort
	tabletest()
	copytest()
	messtest()
}



#  listtest() -- test sorting of lists and sets

procedure listtest() {
	# local n, x, S, L1, L2, L3
	local n
	local x
	local S
	local L1
	local L2
	local L3

	every n := (0 to 10) | 23 | 47 | 91 do {

		L1 := list()
		write(n, ":")
		every !n do
			L1.put(randval())

		L2 := L1.sort()
		L3 := L2.sort()
		check(L2, L3)

		L2 := L1.copy().sort()
		L3 := L2.copy().sort()
		check(L2, L3)

	}
}


#  rectest() -- test sorting of records

record r0()
record r1(a)
record r2(a,b)
record r5(a,b,c,d,e)

procedure rectest() {
	write()
	wlist(r0().sort())
	wlist(r0().copy().sort())
	wlist(r1(12)).sort()
	wlist(r2(5,2)).sort()
	wlist(r5(2,7,1,8,3)).sort()
	wlist(r5(3,1,4,1,6)).sort()
	wlist(r5("t","e","p","a","d").sort())
	wlist(r5("t","e","p","a","d").copy().sort())
	return
}


#  tabletest() -- test sorting of tables

procedure tabletest() {
	# local T, L, t
	local T
	local L
	local t

	T := table()
	T[7] := "h"
	T[2] := "a"
	T[8] := "r"
	T[0] := "e"
	T[3] := "o"
	T[6] := "s"
	T[5] := "n"
	T[1] := "t"
	T[4] := "i"
	T[9] := "d"

	write()
	L := T.sort();  every t:= !L do writes(" ", t.key, " ", t.value); write()
	L := T.sort(1); every t:= !L do writes(" ", t.key, " ", t.value); write()
	L := T.sort(2); every t:= !L do writes(" ", t.key, " ", t.value); write()

	T := T.copy()
	L := T.sort();  every t:= !L do writes(" ", t.key, " ", t.value); write()
	L := T.sort(1); every t:= !L do writes(" ", t.key, " ", t.value); write()
	L := T.sort(2); every t:= !L do writes(" ", t.key, " ", t.value); write()
	return
}



#  randval() -- return random integer, real, or string value

procedure randval() {
	return case ?3 of {
		1:  ?999					# 000 - 999
		2:  ?99 / 10.0				# 0.0 - 9.9
		3:  ?letters || ?letters || ?letters	# "AAA" - "ZZZ"
		}
}


#  check that two lists have identical components
#  and that they are in ascending order

procedure check(a, b) {
	# local i, ai, ai1, bi, d
	local i
	local ai
	local ai1
	local bi
	local d

	if *a ~= *b then
		stop("different sizes: ", image(a), " / ", image(b))
	every i := 1 to *a do {
		ai := a[i]
		bi := b[i]
		ai1 := a[i-1] | nil
		if ai ~=== bi then
			stop("element ", i, " differs")
		if type(ai) === type(ai1) then {
			case type(ai) of {
				"integer":	d := (ai1 > ai) | nil
				"real":	d := (ai1 > ai) | nil
				"string":	d := (ai1 >> ai) | nil
			}
		stop("element ", i, " out of order: ", image(\d))
		}
	}
	return
}


#  write list

procedure wlist(L) {
	writes(*L, ":")
	every writes(" ", !L | "\n")
	return
}



#  test copy(), especially that copies are really distinct

procedure copytest() {
	# local L1, L2, S1, S2, T1, T2, R1, R2
	local L1
	local L2
	local S1
	local S2
	local T1
	local T2
	local R1
	local R2

	write()

	L1 := [1,2,3]
	L1.push(L1)
	L2 := L1.copy()
	L2.pull()
	L2.put(4)
	every writes(" ", "L1:" | image(!L1) | "\n")
	every writes(" ", "L2:" | image(!L2) | "\n")

	T1 := table()
	T1[2] := "j"
	T1[5] := "c"
	T1[8] := "n"
	T1[15] := T1
	T2 := T1.copy()
	T2.delete(5)
	T2[11] := "t"
	every writes(" ", "T1:" | image(!T1.sort()) | "\n")
	every writes(" ", "T2:" | image(!T2.sort()) | "\n")

	return
}



#  sort different types together

procedure messtest() {
	# local L0, L1, L2, L3
	local L0
	local L1
	local L2
	local L3

	write()
	L0 := []
	L1 := [
		"", "0cs", 4.4, 2.2, "a", nil, number, wlist, "epsilons", L0.put,
		r0, "delta", image, "beta", table(5), [], write, "123cs", [3,4], -3^41,
		image, table()[4]:=7, %stdin, 3.3, reverse, r1(1), [], table(),
		r5, r5(1,23), nil,
		create 1 | 2, 5.5,
		"", r2(5,6), -7^23, L0.get,
		"epsilon", [1,2,3], r5(7,8,9), r2, %stdout, 4, , 1,
		r5(1,2,3), r1, check,
		create 3 | 4,
		"XYZcs", 1.1, r1(5), 5^28, L0.push,
		"1234cs", 5, r0(), read, "gamma", r5(4,5,6,7,8), 2,
		create 5 to 7,
		table, r2(1,2), toupper, r0(), "alpha", messtest, %stderr, 11^19,
		listtest, "gamma", main, 3, L0.pop ]
	L1.put(L1)
	L2 := L1.copy()
	every L1.put((!L2).copy())

	write()
	every write(image(!L1.sort()))

	wsortf(L1, 2)
	return
}

procedure wsortf(L, n) {
	# local e, s
	local e
	local s

	write()
	every e := !L.sort(n) do {
		s := image(e)
		if type(e) === (list | type) then
			writes("key=", image(e[n]), " ")	# may fail
		write(s)
		}
	return
}
