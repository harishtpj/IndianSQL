package page

import (
	"encoding/binary"
	"fmt"

	"github.com/harishtpj/indiansql/internal/apperrors"
	"github.com/harishtpj/indiansql/internal/consts"
)

const (
	freeStartOffset  = 1
	freeEndOffset    = freeStartOffset + 2
	cellCountOffset  = freeEndOffset + 2
	rightChildOffset = cellCountOffset + 2
	parentPageOffset = rightChildOffset + 4
	PageHeaderSize   = parentPageOffset + 4
)

type PageHeader struct {
	Type       PageType
	FreeStart  uint16
	FreeEnd    uint16
	CellCount  uint16
	RightChild uint32
	ParentPage uint32
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
	binary.BigEndian.PutUint32(dst[parentPageOffset:], h.ParentPage)
	return nil
}

func DecodePageHeader(src []byte) (hdr *PageHeader, err error) {
	hdr = &PageHeader{}

	if len(src) < PageHeaderSize {
		err = fmt.Errorf("decode page header: %w", apperrors.ErrInvalidHeader)
		return
	}

	if hdr.Type = PageType(src[0]); hdr.Type >= MaxPageType {
		err = fmt.Errorf("decode page header: invalid page type %d: %w", hdr.Type, apperrors.ErrInvalidHeader)
		return
	}

	hdr.FreeStart = binary.BigEndian.Uint16(src[freeStartOffset:])
	if hdr.FreeStart < PageHeaderSize {
		err = fmt.Errorf("decode page header: free start %d within header: %w", hdr.FreeStart, apperrors.ErrInvalidHeader)
		return
	}

	hdr.FreeEnd = binary.BigEndian.Uint16(src[freeEndOffset:])
	if hdr.FreeStart > hdr.FreeEnd {
		err = fmt.Errorf("decode page header: free start %d beyond free end %d: %w", hdr.FreeStart, hdr.FreeEnd, apperrors.ErrInvalidHeader)
		return
	} else if hdr.FreeEnd > uint16(len(src)) {
		err = fmt.Errorf("decode page header: free end %d beyond buffer size %d: %w", hdr.FreeEnd, len(src), apperrors.ErrInvalidHeader)
		return
	}

	hdr.CellCount = binary.BigEndian.Uint16(src[cellCountOffset:])
	hdr.RightChild = binary.BigEndian.Uint32(src[rightChildOffset:])
	hdr.ParentPage = binary.BigEndian.Uint32(src[parentPageOffset:])
	return
}

func (h *PageHeader) Parent() uint32 {
	return h.ParentPage
}

func (h *PageHeader) SetParent(id uint32) {
	h.ParentPage = id
}
