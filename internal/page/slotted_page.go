package page

import (
	"encoding/binary"
	"errors"
)

type SlottedPage struct {
	data []byte
}

func NewSlottedPage(data []byte) *SlottedPage {
	return &SlottedPage{data: data}
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

func (sp *SlottedPage) appendSlot(off uint16) {
	pgHdr := sp.Header()
	binary.BigEndian.PutUint16(sp.data[pgHdr.FreeStart:], off)
	pgHdr.FreeStart += 2
	sp.SetHeader(pgHdr)
}

func (sp *SlottedPage) InsertCell(c *Cell) error {
	buf := make([]byte, c.Size())
	if err := c.Encode(buf); err != nil {
		return err
	}
	off, err := sp.writeCellRaw(buf)
	if err != nil {
		return err
	}
	sp.appendSlot(off)
	pgHdr := sp.Header()
	pgHdr.CellCount++
	sp.SetHeader(pgHdr)
	return nil
}

func (sp *SlottedPage) GetCell(idx uint16) (*Cell, error) {
	return DecodeCell(sp.data[sp.getSlot(idx):])
}
