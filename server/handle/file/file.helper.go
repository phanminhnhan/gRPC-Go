package file

import (
	"bytes"
)
 
type File struct {
	Name   string
	buffer *bytes.Buffer
	fType  string
}
 
func NewFile(name, fType string) *File {
	return &File{
		Name:   name,
		buffer: &bytes.Buffer{},
		fType: fType,
	}
}
 
func (f *File) Write(chunk []byte) error {
	_, err := f.buffer.Write(chunk)
 
	return err
}