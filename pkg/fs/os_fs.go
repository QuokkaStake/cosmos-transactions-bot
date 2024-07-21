package fs

import "os"

type OsFS struct {
}

func (fs *OsFS) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (fs *OsFS) Create(path string) (File, error) {
	return os.Create(path)
}
