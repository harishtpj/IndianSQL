package page

import (
	"encoding/binary"
	"fmt"

	"github.com/harishtpj/indiansql/internal/apperrors"
)

type InternalCell struct {
	Key   uint64
	Child uint32
}

func (ic *InternalCell) Size() uint16 {
	return 8 + 4
}

func (ic *InternalCell) Encode(dst []byte) {
	binary.BigEndian.PutUint64(dst, ic.Key)
	binary.BigEndian.PutUint32(dst[8:], ic.Child)
}

func DecodeInternalCell(src []byte) (*InternalCell, error) {
	if len(src) < 12 {
		return nil, fmt.Errorf("decode internal cell: %w", apperrors.ErrInvalidCell)
	}

	return &InternalCell{
		Key:   binary.BigEndian.Uint64(src),
		Child: binary.BigEndian.Uint32(src[8:]),
	}, nil
}
