// The contents of this file is free and unencumbered software released into
// the public domain. Refer to <http://unlicense.org/> for more information.

package pkg

// note: similar keys have different key codes, like Enter and
//       Keypad Enter

// cross-platform key codes, compatible with SDL 2 and USB HID speccy
// https://hg.libsdl.org/SDL/file/default/include/SDL_scancode.h
// http://www.usb.org/developers/hidpage/Hut1_12v2.pdf (page 53)

// key for codes 0 - 3 are not present on keyboards, they are:
// 0 - Reserved (no event)
// 1 - ErrorRollOver
// 2 - POSTFail
// 3 - ErrorUndefined

// Keycode defines a driver-independent keycode.
type Keycode uint16

// Code returns the key code's value.
func (k Keycode) Code() uint16 {
	return uint16(k)
}

// Define character codes
const (
	KeyA Keycode = 4 + iota
	KeyB
	KeyC
	KeyD
	KeyE
	KeyF
	KeyG
	KeyH
	KeyI
	KeyJ
	KeyK
	KeyL
	KeyM
	KeyN
	KeyO
	KeyP
	KeyQ
	KeyR
	KeyS
	KeyT
	KeyU
	KeyV
	KeyW
	KeyX
	KeyY
	KeyZ
)

// Define number codes.
const (
	Key1 Keycode = 30 + iota
	Key2
	Key3
	Key4
	Key5
	Key6
	Key7
	Key8
	Key9
	Key0
)

// Define special key codes 40-50
const (
	// choice is to use Enter name instead of Return
	// and the key code is different from Keypad Enter
	KeyEnter Keycode = 40 + iota
	KeyEscape
	KeyBackspace
	KeyTab
	KeySpace

	// KeyMinus on keypad has different key code
	KeyMinus
	KeyEquals
	KeyLeftBracket
	KeyRightBracket
	KeyBackslash
)

// key code number 50 is skipped, because it is unclear
// where is the key, and what is its name and function

// Define special key codes 51-57
const (
	// different name from SDL2 for brevity
	KeyColon Keycode = 51 + iota
	KeyApostrophe
	// KeyTilde is an alias
	KeyGrave
	KeyCommad
	// KeyDot is an alias, keypad period is a different
	KeyPeriod
	Slash
	CapsLock
)

// define aliases
const (
	KeyTilde = KeyGrave
	KeyDot   = KeyPeriod
)

// Define special key codes 58-69
const (
	KeyF1 Keycode = 58 + iota
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
)

// Define special key codes 70-82
const (
	KeyPrintScreen Keycode = 70 + iota
	KeyScrollLock
	KeyPause
	KeyInsert
	KeyHome
	KeyPageUp
	KeyDelete
	KeyEnd
	KeyPageDown
	KeyRight
	KeyLeft
	KeyDown
	KeyUp
)

// Define special key codes 83-99
const (
	KeyNumLock Keycode = 83 + iota
	KeyKpDivide
	KeyKpMultiply
	KeyKpMinus
	KeyKpPlus
	KeyKpEnter
	KeyKp1
	KeyKp2
	KeyKp3
	KeyKp4
	KeyKp5
	KeyKp6
	KeyKp7
	KeyKp8
	KeyKp9
	KeyKp0
	KeyKpPeriod
)

// KeyKpDot alias for period
const KeyKpDot = KeyKpPeriod

// key code 100 is skipped, because I can not find the key
// key code 101 is not present on Mac
// key codes 102-223 are not present on PC

// define modifiers
const (
	KeyLCtrl Keycode = 224 + iota
	KeyLShift
	KeyLAlt
	// KeyLWin is an alias
	KeyLGUI
	KeyRCtrl
	KeyRShift
	KeyRAlt
	// KeyRWin is an alias
	KeyRGUI
)

// define modifier aliases
const (
	KeyLWin = KeyLGUI
	KeyRWin = KeyRGUI
)
