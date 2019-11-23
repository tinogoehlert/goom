package drivers

// Window a generic window
type Window interface {
	Close()
	GetInput() InputDriver
	GetSize() (int, int)
	RunGame(func(), func(), func(float64))
}
