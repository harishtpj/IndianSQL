package apperrors

import "errors"

var (
	ErrDuplicateKey = errors.New("duplicate key")
	ErrNodeFull     = errors.New("node full")

	ErrCorruption    = errors.New("database corruption")
	ErrKeyNotFound   = errors.New("key not found")
	ErrInvalidCell   = errors.New("invalid cell")
	ErrInvalidHeader = errors.New("invalid header")

	ErrDBHeaderSmall = errors.New("db header too small")

	ErrInvalidPageType = errors.New("invalid page type")
	ErrCantSplit       = errors.New("node is too small to split")
)
