package sprite

const (
	DIRECTION_RIGHT = 1
	DIRECTION_LEFT  = -1
	DIRECTION_BACK  = 0
	DIRECTION_FRONT = 6
)

type Object interface { // or maybe better call it Sprite
	Start(actionName string) error
	Stop() error
	Show() Object // will be called each gl.draw cycle
	AddAction(Action) Object
	SetDefaultAction(Action) Object
	SetPosition(x, y float64) Object
	X() float64
	Y() float64
	Cx() int
	Cy() int
	SetDirection(direction int) Object
	SetAllowedDirections([]int) Object
}

type Action interface {
	Stop() error // not everything can be stopped
	AllowStopping(bool) Action
	SetFrames(frames []Frame, vectors []Vector) Action
	SetX(float64) Action
	SetY(float64) Action
	SetName(string) Action
	GetName() string
	Reset() error
	Run() (Frame, float64, float64)
	IsComplete() bool
	SetDirection(direction int) Action
	SetAllowedDirections([]int) Action
}

type Frame interface {
	SetImage(pathToImg string) error
	SetSize(width, height int) Frame
	GetWidth() int
	GetHeight() int
	Show(x, y int) Frame
	SetDirection(int) Frame
	GetDirection() int
	SetAllowedDirections([]int) Frame
	GetImgName() string
	SetImgName(string) Frame
}

type Vector interface {
	SetAngle(angle float64) Vector
	SetSpeed(speed float64) Vector
	GetSpeed() float64
	SetDurationInCycles(n int) Vector
	Apply(x, y float64) (float64, float64)
	IsApplied() bool
	Reset() Vector
	SetDirection(direction int) Vector
	SetAllowedDirections([]int) Vector
}
