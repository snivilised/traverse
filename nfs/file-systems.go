package nfs

import (
	"io/fs"
	"os"
)

type nativeFS struct {
	fsys fs.FS
}

func NewNativeFS(path string) fs.ReadDirFS {
	return &nativeFS{
		fsys: os.DirFS(path),
	}
}

func (n *nativeFS) Open(path string) (fs.File, error) {
	return n.fsys.Open(path)
}

func (n *nativeFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return fs.ReadDir(n.fsys, name)
}

type queryStatusFS struct {
	fsys fs.FS
}

func NewQueryStatusFS(fsys fs.FS) fs.StatFS {
	return &queryStatusFS{
		fsys: fsys,
	}
}

func (q *queryStatusFS) Open(name string) (fs.File, error) {
	return q.fsys.Open(name)
}

func (q *queryStatusFS) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(name)
}
