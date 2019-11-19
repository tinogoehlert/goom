package drivers

// Window a generic window
type Window interface {
	Size() (int, int)
	Close()
	Input() InputDriver
	FrameBufferSize() (int, int)
	Run(func(), func(), func(float64))
}
