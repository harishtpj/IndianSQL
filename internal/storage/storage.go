package storage

type Storage interface {
	ReadAt(b []byte, off int64) (n int, err error)
	WriteAt(b []byte, off int64) (n int, err error)
	Size() (int64, error)
	Close() error
}
