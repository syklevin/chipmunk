package chipmunk

import (
	"log"
)

func k_scalar_body(body *Body, r, n Vect) float32 {
	rcn := Cross(r, n)
	return body.m_inv + (body.i_inv * rcn * rcn)
}

func k_scalar(a, b *Body, r1, r2, n Vect) float32 {
	value := k_scalar_body(a, r1, n) + k_scalar_body(b, r2, n)
	if value == 0.0 {
		log.Printf("Warning: Unsolvable collision or constraint.")
	}
	return value
}

func k_scalar2(a, b *Body, r1, r2, n Vect) float32 {
	rcn := (r1.X * n.Y) - (r1.Y * n.X)
	rcn = a.m_inv + (a.i_inv * rcn * rcn)

	rcn2 := (r2.X * n.Y) - (r2.Y * n.X)
	rcn2 = b.m_inv + (b.i_inv * rcn2 * rcn2)

	value := rcn + rcn2
	if value == 0.0 {
		log.Printf("Warning: Unsolvable collision or constraint.")
	}
	return value
}

func relative_velocity2(a, b *Body, r1, r2 Vect) Vect {
	v1 := Add(b.v, Mult(Perp(r2), b.w))
	v2 := Add(a.v, Mult(Perp(r1), a.w))
	return Sub(v1, v2)
}

func relative_velocity(a, b *Body, r1, r2 Vect) Vect {
	return Vect{(-r2.Y*b.w + b.v.X) - (-r1.Y*a.w + a.v.X), (r2.X*b.w + b.v.Y) - (r1.X*a.w + a.v.Y)}
}

func normal_relative_velocity(a, b *Body, r1, r2, n Vect) float32 {
	return Dot(relative_velocity(a, b, r1, r2), n)
}

func apply_impulses(a, b *Body, r1, r2, j Vect) {
	j1 := Vect{-j.X, -j.Y}

	a.v.Add(Mult(j1, a.m_inv))
	a.w += a.i_inv * Cross(r1, j1)

	b.v.Add(Mult(j, b.m_inv))
	b.w += b.i_inv * Cross(r2, j)
}

func apply_bias_impulses(a, b *Body, r1, r2, j Vect) {

	j1 := Vect{-j.X, -j.Y}

	a.v_bias.Add(Mult(j1, a.m_inv))
	a.w_bias += a.i_inv * Cross(r1, j1)

	b.v_bias.Add(Mult(j, b.m_inv))
	b.w_bias += b.i_inv * Cross(r2, j)
}
