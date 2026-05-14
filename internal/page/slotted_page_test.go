package page

import "testing"

func newTestPage() *SlottedPage {
	buf := make([]byte, testPageSize)

	h := PageHeader{
		Type:      PageTypeLeaf,
		FreeStart: PageHeaderSize,
		FreeEnd:   testPageSize,
		CellCount: 0,
	}

	h.Encode(buf)

	return WrapSlottedPage(buf)
}

func newCell(k uint64, v string) *Cell {
	return &Cell{
		Key:   k,
		Value: []byte(v),
	}
}

func TestHeaderInitialization(t *testing.T) {
	sp := newTestPage()

	h := sp.Header()

	if h.Type != PageTypeLeaf {
		t.Fatal("wrong page type")
	}

	if h.FreeStart != PageHeaderSize {
		t.Fatalf("invalid FreeStart: %d", h.FreeStart)
	}

	if h.FreeEnd != testPageSize {
		t.Fatalf("invalid FreeEnd: %d", h.FreeEnd)
	}

	if h.CellCount != 0 {
		t.Fatal("cell count must be zero")
	}
}

func TestSlotReadWrite(t *testing.T) {
	sp := newTestPage()

	sp.setSlot(0, 3000)

	if sp.getSlot(0) != 3000 {
		t.Fatal("slot read/write mismatch")
	}
}

func TestInsertSingleCell(t *testing.T) {
	sp := newTestPage()

	cell := newCell(10, "hello")

	if err := sp.InsertCell(cell); err != nil {
		t.Fatal(err)
	}

	h := sp.Header()

	if h.CellCount != 1 {
		t.Fatalf("expected 1 cell, got %d", h.CellCount)
	}

	got, err := DecodeCell(sp.GetCellRaw(0))
	if err != nil {
		t.Fatal(err)
	}

	if got.Key != 10 {
		t.Fatalf("key mismatch: %d", got.Key)
	}

	if string(got.Value) != "hello" {
		t.Fatalf("value mismatch: %s", got.Value)
	}
}

func TestMultipleInsert(t *testing.T) {
	sp := newTestPage()

	for i := 0; i < 50; i++ {
		err := sp.InsertCell(newCell(uint64(i), "x"))
		if err != nil {
			t.Fatal(err)
		}
	}

	h := sp.Header()

	if h.CellCount != 50 {
		t.Fatalf("expected 50 cells, got %d", h.CellCount)
	}

	for i := 0; i < 50; i++ {
		c, err := DecodeCell(sp.GetCellRaw(uint16(i)))
		if err != nil {
			t.Fatal(err)
		}

		if c.Key != uint64(i) {
			t.Fatalf("wrong key at index %d", i)
		}
	}
}

func TestFreeSpaceReduction(t *testing.T) {
	sp := newTestPage()

	before := sp.FreeSpace()

	err := sp.InsertCell(newCell(1, "abc"))
	if err != nil {
		t.Fatal(err)
	}

	after := sp.FreeSpace()

	if after >= before {
		t.Fatal("free space did not shrink")
	}
}

func TestPageFull(t *testing.T) {
	sp := newTestPage()

	large := make([]byte, 200)

	for {
		err := sp.InsertCell(&Cell{
			Key:   1,
			Value: large,
		})

		if err != nil {
			return // success: page became full
		}
	}
}

func TestSlotOffsetsValid(t *testing.T) {
	sp := newTestPage()

	for i := 0; i < 10; i++ {
		sp.InsertCell(newCell(uint64(i), "data"))
	}

	h := sp.Header()

	for i := uint16(0); i < h.CellCount; i++ {
		off := sp.getSlot(i)

		if off < PageHeaderSize || off >= testPageSize {
			t.Fatalf("slot offset out of bounds: %d", off)
		}
	}
}
