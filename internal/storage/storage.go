package storage

import (
	"github.com/ythosa/bendy/internal/index"
)

type Index interface {
	Get() (index.InvertIndex, error)
	Set(index index.InvertIndex) error
}

type Files interface {
	Get() ([]string, error)
	Put(filename string) error
	Delete(filename string) error
}

type Storage struct {
	Index
	Files
}
