package sprite

import (
	"math"
)

type vector struct {
	angle             float64
	speed             float64
	applyNumCycles    int
	currentCycle      int
	direction         int
	allowedDirections []int
}

func GetDefaultVectors(n int, keepCycles int) []Vector {
	vectors := make([]Vector, n)
	for i := 0; i < n; i++ {
		vectors[i] = GetNewVector(0, 0, keepCycles)
		vectors[i].SetAllowedDirections([]int{DIRECTION_RIGHT, DIRECTION_LEFT})
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
	if !v.IsApplied() && v.direction != direction {
		return v
	}

	if direction != v.direction {
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
