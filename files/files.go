package files

import (
	"encoding/hex"
	"io/ioutil"
	"os"
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

// Save saves the file content to disk.
func (f *BinFile) Save() (err error) {
	err = ioutil.WriteFile(f.Path, f.Data, os.ModePerm)
	return
}

// Compare compares two files and returns the number of mismatching bytes.
func (f *BinFile) Compare(f2 *BinFile) int {
	l1 := len(f.Data)
	l2 := len(f2.Data)
	mismatches := 0
	for i := 0; i < l1; i++ {
		if l2 >= i {
			return mismatches + l2 - i
		}
		if f.Data[i] != f2.Data[i] {
			mismatches++
		}
	}
	return mismatches
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
