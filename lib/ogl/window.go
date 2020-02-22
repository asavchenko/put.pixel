package ogl

import "github.com/go-gl/glfw/v3.3/glfw"

func GetWindowWidth() int {
	return width
}

func GetWindowHeight() int {
	return height
}

func GoFullScreen() {
	if window.GetMonitor() != nil {
		return
	}
	mon := glfw.GetPrimaryMonitor()
	vmode := mon.GetVideoMode()
	window.SetMonitor(mon, 0, 0, width, height, vmode.RefreshRate)
}

func RestoreFromFullScreen() {
	if window.GetMonitor() != nil {
		window.SetMonitor(nil, lastX, lastY, lastWidth, lastHeight, glfw.DontCare)
	}
}

func IsInFullScreen() bool {
	return window.GetMonitor() != nil
}
