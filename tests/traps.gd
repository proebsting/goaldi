#SRC: icon/traps.icn
# test assignments to trapped variables
# original source unknown; found 2013 in an ancient to-do collection

procedure tvtbl_test () {
local T

	#
	# Test to make sure that the table trapped variable returns
	# the correct value.
	#
	# Old Icon Note: "The parameters to write are not de-referenced
	# until all of them are evaluated.  Any line produced by this section
	# that has includes two different values for T [] is therefore incorrect."
	#
	# In Goaldi, the rules are different:
	# Each write() argument is dereferenced as it is produced,
	# so different values ARE to be expected.
	#
	write ( "TVTBL test 1" )
	T := table()

	write ( "Assignment test:  \t", T [], "\t", T [] := "Assigned" )
	write ( "Reassignment test:\t", T [], "\t", T [] := "Reassigned" )
	write ( "Deletion test:    \t", T [], "\t", T.delete(nil) & T [] )
	write ( "Insertion test:   \t", T [], "\t",
		T [] := ( ( T[] := "Assigned" ) & "Reassigned" ) )

	#
	# Test to make sure that the table is getting updated properly by
	# trapped variable assignment.
	#
	# Note: there have been past errors where "T [] :=..." returns the
	# correct value without properly updating the table.
	#
	write ( "\nTVTBL test 2" )
	T.delete()
	T [] := "Assigned";   write ( "Assignment test:  \t",   T [] )
	T [] := "Reassigned"; write ( "Reassignment test:\t", T [] )
	T [] := "Assigned"
	T [] := ( T.delete() & "Reassigned" )
	write ( "Deletion test:    \t" , T [] )
	T.delete()
	T [] := ( ( T[] := "Assigned" ) & "Reassigned" )
	write ( "Insertion test:   \t", T [] )
	write ( )

}


procedure subs_test (  ) {
# local T, s
local T
local s

	write ( "TVSUBS test" )
	T := table()
	T [ 7 ] := "....."
	T [ 7 ] [ 4 ] := "X"
	write ( "Subs of new table elem:      ", T [ 7 ] )
	T [ 7 ] := "....."
	T [ 7 ] [ 4 ] := "X"
	write ( "Subs of existing table elem: ", T [ 7 ] )

	# Lots more cases should be added here.

	return
}


procedure main (  ) {
	tvtbl_test ( )
	subs_test ( )
	return
}
