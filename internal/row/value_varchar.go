package row

import "github.com/harishtpj/indiansql/internal/schema"

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

func (vv *VarcharValue) String() string {
	return vv.value
}

func (vv *VarcharValue) GetString() string {
	if vv.IsNull() {
		panic("trying to extract string from null value")
	}
	return vv.value
}
