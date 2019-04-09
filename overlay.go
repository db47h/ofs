package ofs

import (
	"archive/zip"
	"log"
	"os"
	"path/filepath"
	"syscall"

	"github.com/pkg/errors"
)

// Overlay is a primitive overlay FileSystem. Add directories or zip files to the overlay with
// the Add method.
//
// There is no support for renaming or deleting files. The main reason to support file Creation
// is to implement file caches on top of the overlay.
//
type Overlay struct {
	fs []FileSystem
	dm map[string]int
}

// Add adds the named directories to the overlay. The last directory added takes precedence.
// Paths are expected to use the host os' separators (i.e. os.PathSeparator).
//
// A dir can also be a zip archive. However, this will be a read-only FileSystem and any File
// returned by Open will be read-only and not seek-able.
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
			log.Printf("OverlayFS: %q: trying %q", dir, filepath.Join(exeDir, dir))
			err = o.Add(true, filepath.Join(exeDir, dir))
			if err == nil {
				continue
			}
			log.Printf("OverlayFS: %q: trying %q", dir, abs)
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
		log.Printf("Overlay FS: added %q", abs)
	}
	return nil
}

// Open implements FileSystem.Open.
//
func (o *Overlay) Open(name string) (File, error) {
	for i := len(o.fs) - 1; i >= 0; i-- {
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
	return nil, os.ErrPermission
}
