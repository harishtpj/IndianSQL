package storage

import "io"

type MemoryStorage struct {
	data []byte
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}

func (ms *MemoryStorage) ReadAt(b []byte, off int64) (n int, err error) {
	if size, _ := ms.Size(); off >= size {
		return 0, io.EOF
	}
	n = copy(b, ms.data[off:])
	if n < len(b) {
		return n, io.EOF
	}
	return n, nil
}

func (ms *MemoryStorage) WriteAt(b []byte, off int64) (n int, err error) {
	end := off + int64(len(b))
	size, _ := ms.Size()
	if end > size {
		diff := end - size
		ms.data = append(ms.data, make([]byte, diff)...)
	}
	n = copy(ms.data[off:end], b)
	return n, nil
}

func (ms *MemoryStorage) Size() (int64, error) {
	return int64(len(ms.data)), nil
}

func (ms *MemoryStorage) Close() error {
	return nil
}
