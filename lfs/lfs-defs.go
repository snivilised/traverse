package lfs

import (
	"io/fs"
	"os"
)

// 📦 pkg: lfs - contains local file system abstractions for navigation.
// Since there are no standard write-able file system interfaces,
// we need to define proprietary ones here in this package.
// This is a low level package that should not use anything else in
// traverse.

type (
	// FileSystems contains the logical file systems required
	// for navigation.
	FileSystems struct {
		// T is the file system that contains just the functionality required
		// for traversal. It can also represent other file systems including afero.
		T TraverseFS
	}

	// ExistsInFS contains methods that check the existence of file system items.
	ExistsInFS interface {
		// FileExists does file exist at the path specified
		FileExists(name string) bool

		// DirectoryExists does directory exist at the path specified
		DirectoryExists(name string) bool
	}

	// ReadFileFS file system non streaming reader
	ReadFileFS interface {
		fs.FS
		// Read reads file at path, from file system specified
		ReadFile(name string) ([]byte, error)
	}

	// ReaderFS
	ReaderFS interface {
		fs.StatFS
		fs.ReadDirFS
		ExistsInFS
		ReadFileFS
	}

	// MakeDirFS is a file system with a MkDirAll method.
	MakeDirFS interface {
		ExistsInFS
		MakeDir(name string, perm os.FileMode) error
		MakeDirAll(name string, perm os.FileMode) error
	}

	// CopyFS
	CopyFS interface {
		Copy(from, to string) error
		// CopyFS copies the file system fsys into the directory dir,
		// creating dir if necessary.
		CopyFS(dir string, fsys fs.FS) error
	}

	// MoveFS
	MoveFS interface {
		Move(from, to string) error
	}

	// RemoveFS
	RemoveFS interface {
		Remove(name string) error
		RemoveAll(path string) error
	}

	// RenameFS
	RenameFS interface {
		Rename(from, to string) error
	}

	// WriteFileFS file system non streaming writer
	WriteFileFS interface {
		// Create creates or truncates the named file.
		Create(name string) (*os.File, error)
		// Write writes file at path, to file system specified
		WriteFile(name string, data []byte, perm os.FileMode) error
	}

	// WriterFS
	WriterFS interface {
		CopyFS
		MoveFS
		ExistsInFS
		RemoveFS
		RenameFS
		WriteFileFS
	}

	// TraverseFS non streaming file system with reader and some
	// writer capabilities
	TraverseFS interface {
		MakeDirFS
		ReaderFS
		WriteFileFS
	}

	// UniversalFS the file system that can do it all
	UniversalFS interface {
		CopyFS
		MoveFS
		RemoveFS
		RenameFS
		TraverseFS
	}
)
