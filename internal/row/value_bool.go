package row

import "github.com/harishtpj/indiansql/internal/schema"

type BoolValue struct {
	value bool
	null  bool
}

func NewBoolValue(val bool) *BoolValue {
	return &BoolValue{value: val}
}

func NewNullBoolValue() *BoolValue {
	return &BoolValue{null: true}
}

func (bv *BoolValue) Type() schema.ColumnType {
	return schema.ColumnTypeBoolean
}

func (bv *BoolValue) IsNull() bool {
	return bv.null
}

func (bv *BoolValue) Bytes() []byte {
	enc := make([]byte, 1)
	if bv.IsNull() {
		enc[0] |= 1 << 0
	}
	if bv.value {
		enc[0] |= 1 << 1
	}
	return enc
}

func (bv *BoolValue) String() string {
	if bv.IsNull() {
		return "NULL"
	}
	if bv.value {
		return "TRUE"
	}
	return "FALSE"
}

func (bv *BoolValue) GetBool() bool {
	if bv.IsNull() {
		panic("trying to extract boolean from null value")
	}
	return bv.value
}

func (bv *BoolValue) GetInt() byte {
	if bv.IsNull() {
		panic("trying to extract boolean from null value")
	}
	if bv.value {
		return 1
	}
	return 0
}
