#SRC: icon/sets.icn

# set test

procedure main() {
	local x
	local y
	local z

	wset("empty", x := set())
	write(type(x))
	write(image(?x))			# should fail
	write(image(x.member()))	# should fail
	wset("insert", x.insert(nil))
	write(image(?x))			# should write nil
	write(image(x.member()))	# should write nil
	write(image(x.member(3)))	# should write nil
	wset("insert", x.insert())
	wset("delete", x.delete())
	wset("delete", x.delete())
	write()

	wset("x", x := set([1,2,4]))
	wset("y", y := set([1,2,5]))
	wset("x ++ y", x ++ y)
	wset("y ++ x", y ++ x)
	wset("x -- y", x -- y)
	wset("y -- x", y -- x)
	wset("x ** y", x ** y)
	wset("y ** x", y ** x)
	write()

	wset("empty", x := set(nil))
	wset("+ 1", x.insert(1))	# only inserts 1
	wset("+ 2", x.insert(2))
	wset("+ c", x.insert("c"))
	wset("- 3", x.delete(3))		# deletes nothing
	wset("- 1", x.delete(1))		# only deletes 1
	wset("- 1", x.delete(1))
	wset("+ 2", x.insert(2))
	wset("+ 1", x.insert(1))
	wset("+ 7.0", x.insert(7.0))
	wset("+ 7.0", x.insert(7.0))
	wset(`+ "cs"`, x.insert("cs"))
	wset(`+ "cs"`, x.insert("cs"))
	wset("x =", x)
	write()

	wset("3,a,4", y := set([3,"a",4]))
	wset("y ++ x", y ++ x)
	wset("y ** x", y ** x)
	wset("y -- x", y -- x)
	wset("x -- y", x -- y)
	write()

	every (z := set()).insert(!y)
	wset("z from !y", z)

	write()
	x := set([3,1,4,1,5,9,2,6,5,3,5])
	y := copy(x)
	x.delete(4)
	x.insert(7)
	y.insert(0)
	y.delete(1)
	wset("x", x)
	wset("y", y)
}



#	dump a set, assuming it contains nothing other than:
#	nil, 0 - 9, "", "a" - "e", "cs"

procedure wset(label, S) {
	local x

	writes(right(label, 10), " :", right(*S, 2), " :")
	every x := nil | (0 to 9) | "" | !"abcde" | "cs" do
		writes(" ", image(S.member(x)))
	write()
	return
}
