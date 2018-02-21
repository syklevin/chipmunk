package chipmunk

// Wrapper around []
type Vertices []Vect

// Checks if verts forms a valid polygon.
// The vertices must be convex and winded clockwise.
func (verts Vertices) ValidatePolygon() bool {
	numVerts := len(verts)
	for i := 0; i < numVerts; i++ {
		a := verts[i]
		b := verts[(i+1)%numVerts]
		c := verts[(i+2)%numVerts]

		if Cross(Sub(b, a), Sub(c, b)) > 0.0 {
			return false
		}
	}

	return true
}
