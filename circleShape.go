package chipmunk

type CircleShape struct {
	Shape *Shape
	// Center of the circle. Call Update() on the parent shape if changed.
	Position Vect
	// Radius of the circle. Call Update() on the parent shape if changed.
	Radius float32
	// Global center of the circle. Do not touch!
	Tc Vect
}

// Creates a new CircleShape with the given center and radius.
func NewCircle(pos Vect, radius float32) *Shape {
	shape := newShape()
	circle := &CircleShape{
		Position: pos,
		Radius:   float32(radius),
		Shape:    shape,
	}
	shape.ShapeClass = circle
	return shape
}

// Returns ShapeType_Circle. Needed to implemet the ShapeClass interface.
func (circle *CircleShape) ShapeType() ShapeType {
	return ShapeType_Circle
}

func (circle *CircleShape) Moment(mass float32) float32 {
	return (float32(mass) * (0.5 * (circle.Radius * circle.Radius))) + LengthSqr(circle.Position)
}

// Recalculates the global center of the circle and the the bounding box.
func (circle *CircleShape) update(xf Transform) AABB {
	//global center of the circle
	center := xf.TransformVect(circle.Position)
	circle.Tc = center
	rv := Vect{circle.Radius, circle.Radius}

	return AABB{
		Sub(center, rv),
		Add(center, rv),
	}
}

// Returns ShapeType_Box. Needed to implemet the ShapeClass interface.
func (circle *CircleShape) Clone(s *Shape) ShapeClass {
	clone := *circle
	clone.Shape = s
	return &clone
}

// Returns true if the given point is located inside the circle.
func (circle *CircleShape) TestPoint(point Vect) bool {
	d := Sub(point, circle.Tc)

	return Dot(d, d) <= circle.Radius*circle.Radius
}
