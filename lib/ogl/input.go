package ogl

import (
	"fmt"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	KEY_UNKNOWN       = glfw.KeyUnknown
	kEy_SPACE         = glfw.KeySpace
	KEY_APOSTROPHE    = glfw.KeyApostrophe
	KEY_COMMA         = glfw.KeyComma
	KEY_MINUS         = glfw.KeyMinus
	KEY_PERIOD        = glfw.KeyPeriod
	KEY_SLASH         = glfw.KeySlash
	KEY_0             = glfw.Key0
	KEY_1             = glfw.Key1
	KEY_2             = glfw.Key2
	KEY_3             = glfw.Key3
	KEY_4             = glfw.Key4
	KEY_5             = glfw.Key5
	KEY_6             = glfw.Key6
	KEY_7             = glfw.Key7
	KEY_8             = glfw.Key8
	KEY_9             = glfw.Key9
	KEY_SEMICOLON     = glfw.KeySemicolon
	KEY_EQUAL         = glfw.KeyEqual
	KEY_A             = glfw.KeyA
	KEY_B             = glfw.KeyB
	KEY_C             = glfw.KeyC
	KEY_D             = glfw.KeyD
	KEY_E             = glfw.KeyE
	KEY_F             = glfw.KeyF
	KEY_G             = glfw.KeyG
	KEY_H             = glfw.KeyH
	KEY_I             = glfw.KeyI
	KEY_J             = glfw.KeyJ
	KEY_K             = glfw.KeyK
	KEY_L             = glfw.KeyL
	KEY_M             = glfw.KeyM
	KEY_N             = glfw.KeyN
	KEY_O             = glfw.KeyO
	KEY_P             = glfw.KeyP
	KEY_Q             = glfw.KeyQ
	KEY_R             = glfw.KeyR
	KEY_S             = glfw.KeyS
	KEY_T             = glfw.KeyT
	KEY_U             = glfw.KeyU
	KEY_V             = glfw.KeyV
	KEY_W             = glfw.KeyW
	KEY_X             = glfw.KeyX
	KEY_Y             = glfw.KeyY
	KEY_Z             = glfw.KeyZ
	KEY_LEFT_BRACKET  = glfw.KeyLeftBracket
	KEY_BACKSLASH     = glfw.KeyBackslash
	KEY_RIGHT_BRACKET = glfw.KeyRightBracket
	KEY_GRAVE_ACCENT  = glfw.KeyGraveAccent
	KEY_WORLD1        = glfw.KeyWorld1
	KEY_WORLD2        = glfw.KeyWorld2
	KEY_ESC           = glfw.KeyEscape
	KEY_ENTER         = glfw.KeyEnter
	KEY_TAB           = glfw.KeyTab
	KEY_BACKSPACE     = glfw.KeyBackspace
	KEY_INSERT        = glfw.KeyInsert
	KEY_DELETE        = glfw.KeyDelete
	KEY_RIGHT         = glfw.KeyRight
	KEY_LEFT          = glfw.KeyLeft
	KEY_DOWN          = glfw.KeyDown
	KEY_UP            = glfw.KeyUp
	KEY_PAGEuP        = glfw.KeyPageUp
	KEY_PAGE_DOWN     = glfw.KeyPageDown
	KEY_HOME          = glfw.KeyHome
	KEY_END           = glfw.KeyEnd
	KEY_CAPSlOCK      = glfw.KeyCapsLock
	KEY_SCROLL_LOCK   = glfw.KeyScrollLock
	KEY_NUM_LOCK      = glfw.KeyNumLock
	KEY_PRINT_SCREEN  = glfw.KeyPrintScreen
	KEY_PAUSE         = glfw.KeyPause
	KEY_F1            = glfw.KeyF1
	KEY_F2            = glfw.KeyF2
	KEY_F3            = glfw.KeyF3
	KEY_F4            = glfw.KeyF4
	KEY_F5            = glfw.KeyF5
	KEY_F6            = glfw.KeyF6
	KEY_F7            = glfw.KeyF7
	KEY_F8            = glfw.KeyF8
	KEY_F9            = glfw.KeyF9
	KEY_F10           = glfw.KeyF10
	KEY_F11           = glfw.KeyF11
	KEY_F12           = glfw.KeyF12
	KEY_F13           = glfw.KeyF13
	KEY_F14           = glfw.KeyF14
	KEY_F15           = glfw.KeyF15
	KEY_F16           = glfw.KeyF16
	KEY_F17           = glfw.KeyF17
	KEY_F18           = glfw.KeyF18
	KEY_F19           = glfw.KeyF19
	KEY_F20           = glfw.KeyF20
	KEY_F21           = glfw.KeyF21
	KEY_F22           = glfw.KeyF22
	KEY_F23           = glfw.KeyF23
	KEY_F24           = glfw.KeyF24
	KEY_F25           = glfw.KeyF25
	KEY_KP0           = glfw.KeyKP0
	KEY_KP1           = glfw.KeyKP1
	KEY_KP2           = glfw.KeyKP2
	KEY_KP3           = glfw.KeyKP3
	KEY_KP4           = glfw.KeyKP4
	KEY_KP5           = glfw.KeyKP5
	KEY_KP6           = glfw.KeyKP6
	KEY_KP7           = glfw.KeyKP7
	KEY_KP8           = glfw.KeyKP8
	KEY_KP9           = glfw.KeyKP9
	KEY_KP_DECIMAL    = glfw.KeyKPDecimal
	KEY_KP_DIVIDE     = glfw.KeyKPDivide
	KEY_KP_MULTIPLY   = glfw.KeyKPMultiply
	KEY_KP_SUBTRACT   = glfw.KeyKPSubtract
	KEY_KP_ADD        = glfw.KeyKPAdd
	KEY_KP_ENTER      = glfw.KeyKPEnter
	KEY_KP_EQUAL      = glfw.KeyKPEqual
	KEY_LEFT_SHIFT    = glfw.KeyLeftShift
	KEY_LEFT_CONTROL  = glfw.KeyLeftControl
	KEY_LEFT_ALT      = glfw.KeyLeftAlt
	KEY_LEFT_SUPER    = glfw.KeyLeftSuper
	KEY_RIGHT_SHIFT   = glfw.KeyRightShift
	KEY_RIGHT_CONTROL = glfw.KeyRightControl
	KEY_RIGHT_ALT     = glfw.KeyRightAlt
	KEY_RIGHT_SUPER   = glfw.KeyRightSuper
	KEY_MENU          = glfw.KeyMenu
	KEY_LAST          = glfw.KeyLast
)

const (
	KEY_MODIFIER_SHIFT     = glfw.ModShift
	KEY_MODIFIER_CTRL      = glfw.ModControl
	KEY_MODIFIER_ALT       = glfw.ModAlt
	KEY_MODIFIER_SUPER     = glfw.ModSuper
	KEY_MODIFIER_CAPS_LOCK = glfw.ModCapsLock
	KEY_MODIFIER_NUM_LOCK  = glfw.ModNumLock
)

func processInput(window *glfw.Window) {

}

func OnKeypress(key glfw.Key, callback func()) {
	if _, exists := keyCallbacks["press"]; !exists {
		keyCallbacks["press"] = make(map[glfw.Key][]func(), 0)
	}
	keyCallbacks["press"][key] = append(keyCallbacks["press"][key], callback)
}

func OnKeydown(key glfw.Key, callback func()) {
	if _, exists := keyCallbacks["down"]; !exists {
		keyCallbacks["down"] = make(map[glfw.Key][]func(), 0)
	}
	keyCallbacks["down"][key] = append(keyCallbacks["down"][key], callback)
}

func OnKeyup(key glfw.Key, callback func()) {
	if _, exists := keyCallbacks["up"]; !exists {
		keyCallbacks["up"] = make(map[glfw.Key][]func(), 0)
	}
	keyCallbacks["up"][key] = append(keyCallbacks["up"][key], callback)
}

func OnCombination(keys []interface{}, callback func()) {
	typedKeys := make([]glfw.Key, len(keys))
	for i, ik := range keys {
		typedKeys[i] = ik.(glfw.Key)
	}
	keyCombinationCallbacks[fmt.Sprint(keys)] = map[string]interface{}{
		"keys":     typedKeys,
		"callback": callback,
	}
}

func onKeyPress(w *glfw.Window, keyPressed glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Repeat {
		pressedKeys = append(pressedKeys, keyPressed)
	}
	if action == glfw.Release {
		for i, key := range pressedKeys {
			if key == keyPressed {
				pressedKeys = append(pressedKeys[:i], pressedKeys[i+1:]...)
			}
		}
	}
	for state, keys := range keyCallbacks {
		for key, callbacks := range keys {
			if key != keyPressed {
				continue
			}
			for _, callback := range callbacks {
				switch state {
				case "press":
					if action == glfw.Action(glfw.Press) {
						callback()
					}
				case "up":
					if action == glfw.Action(glfw.Release) {
						callback()
					}
				case "down":
					if action == glfw.Action(glfw.Repeat) {
						callback()
					}
				}

			}
		}
	}
	if action != glfw.Action(glfw.Press) && action != glfw.Action(glfw.Repeat) {
		return
	}

	for _, data := range keyCombinationCallbacks {
		keys := data["keys"].([]glfw.Key)
		if len(keys) < 1 {
			continue
		}
		if keyPressed != keys[len(keys)-1] {
			continue
		}
		callback := data["callback"].(func())
		canFire := true
		for i := 0; i < len(keys)-1; i++ {
			for _, k := range pressedKeys {
				if keys[i] != k {
					canFire = false
					break
				}
			}
		}
		if !canFire {
			continue
		}

		callback()
	}
}
