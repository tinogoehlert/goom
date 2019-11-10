package drivers

// Window a generic window
type Window interface {
	Size() (int, int)
	AspectRatio() float32
	PollEvents()
	Close()
	FrameBufferSize() (int, int)
	ShouldClose() bool
}
