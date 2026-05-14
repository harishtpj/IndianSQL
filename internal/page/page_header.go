package page

import (
	"encoding/binary"
	"errors"

	"github.com/harishtpj/indiansql/internal/consts"
)

const (
	freeStartOffset  = 1
	freeEndOffset    = freeStartOffset + 2
	cellCountOffset  = freeEndOffset + 2
	rightChildOffset = cellCountOffset + 2
	PageHeaderSize   = rightChildOffset + 4
)

type PageHeader struct {
	Type       PageType
	FreeStart  uint16
	FreeEnd    uint16
	CellCount  uint16
	RightChild uint32
}

func InitPage(page []byte, tp PageType) error {
	pgHdr := PageHeader{
		Type:      tp,
		FreeStart: PageHeaderSize,
		FreeEnd:   consts.PageSize,
	}

	return pgHdr.Encode(page)
}

func (h *PageHeader) Encode(dst []byte) error {
	dst[0] = uint8(h.Type)
	binary.BigEndian.PutUint16(dst[freeStartOffset:], h.FreeStart)
	binary.BigEndian.PutUint16(dst[freeEndOffset:], h.FreeEnd)
	binary.BigEndian.PutUint16(dst[cellCountOffset:], h.CellCount)
	binary.BigEndian.PutUint32(dst[rightChildOffset:], h.RightChild)
	return nil
}

func DecodePageHeader(src []byte) (hdr *PageHeader, err error) {
	hdr = &PageHeader{}

	if len(src) < PageHeaderSize {
		err = errors.New("Page Header too small")
		return
	}

	if hdr.Type = PageType(src[0]); hdr.Type >= MaxPageType {
		err = errors.New("Invalid Page Type")
		return
	}

	hdr.FreeStart = binary.BigEndian.Uint16(src[freeStartOffset:])
	if hdr.FreeStart < PageHeaderSize {
		err = errors.New("Page Start within Header!")
		return
	}

	hdr.FreeEnd = binary.BigEndian.Uint16(src[freeEndOffset:])
	if hdr.FreeStart > hdr.FreeEnd {
		err = errors.New("Page Start is out of bound!")
		return
	} else if hdr.FreeEnd > uint16(len(src)) {
		err = errors.New("Page End is out of bounds!")
		return
	}

	hdr.CellCount = binary.BigEndian.Uint16(src[cellCountOffset:])
	hdr.RightChild = binary.BigEndian.Uint32(src[rightChildOffset:])
	return
}
