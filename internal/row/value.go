package row

import (
	"encoding/binary"
	"strconv"

	"github.com/harishtpj/indiansql/internal/schema"
)

type Value interface {
	Type() schema.ColumnType
	Bytes() []byte
	IsNull() bool
	String() string
}

type NumericValue struct {
	value int64
	null  bool
}

func NewNumericValue(val int64) *NumericValue {
	return &NumericValue{val, false}
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
	return binary.BigEndian.AppendUint64(enc, uint64(nv.value))
}

func (nv *NumericValue) IsNull() bool {
	return nv.null
}

type VarcharValue struct {
	value string
	null  bool
}

func NewVarcharValue(val string) *VarcharValue {
	return &VarcharValue{val, false}
}

func NewNullVarcharValue() *VarcharValue {
	return &VarcharValue{null: true}
}

func (vv *VarcharValue) Type() schema.ColumnType {
	return schema.ColumnTypeVarchar
}

func (vv *VarcharValue) Bytes() []byte {
	enc := make([]byte, 1)
	if vv.IsNull() {
		enc[0] = 1
	} else {
		enc[0] = 0
	}
	return append(enc, []byte(vv.value)...)
}

func (vv *VarcharValue) IsNull() bool {
	return vv.null
}

func (nv *NumericValue) String() string {
	return strconv.FormatInt(nv.value, 10)
}

func (vv *VarcharValue) String() string {
	return vv.value
}

func (nv *NumericValue) GetInt64() int64 {
	if nv.IsNull() {
		panic("trying to extract numeric from null value")
	}
	return nv.value
}

func (vv *VarcharValue) GetString() string {
	if vv.IsNull() {
		panic("trying to extract string from null value")
	}
	return vv.value
}
