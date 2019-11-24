package drivers

import (
	sdl_native "github.com/tinogoehlert/go-sdl2/sdl"

	"github.com/tinogoehlert/goom/audio/sfx"
	"github.com/tinogoehlert/goom/drivers/sdl"
)

func init() {
	WindowMakers[SdlWindow] = newSdlWindow
	AudioDrivers[SdlAudio] = newSdlAudio
	Timers[SdlTimer] = sdl.GetTime
	InputProviders[SdlInput] = newSdlInput
}

func newSdlWindow(title string, width, height int) (Window, error) {
	win, err := sdl.NewWindow(title, width, height)

	return Window(win), err
}

func newSdlAudio(sounds *sfx.Sounds, tempFolder string) (Audio, error) {
	audio, err := sdl.NewAudio(sounds, tempFolder)
	return Audio(audio), err
}

func newSdlInput(win Window) Input {
	return Input(sdl.NewInputDriver())
}

var sdlKeyMap = map[sdl_native.Keycode]Keycode{
	sdl_native.K_UP:       KeyUp,
	sdl_native.K_LEFT:     KeyLeft,
	sdl_native.K_RIGHT:    KeyRight,
	sdl_native.K_DOWN:     KeyDown,
	sdl_native.K_SPACE:    KeySpace,
	sdl_native.K_KP_ENTER: KeyEnter,
	sdl_native.K_0:        Key0,
	sdl_native.K_1:        Key1,
	sdl_native.K_2:        Key2,
	sdl_native.K_3:        Key3,
	sdl_native.K_4:        Key4,
	sdl_native.K_5:        Key5,
	sdl_native.K_6:        Key6,
	sdl_native.K_7:        Key7,
	sdl_native.K_8:        Key8,
	sdl_native.K_9:        Key9,
	sdl_native.K_a:        KeyA,
	sdl_native.K_b:        KeyB,
	sdl_native.K_c:        KeyC,
	sdl_native.K_d:        KeyD,
	sdl_native.K_e:        KeyE,
	sdl_native.K_f:        KeyF,
	sdl_native.K_g:        KeyG,
	sdl_native.K_h:        KeyH,
	sdl_native.K_i:        KeyI,
	sdl_native.K_j:        KeyJ,
	sdl_native.K_k:        KeyK,
	sdl_native.K_l:        KeyL,
	sdl_native.K_m:        KeyM,
	sdl_native.K_n:        KeyN,
	sdl_native.K_o:        KeyO,
	sdl_native.K_p:        KeyP,
	sdl_native.K_q:        KeyQ,
	sdl_native.K_r:        KeyR,
	sdl_native.K_s:        KeyS,
	sdl_native.K_t:        KeyT,
	sdl_native.K_u:        KeyU,
	sdl_native.K_v:        KeyV,
	sdl_native.K_w:        KeyW,
	sdl_native.K_x:        KeyX,
	sdl_native.K_y:        KeyY,
	sdl_native.K_z:        KeyZ,
}

var sdlDriversKeyMap = map[Keycode]sdl_native.Keycode{
	KeyUp:     sdl_native.K_UP,
	KeyLeft:   sdl_native.K_LEFT,
	KeyRight:  sdl_native.K_RIGHT,
	KeyDown:   sdl_native.K_DOWN,
	KeySpace:  sdl_native.K_SPACE,
	KeyEnter:  sdl_native.K_KP_ENTER,
	KeyRShift: sdl_native.K_RSHIFT,
	KeyLShift: sdl_native.K_LSHIFT,
	Key0:      sdl_native.K_0,
	Key1:      sdl_native.K_1,
	Key2:      sdl_native.K_2,
	Key3:      sdl_native.K_3,
	Key4:      sdl_native.K_4,
	Key5:      sdl_native.K_5,
	Key6:      sdl_native.K_6,
	Key7:      sdl_native.K_7,
	Key8:      sdl_native.K_8,
	Key9:      sdl_native.K_9,
	KeyA:      sdl_native.K_a,
	KeyB:      sdl_native.K_b,
}
