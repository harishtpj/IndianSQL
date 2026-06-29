package row

import (
	"encoding/binary"
	"math"
	"strconv"

	"github.com/harishtpj/indiansql/internal/schema"
)

type NumericValue struct {
	value float64
	null  bool
}

func NewNumericValue(value float64) *NumericValue {
	return &NumericValue{value: value}
}

func NewNullNumericValue() *NumericValue {
	return &NumericValue{null: true}
}

func (nv *NumericValue) Type() schema.ColumnType {
	return schema.ColumnTypeNumeric
}

func (nv *NumericValue) Bytes() []byte {
	enc := make([]byte, 1)
	if nv.IsNull() {
		enc[0] = 1
	} else {
		enc[0] = 0
	}
	return binary.BigEndian.AppendUint64(enc, math.Float64bits(nv.value))
}

func (nv *NumericValue) IsNull() bool {
	return nv.null
}

func (nv *NumericValue) String() string {
	if nv.IsNull() {
		return "NULL"
	}
	return strconv.FormatFloat(nv.value, 'f', -1, 64)
}

func (nv *NumericValue) GetFloat() float64 {
	if nv.IsNull() {
		panic("trying to extract numeric from null value")
	}
	return nv.value
}
