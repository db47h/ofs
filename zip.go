package ofs

import (
	"archive/zip"
	"io"
	"os"
	"path"

	"github.com/pkg/errors"
)

type zipArchive struct {
	*zip.ReadCloser
	fm map[string]*zipFileHeader
}

type zipFileHeader struct {
	*zip.File
	files []*zipFileHeader
}

type zipFile struct {
	io.ReadCloser
	h        *zipFileHeader
	dirIndex int
}

func (f *zipFile) Readdir(count int) ([]os.FileInfo, error) {
	i := f.dirIndex
	j := len(f.h.files)
	if i == j {
		return nil, io.EOF
	}
	if count > 0 && count < j-i {
		j = i + count
	}
	if i == j {
		return nil, nil
	}
	fis := make([]os.FileInfo, 0, j-i)
	for ; i < j; i++ {
		fis = append(fis, f.h.files[i].FileInfo())
	}
	f.dirIndex = j
	return fis, nil
}

func (f *zipFile) Seek(offset int64, whence int) (int64, error) {
	return 0, os.ErrInvalid
}

func (f *zipFile) Write(p []byte) (n int, err error) {
	return 0, os.ErrPermission
}

func (f *zipFile) Stat() (os.FileInfo, error) {
	return f.h.FileInfo(), nil
}

func (a *zipArchive) Open(name string) (File, error) {
	if a.fm == nil {
		a.fm = make(map[string]*zipFileHeader, len(a.File))
		for _, f := range a.File {
			n := path.Clean(f.Name)
			a.fm[n] = &zipFileHeader{f, nil}
		}
		// loop again to populate directories
		for n, f := range a.fm {
			d := path.Dir(n)
			if d == "." {
				continue
			}
			p, ok := a.fm[d]
			if !ok {
				return nil, errors.Errorf("parent folder entry %q not found for %q", d, name)
			}
			p.files = append(p.files, f)
		}
	}
	f := a.fm[name]
	if f == nil {
		return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
	}
	r, err := f.Open()
	if err != nil {
		return nil, err
	}
	return &zipFile{r, f, 0}, nil
}

func (a *zipArchive) Create(name string) (File, error) {
	return nil, os.ErrPermission
}
