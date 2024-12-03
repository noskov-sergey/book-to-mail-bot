package files

import (
	"io"
	"os"
	"path/filepath"

	"go.uber.org/zap"

	"github.com/noskov-sergey/book-to-mail-bot/lib/e"
	"github.com/noskov-sergey/book-to-mail-bot/storage"
)

const (
	defaultPerm = 0774
)

type Storage struct {
	basePath string

	log *zap.Logger
}

func New(basePath string, log *zap.Logger) *Storage {
	return &Storage{
		basePath: basePath,
		log:      log.Named("data saver"),
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

	s.log.Info("Has been saved to local bin", zap.String("book name", book.Name))

	return fPath, nil
}
