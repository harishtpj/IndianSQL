package row

import (
	"encoding/binary"
	"strconv"

	"github.com/harishtpj/indiansql/internal/schema"
)

type IntegerValue struct {
	value int64
	null  bool
}

func NewIntegerValue(val int64) *IntegerValue {
	return &IntegerValue{val, false}
}

func NewNullIntegerValue() *IntegerValue {
	return &IntegerValue{null: true}
}

func (iv *IntegerValue) Type() schema.ColumnType {
	return schema.ColumnTypeInteger
}

func (iv *IntegerValue) Bytes() []byte {
	enc := make([]byte, 1)
	if iv.IsNull() {
		enc[0] = 1
	} else {
		enc[0] = 0
	}
	return binary.BigEndian.AppendUint64(enc, uint64(iv.value))
}

func (iv *IntegerValue) IsNull() bool {
	return iv.null
}

func (iv *IntegerValue) String() string {
	if iv.IsNull() {
		return "NULL"
	}
	return strconv.FormatInt(iv.value, 10)
}

func (iv *IntegerValue) GetInt64() int64 {
	if iv.IsNull() {
		panic("trying to extract integer from null value")
	}
	return iv.value
}
