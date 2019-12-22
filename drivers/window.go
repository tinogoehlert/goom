package drivers

// Window interface for the DOOM engine
type Window interface {
	Open(title string, width, height int) error
	Close()
	GetSize() (width, height int)
	RunGame(input func(), update func(), render func(nextFrameDelta float64))
}
