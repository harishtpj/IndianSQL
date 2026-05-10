package storage

type MemoryStorage struct {
	data []byte
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}

func (ms *MemoryStorage) ReadAt(b []byte, off int64) (n int, err error) {
	if off >= int64(len(ms.data)) {
		return 0, nil
	}
	n = copy(b, ms.data[off:])
	return n, nil
}

func (ms *MemoryStorage) WriteAt(b []byte, off int64) (n int, err error) {
	// TODO: Handle partial writes
	if l := int64(len(ms.data)); off >= l {
		ms.data = append(ms.data, make([]byte, off-l+1)...)
	}
	n = copy(ms.data[off:], b)
	return n, nil
}

func (ms *MemoryStorage) Size() (int64, error) {
	return int64(len(ms.data)), nil
}

func (ms *MemoryStorage) Close() error {
	return nil
}
