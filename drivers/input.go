package drivers

type KeyState int

const (
	KeyReleased KeyState = 0
	KeyPressed  KeyState = 1
	KeyRepeated KeyState = 2
)

type Key struct {
	Keycode Keycode
	State   KeyState
}

type InputDriver interface {
	KeyStates() chan Key
	IsPressed(keycode Keycode) bool
}
