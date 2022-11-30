package common

import "errors"

var (
	ErrMetadataNotFound = errors.New("common: metadata not found")

	ErrKeyNotFound = errors.New("common: Key not found")
)
