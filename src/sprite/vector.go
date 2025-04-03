package sprite

import (
	"math"
	"strings"

	"assa.com/put.pixel/lib/convertor"
)

type vector struct {
	angle             float64
	speed             float64
	applyNumCycles    int
	currentCycle      int
	direction         int
	allowedDirections []int
}

func GetDefaultVectors(frames []Frame, keepCycles int) []Vector {
	vectors := make([]Vector, len(frames))
	for i := 0; i < len(frames); i++ {
		frameName := frames[i].GetImgName()
		nameParts := strings.Split(frameName, "_")
		speed := float64(0)
		angle := float64(0)
		if len(nameParts) > 2 {
			speed = convertor.ToFloat64IgnoreErrors(nameParts[1])
			angle = convertor.ToFloat64IgnoreErrors(nameParts[2]) * math.Pi / 180
			keepCycles = convertor.ToIntIgnoreErrors(strings.Replace(nameParts[3], ".png", "", -1))
		}
		vectors[i] = GetNewVector(angle, speed, keepCycles)
		vectors[i].SetAllowedDirections([]int{DIRECTION_RIGHT, DIRECTION_LEFT})
		vectors[i].SetDirection(frames[i].GetDirection())
	}

	return vectors
}

func GetNewVector(a float64, s float64, durationInCycles int) Vector {
	v := &vector{}

	v.SetDurationInCycles(durationInCycles)
	v.SetSpeed(s)
	v.SetAngle(a)

	return v
}

func (v *vector) SetAllowedDirections(directions []int) Vector {
	v.allowedDirections = directions

	return v
}

func (v *vector) SetDirection(direction int) Vector {
	found := false
	for _, ad := range v.allowedDirections {
		if ad == direction {
			found = true
			break
		}
	}
	if !found {
		return v
	}

	if direction == DIRECTION_LEFT && v.direction == DIRECTION_RIGHT {
		v.flip()
	} else if direction == DIRECTION_RIGHT && v.direction == DIRECTION_LEFT {
		v.flip()
	}

	v.direction = direction

	return v
}

func (v *vector) SetAngle(angle float64) Vector {
	v.angle = angle

	return v
}

func (v *vector) SetSpeed(speed float64) Vector {
	v.speed = speed

	return v
}

func (v *vector) SetDurationInCycles(n int) Vector {
	v.applyNumCycles = n

	return v
}

func (v *vector) IsApplied() bool {
	return v.currentCycle >= v.applyNumCycles
}

func (v *vector) GetSpeed() float64 {
	return v.speed
}

func (v *vector) Reset() Vector {
	v.currentCycle = 0

	return v
}

func (v *vector) Apply(x, y float64) (float64, float64) {
	x += math.Cos(v.angle) * v.speed
	y += math.Sin(v.angle) * v.speed
	v.currentCycle++

	return x, y
}

func (v *vector) flip() {
	// normalize
	if v.angle >= math.Pi*2 {
		v.angle -= math.Pi * 2 * float64(int(v.angle/(math.Pi*2)))
	}
	if v.angle <= -math.Pi*2 {
		v.angle += math.Pi * 2 * float64(int(-v.angle/(math.Pi*2)))
	}

	switch v.direction {
	case DIRECTION_LEFT:
		v.angle -= math.Pi - v.angle*2
	case DIRECTION_RIGHT:
		v.angle += math.Pi - v.angle*2
	}
}
