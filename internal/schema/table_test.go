package schema

import (
	"testing"
)

// Test NewTable()
func TestNewTable(t *testing.T) {
	tests := []struct {
		name      string
		tableName string
		columns   []*Column
		shouldErr bool
		errMsg    string
	}{
		{
			name:      "Valid table with single column",
			tableName: "users",
			columns: []*Column{
				{Name: "id", Type: ColumnTypeNumeric, IsPrimaryKey: true},
			},
			shouldErr: false,
		},
		{
			name:      "Valid table with multiple columns",
			tableName: "users",
			columns: []*Column{
				{Name: "id", Type: ColumnTypeNumeric, IsPrimaryKey: true},
				{Name: "name", Type: ColumnTypeVarchar, IsPrimaryKey: false},
				{Name: "age", Type: ColumnTypeNumeric, IsPrimaryKey: false},
			},
			shouldErr: false,
		},
		{
			name:      "Empty table name should error",
			tableName: "",
			columns: []*Column{
				{Name: "id", Type: ColumnTypeNumeric, IsPrimaryKey: true},
			},
			shouldErr: true,
			errMsg:    "name",
		},
		{
			name:      "Empty columns list should error",
			tableName: "users",
			columns:   []*Column{},
			shouldErr: true,
			errMsg:    "columns",
		},
		{
			name:      "No primary key should error",
			tableName: "users",
			columns: []*Column{
				{Name: "id", Type: ColumnTypeNumeric, IsPrimaryKey: false},
				{Name: "name", Type: ColumnTypeVarchar, IsPrimaryKey: false},
			},
			shouldErr: true,
			errMsg:    "primary key",
		},
		{
			name:      "Multiple primary keys should error",
			tableName: "users",
			columns: []*Column{
				{Name: "id", Type: ColumnTypeNumeric, IsPrimaryKey: true},
				{Name: "email", Type: ColumnTypeVarchar, IsPrimaryKey: true},
			},
			shouldErr: true,
			errMsg:    "primary key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table, err := NewTable(tt.tableName, tt.columns)

			if (err != nil) != tt.shouldErr {
				t.Errorf("NewTable() error = %v, shouldErr = %v", err, tt.shouldErr)
			}

			if !tt.shouldErr && table == nil {
				t.Error("NewTable() returned nil table when error was nil")
			}

			if !tt.shouldErr && table != nil {
				if table.Name != tt.tableName {
					t.Errorf("NewTable() table name = %v, want %v", table.Name, tt.tableName)
				}
				if len(table.Columns) != len(tt.columns) {
					t.Errorf("NewTable() column count = %v, want %v", len(table.Columns), len(tt.columns))
				}
			}
		})
	}
}

// Test Table.Validate()
func TestTableValidate(t *testing.T) {
	tests := []struct {
		name      string
		table     *Table
		shouldErr bool
	}{
		{
			name: "Valid table with single column",
			table: &Table{
				Name: "users",
				Columns: []*Column{
					{Name: "id", Type: ColumnTypeNumeric, IsPrimaryKey: true},
				},
				PrimaryKeyIndex: 0,
			},
			shouldErr: false,
		},
		{
			name: "Valid table with multiple columns",
			table: &Table{
				Name: "users",
				Columns: []*Column{
					{Name: "id", Type: ColumnTypeNumeric, IsPrimaryKey: true},
					{Name: "name", Type: ColumnTypeVarchar, IsPrimaryKey: false},
				},
				PrimaryKeyIndex: 0,
			},
			shouldErr: false,
		},
		{
			name: "Empty name should error",
			table: &Table{
				Name: "",
				Columns: []*Column{
					{Name: "id", Type: ColumnTypeNumeric, IsPrimaryKey: true},
				},
				PrimaryKeyIndex: 0,
			},
			shouldErr: true,
		},
		{
			name: "No columns should error",
			table: &Table{
				Name:            "users",
				Columns:         []*Column{},
				PrimaryKeyIndex: 0,
			},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.table.Validate()
			if (err != nil) != tt.shouldErr {
				t.Errorf("Validate() error = %v, shouldErr = %v", err, tt.shouldErr)
			}
		})
	}
}

// Test Table.GetColumn()
func TestTableGetColumn(t *testing.T) {
	table := &Table{
		Name: "users",
		Columns: []*Column{
			{Name: "id", Type: ColumnTypeNumeric, IsPrimaryKey: true},
			{Name: "name", Type: ColumnTypeVarchar, IsPrimaryKey: false},
			{Name: "email", Type: ColumnTypeVarchar, IsPrimaryKey: false},
		},
	}

	tests := []struct {
		name          string
		columnName    string
		expectedFound bool
		expectedType  ColumnType
	}{
		{
			name:          "Find first column",
			columnName:    "id",
			expectedFound: true,
			expectedType:  ColumnTypeNumeric,
		},
		{
			name:          "Find middle column",
			columnName:    "name",
			expectedFound: true,
			expectedType:  ColumnTypeVarchar,
		},
		{
			name:          "Find last column",
			columnName:    "email",
			expectedFound: true,
			expectedType:  ColumnTypeVarchar,
		},
		{
			name:          "Column not found",
			columnName:    "nonexistent",
			expectedFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			col := table.GetColumn(tt.columnName)

			if (col != nil) != tt.expectedFound {
				t.Errorf("GetColumn() found = %v, want %v", col != nil, tt.expectedFound)
			}

			if tt.expectedFound && col != nil && col.Type != tt.expectedType {
				t.Errorf("GetColumn() type = %v, want %v", col.Type, tt.expectedType)
			}
		})
	}
}

// Test Table.GetColumnIndex()
func TestTableGetColumnIndex(t *testing.T) {
	table := &Table{
		Name: "users",
		Columns: []*Column{
			{Name: "id", Type: ColumnTypeNumeric, IsPrimaryKey: true},
			{Name: "name", Type: ColumnTypeVarchar, IsPrimaryKey: false},
			{Name: "email", Type: ColumnTypeVarchar, IsPrimaryKey: false},
		},
	}

	tests := []struct {
		name          string
		columnName    string
		expectedIndex int
	}{
		{
			name:          "Find first column",
			columnName:    "id",
			expectedIndex: 0,
		},
		{
			name:          "Find middle column",
			columnName:    "name",
			expectedIndex: 1,
		},
		{
			name:          "Find last column",
			columnName:    "email",
			expectedIndex: 2,
		},
		{
			name:          "Column not found returns -1",
			columnName:    "nonexistent",
			expectedIndex: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := table.GetColumnIndex(tt.columnName)

			if idx != tt.expectedIndex {
				t.Errorf("GetColumnIndex() = %v, want %v", idx, tt.expectedIndex)
			}
		})
	}
}

// Test Table.GetPrimaryKeyColumn()
func TestTableGetPrimaryKeyColumn(t *testing.T) {
	col1 := &Column{Name: "id", Type: ColumnTypeNumeric, IsPrimaryKey: true}
	col2 := &Column{Name: "name", Type: ColumnTypeVarchar, IsPrimaryKey: false}

	table := &Table{
		Name:            "users",
		Columns:         []*Column{col1, col2},
		PrimaryKeyIndex: 0,
	}

	pkCol := table.GetPrimaryKeyColumn()

	if pkCol == nil {
		t.Fatal("GetPrimaryKeyColumn() returned nil")
	}

	if pkCol.Name != "id" {
		t.Errorf("GetPrimaryKeyColumn() name = %v, want 'id'", pkCol.Name)
	}

	if pkCol.Type != ColumnTypeNumeric {
		t.Errorf("GetPrimaryKeyColumn() type = %v, want NUMERIC", pkCol.Type)
	}

	if !pkCol.IsPrimaryKey {
		t.Error("GetPrimaryKeyColumn() IsPrimaryKey = false, want true")
	}
}

// Test Table.ColumnCount()
func TestTableColumnCount(t *testing.T) {
	tests := []struct {
		name          string
		columnCount   int
		expectedCount int
	}{
		{
			name:          "Single column",
			columnCount:   1,
			expectedCount: 1,
		},
		{
			name:          "Multiple columns",
			columnCount:   5,
			expectedCount: 5,
		},
		{
			name:          "No columns",
			columnCount:   0,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cols := make([]*Column, tt.columnCount)
			for i := range cols {
				isPK := (i == 0)
				cols[i] = &Column{
					Name:         "col" + string(rune(i)),
					Type:         ColumnTypeNumeric,
					IsPrimaryKey: isPK,
				}
			}

			table := &Table{Name: "test", Columns: cols, PrimaryKeyIndex: 0}
			count := table.ColumnCount()

			if count != tt.expectedCount {
				t.Errorf("ColumnCount() = %v, want %v", count, tt.expectedCount)
			}
		})
	}
}

// Test Table.String()
func TestTableString(t *testing.T) {
	table := &Table{
		Name: "users",
		Columns: []*Column{
			{Name: "id", Type: ColumnTypeNumeric, IsPrimaryKey: true},
			{Name: "name", Type: ColumnTypeVarchar, IsPrimaryKey: false},
		},
	}

	str := table.String()

	// String should contain table name
	if len(str) == 0 {
		t.Error("String() returned empty string")
	}

	// Should contain table name
	if !contains(str, "users") {
		t.Errorf("String() should contain table name 'users'")
	}
}

// Helper function for testing
func contains(s, substr string) bool {
	for i := 0; i < len(s); i++ {
		if len(s[i:]) < len(substr) {
			return false
		}
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
