package fs

import (
	"io/fs"
)

type FS interface {
	ReadFile(name string) ([]byte, error)
	Create(path string) (fs.File, error)
}
