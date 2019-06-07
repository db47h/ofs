package ofs_test

import (
	"log"
	"net/http"
	"github.com/db47h/ofs"
)

// This example demonstrates how to use a FileSystem as an http.FileSystem
// with a simple wrapper.
//
func ExampleOverlay_httpFileSystem() {
	var (
		ovl = new(ofs.Overlay)
		httpFS = &httpFileSystem{ovl}
	)
	// configure overlay
	err := ovl.Add(true, "foo")
	if err != nil {
		log.Fatal(err)
	}
	// serve
	log.Fatal(http.ListenAndServe(":8080", http.FileServer(httpFS)))
}

// httpFileSystem is a simple wrapper around a FileSystem that 
// implements http.FileSystem.
//
type httpFileSystem struct {
	fs ofs.FileSystem
}

func (fs *httpFileSystem) Open(name string) (http.File, error) {
	return fs.fs.Open(name)
}

func (fs *httpFileSystem) Create(name string) (http.File, error) {
	return fs.fs.Create(name)
}
