package lfs_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
)

func TestLfs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lfs Suite")
}

type (
	ensureTE struct {
		given     string
		should    string
		relative  string
		expected  string
		directory bool
	}

	RPEntry struct {
		given  string
		should string
		path   string
		expect string
	}
)

var (
	fakeHome      = filepath.Join(string(filepath.Separator), "home", "rabbitweed")
	fakeAbsCwd    = filepath.Join(string(filepath.Separator), "home", "rabbitweed", "music", "xpander")
	fakeAbsParent = filepath.Join(string(filepath.Separator), "home", "rabbitweed", "music")
)

func fakeHomeResolver() (string, error) {
	return fakeHome, nil
}

func fakeAbsResolver(path string) (string, error) {
	if strings.HasPrefix(path, "..") {
		return filepath.Join(fakeAbsParent, path[2:]), nil
	}

	if strings.HasPrefix(path, ".") {
		return filepath.Join(fakeAbsCwd, path[1:]), nil
	}

	return path, nil
}

type (
	makeDirMapFS struct {
		mapFS fstest.MapFS
	}
)

func (f *makeDirMapFS) FileExists(path string) bool {
	fi, err := f.mapFS.Stat(path)
	if err != nil {
		return false
	}

	if fi.IsDir() {
		return false
	}

	return true
}

func (f *makeDirMapFS) DirectoryExists(path string) bool {
	if strings.HasPrefix(path, string(filepath.Separator)) {
		path = path[1:]
	}

	fileInfo, err := f.mapFS.Stat(path)
	if err != nil {
		return false
	}

	if !fileInfo.IsDir() {
		return false
	}

	return true
}

func (f *makeDirMapFS) MakeDir(path string, perm os.FileMode) error {
	if exists := f.DirectoryExists(path); !exists {
		f.mapFS[path] = &fstest.MapFile{
			Mode: fs.ModeDir | perm,
		}
	}

	return nil
}

func (f *makeDirMapFS) MakeDirAll(path string, perm os.FileMode) error {
	var current string
	segments := filepath.SplitList(path)

	for _, part := range segments {
		if current == "" {
			current = part
		} else {
			current += string(filepath.Separator) + part
		}

		if exists := f.DirectoryExists(current); !exists {
			f.mapFS[current] = &fstest.MapFile{
				Mode: fs.ModeDir | perm,
			}
		}
	}

	return nil
}
