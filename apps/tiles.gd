#	tiles.gd -- number tile game
#
#	click a tile to move it to an adjacent empty space
#	swipe offscreen left to reset
#	swipe offscreen right to scramble
#	swipe off to upper left to quit

#	tile background stripe colors (not currently used)
global VCOLORS := ["#C00", "#00A", "#870", "#080"]
global HCOLORS := ["#FC0", "#9F4", "#9DF", "#DAF"]
#	number of tiles
global NWIDE := *VCOLORS	# across
global NHIGH := *HCOLORS	# down

global BACKGD := "#776"	# background color
global UPBEVL := "#EEE"	# bevel color (sunny side)
global DNBEVL := "#BBB"	# bevel color (shady side)
global BEVLW := 6		# bevel width

record tile (
	n,		# tile number
	i, j,	# current location (row and column)
	sprite,	# associated sprite
)

global app			# application canvas
global tsize		# size of tile on canvas, in points

global cell := list(NHIGH, list(NWIDE))	# tiles on grid

procedure main() {

	# setup
	randomize()
	app := canvas()
	app.color(BACKGD).Rect(-1000, -1000, 2000, 2000)
	tsize := (app.Width / app.PixPerPt) / NWIDE
	tsize >:= (app.Height / app.PixPerPt) / NHIGH
	makeTiles()

	# animate initial scrambling
	sleep(1)
	every randmove(8 | 5 | 3)
	every !10 do randmove(1)
	scramble()
	every !10 do randmove(1)
	every randmove(3 | 5 | 8)

	# main loop
	local istart
	local jstart
	while ^e := @app.Events do case e.Action of {
		"touch": {
			istart := ti(e.Y)
			jstart := tj(e.X)
		}
		"release": {
			^i := ti(e.Y)
			^j := tj(e.X)
			if /cell[0<i, 0<j] then {
				# ended in vacant cell; use start point instead
				i := istart
				j := jstart
			}
			if ^t := \cell[0<i, 0<j] then {
				if (^di := -1 to 1) & (^dj := -1 to 1) &
					(^n := t.canmove(di,dj)) then {
						t.move(n, di, dj, 7)
				}
			} else if j < 1 && i < 1 then {
				break
			} else if j < 1 then {
				reset()
			} else if j > NWIDE then {
				scramble()
			}
		}
		"stop":
			break
	}
}

procedure reset() {				#: restore initial configuration
	^tlist := [: \!!cell :]
	every !!cell := nil
	every ^t := !tlist do {
		^i := integer(1 + (t.n - 1) / NWIDE)
		^j := integer(1 + (t.n - 1) % NWIDE)
		cell[i, j] := t
		t.slideto(i, j, 3)
		t.i := i
		t.j := j
	}
}

procedure scramble() {			#: rearrange tiles randomly
	# We need a reachable configuration, so cannot arbitrarily place tiles.
	# Instead, make a large number of legal moves randomly.
	# 100 appears to do a thorough job, so we do more to be extra sure.
	every !250 do randmove()
}

procedure randmove(nsteps) {	#: move random tile, but not to undo previous mv
	/static prevdi := 0
	/static prevdj := 0
	repeat {	# just try until one works (#%#% inefficient)
		^t := \cell[?NHIGH+1,?NWIDE+1] | continue
		if (^di := -1 to 1) & (^dj := -1 to 1) & (^n := t.canmove(di,dj)) then {
			if di = -prevdi & dj = -prevdj then {
				continue
			}
			prevdi := di
			prevdj := dj
			return t.move(n, di, dj, nsteps)
		}
	}
}

procedure makeTiles() {			#: create tiles and place in cells
	^z := tsize / 100	# scaling with no gaps
	z *:= 0.98			# make a little gap
	every ^i := 1 to NHIGH do {
		every ^j := 1 to NWIDE do {
			^n := j + NWIDE * (i - 1)
			if i = NHIGH & j = NWIDE then break
			^c := canvas(100, 100, 1)
			# (don't use) background stripes
			# c.color(VCOLORS[j]).Rect(-30,-50,60,100)
			# c.color(HCOLORS[i]).Rect(-50,-25,100,50)
			# draw tile number label
			c.VFont := font("mono", 40)
			c.color("black").Text(-25, 15, right(n, 2))
			# draw bevels around edges
			c.Size := 2 * BEVLW
			c.color(DNBEVL).Goto(50,-50,90)
			every !2 do c.Forward(100).turn(90)
			c.color(UPBEVL)
			every !2 do c.Forward(100).turn(90)
			# create the tile and its sprite
			^t := tile(n:n, i:i, j:j)
			^x := tx(j)
			^y := ty(i)
			t.sprite := app.AddSprite(c.Canvas, x, y, z)
			cell[i,j] := t
		}
	}
}

procedure ti(y) {				#! compute row number from y-coordinate
	return integer(y / tsize + NHIGH / 2 + 1)
}

procedure tj(x) {				#! compute column number from x-coordinate
	return integer(x / tsize + NWIDE / 2 + 1)
}

procedure tx(j) {				#! compute x-coordinate of sprite in column j
	return (j - (NWIDE + 1) / 2) * tsize
}

procedure ty(i) {				#! compute y-coordinate of sprite in column i
	return (i - (NHIGH + 1) / 2) * tsize
}

procedure dumpCells() {			#: print a snapshot of current state
	every ^i := 1 to NHIGH do {
		writes("[",i,",*]:")
		every ^j := 1 to NWIDE do {
			writes(right((\cell[i,j]).n | "--", 3))
		}
		write()
	}
}

#	move ntiles (from vacant cell through self) by (di,dj) in nsteps
procedure tile.move(ntiles, di, dj, nsteps) {
	^vi := self.i + ntiles * di		# vacant cell location
	^vj := self.j + ntiles * dj
	every ^k := !ntiles do {
		^t := cell[vi - k * di, vj - k * dj]
		t.slideto(t.i + di, t.j + dj, nsteps)
		cell[t.i, t.j] := nil
		t.i +:= di
		t.j +:= dj
		cell[t.i, t.j] := t
	}
	return self
}

procedure tile.slideto(i, j, nsteps) {	#: slide tile to [i,j] in n steps
	^oldx := tx(self.j)
	^oldy := ty(self.i)
	^newx := tx(j)
	^newy := ty(i)
	if \nsteps > 0 then every ^d := !nsteps do {
		^f := d / nsteps
		^ff := 1 - f
		self.sprite.MoveTo(ff * oldx + f * newx, ff * oldy + f * newy,
			self.sprite.Scale)
		sleep(0.02)
	} else {
		self.sprite.MoveTo(newx, newy, self.sprite.Scale)
	}
	return self
}

procedure tile.canmove(di, dj) {	#: return number of tiles that can move
	if di ~= 0 & dj ~= 0 then return fail
	if di ~= 0 then {
		every ^n := 1 to NHIGH-1 do {
			^k := self.i + n * di
			if 1 <= k <= NHIGH & /cell[k, self.j] then {
				return n
			}
		}
	}
	if dj ~= 0 then {
		every ^n := 1 to NWIDE-1 do {
			^k := self.j + n * dj
			if 1 <= k <= NWIDE & /cell[self.i, k] then {
				return n
			}
		}
	}
	return fail
}
