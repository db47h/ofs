package ofs

import (
	"io"
	"os"
)

// A File is returned by a FileSystem's Open method.
//
// The methods should behave the same as those on an *os.File.
//
type File interface {
	io.Closer
	io.Reader
	io.Writer
	io.Seeker
	Readdir(count int) ([]os.FileInfo, error)
	Stat() (os.FileInfo, error)
}

// A FileSystem implements access to a collection of named files. The elements
// in a file path are separated by slash ('/', U+002F) characters, regardless of
// host operating system convention.
//
type FileSystem interface {
	Open(name string) (File, error)
	Create(name string) (File, error)
}
