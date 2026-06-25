package row

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/harishtpj/indiansql/internal/schema"
)

type Decoder struct{}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func (d *Decoder) DecodeRow(primaryKey uint64, encoded []byte, tbl *schema.Table) (*Row, error) {
	if tbl == nil {
		return nil, errors.New("table schema is empty!")
	} else if encoded == nil {
		return nil, errors.New("encoded bytes is nil")
	}

	values, err := d.decodeValues(encoded, tbl)
	if err != nil {
		return nil, err
	}
	return NewRow(values, tbl)
}

func (d *Decoder) decodeValues(encoded []byte, tbl *schema.Table) ([]Value, error) {
	if len(encoded) < 2 {
		return nil, errors.New("corrupted row data")
	}

	nVals := binary.BigEndian.Uint16(encoded)
	if nVals != uint16(tbl.ColumnCount()) {
		return nil, fmt.Errorf("column count mismatch: encoded=%d schema=%d", nVals, tbl.ColumnCount())
	}

	offset := 2
	values := make([]Value, nVals)
	for i := range nVals {
		if offset+3 > len(encoded) {
			return nil, errors.New("truncated column header")
		}
		vType := schema.ColumnType(encoded[offset])
		offset++
		byteLen := int(binary.BigEndian.Uint16(encoded[offset:]))
		offset += 2

		if offset+byteLen > len(encoded) {
			return nil, errors.New("truncated column data")
		}
		if vType != tbl.Columns[i].Type {
			return nil, fmt.Errorf(
				"Schema Mismatch for column %s: encoded %s, wanted %s",
				tbl.Columns[i].Name,
				vType.String(),
				tbl.Columns[i].Type.String(),
			)
		}
		switch vType {
		case schema.ColumnTypeNumeric:
			if encoded[offset] == 1 { // IsNull Check
				values[i] = NewNullNumericValue()
			} else {
				values[i] = NewNumericValue(int64(binary.BigEndian.Uint64(encoded[offset+1:])))
			}
		case schema.ColumnTypeVarchar:
			if encoded[offset] == 1 {
				values[i] = NewNullVarcharValue()
			} else {
				values[i] = NewVarcharValue(string(encoded[offset+1 : offset+byteLen]))
			}
		default:
			return nil, errors.New("invalid type supplied for schema")
		}
		offset += byteLen
	}
	return values, nil
}
