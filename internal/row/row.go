package row

import (
	"errors"
	"fmt"
	"strings"

	"github.com/harishtpj/indiansql/internal/schema"
)

type Row struct {
	Values []*Value
	Schema *schema.Table
}

func NewRow(values []Value, tbl *schema.Table) (*Row, error) {
	if tbl == nil {
		return nil, errors.New("table schema is empty!")
	}

	if len(values) != tbl.ColumnCount() {
		return nil, fmt.Errorf("insufficient values: required %d, got %d", tbl.ColumnCount(), len(values))
	}

	valsPtr := make([]*Value, len(values))
	for i, val := range values {
		valsPtr[i] = &val
	}
	return &Row{valsPtr, tbl}, nil
}

func (r *Row) GetValue(colIndex int) (Value, error) {
	if colIndex < 0 || colIndex >= len(r.Values) {
		return nil, errors.New("invalid column index provided")
	}
	return *r.Values[colIndex], nil
}

func (r *Row) GetValueByName(colName string) (Value, error) {
	colIdx := r.Schema.GetColumnIndex(colName)
	if colIdx == -1 {
		return nil, errors.New("invalid column name provided")
	}
	return r.GetValue(colIdx)
}

func (r *Row) PrimaryKeyValue() (Value, error) {
	return r.GetValue(r.Schema.PrimaryKeyIndex)
}

func (r *Row) PrimaryKeyAsInt64() (int64, error) {
	pkVal, err := r.PrimaryKeyValue()
	if err != nil {
		return 0, err
	}

	if pkVal.Type() != schema.ColumnTypeInteger {
		return 0, fmt.Errorf("expected integer primary key, got %v", pkVal.Type())
	}

	val, ok := pkVal.(*IntegerValue)
	if !ok {
		return 0, errors.New("primary key is not integer")
	}

	return val.GetInt64(), nil
}

func (r *Row) String() string {
	var repr strings.Builder
	fmt.Fprintf(&repr, "(")
	for i, val := range r.Values {
		if (*val).Type() == schema.ColumnTypeVarchar {
			fmt.Fprintf(&repr, "'%s'", (*val).String())
		} else {
			fmt.Fprint(&repr, (*val).String())
		}
		if i+1 != r.Schema.ColumnCount() {
			repr.WriteString(", ")
		}
	}
	return repr.String() + ")"
}
