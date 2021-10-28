package storage

import "container/list"

type Index interface {
	Get() (map[string]*list.List, error)
	Set(index map[string]*list.List) error
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
