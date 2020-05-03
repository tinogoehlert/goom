package glvis

// #include "glvis.h"
import "C"
import (
	"unsafe"
)

// GLVis object is used to initialize glvis.
type GLVis struct {
	visPtr unsafe.Pointer
}

// NewGLVis builds a new GLVis object.
func NewGLVis() GLVis {
	return GLVis{
		visPtr: C.GLVisInit(),
	}
}

// Free deallocates all memory for glvis.
func (f GLVis) Free() {
	C.GLVisFree(unsafe.Pointer(f.visPtr))
}

// BuildVis builds the GL_PVS nodes from WAD and glBSP gwa buffers.
func (f GLVis) BuildVis(wadBuff, gwaBuff []byte) []byte {
	C.BuildVis(
		unsafe.Pointer(f.visPtr),
		(*C.uchar)(unsafe.Pointer(&wadBuff[0])),
		(*C.uchar)(unsafe.Pointer(&gwaBuff[0])),
	)

	size := int((C.uint32_t)(C.GetVisSize(unsafe.Pointer(f.visPtr))))

	buffPtr := unsafe.Pointer(C.GetVisData(unsafe.Pointer(f.visPtr)))
	return C.GoBytes(buffPtr, (C.int)(size))
}
