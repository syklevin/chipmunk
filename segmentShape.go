package chipmunk

//If Settings.AutoUpdateShapes is not set, call Update on the parent shape for changes to the A, B and Radius to take effect.
type SegmentShape struct {
	Shape *Shape
	//start/end points of the segment.
	A, B Vect
	//radius of the segment.
	Radius float32

	//local normal. Do not touch!
	N Vect
	//transformed normal. Do not touch!
	Tn Vect
	//transformed start/end points. Do not touch!
	Ta, Tb Vect

	//tangents at the start/end when chained with other segments. Do not touch!
	A_tangent, B_tangent Vect
}

// Creates a new SegmentShape with the given points and radius.
func NewSegment(a, b Vect, r float32) *Shape {
	shape := newShape()
	seg := &SegmentShape{
		A:      a,
		B:      b,
		Radius: r,
		Shape:  shape,
	}
	shape.ShapeClass = seg
	return shape
}

// Returns ShapeType_Segment. Needed to implemet the ShapeClass interface.
func (segment *SegmentShape) ShapeType() ShapeType {
	return ShapeType_Segment
}

func (segment *SegmentShape) Moment(mass float32) float32 {

	offset := Mult(Add(segment.A, segment.B), 0.5)

	return float32(mass) * (DistSqr(segment.B, segment.A)/12.0 + LengthSqr(offset))
}

//Called to update N, Tn, Ta, Tb and the the bounding box.
func (segment *SegmentShape) update(xf Transform) AABB {
	a := xf.TransformVect(segment.A)
	b := xf.TransformVect(segment.B)
	segment.Ta = a
	segment.Tb = b
	segment.N = Perp(Normalize(Sub(segment.B, segment.A)))
	segment.Tn = xf.RotateVect(segment.N)

	rv := Vect{segment.Radius, segment.Radius}

	min := Min(a, b)
	min.Sub(rv)

	max := Max(a, b)
	max.Add(rv)

	return AABB{
		min,
		max,
	}
}

func (segment *SegmentShape) Clone(s *Shape) ShapeClass {
	clone := *segment
	clone.Shape = s
	return &clone
}

// Only returns false for now.
func (segment *SegmentShape) TestPoint(point Vect) bool {
	return false
}
