package page

import (
	"encoding/binary"
	"errors"
)

type InternalCell struct {
	Key   uint64
	Value uint32
}

func (ic *InternalCell) Size() uint16 {
	return 8 + 4
}

func (ic *InternalCell) Encode(dst []byte) {
	binary.BigEndian.PutUint64(dst, ic.Key)
	binary.BigEndian.PutUint32(dst[8:], ic.Value)
}

func DecodeInternalCell(src []byte) (*InternalCell, error) {
	if len(src) < 12 {
		return nil, errors.New("invalid internal cell size!")
	}

	return &InternalCell{
		Key:   binary.BigEndian.Uint64(src),
		Value: binary.BigEndian.Uint32(src[8:]),
	}, nil
}
