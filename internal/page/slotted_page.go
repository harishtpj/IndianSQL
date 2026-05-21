package page

import (
	"encoding/binary"
	"errors"
)

type SlottedPage struct {
	data []byte
}

func WrapSlottedPage(data []byte) *SlottedPage {
	return &SlottedPage{data: data}
}

func NewSlottedPage(data []byte, pgType PageType) *SlottedPage {
	sp := &SlottedPage{data: data}

	h := PageHeader{
		Type:      pgType,
		FreeStart: PageHeaderSize,
		FreeEnd:   uint16(len(data)),
		CellCount: 0,
	}

	sp.SetHeader(h)
	return sp
}

func (sp *SlottedPage) Header() PageHeader {
	pgHdr, _ := DecodePageHeader(sp.data)
	return *pgHdr
}

func (sp *SlottedPage) SetHeader(hdr PageHeader) {
	hdr.Encode(sp.data)
}

func slotOffset(i uint16) uint16 {
	return PageHeaderSize + i*2
}

func (sp *SlottedPage) getSlot(idx uint16) uint16 {
	return binary.BigEndian.Uint16(sp.data[slotOffset(idx):])
}

func (sp *SlottedPage) setSlot(idx uint16, val uint16) {
	binary.BigEndian.PutUint16(sp.data[slotOffset(idx):], val)
}
func (sp *SlottedPage) FreeSpace() uint16 {
	pgHdr := sp.Header()
	return pgHdr.FreeEnd - pgHdr.FreeStart
}

func (sp *SlottedPage) HasSpace(cellSize uint16) bool {
	return 2+cellSize < sp.FreeSpace()
}

func (sp *SlottedPage) writeCellRaw(data []byte) (uint16, error) {
	pgHdr := sp.Header()
	cellSize := uint16(len(data))

	if !sp.HasSpace(cellSize) {
		return 0, errors.New("cell size is too small")
	}

	pgHdr.FreeEnd -= cellSize
	copy(sp.data[pgHdr.FreeEnd:], data)

	sp.SetHeader(pgHdr)
	return pgHdr.FreeEnd, nil
}

func (sp *SlottedPage) insertSlot(idx uint16, off uint16) {
	pgHdr := sp.Header()
	for i := pgHdr.CellCount; i > idx; i-- {
		sp.setSlot(i, sp.getSlot(i-1))
	}
	sp.setSlot(idx, off)
	pgHdr.FreeStart += 2
	sp.SetHeader(pgHdr)
}

func (sp *SlottedPage) appendSlot(off uint16) {
	pgHdr := sp.Header()
	binary.BigEndian.PutUint16(sp.data[pgHdr.FreeStart:], off)
	pgHdr.FreeStart += 2
	sp.SetHeader(pgHdr)
}

func (sp *SlottedPage) InsertCell(idx uint16, c PageCell) error {
	buf := make([]byte, c.Size())
	c.Encode(buf)

	off, err := sp.writeCellRaw(buf)
	if err != nil {
		return err
	}
	sp.insertSlot(idx, off)
	pgHdr := sp.Header()
	pgHdr.CellCount++
	sp.SetHeader(pgHdr)
	return nil
}

func (sp *SlottedPage) AppendCell(c PageCell) error {
	return sp.InsertCell(sp.Header().CellCount, c)
}

func (sp *SlottedPage) GetCellRaw(idx uint16) []byte {
	return sp.data[sp.getSlot(idx):]
}
