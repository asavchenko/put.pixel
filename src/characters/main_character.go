package characters

import "assa.com/put.pixel/src/sprite"

func GetSprite(x, y float64) sprite.Object {
	idleFrames := sprite.GetFrames("src/characters/main")
	idleVectors := sprite.GetDefaultVectors(len(idleFrames), 36)
	defaultAction := sprite.GetNewAction("idle", idleFrames, idleVectors, true)
	defaultAction.SetAllowedDirections([]int{sprite.DIRECTION_LEFT, sprite.DIRECTION_RIGHT})
	mCh := sprite.GetNewObject(x, y, defaultAction)
	mCh.SetAllowedDirections([]int{sprite.DIRECTION_LEFT, sprite.DIRECTION_RIGHT})
	mCh.SetDirection(sprite.DIRECTION_RIGHT)

	return mCh
}
