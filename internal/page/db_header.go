package page

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/harishtpj/indiansql/internal/consts"
)

// DB Header is Fixed to 100 Bytes
const (
	magicStrLen         = len(consts.MagicStr)
	versionOffset       = magicStrLen
	pageSizeOffset      = versionOffset + 2
	pageCountOffset     = pageSizeOffset + 4
	freeListHeadOffset  = pageCountOffset + 4
	freeListCountOffset = freeListHeadOffset + 4
	rootPageOffset      = freeListCountOffset + 4
	reservedOffset      = rootPageOffset + 4
	reservedSize        = 54
)

type DBHeader struct {
	Magic         [magicStrLen]byte
	Version       uint16
	PageSize      uint32
	PageCount     uint32
	FreeListHead  uint32
	FreeListCount uint32
	RootPage      uint32
	Reserved      [reservedSize]byte
}

func InitDBHeader(page []byte, pgSize uint32) error {
	hdr := DBHeader{
		Version:   consts.VersionNum,
		PageSize:  pgSize,
		PageCount: 1,
	}

	copy(hdr.Magic[:], consts.MagicStr)
	return hdr.Encode(page)
}

func (h *DBHeader) Encode(dst []byte) error {
	copy(dst, h.Magic[:])
	binary.BigEndian.PutUint16(dst[versionOffset:], h.Version)
	binary.BigEndian.PutUint32(dst[pageSizeOffset:], h.PageSize)
	binary.BigEndian.PutUint32(dst[pageCountOffset:], h.PageCount)
	binary.BigEndian.PutUint32(dst[freeListHeadOffset:], h.FreeListHead)
	binary.BigEndian.PutUint32(dst[freeListCountOffset:], h.FreeListCount)
	binary.BigEndian.PutUint32(dst[rootPageOffset:], h.RootPage)
	copy(dst[reservedOffset:], h.Reserved[:])
	return nil
}

func DecodeDBHeader(src []byte) (*DBHeader, error) {
	if len(src) < reservedOffset+reservedSize {
		return nil, errors.New("DB Header too small")
	}

	var hdr DBHeader

	copy(hdr.Magic[:], src[:magicStrLen])

	if !bytes.Equal(hdr.Magic[:], []byte(consts.MagicStr)) {
		return nil, errors.New("corrupted File! Magic String not found")
	}

	hdr.Version = binary.BigEndian.Uint16(src[versionOffset:])

	if hdr.Version != consts.VersionNum {
		return nil, errors.New("unsupported Version")
	}

	hdr.PageSize = binary.BigEndian.Uint32(src[pageSizeOffset:])
	hdr.PageCount = binary.BigEndian.Uint32(src[pageCountOffset:])
	hdr.FreeListHead = binary.BigEndian.Uint32(src[freeListHeadOffset:])
	hdr.FreeListCount = binary.BigEndian.Uint32(src[freeListCountOffset:])
	hdr.RootPage = binary.BigEndian.Uint32(src[rootPageOffset:])

	copy(hdr.Reserved[:], src[reservedOffset:reservedOffset+reservedSize])

	return &hdr, nil
}
