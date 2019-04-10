// Copyright 2019 Denis Bernard <db047h@gmail.com>
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ofs

import (
	"io"
	"os"
)

// A File is returned by a FileSystem's Open or Create method.
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
// Directory names are always rooted at the toplevel of the filesystem. i.e. Open("foo")
// and Open("/foo").
//
type FileSystem interface {
	Open(name string) (File, error)
	Create(name string) (File, error)
}
