package ofs_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/db47h/ofs"
)

type testData struct {
	path     string
	contents string
}

var td = []testData{
	{"testdata/foodir/foobar", "foobar"},
	{"testdata/foodir/foofoo", "foofoo_fromdir"},
	{"testdata/foo", "foo_fromdir"},
	{"testdata/bar", "bar_cache"},
}

func Test_Overlay(t *testing.T) {
	var ovl ofs.Overlay

	// create cache
	tmp, err := ioutil.TempDir("", "goovl")
	if err != nil {
		panic(err)
	}
	t.Log("created temp dir ", tmp)
	defer os.RemoveAll(tmp)
	err = os.Mkdir(filepath.Join(tmp, "testdata"), 0700)
	if err != nil {
		t.Fatal("failed to create temp dir:", err)
	}

	// TODO: there is a problem with using "." for current dir since
	// it will take the current directory of the executable (for tests, somewhere under /tmp)
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal("failed to get current directory:", err)
	}
	err = ovl.Add(true, "testdata.zip", wd, tmp)
	if err != nil {
		t.Fatal("failed to create overlay:", err)
	}

	// test create
	f, err := ovl.Create("testdata/bar")
	if err != nil {
		t.Fatal("failed to create file in overlay:", err)
	}
	_, err = f.Write([]byte("bar_cache"))
	f.Close()
	if err != nil {
		t.Fatal("failed to write data to cache:", err)
	}

	// check that the file go created in the cache
	cached := filepath.Join(tmp, "testdata/bar")
	fi, err := os.Stat(cached)
	if err != nil {
		t.Fatal("failed to stat cached file", err)
	}
	if fi.Size() != 9 {
		t.Fatalf("bad cached data: expected size %d, got %d", 9, fi.Size())
	}

	// now test reading
	for _, d := range td {
		f, err := ovl.Open(d.path)
		if err != nil {
			t.Errorf("failed to open %s: %v", d.path, err)
		}
		c, err := ioutil.ReadAll(f)
		if err != nil {
			t.Errorf("failed to read from %s: %v", d.path, err)
		}
		if string(c) != d.contents {
			t.Errorf("%s: bad file contents. Expected %s, got %s", d.path, d.contents, c)
		}
	}
}
