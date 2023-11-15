package cache

import "errors"

var ErrKeyNotFound = errors.New("key not found")

type Cache interface {
	GetKey(string) (int64, error)
	SetKey(string, int64) error
	FlushData() error
}
