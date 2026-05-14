package page

import (
	"bytes"
	"testing"
)

func TestCellSizeMatchesEncoding(t *testing.T) {
	c := &Cell{
		Key:   42,
		Value: []byte("indiansql"),
	}

	sz := c.Size()

	buf := make([]byte, sz)

	c.Encode(buf)

	if uint16(len(buf)) != sz {
		t.Fatalf("encoded size mismatch")
	}
}

func TestCellEncodeDecodeRoundTrip(t *testing.T) {
	orig := &Cell{
		Key:   999,
		Value: []byte("database-engine"),
	}

	sz := orig.Size()
	buf := make([]byte, sz)

	orig.Encode(buf)

	decoded, err := DecodeCell(buf)
	if err != nil {
		t.Fatal(err)
	}

	if decoded.Key != orig.Key {
		t.Fatalf("key mismatch")
	}

	if !bytes.Equal(decoded.Value, orig.Value) {
		t.Fatalf("value mismatch")
	}
}

func TestCellSizeRawMatchesSize(t *testing.T) {
	c := &Cell{
		Key:   7,
		Value: []byte("raw-size-test"),
	}

	expected := c.Size()

	buf := make([]byte, expected)
	c.Encode(buf)

	raw := CellSizeRaw(buf)

	if raw != expected {
		t.Fatalf("raw size mismatch got %d expected %d", raw, expected)
	}
}

func TestMultipleCellsSequentialParsing(t *testing.T) {
	cells := []*Cell{
		{Key: 1, Value: []byte("a")},
		{Key: 2, Value: []byte("bbbb")},
		{Key: 3, Value: []byte("cccccccc")},
	}

	total := 0
	for _, c := range cells {
		s := c.Size()
		total += int(s)
	}

	page := make([]byte, total)

	off := 0
	for _, c := range cells {
		s := c.Size()
		c.Encode(page[off:])
		off += int(s)
	}

	off = 0
	for i := range cells {
		sz := CellSizeRaw(page[off:])

		decoded, err := DecodeCell(page[off : off+int(sz)])
		if err != nil {
			t.Fatal(err)
		}

		if decoded.Key != cells[i].Key {
			t.Fatalf("key mismatch at %d", i)
		}

		if !bytes.Equal(decoded.Value, cells[i].Value) {
			t.Fatalf("value mismatch at %d", i)
		}

		off += int(sz)
	}
}

func TestCellEmptyValue(t *testing.T) {
	c := &Cell{
		Key:   100,
		Value: nil,
	}

	sz := c.Size()

	buf := make([]byte, sz)

	c.Encode(buf)

	decoded, err := DecodeCell(buf)
	if err != nil {
		t.Fatal(err)
	}

	if decoded.Key != c.Key {
		t.Fatalf("key mismatch")
	}

	if len(decoded.Value) != 0 {
		t.Fatalf("expected empty value")
	}
}

func TestCellCorruptedInput(t *testing.T) {
	_, err := DecodeCell([]byte{1, 2, 3})
	if err == nil {
		t.Fatalf("expected error for corrupted cell")
	}
}
