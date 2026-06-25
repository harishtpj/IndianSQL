package row

import (
	"testing"

	"github.com/harishtpj/indiansql/internal/schema"
)

func TestNumericValue(t *testing.T) {
	tests := []struct {
		name string
		val  int64
	}{
		{"Zero", 0},
		{"Positive", 12345},
		{"Negative", -67890},
		{"Max", 9223372036854775807},
		{"Min", -9223372036854775808},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nv := NewNumericValue(tt.val)
			if nv.GetInt64() != tt.val {
				t.Errorf("GetInt64() = %v, want %v", nv.GetInt64(), tt.val)
			}
			if nv.IsNull() {
				t.Error("IsNull() should be false")
			}
			if nv.Type() != schema.ColumnTypeNumeric {
				t.Error("Type() should be NUMERIC")
			}
			if len(nv.Bytes()) != 9 {
				t.Errorf("Bytes() length = %v, want 8", len(nv.Bytes()))
			}
		})
	}
}

func TestNullNumericValue(t *testing.T) {
	nv := NewNullNumericValue()
	if !nv.IsNull() {
		t.Error("IsNull() should be true")
	}
	if nv.Type() != schema.ColumnTypeNumeric {
		t.Error("Type() should still be NUMERIC")
	}
}

func TestVarcharValue(t *testing.T) {
	tests := []struct {
		name string
		val  string
	}{
		{"Empty", ""},
		{"Simple", "hello"},
		{"Long", "this is a longer string with spaces"},
		{"Special", "hello@world#123!"},
		{"Unicode", "こんにちは"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vv := NewVarcharValue(tt.val)
			if vv.GetString() != tt.val {
				t.Errorf("GetString() = %v, want %v", vv.GetString(), tt.val)
			}
			if vv.IsNull() {
				t.Error("IsNull() should be false")
			}
			if vv.Type() != schema.ColumnTypeVarchar {
				t.Error("Type() should be VARCHAR")
			}
		})
	}
}

func TestNullVarcharValue(t *testing.T) {
	vv := NewNullVarcharValue()
	if !vv.IsNull() {
		t.Error("IsNull() should be true")
	}
	if vv.Type() != schema.ColumnTypeVarchar {
		t.Error("Type() should still be VARCHAR")
	}
}

func TestValueInterface(t *testing.T) {
	var v Value

	v = NewNumericValue(42)
	if v.Type() != schema.ColumnTypeNumeric {
		t.Error("NumericValue should implement Value interface")
	}

	v = NewVarcharValue("test")
	if v.Type() != schema.ColumnTypeVarchar {
		t.Error("VarcharValue should implement Value interface")
	}
}

func TestNewRow(t *testing.T) {
	tbl := &schema.Table{
		Name: "test",
		Columns: []*schema.Column{
			{Name: "id", Type: schema.ColumnTypeNumeric, IsPrimaryKey: true},
			{Name: "name", Type: schema.ColumnTypeVarchar, IsPrimaryKey: false},
		},
		PrimaryKeyIndex: 0,
	}

	tests := []struct {
		name      string
		values    []Value
		shouldErr bool
	}{
		{
			name: "Valid row",
			values: []Value{
				NewNumericValue(1),
				NewVarcharValue("Alice"),
			},
			shouldErr: false,
		},
		{
			name: "Wrong number of values",
			values: []Value{
				NewNumericValue(1),
			},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row, err := NewRow(tt.values, tbl)
			if (err != nil) != tt.shouldErr {
				t.Errorf("error = %v, shouldErr = %v", err, tt.shouldErr)
			}
			if !tt.shouldErr && row == nil {
				t.Error("expected row, got nil")
			}
		})
	}
}

func TestRowGetValue(t *testing.T) {
	tbl := &schema.Table{
		Name: "test",
		Columns: []*schema.Column{
			{Name: "id", Type: schema.ColumnTypeNumeric, IsPrimaryKey: true},
			{Name: "name", Type: schema.ColumnTypeVarchar, IsPrimaryKey: false},
		},
		PrimaryKeyIndex: 0,
	}

	row, _ := NewRow([]Value{
		NewNumericValue(1),
		NewVarcharValue("Alice"),
	}, tbl)

	tests := []struct {
		name      string
		colIndex  int
		shouldErr bool
	}{
		{name: "First column", colIndex: 0, shouldErr: false},
		{name: "Second column", colIndex: 1, shouldErr: false},
		{name: "Out of bounds", colIndex: 2, shouldErr: true},
		{name: "Negative", colIndex: -1, shouldErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := row.GetValue(tt.colIndex)
			if (err != nil) != tt.shouldErr {
				t.Errorf("error = %v, shouldErr = %v", err, tt.shouldErr)
			}
			if !tt.shouldErr && v == nil {
				t.Error("expected value, got nil")
			}
		})
	}
}

func TestRowPrimaryKeyValue(t *testing.T) {
	tbl := &schema.Table{
		Name: "test",
		Columns: []*schema.Column{
			{Name: "id", Type: schema.ColumnTypeNumeric, IsPrimaryKey: true},
			{Name: "name", Type: schema.ColumnTypeVarchar, IsPrimaryKey: false},
		},
		PrimaryKeyIndex: 0,
	}

	row, _ := NewRow([]Value{
		NewNumericValue(123),
		NewVarcharValue("Bob"),
	}, tbl)

	_, err := row.PrimaryKeyValue()
	if err != nil {
		t.Fatal(err)
	}

	pkVal, err := row.PrimaryKeyAsInt64()
	if err != nil {
		t.Fatal(err)
	}

	if pkVal != 123 {
		t.Errorf("PrimaryKeyAsInt64() = %v, want 123", pkVal)
	}
}

func TestRowString(t *testing.T) {
	tbl := &schema.Table{
		Name: "test",
		Columns: []*schema.Column{
			{Name: "id", Type: schema.ColumnTypeNumeric, IsPrimaryKey: true},
			{Name: "name", Type: schema.ColumnTypeVarchar, IsPrimaryKey: false},
		},
		PrimaryKeyIndex: 0,
	}

	row, _ := NewRow([]Value{
		NewNumericValue(1),
		NewVarcharValue("Alice"),
	}, tbl)

	str := row.String()
	if len(str) == 0 {
		t.Error("String() returned empty")
	}
}
