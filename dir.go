package ofs

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// A Dir implements FileSystem using the native file system restricted to a
// specific directory tree.
//
// While the FileSystem.Open method takes '/'-separated paths, a Dir's string
// value is an absolute path to a directory on the native file system, not a
// URL, so it is separated by filepath.Separator, which isn't necessarily '/'.
//
// It behaves similarly to http.Dir, with the difference that the string value
// must be an absolute path and that ErrNotExist errors are not resolved to the
// first non-existing directory when Open or Create fails.
//
type Dir string

func (d Dir) check(name string) error {
	// error if dir is not absolute
	if dir := string(d); !filepath.IsAbs(dir) {
		return errors.Errorf("ofs: dir is not absolute %q", dir)
	}
	// error if using local os' path separators instead of '/'
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) {
		return errors.Errorf("ofs: invalid character in file path %q", name)
	}
	return nil
}

// Open implements FileSystem.Open.
//
func (d Dir) Open(name string) (File, error) {
	if err := d.check(name); err != nil {
		return nil, err
	}
	return os.Open(filepath.Join(string(d), filepath.FromSlash(path.Clean("/"+name))))
}

// Create implements FileSystem.Create.
//
func (d Dir) Create(name string) (File, error) {
	if err := d.check(name); err != nil {
		return nil, err
	}
	return os.Create(filepath.Join(string(d), filepath.FromSlash(path.Clean("/"+name))))
}
