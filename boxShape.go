package chipmunk

// Convenience wrapper around PolygonShape.
type BoxShape struct {
	Shape *Shape
	// The polygon that represents this box. Do not touch!
	Polygon *PolygonShape
	verts   [4]Vect
	// The width of the box. Call UpdatePoly() if changed.
	Width float32
	// The height of the box. Call UpdatePoly() if changed.
	Height float32
	// The center of the box. Call UpdatePoly() if changed.
	Position Vect
}

// Creates a new BoxShape with given position, width and height.
func NewBox(pos Vect, w, h float32) *Shape {
	shape := newShape()

	box := &BoxShape{
		Polygon:  &PolygonShape{Shape: shape},
		Width:    w,
		Height:   h,
		Position: pos,
		Shape:    shape,
	}

	hw := w / 2.0
	hh := h / 2.0

	if hw < 0 {
		hw = -hw
	}
	if hh < 0 {
		hh = -hh
	}

	box.verts = [4]Vect{
		{-hw, -hh},
		{-hw, hh},
		{hw, hh},
		{hw, -hh},
	}

	poly := box.Polygon
	poly.SetVerts(box.verts[:], box.Position)

	shape.ShapeClass = box
	return shape
}

func (box *BoxShape) Moment(mass float32) float32 {
	return (float32(mass) * (box.Width*box.Width + box.Height*box.Height) / 12.0)
}

// Recalculates the internal Polygon with the Width, Height and Position.
func (box *BoxShape) UpdatePoly() {
	hw := box.Width / 2.0
	hh := box.Height / 2.0

	if hw < 0 {
		hw = -hw
	}
	if hh < 0 {
		hh = -hh
	}

	box.verts = [4]Vect{
		{-hw, -hh},
		{-hw, hh},
		{hw, hh},
		{hw, -hh},
	}

	poly := box.Polygon
	poly.SetVerts(box.verts[:], box.Position)
}

// Returns ShapeType_Box. Needed to implemet the ShapeClass interface.
func (box *BoxShape) ShapeType() ShapeType {
	return ShapeType_Box
}

// Returns ShapeType_Box. Needed to implemet the ShapeClass interface.
func (box *BoxShape) Clone(s *Shape) ShapeClass {
	clone := *box
	clone.Polygon = &PolygonShape{Shape: s}
	clone.Shape = s
	clone.UpdatePoly()
	return &clone
}

// Recalculates the transformed vertices, axes and the bounding box.
func (box *BoxShape) update(xf Transform) AABB {
	return box.Polygon.update(xf)
}

// Returns true if the given point is located inside the box.
func (box *BoxShape) TestPoint(point Vect) bool {
	return box.Polygon.TestPoint(point)
}
