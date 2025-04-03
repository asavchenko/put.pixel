package sprite

import (
	"fmt"
	"math"
)

type object struct {
	actions           map[string]Action // key action name
	activeAction      Action
	defaultAction     Action
	x, y              float64
	cx, cy            int
	direction         int   // -1 0 +1
	allowedDirections []int // -1 0 +1
}

func GetNewObject(x, y float64, defaultAction Action) Object {
	o := &object{actions: make(map[string]Action), defaultAction: defaultAction}
	o.SetPosition(x, y)

	return o
}

func (obj *object) Start(actionName string) error {
	if obj.activeAction != nil {
		return fmt.Errorf("can't start new action until the previous action is in progress")
	}
	a, exists := obj.actions[actionName]
	if !exists {
		return fmt.Errorf("%s doesn't exist", actionName)
	}

	obj.activeAction = a

	return a.SetX(obj.x).SetY(obj.y).Reset()
}

func (obj *object) Stop() error {
	if obj.activeAction == nil {
		return nil
	}
	if err := obj.activeAction.Stop(); err != nil {
		return err
	}
	obj.activeAction = nil

	return nil
}

func (obj *object) SetDefaultAction(defaultAction Action) Object {
	obj.defaultAction = defaultAction

	return obj
}

func (obj *object) AddAction(action Action) Object {
	obj.actions[action.GetName()] = action

	return obj
}

func (obj *object) Show() Object {
	a := obj.activeAction
	if a == nil {
		a = obj.defaultAction
	}

	a.SetX(obj.x).SetY(obj.y).SetDirection(obj.direction)
	f, x, y := a.Run()
	obj.SetPosition(x, y)
	f.Show(obj.Cx(), obj.Cy())

	return obj
}

func (obj *object) SetPosition(x, y float64) Object {
	obj.x = x
	obj.y = y

	obj.cx = int(math.Round(obj.x))
	obj.cy = int(math.Round(obj.y))

	return obj
}

func (obj *object) X() float64 {
	return obj.x
}

func (obj *object) Y() float64 {
	return obj.y
}

func (obj *object) Cx() int {
	return obj.cx
}

func (obj *object) Cy() int {
	return obj.cy
}

func (obj *object) SetDirection(direction int) Object {
	if obj.direction != direction && obj.activeAction != nil {
		return obj
	}

	obj.direction = direction

	return obj
}

func (obj *object) SetAllowedDirections(directions []int) Object {
	obj.allowedDirections = directions

	return obj
}
