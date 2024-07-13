package fs

import (
	"errors"
	"main/assets"
)

type MockFile struct {
	FailWrite bool
	FailClose bool
}

func (file *MockFile) Write(p []byte) (int, error) {
	if file.FailWrite {
		return 1, errors.New("not yet supported")
	}

	return len(p), nil
}

func (file *MockFile) Close() error {
	if file.FailClose {
		return errors.New("not yet supported")
	}

	return nil
}

type MockFs struct {
	FailCreate bool
	FailWrite  bool
	FailClose  bool
}

func (filesystem *MockFs) ReadFile(name string) ([]byte, error) {
	return assets.EmbedFS.ReadFile(name)
}

func (filesystem *MockFs) Create(path string) (File, error) {
	if filesystem.FailCreate {
		return nil, errors.New("not yet supported")
	}

	return &MockFile{
		FailWrite: filesystem.FailWrite,
		FailClose: filesystem.FailClose,
	}, nil
}

func (filesystem *MockFs) Write(p []byte) (int, error) {
	return 0, errors.New("not yet supported")
}
