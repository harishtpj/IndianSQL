package storage

import (
	"bytes"
	"io"
	"testing"
)

/*
StorageContract validates REQUIRED behaviour
for ALL storage backends.

Every backend must pass these tests.
*/
func StorageContract(t *testing.T, s Storage) {
	t.Helper()

	t.Run("write_then_read", func(t *testing.T) {
		data := []byte("havensql")

		_, err := s.WriteAt(data, 0)
		if err != nil {
			t.Fatalf("write failed: %v", err)
		}

		buf := make([]byte, len(data))
		_, err = s.ReadAt(buf, 0)
		if err != nil && err != io.EOF {
			t.Fatalf("read failed: %v", err)
		}

		if !bytes.Equal(buf, data) {
			t.Fatalf("data mismatch")
		}
	})

	t.Run("read_beyond_eof_returns_zero", func(t *testing.T) {
		buf := make([]byte, 16)

		_, err := s.ReadAt(buf, 1<<20) // far offset
		if err != nil && err != io.EOF {
			t.Fatalf("unexpected error: %v", err)
		}

		if !bytes.Equal(buf, make([]byte, 16)) {
			t.Fatalf("expected zero-filled buffer")
		}
	})

	t.Run("size_updates_after_write", func(t *testing.T) {
		data := []byte("abc")

		_, err := s.WriteAt(data, 100)
		if err != nil {
			t.Fatal(err)
		}

		size, err := s.Size()
		if err != nil {
			t.Fatal(err)
		}

		if size < 103 {
			t.Fatalf("size not updated")
		}
	})

	t.Run("overwrite_existing_data", func(t *testing.T) {
		_, _ = s.WriteAt([]byte("AAAA"), 0)
		_, _ = s.WriteAt([]byte("BBBB"), 0)

		buf := make([]byte, 4)
		_, _ = s.ReadAt(buf, 0)

		if !bytes.Equal(buf, []byte("BBBB")) {
			t.Fatalf("overwrite failed")
		}
	})

	t.Run("random_offset_writes", func(t *testing.T) {
		offsets := []int64{0, 32, 128, 1024}

		for _, off := range offsets {
			data := []byte{byte(off % 255)}

			_, err := s.WriteAt(data, off)
			if err != nil {
				t.Fatal(err)
			}

			buf := make([]byte, 1)
			_, _ = s.ReadAt(buf, off)

			if buf[0] != data[0] {
				t.Fatalf("random write failed at %d", off)
			}
		}
	})
}
