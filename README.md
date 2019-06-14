[godoc]: https://godoc.org/github.com/db47h/ofs
[godoc-img]: https://godoc.org/github.com/db47h/ofs?status.svg
[goreport]: https://goreportcard.com/report/github.com/db47h/ofs
[goreport-img]: https://goreportcard.com/badge/github.com/db47h/ofs
[license]: https://img.shields.io/github/license/db47h/ofs.svg

# ofs

[![GoDoc][godoc-img]][godoc]
[![GoReportCard][goreport-img]][goreport]
![MIT License][license]

`import "github.com/db47h/ofs"`
* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
Package ofs provides a primitive overlay FileSystem compatible with go/http. It
has some write support and transparent zip file access (read-only).

This was designed primarily to handle asset loading where we want transpartent
support for patches and mods. For those looking for something more advanced,
there's <a href="https://github.com/spf13/afero/">https://github.com/spf13/afero/</a>.

For example, suppose we have an assets directory:

	assets/
	    shaders/
	        basic.glsl
	    sprites/
	        gear.png

The application will be shipped as an executable along with the assets packaged
in a single zip file assets.zip. We also want transparent modding support, so we
use an overlay filesystem that will look for files in the mods/assets directory
then fallback to the assets.zip archive:

	var ovl ofs.Overlay
	err := ovl.Add(false, "assets.zip", "mods")
	shader, err := ovl.Open("assets/shaders/basic.glsl")

The file "assets/shaders/basic.glsl" will be looked up in
"mods/assets/shaders/basic.glsl" then "assets/shaders/basic.glsl" within the
assets.zip file.

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

## <a name="pkg-index">Index</a>
* [type Dir](#Dir)
  * [func (d Dir) Create(name string) (File, error)](#Dir.Create)
  * [func (d Dir) Open(name string) (File, error)](#Dir.Open)
* [type File](#File)
* [type FileSystem](#FileSystem)
* [type Overlay](#Overlay)
  * [func (o *Overlay) Add(mustExist bool, dirs ...string) error](#Overlay.Add)
  * [func (o *Overlay) Create(name string) (File, error)](#Overlay.Create)
  * [func (o *Overlay) Open(name string) (File, error)](#Overlay.Open)

#### <a name="pkg-files">Package files</a>
[dir.go](/src/target/dir.go) [fs.go](/src/target/fs.go) [overlay.go](/src/target/overlay.go) [zip.go](/src/target/zip.go) 


## <a name="Dir">type</a> [Dir](/src/target/dir.go?s=1441:1456#L37)
``` go
type Dir string
```
A Dir implements FileSystem using the native file system restricted to a
specific directory tree.

While the FileSystem.Open method takes '/'-separated paths, a Dir's string
value is an absolute path to a directory on the native file system, not a
URL, so it is separated by filepath.Separator, which isn't necessarily '/'.

It behaves similarly to http.Dir, with the difference that the string value
must be an absolute path and that ErrNotExist errors are not resolved to the
first non-existing directory when Open or Create fails.










### <a name="Dir.Create">func</a> (Dir) [Create](/src/target/dir.go?s=2147:2193#L62)
``` go
func (d Dir) Create(name string) (File, error)
```
Create implements FileSystem.Create.




### <a name="Dir.Open">func</a> (Dir) [Open](/src/target/dir.go?s=1910:1954#L53)
``` go
func (d Dir) Open(name string) (File, error)
```
Open implements FileSystem.Open.




## <a name="File">type</a> [File](/src/target/fs.go?s=953:1092#L26)
``` go
type File interface {
    io.Closer
    io.Reader
    io.Writer
    io.Seeker
    Readdir(count int) ([]os.FileInfo, error)
    Stat() (os.FileInfo, error)
}
```
A File is returned by a FileSystem's Open or Create method.

The methods should behave the same as those on an *os.File.










## <a name="FileSystem">type</a> [FileSystem](/src/target/fs.go?s=1407:1504#L42)
``` go
type FileSystem interface {
    Open(name string) (File, error)
    Create(name string) (File, error)
}
```
A FileSystem implements access to a collection of named files. The elements
in a file path are separated by slash ('/', U+002F) characters, regardless of
host operating system convention.

Directory names are always rooted at the toplevel of the filesystem. i.e. Open("foo")
and Open("/foo").










## <a name="Overlay">type</a> [Overlay](/src/target/overlay.go?s=2841:3099#L89)
``` go
type Overlay struct {
    // If ResolveExecDir is true, Add will try to resolve non-absolute paths
    // relative to the path of the executable before trying the current directory.
    ResolveExecDir bool
    // contains filtered or unexported fields
}
```
Overlay is a primitive overlay FileSystem. Add directories or zip files to the overlay with
the Add method.










### <a name="Overlay.Add">func</a> (\*Overlay) [Add](/src/target/overlay.go?s=3995:4054#L115)
``` go
func (o *Overlay) Add(mustExist bool, dirs ...string) error
```
Add adds the named directories to the overlay. The last directory added takes
precedence.

A dir can also be a zip archive. However, this will be a read-only FileSystem
and any File returned by Open will be read-only and not seek-able.

While the FileSystem.Open method takes '/'-separated paths, dir string values
are filenames on the native file system, so it is separated by
filepath.Separator, which isn't necessarily '/'.

In order to allow client code to safely call os.Chdir without interference,
Overlay only keeps track of absolute paths. When a directory is added to the
overlay, non-absolute paths are resolved relative to the path of executable
first if ResolveRelative, then in the current directory.

Add will silently ignore non-existing directories if mustExist is false, and
Open and Create will never look for files in these.




### <a name="Overlay.Create">func</a> (\*Overlay) [Create](/src/target/overlay.go?s=5654:5705#L185)
``` go
func (o *Overlay) Create(name string) (File, error)
```
Create implements FileSystem.Create.




### <a name="Overlay.Open">func</a> (\*Overlay) [Open](/src/target/overlay.go?s=5182:5231#L165)
``` go
func (o *Overlay) Open(name string) (File, error)
```
Open implements FileSystem.Open.








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
