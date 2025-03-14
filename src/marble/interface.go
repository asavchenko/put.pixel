package marble

type Object interface {
	Show()
	Rotate(float64)
	RotateTo(x, y int)
	Move()
	SetSpeed(float64)
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
	WillIntersectsWith(Object) (bool, float64, float64)
	GetDistanceTo(...interface{}) float64
	GetSpeed() float64
	IsInside(float64, float64) bool
	Contains(x, y float64) bool
}

type Group interface {
	GetMarbles() []Object
	Add(Object)
	Contains(x, y float64) bool
	IntersectsWith(Object) bool
	GetLeftBorder() int
	GetRightBorder() int
	GetBottomBorder() int
	MoveDown()
	GetIntersection(Object) (bool, float64, float64)
	Show()
	RemoveMatched(Object)
}
