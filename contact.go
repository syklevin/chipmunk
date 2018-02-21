package chipmunk

type Contact struct {
	p, n Vect
	dist float32

	r1, r2               Vect
	nMass, tMass, bounce float32

	jnAcc, jtAcc, jBias float32
	bias                float32

	hash HashValue
}

func (con *Contact) reset(pos, norm Vect, dist float32, hash HashValue) {
	con.p = pos
	con.n = norm
	con.dist = dist
	con.hash = hash

	con.jnAcc = 0.0
	con.jtAcc = 0.0
	con.jBias = 0.0
}

func (con *Contact) Normal() Vect {
	return con.n
}

func (con *Contact) Position() Vect {
	return con.p
}
