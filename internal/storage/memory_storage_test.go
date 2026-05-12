package storage

import "testing"

func TestMemoryStorage(t *testing.T) {
	s := NewMemoryStorage()
	defer s.Close()

	StorageContract(t, s)
}
