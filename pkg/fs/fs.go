package fs

import (
	"io"
)

type File interface {
	io.WriteCloser
}

type FS interface {
	ReadFile(name string) ([]byte, error)
	Create(path string) (File, error)
}
