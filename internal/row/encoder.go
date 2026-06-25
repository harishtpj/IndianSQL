package row

import (
	"encoding/binary"
	"errors"
	"math"
)

type Encoder struct{}

func NewEncoder() *Encoder {
	return &Encoder{}
}

func (e *Encoder) EncodeRow(row *Row) (uint64, []byte, error) {
	if row == nil {
		return 0, nil, errors.New("given row is nil")
	}
	encoded := binary.BigEndian.AppendUint16(nil, uint16(row.Schema.ColumnCount()))
	for _, val := range row.Values {
		value := *val
		encoded = append(encoded, byte(value.Type()))
		asBytes := value.Bytes()
		encoded = binary.BigEndian.AppendUint16(encoded, uint16(len(asBytes)))
		encoded = append(encoded, asBytes...)
	}

	pk, err := row.PrimaryKeyAsInt64()
	return uint64(pk), encoded, err
}

func (e *Encoder) EncodeRowWithoutPK(row *Row) ([]byte, error) {
	if row == nil {
		return nil, errors.New("given row is nil")
	}

	encoded := binary.BigEndian.AppendUint16(nil, uint16(row.Schema.ColumnCount()))
	for _, val := range row.Values {
		value := *val
		encoded = append(encoded, byte(value.Type()))
		asBytes := value.Bytes()
		if len(asBytes) > math.MaxUint16 {
			return nil, errors.New("value too large")
		}
		encoded = binary.BigEndian.AppendUint16(encoded, uint16(len(asBytes)))
		encoded = append(encoded, asBytes...)
	}
	return encoded, nil
}
