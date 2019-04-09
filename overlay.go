/*
Package ofs provides a primitive overlay FileSystem compatible with go/http.
It has some write support and transparent zip file access (read-only).

For example, suppose we have an assets directory:

	assets/
		shaders/
			basic.glsl
		sprites/
			gear.png

The application will be shipped as an executable along with the assets packaged
in a single zip file assets.zip. We also want transparent modding support, so
we use an overlay filesystem that will look for files in the mods/assets
directory then fallback to the assets.zip archive:

	var ovl ofs.Overlay
	err := ovl.Add(false, "assets.zip", "mods")
	shader, err := ovl.Open("assets/shaders/basic.glsl")

The file "assets/shaders/basic.glsl" will be looked up in "mods/assets/shaders/basic.glsl"
then "assets/shaders/basic.glsl" within the assets.zip file.

One could also add a local cache directory on top of the overlay for all write
operations:

	// fallback to some temp dir if any of these fail
	cache, err := os.UserCacheDir()
	cache = filepath.Join(cache, "myApp")
	err = os.MkDir(cache)
	err = ovl.Add(true, cache)

Note that there is no support to remove files. However, the Overlay FileSystem
does not cache any information, client code can therefore use regular os calls
to remove files without interference with the overlay.
*/
package ofs

import (
	"archive/zip"
	"os"
	"path/filepath"
	"syscall"

	"github.com/pkg/errors"
)

// Overlay is a primitive overlay FileSystem. Add directories or zip files to the overlay with
// the Add method.
//
type Overlay struct {
	fs []FileSystem
	dm map[string]int
}

// Add adds the named directories to the overlay. The last directory added takes
// precedence.
//
// A dir can also be a zip archive. However, this will be a read-only FileSystem
// and any File returned by Open will be read-only and not seek-able.
//
// While the FileSystem.Open method takes '/'-separated paths, dir string values
// are filenames on the native file system, so it is separated by
// filepath.Separator, which isn't necessarily '/'.
//
// In order to allow client code to safely call os.Chdir without interference,
// Overlay only keeps track of absolute paths. When a directory is added to the
// overlay, non-absolute paths are resolved relative to the path of executable
// first, then in the current directory.
//
// Add will silently ignore non-existing directories if mustExist is false, and
// Open and Create will never look for files in these.
//
func (o *Overlay) Add(mustExist bool, dirs ...string) error {
	if o.dm == nil {
		o.dm = make(map[string]int)
	}
	exeDir := filepath.Dir(os.Args[0])
	exeDir, err := filepath.Abs(exeDir)
	if err != nil {
		exeDir = ""
	}
	for _, dir := range dirs {
		dir = filepath.Clean(dir)
		abs, err := filepath.Abs(dir)
		// sanitize and resolve to an actual full path
		if err != nil {
			return errors.Wrapf(err, "failed to get absolute path for %q", dir)
		}
		// for relative paths, try exec directory first
		if !filepath.IsAbs(dir) && exeDir != "" {
			err = o.Add(true, filepath.Join(exeDir, dir))
			if err == nil {
				continue
			}
		}

		// skip if this resolves to an existing dir
		if _, ok := o.dm[abs]; ok {
			continue
		}
		// make sure this is an actual directory
		fi, err := os.Stat(abs)
		if err != nil {
			if !mustExist && os.IsNotExist(err) {
				continue
			}
			return errors.Wrapf(err, "failed to stat %q", dir)
		}
		var fs FileSystem
		if !fi.IsDir() {
			// attempt to open as a zip archive
			zr, err := zip.OpenReader(abs)
			if err != nil {
				return syscall.ENOTDIR
			}
			fs = &zipArchive{zr, nil}
		} else {
			fs = Dir(abs)
		}
		o.dm[abs] = len(o.fs)
		o.fs = append(o.fs, fs)
	}
	return nil
}

// Open implements FileSystem.Open.
//
func (o *Overlay) Open(name string) (File, error) {
	i := len(o.fs)
	if i == 0 {
		return nil, &os.PathError{Op: "open", Path: name, Err: errors.Errorf("ofs: no filesystems configured")}
	}
	for i = i - 1; i >= 0; i-- {
		f, err := o.fs[i].Open(name)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		return f, err
	}
	return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
}

// Create implements FileSystem.Create.
//
func (o *Overlay) Create(name string) (File, error) {
	if i := len(o.fs) - 1; i >= 0 {
		return o.fs[i].Create(name)
	}
	return nil, &os.PathError{Op: "create", Path: name, Err: errors.Errorf("ofs: no filesystems configured")}
}
