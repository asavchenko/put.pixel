package characters

import "assa.com/put.pixel/src/sprite"

func GetSprite(x, y float64) sprite.Object {
	idleVectors
	defaultAction := sprite.GetNewAction("idle", idleFrames, idleVectors, true)
	pathToDefaultFrameImg := "main/default.png"
	defaultFrame := sprite.GetNewFrame(pathToDefaultFrameImg, 64, 64)
	mCh := sprite.GetNewObject(x, y, defaultFrame)
	mCh.AddAction()
	return mCh
}
