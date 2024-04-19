package files

import (
	"book-to-mail-bot/lib/e"
	"book-to-mail-bot/storage"
	"io"
	"log"
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

func (s *Storage) Save(book *storage.Book, data io.ReadCloser) (path string, err error) {
	defer func() { err = e.WrapIfErr("can't save book", err) }()

	if err := os.MkdirAll(s.basePath, defaultPerm); err != nil {
		return "", err
	}

	fPath := filepath.Join(s.basePath, book.Name)

	file, err := os.Create(fPath)
	if err != nil {
		return "", err
	}

	defer func() { _ = file.Close() }()

	_, err = io.Copy(file, data)
	if err != nil {
		return "", err
	}

	log.Printf("book - '%s' has been saved to local bin", book.Name)

	return fPath, nil
}
