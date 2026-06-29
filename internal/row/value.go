package row

import (
	"github.com/harishtpj/indiansql/internal/schema"
)

type Value interface {
	Type() schema.ColumnType
	Bytes() []byte
	IsNull() bool
	String() string
}
