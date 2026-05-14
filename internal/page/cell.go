package page

import (
	"encoding/binary"
	"errors"
)

const (
	valSizeOffset = 8
	valOffset     = valSizeOffset + 2
)

type Cell struct {
	Key   uint64
	Value []byte
}

func (c *Cell) Size() uint16 {
	return uint16(valOffset + len(c.Value))
}

func (c *Cell) Encode(dst []byte) {
	binary.BigEndian.PutUint64(dst, c.Key)
	binary.BigEndian.PutUint16(dst[valSizeOffset:], uint16(len(c.Value)))
	copy(dst[valOffset:], c.Value)
}

func DecodeCell(src []byte) (*Cell, error) {
	if len(src) < valSizeOffset {
		return nil, errors.New("invalid cell size")
	}

	var c Cell
	c.Key = binary.BigEndian.Uint64(src)
	valLen := binary.BigEndian.Uint16(src[valSizeOffset:])
	c.Value = make([]byte, valLen)
	copy(c.Value, src[valOffset:valLen+valOffset])

	return &c, nil
}

func CellSizeRaw(src []byte) uint16 {
	return valOffset + binary.BigEndian.Uint16(src[valSizeOffset:])
}
