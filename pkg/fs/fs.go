package fs

import (
	"io/fs"
	"os"
)

type FS interface {
	ReadFile(name string) ([]byte, error)
	Create(path string) (fs.File, error)
}

type OsFS struct {
}

func (fs *OsFS) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (fs *OsFS) Create(path string) (fs.File, error) {
	return os.Create(path)
}
