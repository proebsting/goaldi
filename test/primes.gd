#SRC: icon/primes.icn

# a simple and slow prime number generator

procedure main() {
	local i
	every i := 2 to 100 do {
		if i % (2 to i - 1) = 0 then
			continue
		write(i)
	}
}
