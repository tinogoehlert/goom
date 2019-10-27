package files

import (
	"encoding/hex"
	"io/ioutil"
	"path"
)

// BinFile stores a loaded binary file.
type BinFile struct {
	Path string
	Data []byte
}

// Load loads the file content into memory.
func (f *BinFile) Load() (err error) {
	f.Data, err = ioutil.ReadFile(f.Path)
	return
}

// Hex returns the file content as hexadecimal bytes.
func (f *BinFile) Hex() string {
	return hex.EncodeToString(f.Data)
}

// Dump dumps the file in a readable line-numbered hex-format.
func (f *BinFile) Dump() string {
	return hex.Dump(f.Data)
}

// FromHex load file content from hexadecimal bytes.
func (f *BinFile) FromHex(hexString string) (err error) {
	f.Data, err = hex.DecodeString(hexString)
	return
}

// NewBinFile returns a binary file container.
func NewBinFile(location ...string) *BinFile {
	return &BinFile{path.Join(location...), nil}
}
