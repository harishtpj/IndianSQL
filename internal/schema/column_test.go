package schema

import (
	"testing"
)

// Test ColumnType.String()
func TestColumnTypeString(t *testing.T) {
	tests := []struct {
		name        string
		colType     ColumnType
		expectedStr string
	}{
		{
			name:        "INTEGER type",
			colType:     ColumnTypeInteger,
			expectedStr: "INTEGER",
		},
		{
			name:        "VARCHAR type",
			colType:     ColumnTypeVarchar,
			expectedStr: "VARCHAR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.colType.String()
			if got != tt.expectedStr {
				t.Errorf("String() = %v, want %v", got, tt.expectedStr)
			}
		})
	}
}

// Test ColumnType.IsValid()
func TestColumnTypeIsValid(t *testing.T) {
	tests := []struct {
		name     string
		colType  ColumnType
		expected bool
	}{
		{
			name:     "INTEGER is valid",
			colType:  ColumnTypeInteger,
			expected: true,
		},
		{
			name:     "VARCHAR is valid",
			colType:  ColumnTypeVarchar,
			expected: true,
		},
		{
			name:     "MaxColumnType (sentinel) is invalid",
			colType:  MaxColumnType,
			expected: false,
		},
		{
			name:     "Out of range value is invalid",
			colType:  ColumnType(255),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.colType.IsValid()
			if got != tt.expected {
				t.Errorf("IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Test NewColumn()
func TestNewColumn(t *testing.T) {
	tests := []struct {
		name      string
		colName   string
		colType   ColumnType
		isPK      bool
		shouldErr bool
		errMsg    string
	}{
		{
			name:      "Valid INTEGER column",
			colName:   "id",
			colType:   ColumnTypeInteger,
			isPK:      true,
			shouldErr: false,
		},
		{
			name:      "Valid VARCHAR column",
			colName:   "name",
			colType:   ColumnTypeVarchar,
			isPK:      false,
			shouldErr: false,
		},
		{
			name:      "Empty column name should error",
			colName:   "",
			colType:   ColumnTypeInteger,
			isPK:      false,
			shouldErr: true,
			errMsg:    "name",
		},
		{
			name:      "Invalid column type should error",
			colName:   "col",
			colType:   MaxColumnType,
			isPK:      false,
			shouldErr: true,
			errMsg:    "type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			col, err := NewColumn(tt.colName, tt.colType, tt.isPK)

			if (err != nil) != tt.shouldErr {
				t.Errorf("NewColumn() error = %v, shouldErr = %v", err, tt.shouldErr)
			}

			if !tt.shouldErr && col == nil {
				t.Error("NewColumn() returned nil column when error was nil")
			}

			if !tt.shouldErr && col != nil {
				if col.Name != tt.colName || col.Type != tt.colType || col.IsPrimaryKey != tt.isPK {
					t.Errorf("NewColumn() created column with wrong values")
				}
			}
		})
	}
}

// Test Column.Validate()
func TestColumnValidate(t *testing.T) {
	tests := []struct {
		name      string
		column    *Column
		shouldErr bool
	}{
		{
			name: "Valid INTEGER column",
			column: &Column{
				Name:         "id",
				Type:         ColumnTypeInteger,
				IsPrimaryKey: true,
			},
			shouldErr: false,
		},
		{
			name: "Valid VARCHAR column",
			column: &Column{
				Name:         "name",
				Type:         ColumnTypeVarchar,
				IsPrimaryKey: false,
			},
			shouldErr: false,
		},
		{
			name: "Empty name should error",
			column: &Column{
				Name:         "",
				Type:         ColumnTypeInteger,
				IsPrimaryKey: false,
			},
			shouldErr: true,
		},
		{
			name: "Invalid type should error",
			column: &Column{
				Name:         "col",
				Type:         MaxColumnType,
				IsPrimaryKey: false,
			},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.column.Validate()
			if (err != nil) != tt.shouldErr {
				t.Errorf("Validate() error = %v, shouldErr = %v", err, tt.shouldErr)
			}
		})
	}
}
