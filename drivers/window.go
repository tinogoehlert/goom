package drivers

// Window a generic window
type Window interface {
	Size() (int, int)
	PollEvents()
	Close()
	FrameBufferSize() (int, int)
	ShouldClose() bool
}
