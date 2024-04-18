package files

import (
	"book-to-mail-bot/lib/e"
	"book-to-mail-bot/storage"
	"io"
	"os"
	"path/filepath"
)

const (
	defaultPerm = 0774
)

type Storage struct {
	basePath string
}

func New(basePath string) *Storage {
	return &Storage{
		basePath: basePath,
	}
}

func (s *Storage) Save(page *storage.Book, data io.ReadCloser) (path string, err error) {
	defer func() { err = e.WrapIfErr("can't save book", err) }()

	if err := os.MkdirAll(s.basePath, defaultPerm); err != nil {
		return "", err
	}

	fPath := filepath.Join(s.basePath, page.Name)

	file, err := os.Create(fPath)
	if err != nil {
		return "", err
	}

	defer func() { _ = file.Close() }()

	_, err = io.Copy(file, data)
	if err != nil {
		return "", err
	}

	return fPath, nil
}
