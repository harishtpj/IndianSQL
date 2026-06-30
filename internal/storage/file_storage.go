package storage

import "os"

type FileStorage struct {
	file *os.File
}

func NewFileStorage(path string) (*FileStorage, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return &FileStorage{file}, nil
}

func (fs *FileStorage) ReadAt(b []byte, off int64) (n int, err error) {
	return fs.file.ReadAt(b, off)
}

func (fs *FileStorage) WriteAt(b []byte, off int64) (n int, err error) {
	return fs.file.WriteAt(b, off)
}

func (fs *FileStorage) Size() (int64, error) {
	info, err := fs.file.Stat()
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

func (fs *FileStorage) Flush() error {
	return fs.file.Sync()
}

func (fs *FileStorage) Close() error {
	return fs.file.Close()
}
