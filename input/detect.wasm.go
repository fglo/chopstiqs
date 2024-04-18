//go:build wasm
// +build wasm

package input

import (
	"runtime"
	"strings"
	"syscall/js"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	lastPressedKeys [ebiten.KeyMax + 1]bool
	keyState        map[ebiten.Key]int = make(map[ebiten.Key]int)
)

func init() {
	js.Global().Get("document").Call("addEventListener", "keyup", js.FuncOf(func(this js.Value, args []js.Value) any {
		e := args[0]
		key := jsKeyToID(e.Get("code"))
		keyState[key]--

		if key == ebiten.KeyMetaLeft || key == ebiten.KeyMetaRight {
			for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
				keyState[k] = 0
			}
		}

		switch {
		case key == ebiten.KeyMetaLeft && keyState[ebiten.KeyMetaRight] == 0 || key == ebiten.KeyMetaRight && keyState[ebiten.KeyMetaLeft] == 0:
			keyState[ebiten.KeyMeta] = 0
		case key == ebiten.KeyAltLeft && keyState[ebiten.KeyAltRight] == 0 || key == ebiten.KeyAltRight && keyState[ebiten.KeyAltLeft] == 0:
			keyState[ebiten.KeyAlt] = 0
		case key == ebiten.KeyControlLeft && keyState[ebiten.KeyControlRight] == 0 || key == ebiten.KeyControlRight && keyState[ebiten.KeyControlLeft] == 0:
			keyState[ebiten.KeyControl] = 0
		case key == ebiten.KeyShiftLeft && keyState[ebiten.KeyShiftRight] == 0 || key == ebiten.KeyShiftRight && keyState[ebiten.KeyShiftLeft] == 0:
			keyState[ebiten.KeyShift] = 0
		}

		return nil
	}))

	js.Global().Get("document").Call("addEventListener", "keydown", js.FuncOf(func(this js.Value, args []js.Value) any {
		e := args[0]
		key := jsKeyToID(e.Get("code"))
		meta := e.Get("metaKey").Bool()

		if !meta && keyState[ebiten.KeyMeta] > 0 {
			for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
				keyState[k] = 0
			}
		}

		keyState[key]++

		switch {
		case key == ebiten.KeyMetaLeft || key == ebiten.KeyMetaRight:
			keyState[ebiten.KeyMeta]++
		case key == ebiten.KeyAltLeft || key == ebiten.KeyAltRight:
			keyState[ebiten.KeyAlt]++
		case key == ebiten.KeyControlLeft || key == ebiten.KeyControlRight:
			keyState[ebiten.KeyControl]++
		case key == ebiten.KeyShiftLeft || key == ebiten.KeyShiftRight:
			keyState[ebiten.KeyShift]++
		}

		return nil
	}))

	for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
		keyState[k] = 0
	}
}

func DetectSystem() {
	switch runtime.GOOS {
	case "windows":
		OS = Windows
	case "linux":
		OS = Linux
	case "darwin":
		OS = MacOS
	case "js":
		platform := strings.ToLower(js.Global().Get("navigator").Get("platform").String())
		switch {
		case strings.Contains(platform, "mac"):
			OS = MacOS
		case strings.Contains(platform, "win"):
			OS = Windows
		case strings.Contains(platform, "linux"):
			OS = Linux
		}
	}
}

func DetectPressedKeys() {
	for k, p := range KeyPressed {
		lastPressedKeys[k] = p
	}

	if OSMacOS() {
		for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
			KeyPressed[k] = keyState[k] > 0
			if KeyPressed[k] {
				AnyKeyPressed = true
			}

			KeyJustPressed[k] = inpututil.IsKeyJustPressed(k)
			if KeyJustPressed[k] {
				AnyJustKeyPressed = true
			}
		}

	} else {
		for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
			KeyPressed[k] = ebiten.IsKeyPressed(k)
			if KeyPressed[k] {
				AnyKeyPressed = true
			}

			KeyJustPressed[k] = inpututil.IsKeyJustPressed(k)
			if KeyJustPressed[k] {
				AnyJustKeyPressed = true
			}
		}
	}
}

func jsKeyToID(key js.Value) ebiten.Key {
	// js.Value cannot be used as a map key.
	// As the number of keys is around 100, just a dumb loop should work.
	for uiKey, jsKey := range uiKeyToJSKey {
		if jsKey.Equal(key) {
			return uiKey
		}
	}
	return -1
}

var uiKeyToJSKey = map[ebiten.Key]js.Value{
	ebiten.KeyA:              js.ValueOf("KeyA"),
	ebiten.KeyAltLeft:        js.ValueOf("AltLeft"),
	ebiten.KeyAltRight:       js.ValueOf("AltRight"),
	ebiten.KeyArrowDown:      js.ValueOf("ArrowDown"),
	ebiten.KeyArrowLeft:      js.ValueOf("ArrowLeft"),
	ebiten.KeyArrowRight:     js.ValueOf("ArrowRight"),
	ebiten.KeyArrowUp:        js.ValueOf("ArrowUp"),
	ebiten.KeyB:              js.ValueOf("KeyB"),
	ebiten.KeyBackquote:      js.ValueOf("Backquote"),
	ebiten.KeyBackslash:      js.ValueOf("Backslash"),
	ebiten.KeyBackspace:      js.ValueOf("Backspace"),
	ebiten.KeyBracketLeft:    js.ValueOf("BracketLeft"),
	ebiten.KeyBracketRight:   js.ValueOf("BracketRight"),
	ebiten.KeyC:              js.ValueOf("KeyC"),
	ebiten.KeyCapsLock:       js.ValueOf("CapsLock"),
	ebiten.KeyComma:          js.ValueOf("Comma"),
	ebiten.KeyContextMenu:    js.ValueOf("ContextMenu"),
	ebiten.KeyControlLeft:    js.ValueOf("ControlLeft"),
	ebiten.KeyControlRight:   js.ValueOf("ControlRight"),
	ebiten.KeyD:              js.ValueOf("KeyD"),
	ebiten.KeyDelete:         js.ValueOf("Delete"),
	ebiten.KeyDigit0:         js.ValueOf("Digit0"),
	ebiten.KeyDigit1:         js.ValueOf("Digit1"),
	ebiten.KeyDigit2:         js.ValueOf("Digit2"),
	ebiten.KeyDigit3:         js.ValueOf("Digit3"),
	ebiten.KeyDigit4:         js.ValueOf("Digit4"),
	ebiten.KeyDigit5:         js.ValueOf("Digit5"),
	ebiten.KeyDigit6:         js.ValueOf("Digit6"),
	ebiten.KeyDigit7:         js.ValueOf("Digit7"),
	ebiten.KeyDigit8:         js.ValueOf("Digit8"),
	ebiten.KeyDigit9:         js.ValueOf("Digit9"),
	ebiten.KeyE:              js.ValueOf("KeyE"),
	ebiten.KeyEnd:            js.ValueOf("End"),
	ebiten.KeyEnter:          js.ValueOf("Enter"),
	ebiten.KeyEqual:          js.ValueOf("Equal"),
	ebiten.KeyEscape:         js.ValueOf("Escape"),
	ebiten.KeyF:              js.ValueOf("KeyF"),
	ebiten.KeyF1:             js.ValueOf("F1"),
	ebiten.KeyF10:            js.ValueOf("F10"),
	ebiten.KeyF11:            js.ValueOf("F11"),
	ebiten.KeyF12:            js.ValueOf("F12"),
	ebiten.KeyF2:             js.ValueOf("F2"),
	ebiten.KeyF3:             js.ValueOf("F3"),
	ebiten.KeyF4:             js.ValueOf("F4"),
	ebiten.KeyF5:             js.ValueOf("F5"),
	ebiten.KeyF6:             js.ValueOf("F6"),
	ebiten.KeyF7:             js.ValueOf("F7"),
	ebiten.KeyF8:             js.ValueOf("F8"),
	ebiten.KeyF9:             js.ValueOf("F9"),
	ebiten.KeyG:              js.ValueOf("KeyG"),
	ebiten.KeyH:              js.ValueOf("KeyH"),
	ebiten.KeyHome:           js.ValueOf("Home"),
	ebiten.KeyI:              js.ValueOf("KeyI"),
	ebiten.KeyInsert:         js.ValueOf("Insert"),
	ebiten.KeyJ:              js.ValueOf("KeyJ"),
	ebiten.KeyK:              js.ValueOf("KeyK"),
	ebiten.KeyL:              js.ValueOf("KeyL"),
	ebiten.KeyM:              js.ValueOf("KeyM"),
	ebiten.KeyMetaLeft:       js.ValueOf("MetaLeft"),
	ebiten.KeyMetaRight:      js.ValueOf("MetaRight"),
	ebiten.KeyMinus:          js.ValueOf("Minus"),
	ebiten.KeyN:              js.ValueOf("KeyN"),
	ebiten.KeyNumLock:        js.ValueOf("NumLock"),
	ebiten.KeyNumpad0:        js.ValueOf("Numpad0"),
	ebiten.KeyNumpad1:        js.ValueOf("Numpad1"),
	ebiten.KeyNumpad2:        js.ValueOf("Numpad2"),
	ebiten.KeyNumpad3:        js.ValueOf("Numpad3"),
	ebiten.KeyNumpad4:        js.ValueOf("Numpad4"),
	ebiten.KeyNumpad5:        js.ValueOf("Numpad5"),
	ebiten.KeyNumpad6:        js.ValueOf("Numpad6"),
	ebiten.KeyNumpad7:        js.ValueOf("Numpad7"),
	ebiten.KeyNumpad8:        js.ValueOf("Numpad8"),
	ebiten.KeyNumpad9:        js.ValueOf("Numpad9"),
	ebiten.KeyNumpadAdd:      js.ValueOf("NumpadAdd"),
	ebiten.KeyNumpadDecimal:  js.ValueOf("NumpadDecimal"),
	ebiten.KeyNumpadDivide:   js.ValueOf("NumpadDivide"),
	ebiten.KeyNumpadEnter:    js.ValueOf("NumpadEnter"),
	ebiten.KeyNumpadEqual:    js.ValueOf("NumpadEqual"),
	ebiten.KeyNumpadMultiply: js.ValueOf("NumpadMultiply"),
	ebiten.KeyNumpadSubtract: js.ValueOf("NumpadSubtract"),
	ebiten.KeyO:              js.ValueOf("KeyO"),
	ebiten.KeyP:              js.ValueOf("KeyP"),
	ebiten.KeyPageDown:       js.ValueOf("PageDown"),
	ebiten.KeyPageUp:         js.ValueOf("PageUp"),
	ebiten.KeyPause:          js.ValueOf("Pause"),
	ebiten.KeyPeriod:         js.ValueOf("Period"),
	ebiten.KeyPrintScreen:    js.ValueOf("PrintScreen"),
	ebiten.KeyQ:              js.ValueOf("KeyQ"),
	ebiten.KeyQuote:          js.ValueOf("Quote"),
	ebiten.KeyR:              js.ValueOf("KeyR"),
	ebiten.KeyS:              js.ValueOf("KeyS"),
	ebiten.KeyScrollLock:     js.ValueOf("ScrollLock"),
	ebiten.KeySemicolon:      js.ValueOf("Semicolon"),
	ebiten.KeyShiftLeft:      js.ValueOf("ShiftLeft"),
	ebiten.KeyShiftRight:     js.ValueOf("ShiftRight"),
	ebiten.KeySlash:          js.ValueOf("Slash"),
	ebiten.KeySpace:          js.ValueOf("Space"),
	ebiten.KeyT:              js.ValueOf("KeyT"),
	ebiten.KeyTab:            js.ValueOf("Tab"),
	ebiten.KeyU:              js.ValueOf("KeyU"),
	ebiten.KeyV:              js.ValueOf("KeyV"),
	ebiten.KeyW:              js.ValueOf("KeyW"),
	ebiten.KeyX:              js.ValueOf("KeyX"),
	ebiten.KeyY:              js.ValueOf("KeyY"),
	ebiten.KeyZ:              js.ValueOf("KeyZ"),
}
