package row

import (
	"testing"

	"github.com/harishtpj/indiansql/internal/schema"
)

func TestEncodeRow(t *testing.T) {
	encoder := NewEncoder()
	tbl := &schema.Table{
		Name: "users",
		Columns: []*schema.Column{
			{Name: "id", Type: schema.ColumnTypeNumeric, IsPrimaryKey: true},
			{Name: "name", Type: schema.ColumnTypeVarchar, IsPrimaryKey: false},
		},
		PrimaryKeyIndex: 0,
	}

	row, _ := NewRow([]Value{
		NewNumericValue(42),
		NewVarcharValue("Alice"),
	}, tbl)

	pk, encoded, err := encoder.EncodeRow(row)
	if err != nil {
		t.Fatal(err)
	}

	if pk != 42 {
		t.Errorf("PrimaryKey = %v, want 42", pk)
	}

	if len(encoded) == 0 {
		t.Error("Encoded data is empty")
	}
}

func TestEncodeRowVariousValues(t *testing.T) {
	encoder := NewEncoder()
	tbl := &schema.Table{
		Name: "test",
		Columns: []*schema.Column{
			{Name: "id", Type: schema.ColumnTypeNumeric, IsPrimaryKey: true},
			{Name: "value", Type: schema.ColumnTypeNumeric, IsPrimaryKey: false},
			{Name: "text", Type: schema.ColumnTypeVarchar, IsPrimaryKey: false},
		},
		PrimaryKeyIndex: 0,
	}

	tests := []struct {
		name string
		pk   int64
		val  int64
		text string
	}{
		{"Simple", 1, 100, "hello"},
		{"Long text", 2, 200, "this is a longer string"},
		{"Empty string", 3, 300, ""},
		{"Special chars", 4, 400, "!@#$%^&*()"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row, _ := NewRow([]Value{
				NewNumericValue(tt.pk),
				NewNumericValue(tt.val),
				NewVarcharValue(tt.text),
			}, tbl)

			pk, encoded, err := encoder.EncodeRow(row)
			if err != nil {
				t.Fatal(err)
			}

			if pk != uint64(tt.pk) {
				t.Errorf("PrimaryKey = %v, want %v", pk, tt.pk)
			}

			if len(encoded) == 0 {
				t.Error("Encoded is empty")
			}
		})
	}
}

func TestDecodeRow(t *testing.T) {
	encoder := NewEncoder()
	decoder := NewDecoder()
	tbl := &schema.Table{
		Name: "users",
		Columns: []*schema.Column{
			{Name: "id", Type: schema.ColumnTypeNumeric, IsPrimaryKey: true},
			{Name: "name", Type: schema.ColumnTypeVarchar, IsPrimaryKey: false},
		},
		PrimaryKeyIndex: 0,
	}

	original, _ := NewRow([]Value{
		NewNumericValue(42),
		NewVarcharValue("Alice"),
	}, tbl)

	pk, encoded, _ := encoder.EncodeRow(original)

	decoded, err := decoder.DecodeRow(pk, encoded, tbl)
	if err != nil {
		t.Fatal(err)
	}

	if decoded == nil {
		t.Fatal("Decoded row is nil")
	}

	v1, _ := decoded.GetValue(0)
	if v1.Type() != schema.ColumnTypeNumeric {
		t.Error("First column should be NUMERIC")
	}

	v2, _ := decoded.GetValue(1)
	if v2.Type() != schema.ColumnTypeVarchar {
		t.Error("Second column should be VARCHAR")
	}
}

func TestEncodeDecodeRoundTrip(t *testing.T) {
	encoder := NewEncoder()
	decoder := NewDecoder()
	tbl := &schema.Table{
		Name: "products",
		Columns: []*schema.Column{
			{Name: "id", Type: schema.ColumnTypeNumeric, IsPrimaryKey: true},
			{Name: "price", Type: schema.ColumnTypeNumeric, IsPrimaryKey: false},
			{Name: "title", Type: schema.ColumnTypeVarchar, IsPrimaryKey: false},
		},
		PrimaryKeyIndex: 0,
	}

	original, _ := NewRow([]Value{
		NewNumericValue(999),
		NewNumericValue(19999),
		NewVarcharValue("Laptop"),
	}, tbl)

	pk, encoded, _ := encoder.EncodeRow(original)
	decoded, _ := decoder.DecodeRow(pk, encoded, tbl)

	if decoded == nil {
		t.Fatal("Decoded is nil")
	}

	v1, _ := decoded.GetValue(0)
	nv1 := v1.(*NumericValue)
	if nv1.GetInt64() != 999 {
		t.Errorf("Column 0 = %v, want 999", nv1.GetInt64())
	}

	v2, _ := decoded.GetValue(1)
	nv2 := v2.(*NumericValue)
	if nv2.GetInt64() != 19999 {
		t.Errorf("Column 1 = %v, want 19999", nv2.GetInt64())
	}

	v3, _ := decoded.GetValue(2)
	vv3 := v3.(*VarcharValue)
	if vv3.GetString() != "Laptop" {
		t.Errorf("Column 2 = %v, want Laptop", vv3.GetString())
	}
}

func TestDecodeRowInvalidData(t *testing.T) {
	decoder := NewDecoder()
	tbl := &schema.Table{
		Name: "test",
		Columns: []*schema.Column{
			{Name: "id", Type: schema.ColumnTypeNumeric, IsPrimaryKey: true},
		},
		PrimaryKeyIndex: 0,
	}

	tests := []struct {
		name      string
		pk        uint64
		encoded   []byte
		shouldErr bool
	}{
		{"Empty data", 1, []byte{}, true},
		{"Truncated data", 1, []byte{0x01}, true},
		{"Valid minimal", 1, []byte{
			0x00, 0x01, // column count = 1

			0x00,       // ColumnTypeNumeric (assuming Numeric == 0)
			0x00, 0x09, // length = 9 bytes

			0x00, // not null

			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x04, // uint64(4)
		}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := decoder.DecodeRow(tt.pk, tt.encoded, tbl)
			if (err != nil) != tt.shouldErr {
				t.Errorf("error = %v, shouldErr = %v", err, tt.shouldErr)
			}
		})
	}
}
