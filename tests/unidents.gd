#SRC: Goaldi original
#
#  Using Unicode identifiers in various situations

global Dvořák := 1841

record rtype(Σ, Ω)

procedure main(args[]) {
	local Strauß := 1825
	^π := %pi
	with %ϕ := %phi do {
		println(Dvořák, Strauß, π, %ϕ)
	}
	^R := rtype(123,456)
	println(image(R), R.Σ, R.Ω)
	R := rtype(Ω:999, Σ:888)
	println(image(R), R.Σ, R.Ω)
	^C := constructor("greeks", "α", "β", "γ")
	println(image(C))
	println(image(C(1,2,3)))
	^T := tuple(Ψ:"saguaro", Ξ:"gate", Π:"wicket",)
	println(image(T))
	^t1 := 1017
	^t2 := 1023
	^Δt := t2 - t1
	println(t1, t2, Δt)
}
