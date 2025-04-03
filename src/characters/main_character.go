package characters

import (
	"assa.com/put.pixel/src/sprite"
	"math"
)

func GetSprite(x, y float64) sprite.Object {
	idleFrames := sprite.GetFrames("src/characters/main", []int{sprite.DIRECTION_LEFT, sprite.DIRECTION_RIGHT}, sprite.DIRECTION_RIGHT)
	idleVectors := sprite.GetDefaultVectors(idleFrames, 18)
	defaultAction := sprite.GetNewAction("idle", idleFrames, idleVectors, true)
	defaultAction.SetAllowedDirections([]int{sprite.DIRECTION_LEFT, sprite.DIRECTION_RIGHT})
	mCh := sprite.GetNewObject(x, y, defaultAction)
	mCh.SetAllowedDirections([]int{sprite.DIRECTION_LEFT, sprite.DIRECTION_RIGHT})

	walkFrames := sprite.GetFrames("src/characters/main/walk", []int{sprite.DIRECTION_LEFT, sprite.DIRECTION_RIGHT}, sprite.DIRECTION_RIGHT)
	walkVectors := sprite.GetDefaultVectors(walkFrames, 18)
	for _, v := range walkVectors {
		v.SetSpeed(v.GetSpeed() + math.Pi/4)
		v.SetDurationInCycles(9)
	}
	walkAction := sprite.GetNewAction("walk", walkFrames, walkVectors, true)
	walkAction.SetAllowedDirections([]int{sprite.DIRECTION_LEFT, sprite.DIRECTION_RIGHT})
	mCh.AddAction(walkAction)

	mCh.SetDirection(sprite.DIRECTION_RIGHT)

	return mCh
}
