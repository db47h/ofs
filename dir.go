package ofs

import (
	"os"
	"path"
	"path/filepath"
)

// A Dir implements FileSystem.
//
type Dir string

// Open implements FileSystem.Open.
//
func (d Dir) Open(name string) (File, error) {
	name = path.Clean("/" + name)

	return os.Open(filepath.Join(string(d), filepath.FromSlash(name)))
}

// Create implements FileSystem.Create.
//
func (d Dir) Create(name string) (File, error) {
	name = path.Clean("/" + name)
	return os.Create(filepath.Join(string(d), filepath.FromSlash(name)))
}
