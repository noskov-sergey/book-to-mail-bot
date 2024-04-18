package storage

import "io"

type Storage interface {
	Save(page *Book, data io.ReadCloser) (path string, err error)
}

type Book struct {
	Name string
}
