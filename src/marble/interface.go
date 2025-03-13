package marble

type Object interface {
	Show()
	Rotate(float64)
	RotateTo(x, y int)
	Move(int)
	MoveTo(float64, float64)
	X() float64
	Y() float64
	CX() int
	CY() int
	A() float64
	R() float64
	CR() int
	D() float64
	CD() int
	SetAngle(float64)
	Color() byte
	IntersectsWith(Object) bool
	GetDistanceTo(Object) float64
	IsInside(float64, float64) bool
}
