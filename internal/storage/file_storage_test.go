package storage

import (
	"os"
	"testing"
)

func TestFileStorage(t *testing.T) {
	tmp, err := os.CreateTemp("", "havensql-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())

	s, err := NewFileStorage(tmp.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	StorageContract(t, s)
}
