package sprite

import (
	"fmt"
)

type action struct {
	x, y              float64 // current position
	frames            []Frame
	vectors           []Vector // the length is equal to the numCycles
	name              string   // move_left, move_right, jump_forward, idle, smoke... etc
	currentFrame      Frame
	currentFrameIdx   int // defines the next frame, when vectors from the current frame all applied
	canBeStopped      bool
	direction         int
	allowedDirections []int
}

func GetNewAction(name string, frames []Frame, vectors []Vector, isStoppingAllowed bool) Action {
	a := &action{}
	a.AllowStopping(isStoppingAllowed)
	a.SetFrames(frames, vectors)
	a.SetName(name)

	return a
}

func (a *action) AllowStopping(b bool) Action {
	a.canBeStopped = b

	return a
}

func (a *action) SetFrames(frames []Frame, vectors []Vector) Action {
	a.frames = frames
	a.vectors = vectors

	return a
}

func (a *action) SetX(f float64) Action {
	a.x = f

	return a
}

func (a *action) SetY(f float64) Action {
	a.y = f

	return a
}

func (a *action) SetName(s string) Action {
	a.name = s

	return a
}

func (a *action) GetName() string {
	return a.name
}

func (a *action) Reset() error {
	if !a.canBeStopped && a.currentFrame != nil {
		return fmt.Errorf("%s can't be reset", a.name)
	}
	a.currentFrameIdx = 0
	a.currentFrame = nil

	return nil
}

func (a *action) Stop() error {
	if !a.canBeStopped && !a.IsComplete() {
		return fmt.Errorf("%s can't be stopped", a.name)
	}

	a.currentFrameIdx = 0
	a.currentFrame = nil

	return nil
}

func (a *action) Run() (Frame, float64, float64) {
	v := a.vectors[a.currentFrameIdx]
	v.SetDirection(a.direction)
	if v.IsApplied() {
		v.Reset()
		a.currentFrameIdx++
		if a.currentFrameIdx >= len(a.frames) {
			a.currentFrameIdx = 0
		}
		v = a.vectors[a.currentFrameIdx]
	}
	a.currentFrame = a.frames[a.currentFrameIdx]

	a.x, a.y = v.Apply(a.x, a.y)

	return a.currentFrame.SetDirection(a.direction), a.x, a.y
}

func (a *action) IsComplete() bool {
	return a.currentFrameIdx >= len(a.frames)-1
}

func (a *action) SetDirection(direction int) Action {
	found := false
	for _, ad := range a.allowedDirections {
		if ad != direction {
			continue
		}

		found = true
		break
	}

	if !found {
		return a
	}

	a.direction = direction

	return a
}

func (a *action) SetAllowedDirections(directions []int) Action {
	a.allowedDirections = directions

	return a
}
