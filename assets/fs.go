package assets

import (
	"embed"
	"errors"
	"io/fs"
)

//go:embed *
var EmbedFS embed.FS

type TmpFSInterface struct {
}

func (filesystem *TmpFSInterface) ReadFile(name string) ([]byte, error) {
	return EmbedFS.ReadFile(name)
}

func (filesystem *TmpFSInterface) Create(path string) (fs.File, error) {
	return nil, errors.New("not yet supported")
}

var FS = &TmpFSInterface{}
