package chipmunk

import (
	"log"
	//"fmt"
	"math"
)

type PolygonAxis struct {
	// The axis normal.
	N Vect
	D float32
}

type PolygonShape struct {
	Shape *Shape
	// The raw vertices of the polygon. Do not touch!
	// Use polygon.SetVerts() to change this.
	Verts Vertices
	// The transformed vertices. Do not touch!
	TVerts Vertices
	// The axes of the polygon. Do not touch!
	Axes []PolygonAxis
	// The transformed axes of the polygon Do not touch!
	TAxes []PolygonAxis
	// The number of vertices. Do not touch!
	NumVerts int
}

// Creates a new PolygonShape with the given vertices offset by offset.
// Returns nil if the given vertices are not valid.
func NewPolygon(verts Vertices, offset Vect) *Shape {
	if verts == nil {
		log.Printf("Error: no vertices passed!")
		return nil
	}

	shape := newShape()
	poly := &PolygonShape{Shape: shape}

	poly.SetVerts(verts, offset)

	shape.ShapeClass = poly
	return shape
}

func (poly *PolygonShape) Moment(mass float32) float32 {

	sum1 := float32(0)
	sum2 := float32(0)

	println("using bad Moment calculation")
	offset := Vect{0, 0}

	for i := 0; i < poly.NumVerts; i++ {

		v1 := Add(poly.Verts[i], offset)
		v2 := Add(poly.Verts[(i+1)%poly.NumVerts], offset)

		a := Cross(v2, v1)
		b := Dot(v1, v1) + Dot(v1, v2) + Dot(v2, v2)

		sum1 += a * b
		sum2 += a
	}

	return (float32(mass) * sum1) / (6.0 * sum2)
}

// Sets the vertices offset by the offset and calculates the PolygonAxes.
func (poly *PolygonShape) SetVerts(verts Vertices, offset Vect) {

	if verts == nil {
		log.Printf("Error: no vertices passed!")
		return
	}

	if verts.ValidatePolygon() == false {
		log.Printf("Warning: vertices not valid")
	}

	numVerts := len(verts)
	oldnumVerts := len(poly.Verts)
	poly.NumVerts = numVerts

	if oldnumVerts < numVerts {
		//create new slices
		poly.Verts = make(Vertices, numVerts)
		poly.TVerts = make(Vertices, numVerts)
		poly.Axes = make([]PolygonAxis, numVerts)
		poly.TAxes = make([]PolygonAxis, numVerts)

	} else {
		//reuse old slices
		poly.Verts = poly.Verts[:numVerts]
		poly.TVerts = poly.TVerts[:numVerts]
		poly.Axes = poly.Axes[:numVerts]
		poly.TAxes = poly.TAxes[:numVerts]
	}

	for i := 0; i < numVerts; i++ {
		a := Add(offset, verts[i])
		b := Add(offset, verts[(i+1)%numVerts])
		n := Normalize(Perp(Sub(b, a)))

		poly.Verts[i] = a
		poly.Axes[i].N = n
		poly.Axes[i].D = Dot(n, a)
	}
}

// Returns ShapeType_Polygon. Needed to implemet the ShapeClass interface.
func (poly *PolygonShape) ShapeType() ShapeType {
	return ShapeType_Polygon
}

func (poly *PolygonShape) Clone(s *Shape) ShapeClass {
	return poly.Clone2(s)
}

func (poly *PolygonShape) Clone2(s *Shape) *PolygonShape {
	clone := *poly
	clone.Verts = make(Vertices, len(poly.Verts))
	clone.TVerts = make(Vertices, len(poly.TVerts))
	clone.Axes = make([]PolygonAxis, len(poly.Axes))
	clone.TAxes = make([]PolygonAxis, len(poly.TAxes))

	clone.Verts = append(clone.Verts, poly.Verts...)
	clone.TVerts = append(clone.TVerts, poly.TVerts...)
	clone.Axes = append(clone.Axes, poly.Axes...)
	clone.TAxes = append(clone.TAxes, poly.TAxes...)

	clone.Shape = s

	return &clone
}

// Calculates the transformed vertices and axes and the bounding box.
func (poly *PolygonShape) update(xf Transform) AABB {
	//transform axes
	{
		src := poly.Axes
		dst := poly.TAxes

		for i := 0; i < poly.NumVerts; i++ {
			n := xf.RotateVect(src[i].N)
			dst[i].N = n
			dst[i].D = Dot(xf.Position, n) + src[i].D
		}
		/*
			fmt.Println("")
			fmt.Println("Started Axes")
			fmt.Println(xf.Rotation, xf.Position)
			for i:=0;i<poly.NumVerts;i++ {
				fmt.Println(src[i], dst[i])
			}
		*/
	}
	//transform verts
	{
		inf := float32(math.Inf(1))
		aabb := AABB{
			Lower: Vect{inf, inf},
			Upper: Vect{-inf, -inf},
		}

		src := poly.Verts
		dst := poly.TVerts

		for i := 0; i < poly.NumVerts; i++ {
			v := xf.TransformVect(src[i])

			dst[i] = v
			aabb.Lower.X = FMin(aabb.Lower.X, v.X)
			aabb.Upper.X = FMax(aabb.Upper.X, v.X)
			aabb.Lower.Y = FMin(aabb.Lower.Y, v.Y)
			aabb.Upper.Y = FMax(aabb.Upper.Y, v.Y)
		}

		/*
			fmt.Println("Verts")
			for i:=0;i<poly.NumVerts;i++ {
				fmt.Println(src[i], dst[i])
			}
		*/
		return aabb
	}
}

// Returns true if the given point is located inside the box.
func (poly *PolygonShape) TestPoint(point Vect) bool {
	return poly.ContainsVert(point)
}

func (poly *PolygonShape) ContainsVert(v Vect) bool {
	for _, axis := range poly.TAxes {
		dist := Dot(axis.N, v) - axis.D
		if dist > 0.0 {
			return false
		}
	}

	return true
}

func (poly *PolygonShape) ContainsVertPartial(v, n Vect) bool {
	for _, axis := range poly.TAxes {
		if Dot(axis.N, n) < 0.0 {
			continue
		}
		dist := Dot(axis.N, v) - axis.D
		if dist > 0.0 {
			return false
		}
	}

	return true
}

func (poly *PolygonShape) ValueOnAxis(n Vect, d float32) float32 {
	verts := poly.TVerts
	min := Dot(n, verts[0])

	for i := 1; i < poly.NumVerts; i++ {
		min = FMin(min, Dot(n, verts[i]))
	}

	return min - d
}
